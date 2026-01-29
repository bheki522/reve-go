package transport

import (
	"context"
	"math"
	"math/rand/v2"
	"net/http"
	"time"
)

// Retrier handles retry logic with exponential backoff.
type Retrier struct {
	maxRetries int
	minWait    time.Duration
	maxWait    time.Duration
}

// NewRetrier creates a new retrier.
func NewRetrier(maxRetries int, minWait, maxWait time.Duration) *Retrier {
	return &Retrier{
		maxRetries: maxRetries,
		minWait:    minWait,
		maxWait:    maxWait,
	}
}

// Do executes a function with retry logic.
func (r *Retrier) Do(ctx context.Context, fn func() (*Response, error)) (*Response, error) {
	var lastErr error

	for attempt := 0; attempt <= r.maxRetries; attempt++ {
		if attempt > 0 {
			if err := r.wait(ctx, attempt); err != nil {
				return nil, err
			}
		}

		resp, err := fn()
		if err == nil {
			return resp, nil
		}

		lastErr = err

		if !r.shouldRetry(err) {
			return nil, err
		}

		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
	}

	return nil, lastErr
}

// DoRaw executes a function returning raw response with retry logic.
func (r *Retrier) DoRaw(ctx context.Context, fn func() (*RawResponse, error)) (*RawResponse, error) {
	var lastErr error

	for attempt := 0; attempt <= r.maxRetries; attempt++ {
		if attempt > 0 {
			if err := r.wait(ctx, attempt); err != nil {
				return nil, err
			}
		}

		resp, err := fn()
		if err == nil {
			return resp, nil
		}

		lastErr = err

		if !r.shouldRetry(err) {
			return nil, err
		}

		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
	}

	return nil, lastErr
}

func (r *Retrier) wait(ctx context.Context, attempt int) error {
	backoff := r.calculateBackoff(attempt)
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(backoff):
		return nil
	}
}

func (r *Retrier) calculateBackoff(attempt int) time.Duration {
	backoff := float64(r.minWait) * math.Pow(2, float64(attempt-1))
	jitter := backoff * 0.25 * (rand.Float64()*2 - 1)
	backoff += jitter

	if backoff > float64(r.maxWait) {
		backoff = float64(r.maxWait)
	}

	return time.Duration(backoff)
}

func (r *Retrier) shouldRetry(err error) bool {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.Retryable()
	}
	return false
}

// isRetryableStatus checks if HTTP status code is retryable.
func isRetryableStatus(code int) bool {
	switch code {
	case http.StatusTooManyRequests,
		http.StatusInternalServerError,
		http.StatusBadGateway,
		http.StatusServiceUnavailable,
		http.StatusGatewayTimeout:
		return true
	}
	return false
}
