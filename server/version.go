package server

const (
	appName        = "identity-server"
	appDisplayName = "Identity Server"
	appCopyright   = "Copyright \u00A9 2023-2024 energimind.com"
)

//nolint:gochecknoglobals
var (
	appVersion = "v0.0.0"
	appBuildID = "devel"
	Version    = VersionInfo{
		Name:        appName,
		DisplayName: appDisplayName,
		Copyright:   appCopyright,
		Signature:   appDisplayName + " " + appVersion + ", build " + appBuildID,
		Version:     appVersion,
		BuildID:     appBuildID,
	}
)

// VersionInfo represents the identity server version information.
type VersionInfo struct {
	Name        string
	DisplayName string
	Copyright   string
	Signature   string
	Version     string
	BuildID     string
}
