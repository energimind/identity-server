package server

const (
	realmName        = "identity-server"
	realmDisplayName = "Identity Server"
	realmCopyright   = "Copyright \u00A9 2023-2024 energimind.com"
)

//nolint:gochecknoglobals
var (
	realmVersion = "v0.0.0"
	realmBuildID = "devel"
	Version      = VersionInfo{
		Name:        realmName,
		DisplayName: realmDisplayName,
		Copyright:   realmCopyright,
		Signature:   realmDisplayName + " " + realmVersion + ", build " + realmBuildID,
		Version:     realmVersion,
		BuildID:     realmBuildID,
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
