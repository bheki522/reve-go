package reve_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	reve "github.com/shamspias/reve-go"
	"github.com/shamspias/reve-go/image"
	"github.com/shamspias/reve-go/internal/transport"
	"github.com/shamspias/reve-go/internal/validator"
	"github.com/shamspias/reve-go/types"
)

func TestNewClient(t *testing.T) {
	client := reve.NewClient("test-key")
	if client == nil {
		t.Fatal("NewClient returned nil")
	}
	if client.Images == nil {
		t.Fatal("Images service is nil")
	}
}

func TestClientOptions(t *testing.T) {
	client := reve.NewClient("test-key",
		reve.WithTimeout(30*time.Second),
		reve.WithRetry(5, time.Second, 30*time.Second),
		reve.WithUserAgent("test-agent"),
		reve.WithDebug(true),
	)

	cfg := client.Config()
	if cfg.Timeout != 30*time.Second {
		t.Errorf("Expected timeout 30s, got %v", cfg.Timeout)
	}
	if cfg.MaxRetries != 5 {
		t.Errorf("Expected 5 retries, got %d", cfg.MaxRetries)
	}
	if cfg.UserAgent != "test-agent" {
		t.Errorf("Expected test-agent, got %s", cfg.UserAgent)
	}
}

func TestAspectRatio(t *testing.T) {
	tests := []struct {
		ratio types.AspectRatio
		valid bool
	}{
		{types.Ratio16x9, true},
		{types.Ratio9x16, true},
		{types.Ratio3x2, true},
		{types.Ratio2x3, true},
		{types.Ratio4x3, true},
		{types.Ratio3x4, true},
		{types.Ratio1x1, true},
		{types.RatioAuto, true},
		{"", true},
		{"invalid", false},
	}

	for _, tt := range tests {
		if tt.ratio.Valid() != tt.valid {
			t.Errorf("AspectRatio(%s).Valid() = %v, want %v", tt.ratio, tt.ratio.Valid(), tt.valid)
		}
	}
}

func TestOutputFormat(t *testing.T) {
	tests := []struct {
		format types.OutputFormat
		ext    string
	}{
		{types.FormatPNG, ".png"},
		{types.FormatJPEG, ".jpeg"},
		{types.FormatWebP, ".webp"},
		{types.FormatJSON, ".png"},
	}

	for _, tt := range tests {
		if tt.format.Extension() != tt.ext {
			t.Errorf("Extension() = %s, want %s", tt.format.Extension(), tt.ext)
		}
	}
}

func TestDetectFormat(t *testing.T) {
	tests := []struct {
		path string
		want types.OutputFormat
	}{
		{"image.png", types.FormatPNG},
		{"image.jpg", types.FormatJPEG},
		{"image.webp", types.FormatWebP},
		{"image.unknown", types.FormatPNG},
	}

	for _, tt := range tests {
		if got := types.DetectFormat(tt.path); got != tt.want {
			t.Errorf("DetectFormat(%s) = %v, want %v", tt.path, got, tt.want)
		}
	}
}

func TestRef(t *testing.T) {
	tests := []struct {
		index int
		want  string
	}{
		{0, "<img>0</img>"},
		{1, "<img>1</img>"},
		{5, "<img>5</img>"},
	}

	for _, tt := range tests {
		if got := types.Ref(tt.index); got != tt.want {
			t.Errorf("Ref(%d) = %s, want %s", tt.index, got, tt.want)
		}
	}
}

func TestImage(t *testing.T) {
	data := []byte("test image data")
	img := types.NewImage(data)

	encoded := img.Base64()
	if encoded == "" {
		t.Error("Base64() returned empty")
	}

	decoded, err := img.Bytes()
	if err != nil {
		t.Fatalf("Bytes() error: %v", err)
	}
	if string(decoded) != string(data) {
		t.Errorf("Bytes() = %s, want %s", decoded, data)
	}
}

