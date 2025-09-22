// semantic_datatypes_demo.go - Demonstrates the new semantic dataTypes support

package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/beyondcivic/gocroissant/pkg/croissant"
)

func main() {
	fmt.Println("=== Croissant Semantic DataTypes Demo ===\n")

	// Create a dataset with enumeration support
	metadata := createDatasetWithSemanticTypes()

	// Serialize to JSON
	jsonData, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Generated Croissant metadata with semantic dataTypes:")
	fmt.Println(string(jsonData))

	// Validate the metadata
	fmt.Println("\n=== Validation Results ===")
	validateMetadata(metadata)
}

func createDatasetWithSemanticTypes() croissant.Metadata {
	// Create context with all necessary namespaces
	context := croissant.CreateDefaultContext()

	// Create split enumeration RecordSet
	splitRecordSet := croissant.CreateSplitRecordSet()

	// Create a gender enumeration RecordSet (like in Titanic example)
	genderRecordSet := croissant.RecordSet{
		ID:       "genders",
		Type:     "cr:RecordSet",
		Name:     "genders",
		DataType: croissant.NewArrayDataType("sc:Enumeration", "wd:Q48277"), // Gender concept
		Key:      croissant.NewSingleKey("genders/name"),
		Fields: []croissant.Field{
			{
				ID:       "genders/name",
				Type:     "cr:Field",
				Name:     "genders/name",
				DataType: croissant.NewSingleDataType("sc:Text"),
			},
			{
				ID:       "genders/url",
				Type:     "cr:Field",
				Name:     "genders/url",
				DataType: croissant.NewSingleDataType("sc:URL"),
			},
		},
		Data: []map[string]interface{}{
			{"genders/name": "female", "genders/url": "wd:Q6581072"},
			{"genders/name": "male", "genders/url": "wd:Q6581097"},
		},
	}

	// Create main data RecordSet with semantic field types
	mainRecordSet := croissant.RecordSet{
		ID:          "main_data",
		Type:        "cr:RecordSet",
		Name:        "main_data",
		Description: "Main dataset with various semantic datatypes",
		Fields: []croissant.Field{
			// Basic image field
			{
				ID:          "main_data/image",
				Type:        "cr:Field",
				Name:        "image",
				Description: "Image content",
				DataType:    croissant.NewSingleDataType("sc:ImageObject"),
				Source: croissant.FieldSource{
					Extract: croissant.Extract{
						Column: "image_path",
					},
					FileObject: croissant.FileObject{
						ID: "demo_data.csv",
					},
				},
			},
			// Label field
			{
				ID:          "main_data/category",
				Type:        "cr:Field",
				Name:        "category",
				Description: "Object category label",
				DataType:    croissant.NewArrayDataType("sc:Text", "cr:Label"),
				Source: croissant.FieldSource{
					Extract: croissant.Extract{
						Column: "category",
					},
					FileObject: croissant.FileObject{
						ID: "demo_data.csv",
					},
				},
			},
			// Split reference field
			{
				ID:          "main_data/split",
				Type:        "cr:Field",
				Name:        "split",
				Description: "Data split assignment",
				DataType:    croissant.NewArrayDataType("cr:Split", "sc:Text"),
				References:  &croissant.FieldRef{Field: "splits/name"},
				Source: croissant.FieldSource{
					Extract: croissant.Extract{
						Column: "split",
					},
					FileObject: croissant.FileObject{
						ID: "demo_data.csv",
					},
				},
			},
			// Gender reference field
			{
				ID:          "main_data/gender",
				Type:        "cr:Field",
				Name:        "gender",
				Description: "Gender classification",
				DataType:    croissant.NewArrayDataType("sc:Enumeration", "sc:Text"),
				References:  &croissant.FieldRef{Field: "genders/name"},
				Source: croissant.FieldSource{
					Extract: croissant.Extract{
						Column: "gender",
					},
					FileObject: croissant.FileObject{
						ID: "demo_data.csv",
					},
				},
			},
			// Bounding box field
			{
				ID:          "main_data/bbox",
				Type:        "cr:Field",
				Name:        "bbox",
				Description: "Object bounding box",
				DataType:    croissant.NewSingleDataType("cr:BoundingBox"),
				Source: croissant.FieldSource{
					Extract: croissant.Extract{
						Column: "bbox",
					},
					FileObject: croissant.FileObject{
						ID: "demo_data.csv",
					},
					Format: "CENTER_XYWH",
				},
			},
		},
	}

	// Create the complete metadata
	metadata := croissant.Metadata{
		Context:     context,
		Type:        "sc:Dataset",
		Name:        "semantic_dataset_demo",
		Description: "Demonstration dataset with semantic dataTypes",
		ConformsTo:  "http://mlcommons.org/croissant/1.0",
		Version:     "1.0.0",
		Distributions: []croissant.Distribution{
			{
				ID:             "demo_data.csv",
				Type:           "cr:FileObject",
				Name:           "demo_data.csv",
				ContentURL:     "demo_data.csv",
				EncodingFormat: "text/csv",
			},
		},
		RecordSets: []croissant.RecordSet{
			splitRecordSet,
			genderRecordSet,
			mainRecordSet,
		},
	}

	return metadata
}

func validateMetadata(metadata croissant.Metadata) {
	// Convert to node for validation
	node := croissant.FromMetadata(metadata)
	node.Validate(node.Issues)

	if node.Issues.HasErrors() || node.Issues.HasWarnings() {
		fmt.Println("Validation Report:")
		fmt.Println(node.Issues.Report())
	} else {
		fmt.Println("✅ No validation errors or warnings")
	}

	// Test dataType validation
	fmt.Println("\n=== DataType Validation Examples ===")
	testDataTypes := []string{
		"sc:Text",
		"sc:ImageObject",
		"cr:Label",
		"cr:BoundingBox",
		"cr:Split",
		"sc:Enumeration",
		"wd:Q48277", // Gender
		"invalid:Type",
	}

	for _, dataType := range testDataTypes {
		isValid := croissant.IsValidDataType(dataType)
		status := "✅"
		if !isValid {
			status = "❌"
		}
		fmt.Printf("%s %s\n", status, dataType)
	}
}
