package version

// These variables should be set at compile time
var (
	// AppName is the name of the application
	AppName = "gocroissant"

	// Version is the service version
	Version = "dev"

	// GitHash is the hash of git commit the service is built from
	GitHash = "dev"

	// BuildTime build time in RFC3339 format
	BuildTime = "now"
)