func TestCreateParamsValidation(t *testing.T) {
	tests := []struct {
		name    string
		params  *image.CreateParams
		wantErr error
	}{
		{"valid", &image.CreateParams{Prompt: "test"}, nil},
		{"empty prompt", &image.CreateParams{}, validator.ErrEmptyPrompt},
		{"too long", &image.CreateParams{Prompt: strings.Repeat("a", 2561)}, validator.ErrPromptTooLong},
		{"invalid ratio", &image.CreateParams{Prompt: "test", AspectRatio: "bad"}, validator.ErrInvalidAspectRatio},
		{"invalid scaling", &image.CreateParams{Prompt: "test", TestTimeScaling: 20}, validator.ErrInvalidScaling},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.params.Validate()
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Validate() = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestEditParamsValidation(t *testing.T) {
	tests := []struct {
		name    string
		params  *image.EditParams
		wantErr error
	}{
		{"valid", &image.EditParams{Instruction: "test", ReferenceImage: "base64"}, nil},
		{"empty instruction", &image.EditParams{ReferenceImage: "base64"}, validator.ErrEmptyInstruction},
		{"empty image", &image.EditParams{Instruction: "test"}, validator.ErrEmptyReferenceImage},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.params.Validate()
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Validate() = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestRemixParamsValidation(t *testing.T) {
	tests := []struct {
		name    string
		params  *image.RemixParams
		wantErr error
	}{
		{"valid", &image.RemixParams{Prompt: "test", ReferenceImages: []string{"img1"}}, nil},
		{"empty prompt", &image.RemixParams{ReferenceImages: []string{"img1"}}, validator.ErrEmptyPrompt},
		{"no images", &image.RemixParams{Prompt: "test"}, validator.ErrNoReferenceImages},
		{"too many", &image.RemixParams{Prompt: "test", ReferenceImages: make([]string, 7)}, validator.ErrTooManyReferenceImages},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.params.Validate()
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Validate() = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/v1/image/create" {
			t.Errorf("Expected /v1/image/create, got %s", r.URL.Path)
		}

		resp := types.Result{
			Image:       "base64data",
			Version:     "test-version",
			CreditsUsed: 18,
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := reve.NewClient("test-key", reve.WithBaseURL(server.URL), reve.WithNoRetry())

	result, err := client.Images.Create(context.Background(), &image.CreateParams{
		Prompt: "test",
	})

	if err != nil {
		t.Fatalf("Create() error: %v", err)
	}
	if result.Image != "base64data" {
		t.Errorf("Image = %s, want base64data", result.Image)
	}
	if result.CreditsUsed != 18 {
		t.Errorf("CreditsUsed = %d, want 18", result.CreditsUsed)
	}
}

func TestAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error_code": "MISSING_REQUIRED_PARAMETER",
			"message":    "Prompt is required",
		})
	}))
	defer server.Close()

	client := reve.NewClient("test-key", reve.WithBaseURL(server.URL), reve.WithNoRetry())

	_, err := client.Images.Create(context.Background(), &image.CreateParams{Prompt: "test"})

	if err == nil {
		t.Fatal("Expected error")
	}

	var apiErr *transport.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("Expected APIError, got %T", err)
	}
	if apiErr.StatusCode != http.StatusBadRequest {
		t.Errorf("StatusCode = %d, want 400", apiErr.StatusCode)
	}
}

func TestAPIErrorMethods(t *testing.T) {
	tests := []struct {
		name   string
		err    *transport.APIError
		method func(*transport.APIError) bool
		want   bool
	}{
		{"IsRateLimit", &transport.APIError{Code: transport.ErrCodeRateLimit}, (*transport.APIError).IsRateLimit, true},
		{"IsInsufficientFunds", &transport.APIError{Code: transport.ErrCodeInsufficientFunds}, (*transport.APIError).IsInsufficientFunds, true},
		{"IsContentViolation", &transport.APIError{Code: transport.ErrCodeContentViolation}, (*transport.APIError).IsContentViolation, true},
		{"IsAuthError", &transport.APIError{Code: transport.ErrCodeInvalidAPIKey}, (*transport.APIError).IsAuthError, true},
		{"Retryable 429", &transport.APIError{StatusCode: 429}, (*transport.APIError).Retryable, true},
		{"Retryable 500", &transport.APIError{StatusCode: 500}, (*transport.APIError).Retryable, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.method(tt.err); got != tt.want {
				t.Errorf("%s() = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestRetry(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < 3 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		json.NewEncoder(w).Encode(types.Result{Image: "success"})
	}))
	defer server.Close()

	client := reve.NewClient("test-key",
		reve.WithBaseURL(server.URL),
		reve.WithRetry(3, 10*time.Millisecond, 100*time.Millisecond),
	)

	result, err := client.Images.Create(context.Background(), &image.CreateParams{Prompt: "test"})

	if err != nil {
		t.Fatalf("Create() error: %v", err)
	}
	if result.Image != "success" {
		t.Errorf("Image = %s, want success", result.Image)
	}
	if attempts != 3 {
		t.Errorf("attempts = %d, want 3", attempts)
	}
}

func TestCostEstimation(t *testing.T) {
	cost := image.EstimateCreate(1, nil)
	if cost.BaseCredits != 18 {
		t.Errorf("BaseCredits = %d, want 18", cost.BaseCredits)
	}

	cost = image.EstimateEdit(false, 1, nil)
	if cost.BaseCredits != 30 {
		t.Errorf("Edit BaseCredits = %d, want 30", cost.BaseCredits)
	}

	cost = image.EstimateEdit(true, 1, nil)
	if cost.BaseCredits != 5 {
		t.Errorf("Edit Fast BaseCredits = %d, want 5", cost.BaseCredits)
	}

	cost = image.EstimateCreate(2, nil)
	if cost.TotalCredits != 36 {
		t.Errorf("Scaled TotalCredits = %d, want 36", cost.TotalCredits)
	}
}

func TestBatchHelpers(t *testing.T) {
	results := []image.BatchResult{
		{Index: 0, Result: &types.Result{Image: "1"}, Error: nil},
		{Index: 1, Result: nil, Error: validator.ErrEmptyPrompt},
		{Index: 2, Result: &types.Result{Image: "2"}, Error: nil},
	}

	if image.SuccessCount(results) != 2 {
		t.Errorf("SuccessCount = %d, want 2", image.SuccessCount(results))
	}
	if image.ErrorCount(results) != 1 {
		t.Errorf("ErrorCount = %d, want 1", image.ErrorCount(results))
	}
	if len(image.Successful(results)) != 2 {
		t.Errorf("Successful len = %d, want 2", len(image.Successful(results)))
	}
	if len(image.Errors(results)) != 1 {
		t.Errorf("Errors len = %d, want 1", len(image.Errors(results)))
	}
}

func BenchmarkCreate(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(types.Result{Image: "test"})
	}))
	defer server.Close()

	client := reve.NewClient("test-key", reve.WithBaseURL(server.URL), reve.WithNoRetry())
	ctx := context.Background()
	params := &image.CreateParams{Prompt: "test"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = client.Images.Create(ctx, params)
	}
}
