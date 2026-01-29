package types

import (
	"encoding/base64"
	"fmt"
	"os"
)

// Image represents an image for API operations.
type Image struct {
	data   []byte
	base64 string
}

// NewImage creates an Image from raw bytes.
//
// Example:
//
//	data, _ := os.ReadFile("photo.png")
//	img := types.NewImage(data)
func NewImage(data []byte) *Image {
	return &Image{data: data}
}

// NewImageFromBase64 creates an Image from base64 string.
//
// Example:
//
//	img := types.NewImageFromBase64("iVBORw0KGgo...")
func NewImageFromBase64(encoded string) *Image {
	return &Image{base64: encoded}
}

// NewImageFromFile loads an Image from a file.
//
// Example:
//
//	img, err := types.NewImageFromFile("photo.png")
func NewImageFromFile(path string) (*Image, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read image: %w", err)
	}
	return &Image{data: data}, nil
}

// Bytes returns the raw image bytes.
func (img *Image) Bytes() ([]byte, error) {
	if len(img.data) > 0 {
		return img.data, nil
	}
	if img.base64 != "" {
		return base64.StdEncoding.DecodeString(img.base64)
	}
	return nil, fmt.Errorf("image is empty")
}

// Base64 returns the base64 encoded image.
func (img *Image) Base64() string {
	if img.base64 != "" {
		return img.base64
	}
	return base64.StdEncoding.EncodeToString(img.data)
}

// SaveTo saves the image to a file.
func (img *Image) SaveTo(path string) error {
	data, err := img.Bytes()
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// Size returns the size in bytes.
func (img *Image) Size() int {
	if len(img.data) > 0 {
		return len(img.data)
	}
	return len(img.base64) * 3 / 4
}

// Ref creates an image reference tag for remix prompts.
//
// Example:
//
//	prompt := fmt.Sprintf("Apply style from %s to %s", types.Ref(0), types.Ref(1))
func Ref(index int) string {
	return fmt.Sprintf("<img>%d</img>", index)
}
