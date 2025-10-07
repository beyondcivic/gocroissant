// commands.go
// Contains cobra command definitions
//
//nolint:funlen,mnd
package cmd

import (
	"fmt"
	"os"

	"github.com/beyondcivic/gocroissant/pkg/croissant"
	"github.com/beyondcivic/gocroissant/pkg/version"
	"github.com/spf13/cobra"
)

// Version Command.
// Displays tool version and build information.
func versionCmd() *cobra.Command {
	return &cobra.Command{
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
}

// Generate command.
func generateCmd() *cobra.Command {
	var generateCmd = &cobra.Command{
		Use:   "generate [csvPath]",
		Short: "Generate Croissant metadata from a CSV file",
		Long: `Generate Croissant metadata from a CSV file, automatically inferring data types 
		and creating a structured JSON-LD output that complies with the ML Commons Croissant specification.`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			csvPath := args[0]
			flagOutputPath, _ := cmd.Flags().GetString("output")
			flagValidate, _ := cmd.Flags().GetBool("validate")
			flagStrict, _ := cmd.Flags().GetBool("strict")
			flagCheckFiles, _ := cmd.Flags().GetBool("check-files")

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
			outputPath := determineOutputPath(flagOutputPath, csvPath)

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

			fmt.Printf("✓ Croissant metadata generated successfully")
			if outputPath != "" {
				fmt.Printf(" and saved to: %s\n", outputPath)
			}

			// Set validation options
			if flagValidate || flagStrict || flagCheckFiles {
				options := commonValidationCmd(flagStrict, flagCheckFiles, false)
				metadata.ValidateWithOptions(options)

				analyzeMetadataIssues(metadata.GetIssues())
			}
		},
	}
	generateCmd.Flags().StringP("output", "o", "", "Output path for the metadata JSON file")
	generateCmd.Flags().BoolP("validate", "v", false, "Validate the generated metadata and print issues")
	generateCmd.Flags().Bool("strict", false, "Enable strict validation mode")
	generateCmd.Flags().Bool("check-files", false, "Check if referenced files exist")

	return generateCmd
}

// Validate command - validate a croissant jsonld file.
func validateCmd() *cobra.Command {
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

			fmt.Printf("Validating Croissant metadata file '%s'...\n", jsonldPath)

			// Read and parse the file manually to use validation options
			data, err := os.ReadFile(jsonldPath)
			if err != nil {
				fmt.Printf("Error reading file: %v\n", err)
				os.Exit(1)
			}
			// Set validation options
			options := commonValidationCmd(strict, checkFiles, checkUrls)

			issues, err := croissant.ValidateJSONWithOptions(data, options)
			if err != nil {
				fmt.Printf("Error validating metadata: %v\n", err)
				os.Exit(1)
			}

			analyzeMetadataIssues(issues)
		},
	}
	validateCmd.Flags().Bool("strict", false, "Enable strict validation mode")
	validateCmd.Flags().Bool("check-files", false, "Check if referenced files exist")
	validateCmd.Flags().Bool("check-urls", false, "Validate URLs by making HTTP requests")

	return validateCmd
}

// Info command - analyze CSV files.
func infoCmd() *cobra.Command {
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

	return infoCmd
}

