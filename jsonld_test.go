package main
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/beyondcivic/gocroissant/pkg/croissant"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run jsonld_test.go <jsonld_file>")
		os.Exit(1)
	}

	filePath := os.Args[1]

	// Read the file
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	// Create a JSON-LD processor
	processor := croissant.NewJSONLDProcessor()

	// Validate JSON-LD
	fmt.Printf("Validating JSON-LD document: %s\n", filePath)
	if err := processor.ValidateJSONLD(data); err != nil {
		fmt.Printf("❌ Invalid JSON-LD: %v\n", err)
	} else {
		fmt.Printf("✅ Valid JSON-LD document\n")
	}

	// Parse and expand JSON-LD
	fmt.Println("\nExpanding JSON-LD document...")
	expanded, err := processor.ParseJSONLD(data)
	if err != nil {
		log.Fatalf("Error expanding JSON-LD: %v", err)
	}

	// Extract common Croissant properties
	fmt.Println("\nExtracting Croissant properties...")
	properties := croissant.ExtractCroissantProperties(expanded)

	fmt.Printf("Dataset Name: %s\n", properties["name"])
	fmt.Printf("Description: %s\n", properties["description"])
	fmt.Printf("Type: %s\n", properties["@type"])
	fmt.Printf("ConformsTo: %s\n", properties["conformsTo"])

	// Parse as Croissant metadata
	fmt.Println("\nParsing as Croissant metadata...")
	metadata, err := processor.ParseCroissantMetadata(data)
	if err != nil {
		log.Fatalf("Error parsing Croissant metadata: %v", err)
	}

	fmt.Printf("✅ Successfully parsed Croissant metadata for dataset: %s\n", metadata.Name)
	fmt.Printf("   - Version: %s\n", metadata.Version)
	fmt.Printf("   - Distributions: %d\n", len(metadata.Distributions))
	fmt.Printf("   - Record Sets: %d\n", len(metadata.RecordSets))
}