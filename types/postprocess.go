package types

// ProcessType represents postprocessing operation types.
type ProcessType string

// Supported postprocessing types.
const (
	ProcessUpscale          ProcessType = "upscale"
	ProcessRemoveBackground ProcessType = "remove_background"
)

// Postprocess represents a postprocessing operation.
type Postprocess struct {
	Process       ProcessType `json:"process"`
	UpscaleFactor int         `json:"upscale_factor,omitempty"`
}

// Upscale creates an upscale postprocessing operation.
//
// Example:
//
//	params := &image.CreateParams{
//		Prompt:      "A sunset",
//		Postprocess: []types.Postprocess{types.Upscale(2)},
//	}
func Upscale(factor int) Postprocess {
	return Postprocess{
		Process:       ProcessUpscale,
		UpscaleFactor: factor,
	}
}

// RemoveBackground creates a background removal operation.
//
// Example:
//
//	params := &image.CreateParams{
//		Prompt:      "A product photo",
//		Postprocess: []types.Postprocess{types.RemoveBackground()},
//	}
func RemoveBackground() Postprocess {
	return Postprocess{
		Process: ProcessRemoveBackground,
	}
}

// Validate validates the postprocessing operation.
func (p Postprocess) Validate() error {
	if p.Process == ProcessUpscale {
		if p.UpscaleFactor < 2 || p.UpscaleFactor > 4 {
			return ErrInvalidUpscale{}
		}
	}
	return nil
}

// ErrInvalidUpscale is returned for invalid upscale factors.
type ErrInvalidUpscale struct{}

func (e ErrInvalidUpscale) Error() string {
	return "upscale factor must be 2, 3, or 4"
}
