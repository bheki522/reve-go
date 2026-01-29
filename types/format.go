package types

import (
	"path/filepath"
	"strings"
)

// OutputFormat represents the response format.
type OutputFormat string

// Supported output formats.
const (
	FormatJSON OutputFormat = "application/json"
	FormatPNG  OutputFormat = "image/png"
	FormatJPEG OutputFormat = "image/jpeg"
	FormatWebP OutputFormat = "image/webp"
)

// String returns the string representation.
func (f OutputFormat) String() string {
	return string(f)
}

// ContentType returns the MIME type.
func (f OutputFormat) ContentType() string {
	return string(f)
}

// Extension returns the file extension.
func (f OutputFormat) Extension() string {
	switch f {
	case FormatPNG:
		return ".png"
	case FormatJPEG:
		return ".jpeg"
	case FormatWebP:
		return ".webp"
	default:
		return ".png"
	}
}

// DetectFormat detects format from file path.
func DetectFormat(path string) OutputFormat {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".png":
		return FormatPNG
	case ".jpg", ".jpeg":
		return FormatJPEG
	case ".webp":
		return FormatWebP
	default:
		return FormatPNG
	}
}