func matchCmd() *cobra.Command {
	// Match command - compare two Croissant metadata files
	var matchCmd = &cobra.Command{
		Use:   "match [reference] [candidate]",
		Short: "Compare two Croissant metadata files for compatibility",
		Long: `Compare two Croissant metadata JSON-LD files to check if the candidate is compatible 
		with the reference. The candidate can have additional fields, but all reference fields 
		must exist in the candidate with matching data types and names.`,
		Args: cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			referencePath := args[0]
			candidatePath := args[1]
			verbose, _ := cmd.Flags().GetBool("verbose")

			// Validate input files
			if !fileExists(referencePath) {
				fmt.Printf("Error: Reference file '%s' does not exist.\n", referencePath)
				os.Exit(1)
			}

			if !fileExists(candidatePath) {
				fmt.Printf("Error: Candidate file '%s' does not exist.\n", candidatePath)
				os.Exit(1)
			}

			fmt.Printf("Comparing metadata files...\n")
			fmt.Printf("Reference: %s\n", referencePath)
			fmt.Printf("Candidate: %s\n", candidatePath)
			fmt.Println()

			// Load reference metadata
			referenceMetadata, err := croissant.LoadMetadataFromFile(referencePath)
			if err != nil {
				fmt.Printf("Error loading reference metadata: %v\n", err)
				os.Exit(1)
			}

			// Load candidate metadata
			candidateMetadata, err := croissant.LoadMetadataFromFile(candidatePath)
			if err != nil {
				fmt.Printf("Error loading candidate metadata: %v\n", err)
				os.Exit(1)
			}

			// Perform the match
			result := croissant.MatchMetadata(*referenceMetadata, *candidateMetadata)

			// Display results
			if result.IsMatch {
				fmt.Printf("✓ Compatibility check PASSED\n")
				fmt.Printf("The candidate metadata is compatible with the reference.\n")
			} else {
				fmt.Printf("✗ Compatibility check FAILED\n")
				fmt.Printf("The candidate metadata is NOT compatible with the reference.\n")
			}
			fmt.Println()

			// Display detailed results
			matchPrintResultAnalysis(result, verbose)

			// Summary
			fmt.Printf("Summary:\n")
			fmt.Printf("  Matched: %d\n", len(result.MatchedFields))
			fmt.Printf("  Missing: %d\n", len(result.MissingFields))
			fmt.Printf("  Type mismatches: %d\n", len(result.TypeMismatches))
			fmt.Printf("  Extra fields: %d\n", len(result.ExtraFields))

			if !result.IsMatch {
				os.Exit(1)
			}
		},
	}
	matchCmd.Flags().BoolP("verbose", "v", false, "Show detailed information including extra fields in candidate")

	return matchCmd
}

// Prints information about matched fields.
// Lists matched, missing, type mismatched, and extra fields.
// nolint:cyclop
func matchPrintResultAnalysis(result *croissant.MatchResult, verbose bool) {
	if len(result.MatchedFields) > 0 {
		fmt.Printf("Matched fields (%d):\n", len(result.MatchedFields))
		for _, field := range result.MatchedFields {
			fmt.Printf("  ✓ %s\n", field)
		}
		fmt.Println()
	}

	if len(result.MissingFields) > 0 {
		fmt.Printf("Missing fields (%d):\n", len(result.MissingFields))
		for _, field := range result.MissingFields {
			fmt.Printf("  ✗ %s (required by reference but not found in candidate)\n", field)
		}
		fmt.Println()
	}

	if len(result.TypeMismatches) > 0 {
		fmt.Printf("Type mismatches (%d):\n", len(result.TypeMismatches))
		for _, mismatch := range result.TypeMismatches {
			fmt.Printf("  ✗ %s: reference expects '%s', candidate has '%s'\n",
				mismatch.FieldName, mismatch.ReferenceType, mismatch.CandidateType)
		}
		fmt.Println()
	}

	if len(result.ExtraFields) > 0 && verbose {
		fmt.Printf("Extra fields in candidate (%d):\n", len(result.ExtraFields))
		for _, field := range result.ExtraFields {
			fmt.Printf("  + %s (additional field in candidate)\n", field)
		}
		fmt.Println()
	} else if len(result.ExtraFields) > 0 {
		fmt.Printf("Candidate has %d additional field(s) (use --verbose to see details)\n", len(result.ExtraFields))
		fmt.Println()
	}
}

// Common configuration of validation options.
func commonValidationCmd(flagStrict bool, flagCheckFiles bool, flagCheckUrls bool) croissant.ValidationOptions {
	options := croissant.DefaultValidationOptions()
	options.StrictMode = flagStrict
	options.CheckFileExists = flagCheckFiles
	options.ValidateURLs = flagCheckUrls

	return options
}

// Pretty prints issues for command output.
// Does not return, calls os.Exit().
func analyzeMetadataIssues(issues *croissant.Issues) {
	report := issues.Report()
	if report != "" {
		fmt.Println(report)
	} else {
		fmt.Println("✓ Validation passed with no issues.")
	}

	if issues.HasErrors() {
		fmt.Printf("\n Validation failed with errors\n")
		os.Exit(1)
	}

	// If there's only warnings, return a safe exit code.
	os.Exit(0)
}
