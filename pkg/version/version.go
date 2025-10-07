// Package version provides version information for gocroissant.
//
// This package contains build-time version information including version number,
// build time, and Git commit hash. The version information is typically set during
// the build process using Go build flags.
//
// # Usage
//
//	fmt.Printf("Version: %s\n", version.Version)
//	stamp := version.RetrieveStamp()
//	fmt.Printf("  Built with %s on %s\n", stamp.InfoGoCompiler, stamp.InfoBuildTime)
//
// The GetBuildInfo function provides detailed build information including
// version, build time, Git commit, and Go version used for compilation.
package version

import (
	"fmt"
	"os"
	"runtime/debug"
)

// These variables should be set at compile time.
//
// nolint:gochecknoglobals
var (
	// AppName is the name of the application.
	AppName = "gocroissant"

	// Version is the service version.
	Version = "dev"
)

type Stamp struct {
	InfoGoVersion  string
	InfoGoCompiler string
	InfoGOARCH     string
	InfoGOOS       string
	InfoBuildTime  string
	VCSRevision    string
}

func RetrieveStamp() *Stamp {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		fmt.Printf("Error: could not read build info")
		os.Exit(1)
	}
	//nolint:exhaustruct
	stamp := Stamp{}
	for _, setting := range info.Settings {
		switch setting.Key {
		case "info.goVersion":
			stamp.InfoGoVersion = setting.Value
		case "-compiler":
			stamp.InfoGoCompiler = setting.Value
		case "GOARCH":
			stamp.InfoGOARCH = setting.Value
		case "GOOS":
			stamp.InfoGOOS = setting.Value
		case "vcs.time":
			stamp.InfoBuildTime = setting.Value
		case "vcs.revision":
			stamp.VCSRevision = setting.Value
		}
	}

	return &stamp
}
