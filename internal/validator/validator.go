// Package validator provides request validation.
package validator

import "errors"

// Validation errors.
var (
	ErrEmptyPrompt            = errors.New("prompt cannot be empty")
	ErrPromptTooLong          = errors.New("prompt exceeds 2560 characters")
	ErrEmptyInstruction       = errors.New("edit instruction cannot be empty")
	ErrEmptyReferenceImage    = errors.New("reference image cannot be empty")
	ErrNoReferenceImages      = errors.New("at least one reference image required")
	ErrTooManyReferenceImages = errors.New("maximum 6 reference images allowed")
	ErrInvalidAspectRatio     = errors.New("invalid aspect ratio")
	ErrInvalidUpscaleFactor   = errors.New("upscale factor must be 2, 3, or 4")
	ErrInvalidScaling         = errors.New("test time scaling must be 1-15")
)

// Constants
const (
	MaxPromptLength    = 2560
	MaxReferenceImages = 6
	MinScaling         = 1.0
	MaxScaling         = 15.0
)

// ValidatePrompt validates a prompt string.
func ValidatePrompt(prompt string) error {
	if prompt == "" {
		return ErrEmptyPrompt
	}
	if len(prompt) > MaxPromptLength {
		return ErrPromptTooLong
	}
	return nil
}

// ValidateInstruction validates an edit instruction.
func ValidateInstruction(instruction string) error {
	if instruction == "" {
		return ErrEmptyInstruction
	}
	if len(instruction) > MaxPromptLength {
		return ErrPromptTooLong
	}
	return nil
}

// ValidateReferenceImage validates a single reference image.
func ValidateReferenceImage(image string) error {
	if image == "" {
		return ErrEmptyReferenceImage
	}
	return nil
}

// ValidateReferenceImages validates multiple reference images.
func ValidateReferenceImages(images []string) error {
	if len(images) == 0 {
		return ErrNoReferenceImages
	}
	if len(images) > MaxReferenceImages {
		return ErrTooManyReferenceImages
	}
	return nil
}

// ValidateAspectRatio validates an aspect ratio string.
func ValidateAspectRatio(ratio string) error {
	if ratio == "" {
		return nil
	}
	valid := map[string]bool{
		"16:9": true, "9:16": true,
		"3:2": true, "2:3": true,
		"4:3": true, "3:4": true,
		"1:1": true, "auto": true,
	}
	if !valid[ratio] {
		return ErrInvalidAspectRatio
	}
	return nil
}

// ValidateUpscaleFactor validates an upscale factor.
func ValidateUpscaleFactor(factor int) error {
	if factor < 2 || factor > 4 {
		return ErrInvalidUpscaleFactor
	}
	return nil
}

// ValidateScaling validates test time scaling.
func ValidateScaling(scaling float64) error {
	if scaling == 0 {
		return nil
	}
	if scaling < MinScaling || scaling > MaxScaling {
		return ErrInvalidScaling
	}
	return nil
}
