// Package transport provides HTTP transport functionality for the Reve API.
package transport

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/net/proxy"
)

// Client handles HTTP communication with the Reve API.
type Client struct {
	httpClient *http.Client
	baseURL    string
	apiKey     string
	userAgent  string
	debug      bool
	logger     Logger
	retrier    *Retrier
}

// Logger is a function type for logging.
type Logger func(format string, args ...any)

// Config holds transport configuration.
type Config struct {
	BaseURL      string
	APIKey       string
	UserAgent    string
	Timeout      time.Duration
	MaxRetries   int
	RetryMinWait time.Duration
	RetryMaxWait time.Duration
	Debug        bool
	Logger       Logger
	Transport    http.RoundTripper
}

// New creates a new transport client.
func New(cfg *Config) *Client {
	httpClient := &http.Client{
		Timeout: cfg.Timeout,
	}

	if cfg.Transport != nil {
		httpClient.Transport = cfg.Transport
	}

	return &Client{
		httpClient: httpClient,
		baseURL:    cfg.BaseURL,
		apiKey:     cfg.APIKey,
		userAgent:  cfg.UserAgent,
		debug:      cfg.Debug,
		logger:     cfg.Logger,
		retrier:    NewRetrier(cfg.MaxRetries, cfg.RetryMinWait, cfg.RetryMaxWait),
	}
}

// Request represents an API request.
type Request struct {
	Method     string
	Path       string
	Body       any
	Accept     string
	Breadcrumb string
}

// Response represents a JSON response.
type Response struct {
	Body      []byte
	Status    int
	RequestID string
}

// RawResponse represents a binary response.
type RawResponse struct {
	Data             []byte
	ContentType      string
	Version          string
	ContentViolation bool
	RequestID        string
	CreditsUsed      int
	CreditsRemaining int
}

// Do executes a request and returns JSON response.
func (c *Client) Do(ctx context.Context, req *Request) (*Response, error) {
	return c.retrier.Do(ctx, func() (*Response, error) {
		return c.execute(ctx, req)
	})
}

// DoRaw executes a request and returns raw binary response.
func (c *Client) DoRaw(ctx context.Context, req *Request) (*RawResponse, error) {
	return c.retrier.DoRaw(ctx, func() (*RawResponse, error) {
		return c.executeRaw(ctx, req)
	})
}

func (c *Client) execute(ctx context.Context, req *Request) (*Response, error) {
	httpReq, err := c.buildRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	c.log("Request: %s %s", httpReq.Method, httpReq.URL)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, &RequestError{Op: "http", Err: err}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &RequestError{Op: "read response", Err: err}
	}

	c.log("Response: status=%d", resp.StatusCode)

	if resp.StatusCode >= 400 {
		return nil, ParseError(resp, body)
	}

	return &Response{
		Body:      body,
		Status:    resp.StatusCode,
		RequestID: resp.Header.Get("X-Reve-Request-Id"),
	}, nil
}

func (c *Client) executeRaw(ctx context.Context, req *Request) (*RawResponse, error) {
	httpReq, err := c.buildRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	c.log("Request (raw): %s %s", httpReq.Method, httpReq.URL)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, &RequestError{Op: "http", Err: err}
	}
	defer resp.Body.Close()

	if errCode := resp.Header.Get("X-Reve-Error-Code"); errCode != "" {
		body, _ := io.ReadAll(resp.Body)
		return nil, ParseError(resp, body)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &RequestError{Op: "read response", Err: err}
	}

	c.log("Response (raw): status=%d, size=%d", resp.StatusCode, len(data))

	return &RawResponse{
		Data:             data,
		ContentType:      resp.Header.Get("Content-Type"),
		Version:          resp.Header.Get("X-Reve-Version"),
		ContentViolation: resp.Header.Get("X-Reve-Content-Violation") == "true",
		RequestID:        resp.Header.Get("X-Reve-Request-Id"),
		CreditsUsed:      parseIntHeader(resp, "X-Reve-Credits-Used"),
		CreditsRemaining: parseIntHeader(resp, "X-Reve-Credits-Remaining"),
	}, nil
}

func (c *Client) buildRequest(ctx context.Context, req *Request) (*http.Request, error) {
	url := c.baseURL + req.Path
	if req.Breadcrumb != "" {
		url += "?breadcrumb=" + req.Breadcrumb
	}

	var bodyReader io.Reader
	var getBody func() (io.ReadCloser, error)

	if req.Body != nil {
		data, err := json.Marshal(req.Body)
		if err != nil {
			return nil, &RequestError{Op: "marshal", Err: err}
		}
		bodyReader = bytes.NewReader(data)
		bodyData := data
		getBody = func() (io.ReadCloser, error) {
			return io.NopCloser(bytes.NewReader(bodyData)), nil
		}
	}

	httpReq, err := http.NewRequestWithContext(ctx, req.Method, url, bodyReader)
	if err != nil {
		return nil, &RequestError{Op: "create request", Err: err}
	}

	httpReq.GetBody = getBody
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("User-Agent", c.userAgent)

	accept := "application/json"
	if req.Accept != "" {
		accept = req.Accept
	}
	httpReq.Header.Set("Accept", accept)

	return httpReq, nil
}

func (c *Client) log(format string, args ...any) {
	if !c.debug {
		return
	}
	if c.logger != nil {
		c.logger(format, args...)
	} else {
		fmt.Printf("[reve] "+format+"\n", args...)
	}
}

func parseIntHeader(resp *http.Response, key string) int {
	val := resp.Header.Get(key)
	if val == "" {
		return 0
	}
	var n int
	fmt.Sscanf(val, "%d", &n)
	return n
}

// CreateHTTPProxyTransport creates a transport with HTTP proxy.
func CreateHTTPProxyTransport(proxyURL string) (http.RoundTripper, error) {
	parsed, err := url.Parse(proxyURL)
	if err != nil {
		return nil, err
	}

	return &http.Transport{
		Proxy: http.ProxyURL(parsed),
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		TLSClientConfig:       &tls.Config{MinVersion: tls.VersionTLS12},
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}, nil
}

// CreateSOCKS5ProxyTransport creates a transport with SOCKS5 proxy.
func CreateSOCKS5ProxyTransport(addr, username, password string) (http.RoundTripper, error) {
	var auth *proxy.Auth
	if username != "" {
		auth = &proxy.Auth{User: username, Password: password}
	}

	dialer, err := proxy.SOCKS5("tcp", addr, auth, proxy.Direct)
	if err != nil {
		return nil, err
	}

	return &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return dialer.Dial(network, addr)
		},
		TLSClientConfig:       &tls.Config{MinVersion: tls.VersionTLS12},
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}, nil
}

// CreateEnvProxyTransport creates a transport using environment proxy settings.
func CreateEnvProxyTransport() http.RoundTripper {
	return &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		TLSClientConfig:       &tls.Config{MinVersion: tls.VersionTLS12},
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
}
