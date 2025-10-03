// Package cmd provides the command-line interface for gocroissant.
//
// Gocroissant is a Go implementation for working with the ML Commons Croissant
// metadata format - a standardized way to describe machine learning datasets using JSON-LD.
//
// The command-line tool provides functionality to:
//   - Generate Croissant metadata from CSV files with automatic type inference
//   - Validate existing Croissant metadata files for specification compliance
//   - Compare metadata files for schema compatibility
//   - Analyze CSV file structure and display column information
//   - Display version and build information
//
// # Command Reference
//
// Generate metadata with default output path:
//
//	gocroissant generate data.csv
//
// Generate metadata with custom output path:
//
//	gocroissant generate data.csv -o output.jsonld
//
// Generate and validate metadata:
//
//	gocroissant generate data.csv -o metadata.jsonld --validate
//
// Validate existing metadata:
//
//	gocroissant validate metadata.jsonld
//
// Compare two metadata files for compatibility:
//
//	gocroissant match reference.jsonld candidate.jsonld
//
// Analyze CSV file structure:
//
//	gocroissant info data.csv --sample-size 20
//
// Show version information:
//
//	gocroissant version
//
// # Features
//
// Metadata Generation:
//   - Automatic data type inference from CSV content
//   - SHA-256 hash calculation for file verification
//   - Configurable output paths and validation options
//   - Support for environment variable configuration
//
// Validation:
//   - JSON-LD structure validation
//   - Croissant specification compliance checking
//   - Configurable validation modes (standard, strict)
//   - Optional file existence and URL accessibility checking
//
// Schema Comparison:
//   - Field-by-field compatibility analysis
//   - Intelligent type compatibility (numeric type flexibility)
//   - Support for schema evolution (additional fields allowed)
//   - Detailed reporting of matches, mismatches, and missing fields
//
// File Analysis:
//
//   - CSV structure validation and statistics
//
//   - Column type inference with configurable sample sizes
//
//   - File size and row count analysis
//
//     gocroissant version
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/beyondcivic/gocroissant/pkg/croissant"
	"github.com/beyondcivic/gocroissant/pkg/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var RootCmd = &cobra.Command{
	Use:   "gocroissant",
	Short: "Croissant metadata tools",
	Long: `A Go implementation for working with the ML Commons Croissant metadata format.
Croissant is a standardized way to describe machine learning datasets using JSON-LD.`,
	Version: version.Version,
}

func Init() {

	// Initialize viper for configuration
	viper.SetEnvPrefix("CROISSANT")
	viper.AutomaticEnv()

	// Add child commands
	RootCmd.AddCommand(versionCmd())
	RootCmd.AddCommand(generateCmd())
	RootCmd.AddCommand(validateCmd())
	RootCmd.AddCommand(infoCmd())
	RootCmd.AddCommand(matchCmd())
}

func Execute() {
	// Execute the command
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// Helper functions

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func isCSVFile(filename string) bool {
	return croissant.IsCSVFile(filename)
}

func determineOutputPath(providedPath, csvPath string) string {
	if providedPath != "" {
		return providedPath
	}

	// Check environment variable
	envOutputPath := os.Getenv("CROISSANT_OUTPUT_PATH")
	if envOutputPath != "" {
		return envOutputPath
	}

	// Generate default path based on CSV filename
	baseName := strings.TrimSuffix(filepath.Base(csvPath), filepath.Ext(csvPath))
	return baseName + "_metadata.jsonld"
}
