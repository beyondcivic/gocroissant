// main.go
package main

import (
	"fmt"
	"os"

	"github.com/beyondcivic/gocroissant/pkg/croissant"
	"github.com/beyondcivic/gocroissant/pkg/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func main() {
	// Initialize viper for configuration
	viper.SetEnvPrefix("CROISSANT")
	viper.AutomaticEnv()

	var rootCmd = &cobra.Command{
		Use:   "croissant",
		Short: "Croissant metadata tools",
		Long:  `Tools for generating and validating Croissant metadata.`,
	}

	// Version command
	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version information",
		Long:  `Print the version, git hash, and build time information of the gocroissant tool.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("%s version %s\n", version.AppName, version.Version)
			fmt.Printf("Git commit: %s\n", version.GitHash)
			fmt.Printf("Built on: %s\n", version.BuildTime)
		},
	}
	rootCmd.AddCommand(versionCmd)

	// Generate command
	var generateCmd = &cobra.Command{
		Use:   "generate [csvPath]",
		Short: "Generate Croissant metadata from a CSV file",
		Long:  `Generate Croissant metadata from a CSV file, inferring data types and creating a structured JSON output.`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			csvPath := args[0]
			outputPath, _ := cmd.Flags().GetString("output")
			validate, _ := cmd.Flags().GetBool("validate")

			metadata, err := croissant.GenerateMetadataWithValidation(csvPath, outputPath)
			if err != nil {
				fmt.Printf("Error generating metadata: %v\n", err)
				os.Exit(1)
			}

			if validate || cmd.Flags().Changed("validate") {
				report := metadata.Report()
				if report != "" {
					fmt.Println(report)
				} else {
					fmt.Println("Validation passed with no issues.")
				}

				if metadata.HasErrors() {
					os.Exit(1)
				}
			}

			if outputPath != "" {
				fmt.Printf("Croissant metadata generated and saved to: %s\n", outputPath)
			}
		},
	}
	generateCmd.Flags().StringP("output", "o", "", "Output path for the metadata JSON file")
	generateCmd.Flags().BoolP("validate", "v", false, "Validate the generated metadata and print issues")
	rootCmd.AddCommand(generateCmd)

	// Validate command
	var validateCmd = &cobra.Command{
		Use:   "validate [jsonldPath]",
		Short: "Validate an existing Croissant metadata file",
		Long:  `Validate an existing Croissant metadata JSON-LD file and report any issues found.`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			jsonldPath := args[0]

			issues, err := croissant.ValidateFile(jsonldPath)
			if err != nil {
				fmt.Printf("Error validating metadata: %v\n", err)
				os.Exit(1)
			}

			report := issues.Report()
			if report != "" {
				fmt.Println(report)
			} else {
				fmt.Println("Validation passed with no issues.")
			}

			if issues.HasErrors() {
				os.Exit(1)
			}
		},
	}
	rootCmd.AddCommand(validateCmd)

	// Execute the command
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
