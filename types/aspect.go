// Package types provides shared type definitions for the Reve SDK.
package types

// AspectRatio represents supported image aspect ratios.
type AspectRatio string

// Supported aspect ratios.
const (
	Ratio16x9 AspectRatio = "16:9" // Landscape widescreen
	Ratio9x16 AspectRatio = "9:16" // Portrait (mobile)
	Ratio3x2  AspectRatio = "3:2"  // Classic photo landscape
	Ratio2x3  AspectRatio = "2:3"  // Classic photo portrait
	Ratio4x3  AspectRatio = "4:3"  // Standard landscape
	Ratio3x4  AspectRatio = "3:4"  // Standard portrait
	Ratio1x1  AspectRatio = "1:1"  // Square
	RatioAuto AspectRatio = "auto" // Auto-detect
)

// String returns the string representation.
func (r AspectRatio) String() string {
	return string(r)
}

// Valid returns true if the aspect ratio is valid.
func (r AspectRatio) Valid() bool {
	switch r {
	case Ratio16x9, Ratio9x16, Ratio3x2, Ratio2x3,
		Ratio4x3, Ratio3x4, Ratio1x1, RatioAuto, "":
		return true
	}
	return false
}

// Dimensions returns approximate width:height for the ratio.
func (r AspectRatio) Dimensions() (int, int) {
	switch r {
	case Ratio16x9:
		return 16, 9
	case Ratio9x16:
		return 9, 16
	case Ratio3x2:
		return 3, 2
	case Ratio2x3:
		return 2, 3
	case Ratio4x3:
		return 4, 3
	case Ratio3x4:
		return 3, 4
	case Ratio1x1:
		return 1, 1
	default:
		return 0, 0
	}
}
