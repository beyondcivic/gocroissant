package main

import (
	"fmt"
	"log"
	"os"

	"github.com/beyondcivic/gocroissant/pkg/croissant"
)

// This file demonstrates basic usage of the gocroissant library.
// It provides practical examples of how to use the croissant package
// for generating, validating, and comparing Croissant metadata.

func main() {
	// Example 1: Generate metadata from CSV
	fmt.Println("=== Example 1: Generate Metadata ===")
	generateMetadataExample()

	// Example 2: Validate metadata
	fmt.Println("\n=== Example 2: Validate Metadata ===")
	validateMetadataExample()

	// Example 3: Compare metadata files
	fmt.Println("\n=== Example 3: Compare Metadata ===")
	compareMetadataExample()

	// Example 4: Advanced validation with options
	fmt.Println("\n=== Example 4: Advanced Validation ===")
	advancedValidationExample()
}

// generateMetadataExample shows how to generate metadata from a CSV file
func generateMetadataExample() {
	// Generate metadata with automatic type inference
	outputPath, err := croissant.GenerateMetadata("sample_csv/data.csv", "generated_metadata.jsonld")
	if err != nil {
		log.Printf("Error generating metadata: %v", err)
		return
	}
	fmt.Printf("✓ Metadata generated and saved to: %s\n", outputPath)

	// Alternative: Generate and get metadata object for further processing
	metadata, err := croissant.GenerateMetadataWithValidation("sample_csv/data.csv", "")
	if err != nil {
		log.Printf("Error generating metadata with validation: %v", err)
		return
	}
	fmt.Printf("✓ Generated metadata for dataset: %s\n", metadata.Name)
	fmt.Printf("  Record sets: %d\n", len(metadata.RecordSets))
	fmt.Printf("  Distributions: %d\n", len(metadata.Distributions))
}

// validateMetadataExample demonstrates metadata validation
func validateMetadataExample() {
	// Basic validation
	issues, err := croissant.ValidateFile("samples_jsonld/ok1.jsonld")
	if err != nil {
		log.Printf("Validation error: %v", err)
		return
	}

	if issues.HasErrors() {
		fmt.Println("✗ Validation failed:")
		fmt.Println(issues.Report())
	} else {
		fmt.Println("✓ Validation passed!")
		if issues.HasWarnings() {
			fmt.Println("Warnings found:")
			fmt.Println(issues.Report())
		}
	}
}

// compareMetadataExample shows how to compare two metadata files
func compareMetadataExample() {
	// Load reference metadata
	reference, err := croissant.LoadMetadataFromFile("samples_jsonld/ok1.jsonld")
	if err != nil {
		log.Printf("Error loading reference: %v", err)
		return
	}

	// Load candidate metadata
	candidate, err := croissant.LoadMetadataFromFile("samples_jsonld/ok2.jsonld")
	if err != nil {
		log.Printf("Error loading candidate: %v", err)
		return
	}

	// Compare the metadata
	result := croissant.MatchMetadata(*reference, *candidate)

	// Display results
	if result.IsMatch {
		fmt.Printf("✓ Schemas are compatible!\n")
		fmt.Printf("  Matched fields: %d\n", len(result.MatchedFields))

		if len(result.ExtraFields) > 0 {
			fmt.Printf("  Additional fields in candidate: %d\n", len(result.ExtraFields))
		}
	} else {
		fmt.Printf("✗ Schemas are incompatible:\n")

		if len(result.MissingFields) > 0 {
			fmt.Printf("  Missing fields: %d\n", len(result.MissingFields))
			for _, field := range result.MissingFields {
				fmt.Printf("    - %s\n", field)
			}
		}

		if len(result.TypeMismatches) > 0 {
			fmt.Printf("  Type mismatches: %d\n", len(result.TypeMismatches))
			for _, mismatch := range result.TypeMismatches {
				fmt.Printf("    - %s: expected %s, got %s\n",
					mismatch.FieldName, mismatch.ReferenceType, mismatch.CandidateType)
			}
		}
	}
}

// advancedValidationExample shows validation with custom options
func advancedValidationExample() {
	// Read file content
	data, err := os.ReadFile("samples_jsonld/ok1.jsonld")
	if err != nil {
		log.Printf("Error reading file: %v", err)
		return
	}

	// Configure validation options
	options := croissant.ValidationOptions{
		StrictMode:      true,  // Enable additional warnings
		CheckDataTypes:  true,  // Validate data type specifications
		ValidateURLs:    false, // Skip network calls for URL validation
		CheckFileExists: false, // Skip file existence checking for this example
	}

	// Validate with custom options
	issues, err := croissant.ValidateJSONWithOptions(data, options)
	if err != nil {
		log.Printf("Validation error: %v", err)
		return
	}

	fmt.Printf("Advanced validation completed:\n")
	fmt.Printf("  Errors: %d\n", issues.ErrorCount())
	fmt.Printf("  Warnings: %d\n", issues.WarningCount())

	if issues.HasErrors() || issues.HasWarnings() {
		fmt.Println("\nDetailed report:")
		fmt.Println(issues.Report())
	} else {
		fmt.Println("✓ No issues found!")
	}
}
