// Package version provides the identity server version information.
package version

const (
	appName        = "identity-server"
	appDisplayName = "Identity Server"
	appCopyright   = "Copyright \u00A9 2023-2024 energimind.com"
)

//nolint:gochecknoglobals
var (
	appVersion = "v0.0.0"
	appBuildID = "devel"
)

// Version represents the identity server version information.
type Version struct {
	Name        string
	DisplayName string
	Copyright   string
	Signature   string
	Version     string
	BuildID     string
}

// Get returns the identity server version information.
func Get() Version {
	return Version{
		Name:        appName,
		DisplayName: appDisplayName,
		Copyright:   appCopyright,
		Signature:   appDisplayName + " " + appVersion + ", build " + appBuildID,
		Version:     appVersion,
		BuildID:     appBuildID,
	}
}
