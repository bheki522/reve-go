package types

import (
	"encoding/base64"
	"os"
)

// Result represents an image generation result.
type Result struct {
	// Image is the base64 encoded image data.
	Image string `json:"image"`

	// Version is the model version used.
	Version string `json:"version"`

	// ContentViolation indicates if content policy was violated.
	ContentViolation bool `json:"content_violation"`

	// RequestID is the unique request identifier.
	RequestID string `json:"request_id"`

	// CreditsUsed is the number of credits consumed.
	CreditsUsed int `json:"credits_used"`

	// CreditsRemaining is the remaining credits.
	CreditsRemaining int `json:"credits_remaining"`
}

// Bytes returns the raw image bytes.
func (r *Result) Bytes() ([]byte, error) {
	return base64.StdEncoding.DecodeString(r.Image)
}

// SaveTo saves the image to a file.
//
// Example:
//
//	result, _ := client.Images.Create(ctx, params)
//	err := result.SaveTo("output.png")
func (r *Result) SaveTo(path string) error {
	data, err := r.Bytes()
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// RawResult represents a raw binary response.
type RawResult struct {
	// Data is the raw image bytes.
	Data []byte

	// ContentType is the MIME type.
	ContentType string

	// Version is the model version used.
	Version string

	// ContentViolation indicates if content policy was violated.
	ContentViolation bool

	// RequestID is the unique request identifier.
	RequestID string

	// CreditsUsed is the number of credits consumed.
	CreditsUsed int

	// CreditsRemaining is the remaining credits.
	CreditsRemaining int
}

// SaveTo saves the raw image to a file.
func (r *RawResult) SaveTo(path string) error {
	return os.WriteFile(path, r.Data, 0644)
}

// Size returns the size in bytes.
func (r *RawResult) Size() int {
	return len(r.Data)
}
