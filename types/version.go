package types

// ModelVersion represents the model version to use.
type ModelVersion string

// Supported model versions.
const (
	// VersionLatest uses the most recent model.
	VersionLatest ModelVersion = "latest"

	// VersionLatestFast uses the fast variant.
	VersionLatestFast ModelVersion = "latest-fast"

	// Create versions
	VersionCreate20250915 ModelVersion = "reve-create@20250915"

	// Edit versions
	VersionEdit20250915     ModelVersion = "reve-edit@20250915"
	VersionEditFast20251030 ModelVersion = "reve-edit-fast@20251030"

	// Remix versions
	VersionRemix20250915     ModelVersion = "reve-remix@20250915"
	VersionRemixFast20251030 ModelVersion = "reve-remix-fast@20251030"
)

// String returns the string representation.
func (v ModelVersion) String() string {
	return string(v)
}

// IsFast returns true if this is a fast model variant.
func (v ModelVersion) IsFast() bool {
	return v == VersionLatestFast ||
		v == VersionEditFast20251030 ||
		v == VersionRemixFast20251030
}
