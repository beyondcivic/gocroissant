// Package cmd provides the command-line interface for gocroissant.
//
// Gocroissant is a Go implementation for working with the ML Commons Croissant
// metadata format - a standardized way to describe machine learning datasets using JSON-LD.
//
// The command-line tool provides functionality to:
//   - Generate Croissant metadata from CSV files
//   - Validate existing Croissant metadata files
//   - Display version information
//
// # Usage Examples
//
// Generate metadata with default output path:
//
//	gocroissant generate data.csv
//
// Generate metadata with custom output path:
//
//	gocroissant generate data.csv -o output.json
//
// Generate and validate metadata:
//
//	gocroissant generate data.csv -o metadata.jsonld -v
//
// Validate existing metadata:
//
//	gocroissant validate metadata.json
//
// Show version information:
//
//	gocroissant version
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

	// Version command
	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version information",
		Long:  `Print the version, git hash, and build time information of the gocroissant tool.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("%s version %s\n", version.AppName, version.Version)
			stamp := version.RetrieveStamp()
			fmt.Printf("  Built with %s on %s\n", stamp.InfoGoCompiler, stamp.InfoBuildTime)
			fmt.Printf("  Git ref: %s\n", stamp.VCSRevision)
			fmt.Printf("  Go version %s, GOOS %s, GOARCH %s\n", stamp.InfoGoVersion, stamp.InfoGOOS, stamp.InfoGOARCH)
		},
	}
	RootCmd.AddCommand(versionCmd)

	// Generate command
	var generateCmd = &cobra.Command{
		Use:   "generate [csvPath]",
		Short: "Generate Croissant metadata from a CSV file",
		Long: `Generate Croissant metadata from a CSV file, automatically inferring data types 
		and creating a structured JSON-LD output that complies with the ML Commons Croissant specification.`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			csvPath := args[0]
			outputPath, _ := cmd.Flags().GetString("output")
			validate, _ := cmd.Flags().GetBool("validate")
			strict, _ := cmd.Flags().GetBool("strict")
			checkFiles, _ := cmd.Flags().GetBool("check-files")

			// Validate input file
			if !fileExists(csvPath) {
				fmt.Printf("Error: CSV file '%s' does not exist.\n", csvPath)
				os.Exit(1)
			}

			if !isCSVFile(csvPath) {
				fmt.Printf("Error: File '%s' does not appear to be a CSV file.\n", csvPath)
				os.Exit(1)
			}

			// Determine output path
			outputPath = determineOutputPath(outputPath, csvPath)

			// Validate output path
			if err := croissant.ValidateOutputPath(outputPath); err != nil {
				fmt.Printf("Error: Invalid output path: %v\n", err)
				os.Exit(1)
			}

			// Generate metadata
			fmt.Printf("Generating Croissant metadata for '%s'...\n", csvPath)
			metadata, err := croissant.GenerateMetadataWithValidation(csvPath, outputPath)
			if err != nil {
				fmt.Printf("Error generating metadata: %v\n", err)
				os.Exit(1)
			}

			// Set validation options
			if validate || strict || checkFiles {
				options := croissant.DefaultValidationOptions()
				options.StrictMode = strict
				options.CheckFileExists = checkFiles
				metadata.ValidateWithOptions(options)

				report := metadata.Report()
				if report != "" {
					fmt.Println(report)
				} else {
					fmt.Println("✓ Validation passed with no issues.")
				}

				if metadata.HasErrors() {
					fmt.Printf("\nMetadata generation completed but with validation errors.\n")
					if outputPath != "" {
						fmt.Printf("Metadata saved to: %s\n", outputPath)
					}
					os.Exit(1)
				}
			}

			fmt.Printf("✓ Croissant metadata generated successfully")
			if outputPath != "" {
				fmt.Printf(" and saved to: %s\n", outputPath)
			} else {
				fmt.Println()
			}
		},
	}
	generateCmd.Flags().StringP("output", "o", "", "Output path for the metadata JSON file")
	generateCmd.Flags().BoolP("validate", "v", false, "Validate the generated metadata and print issues")
	generateCmd.Flags().Bool("strict", false, "Enable strict validation mode")
	generateCmd.Flags().Bool("check-files", false, "Check if referenced files exist")
	RootCmd.AddCommand(generateCmd)

	// Validate command
	var validateCmd = &cobra.Command{
		Use:   "validate [jsonldPath]",
		Short: "Validate an existing Croissant metadata file",
		Long: `Validate an existing Croissant metadata JSON-LD file and report any issues found.
		This command checks compliance with the ML Commons Croissant specification.`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			jsonldPath := args[0]
			strict, _ := cmd.Flags().GetBool("strict")
			checkFiles, _ := cmd.Flags().GetBool("check-files")
			checkUrls, _ := cmd.Flags().GetBool("check-urls")

			// Validate input file
			if !fileExists(jsonldPath) {
				fmt.Printf("Error: Metadata file '%s' does not exist.\n", jsonldPath)
				os.Exit(1)
			}

			// Set validation options
			options := croissant.DefaultValidationOptions()
			options.StrictMode = strict
			options.CheckFileExists = checkFiles
			options.ValidateURLs = checkUrls

			fmt.Printf("Validating Croissant metadata file '%s'...\n", jsonldPath)

			// Read and parse the file manually to use validation options
			data, err := os.ReadFile(filepath.Clean(jsonldPath))
			if err != nil {
				fmt.Printf("Error reading file: %v\n", err)
				os.Exit(1)
			}

			issues, err := croissant.ValidateJSONWithOptions(data, options)
			if err != nil {
				fmt.Printf("Error validating metadata: %v\n", err)
				os.Exit(1)
			}

			report := issues.Report()
			if report != "" {
				fmt.Println(report)
			} else {
				fmt.Println("✓ Validation passed with no issues.")
			}

			if issues.HasErrors() {
				os.Exit(1)
			}
		},
	}
	validateCmd.Flags().Bool("strict", false, "Enable strict validation mode")
	validateCmd.Flags().Bool("check-files", false, "Check if referenced files exist")
	validateCmd.Flags().Bool("check-urls", false, "Validate URLs by making HTTP requests")
	RootCmd.AddCommand(validateCmd)

	// Info command - analyze CSV files
	var infoCmd = &cobra.Command{
		Use:   "info [csvPath]",
		Short: "Display information about a CSV file",
		Long:  `Analyze a CSV file and display information about its structure, columns, and data types.`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			csvPath := args[0]
			sampleSize, _ := cmd.Flags().GetInt("sample-size")

			if !fileExists(csvPath) {
				fmt.Printf("Error: CSV file '%s' does not exist.\n", csvPath)
				os.Exit(1)
			}

			// Validate CSV structure
			if err := croissant.ValidateCSVStructure(csvPath); err != nil {
				fmt.Printf("CSV validation error: %v\n", err)
				os.Exit(1)
			}

			// Get file stats
			stats, err := croissant.GetFileStats(csvPath)
			if err != nil {
				fmt.Printf("Error getting file stats: %v\n", err)
				os.Exit(1)
			}

			// Count total rows
			totalRows, err := croissant.CountCSVRows(csvPath)
			if err != nil {
				fmt.Printf("Error counting rows: %v\n", err)
				os.Exit(1)
			}

			// Get column information with enhanced type detection
			headers, columnTypes, err := croissant.GetCSVColumnTypes(csvPath, sampleSize)
			if err != nil {
				fmt.Printf("Error analyzing CSV: %v\n", err)
				os.Exit(1)
			}

			// Display information
			fmt.Printf("CSV File Information: %s\n", csvPath)
			fmt.Printf("=====================================\n")
			fmt.Printf("File Size: %v bytes\n", stats["size"])
			fmt.Printf("Total Rows: %d (including header)\n", totalRows)
			fmt.Printf("Data Rows: %d\n", totalRows-1)
			fmt.Printf("Columns: %d\n", len(headers))
			fmt.Printf("Sample Size: %d rows\n", sampleSize)
			fmt.Println()

			fmt.Printf("Column Information:\n")
			fmt.Printf("-------------------\n")
			for i, header := range headers {
				fmt.Printf("%d. %s (%s)\n", i+1, header, columnTypes[i])
			}
		},
	}
	infoCmd.Flags().Int("sample-size", 10, "Number of rows to sample for type inference")
	RootCmd.AddCommand(infoCmd)
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
