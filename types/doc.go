// Package types provides shared type definitions for the Reve SDK.
//
// This package contains all the common types used across the SDK:
//   - AspectRatio: Image aspect ratios (16:9, 9:16, etc.)
//   - ModelVersion: Model versions (latest, fast, specific versions)
//   - OutputFormat: Response formats (JSON, PNG, JPEG, WebP)
//   - Postprocess: Post-processing operations (upscale, remove background)
//   - Image: Image handling utilities
//   - Result: API response types
//
// # Usage
//
//	import "github.com/shamspias/reve-go/types"
//
//	// Aspect ratios
//	ratio := types.Ratio16x9
//
//	// Model versions
//	version := types.VersionLatestFast
//
//	// Postprocessing
//	pp := types.Upscale(2)
//
//	// Image handling
//	img, _ := types.NewImageFromFile("photo.png")
//	base64 := img.Base64()
//
//	// Reference tags for remix
//	prompt := fmt.Sprintf("Apply %s to %s", types.Ref(0), types.Ref(1))
package types
