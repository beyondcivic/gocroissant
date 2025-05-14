package main

import (
	"fmt"
	"os"

	"github.com/beyondcivic/gocroissant/pkg/croissant"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func main() {
	// Initialize viper for configuration
	viper.SetEnvPrefix("CROISSANT")
	viper.AutomaticEnv()

	var rootCmd = &cobra.Command{
		Use:   "croissant-metadata [csvPath]",
		Short: "Generate Croissant metadata from a CSV file",
		Long:  `A tool to generate Croissant metadata from a CSV file, inferring data types and creating a structured JSON output.`,
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			csvPath := args[0]
			outputPath, _ := cmd.Flags().GetString("output")

			metadataPath, err := croissant.GenerateMetadata(csvPath, outputPath)
			if err != nil {
				fmt.Printf("Error generating metadata: %v\n", err)
				os.Exit(1)
			}

			fmt.Printf("Croissant metadata generated and saved to: %s\n", metadataPath)
		},
	}

	// Add flags
	rootCmd.Flags().StringP("output", "o", "", "Output path for the metadata JSON file")

	// Execute the command
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
