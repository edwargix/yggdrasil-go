package version

var buildName string
var buildVersion string

type BuildInfo interface {
	BuildName() string
	BuildVersion() string
}

type BuildTimeMetadata struct{}

// BuildName gets the current build name. This is usually injected if built
// from git, or returns "unknown" otherwise.
func (b *BuildTimeMetadata) BuildName() string {
	if buildName == "" {
		return "unknown"
	}
	return buildName
}

// BuildVersion gets the current build version. This is usually injected if
// built from git, or returns "unknown" otherwise.
func (b *BuildTimeMetadata) BuildVersion() string {
	if buildVersion == "" {
		return "unknown"
	}
	return buildVersion
}
