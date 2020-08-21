package version

var (
	// ProviderVersion is set during the release process to the release version of the binary
	// ProviderVersion is a version string populated by the build using -ldflags "-X
	// ${PKG}/pkg/version.Version=${VERSION}".
	ProviderVersion = "dev"
	// GitCommit is the latest git commit hash populated by the build using
	// -ldflags "-X ${PKG}/pkg/version.GitCommit=${GIT_COMMIT}".
	GitCommit = "UNKNOWN"
)
