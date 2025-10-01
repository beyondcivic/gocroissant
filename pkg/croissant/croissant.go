// Package croissant provides functionality for working with the ML Commons Croissant
// metadata format - a standardized way to describe machine learning datasets using JSON-LD.
//
// This package simplifies the creation of Croissant-compatible metadata from CSV data sources by
// automatically inferring schema types from dataset content, generating complete valid JSON-LD
// metadata, providing validation tools to ensure compatibility, and supporting the full
// Croissant specification.
//
// # Basic Usage
//
// Generate metadata from a CSV file:
//
//	outputPath, err := croissant.GenerateMetadata("data.csv", "dataset.jsonld")
//	if err != nil {
//		log.Fatalf("Error generating metadata: %v", err)
//	}
//	fmt.Printf("Metadata saved to: %s\n", outputPath)
//
// # Data Type Inference
//
// The package automatically infers schema.org data types from CSV content:
//   - Boolean values (true/false)
//   - Integer numbers
//   - Floating-point numbers
//   - Dates in various formats
//   - URLs
//   - Default to Text for other content
//
// # Validation
//
// Validate existing Croissant metadata:
//
//	issues, err := croissant.ValidateMetadata("metadata.jsonld")
//	if err != nil {
//		log.Fatalf("Validation error: %v", err)
//	}
//	if len(issues) == 0 {
//		fmt.Println("Validation passed")
//	}
package croissant

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// CreateEnumerationRecordSet creates a RecordSet for categorical/enumeration data
func CreateEnumerationRecordSet(id, name string, values []string, urls []string) RecordSet {
	fields := []Field{
		{
			ID:       fmt.Sprintf("%s/name", id),
			Type:     "cr:Field",
			Name:     fmt.Sprintf("%s/name", id),
			DataType: NewSingleDataType("sc:Text"),
		},
	}

	// Add URL field if URLs are provided
	if len(urls) > 0 {
		urlField := Field{
			ID:       fmt.Sprintf("%s/url", id),
			Type:     "cr:Field",
			Name:     fmt.Sprintf("%s/url", id),
			DataType: NewSingleDataType("sc:URL"),
		}
		fields = append(fields, urlField)
	}

	// Create data entries
	data := make([]map[string]interface{}, len(values))
	for i, value := range values {
		entry := map[string]interface{}{
			fmt.Sprintf("%s/name", id): value,
		}
		if i < len(urls) && urls[i] != "" {
			entry[fmt.Sprintf("%s/url", id)] = urls[i]
		}
		data[i] = entry
	}

	recordSet := RecordSet{
		ID:          id,
		Type:        "cr:RecordSet",
		Name:        name,
		Description: fmt.Sprintf("Enumeration values for %s", name),
		DataType:    NewNullableSingleDataType("sc:Enumeration"),
		Fields:      fields,
		Key:         NewSingleKey(fmt.Sprintf("%s/name", id)),
		Data:        data,
	}

	return recordSet
}

// CreateSplitRecordSet creates a standard ML split RecordSet
func CreateSplitRecordSet() RecordSet {
	recordSet := CreateEnumerationRecordSet(
		"splits",
		"splits",
		[]string{"train", "val", "test"},
		[]string{"cr:TrainingSplit", "cr:ValidationSplit", "cr:TestSplit"},
	)

	// Set the dataType to cr:Split for splits
	recordSet.DataType = NewNullableSingleDataType("cr:Split")

	return recordSet
}

// CreateDefaultContext creates the ML Commons Croissant 1.0 compliant context
func CreateDefaultContext() Context {
	return Context{
		Language:   "en",
		Vocab:      "https://schema.org/",
		CiteAs:     "cr:citeAs",
		Column:     "cr:column",
		ConformsTo: "dct:conformsTo",
		CR:         "http://mlcommons.org/croissant/",
		DCT:        "http://purl.org/dc/terms/",
		WD:         "https://www.wikidata.org/wiki/",
		Data: DataContext{
			ID:   "cr:data",
			Type: "@json",
		},
		DataType: DataTypeContext{
			ID:   "cr:dataType",
			Type: "@vocab",
		},
		Examples: DataContext{
			ID:   "cr:examples",
			Type: "@json",
		},
		Extract:       "cr:extract",
		Field:         "cr:field",
		FileObject:    "cr:fileObject",
		FileProperty:  "cr:fileProperty",
		FileSet:       "cr:fileSet",
		Format:        "cr:format",
		Includes:      "cr:includes",
		IsLiveDataset: "cr:isLiveDataset",
		JSONPath:      "cr:jsonPath",
		Key:           "cr:key",
		MD5:           "cr:md5",
		ParentField:   "cr:parentField",
		Path:          "cr:path",
		RecordSet:     "cr:recordSet",
		References:    "cr:references",
		Regex:         "cr:regex",
		Repeated:      "cr:repeated",
		Replace:       "cr:replace",
		SC:            "https://schema.org/",
		Separator:     "cr:separator",
		Source:        "cr:source",
		SubField:      "cr:subField",
		Transform:     "cr:transform",
	}
}

// GenerateMetadata generates Croissant metadata from a CSV file (simple API)
func GenerateMetadata(csvPath string, outputPath string) (string, error) {
	metadata, err := GenerateMetadataWithValidation(csvPath, outputPath)
	if err != nil {
		return "", err
	}

	// Check if there are validation errors
	if metadata.HasErrors() {
		return "", CroissantError{Message: "validation failed", Value: metadata.Report()}
	}

	return outputPath, nil
}

// GenerateMetadataWithValidation generates Croissant metadata with validation from a CSV file
func GenerateMetadataWithValidation(csvPath string, outputPath string) (*MetadataWithValidation, error) {
	// Get file information
	fileName := filepath.Base(csvPath)
	fileInfo, err := os.Stat(csvPath)
	if err != nil {
		return nil, CroissantError{Message: "failed to get file info", Value: err}
	}
	fileSize := fileInfo.Size()

	// Calculate SHA-256 hash
	fileSHA256, err := CalculateSHA256(csvPath)
	if err != nil {
		return nil, CroissantError{Message: "failed to calculate SHA-256", Value: err}
	}

	// Get column information
	headers, firstRow, err := GetCSVColumns(csvPath)
	if err != nil {
		return nil, CroissantError{Message: "failed to read CSV", Value: err}
	}

	// Create fields based on CSV columns with data type inference
	fields := make([]Field, 0, len(headers))
	for i, header := range headers {
		fieldID := fmt.Sprintf("main/%s", cleanFieldName(header))
		dataType := "sc:Text" // Default

		// Infer data type from first row if available
		if firstRow != nil && i < len(firstRow) {
			dataType = InferDataType(firstRow[i])
		}

		field := Field{
			ID:          fieldID,
			Type:        "cr:Field",
			Name:        header,
			Description: fmt.Sprintf("Field for %s", header),
			DataType:    NewSingleDataType(dataType),
			Source: FieldSource{
				Extract: &Extract{
					Column: header,
				},
				FileObject: &FileObject{
					ID: fileName,
				},
			},
		}

		fields = append(fields, field)
	}

	// Create metadata structure
	datasetName := strings.TrimSuffix(fileName, filepath.Ext(fileName))
	metadata := Metadata{
		Context:       CreateDefaultContext(),
		Type:          "sc:Dataset",
		Name:          fmt.Sprintf("%s_dataset", datasetName),
		Description:   fmt.Sprintf("Dataset created from %s", fileName),
		ConformsTo:    "http://mlcommons.org/croissant/1.0",
		DatePublished: time.Now().Format("2006-01-02"),
		Version:       "1.0.0",
		Distributions: []Distribution{
			{
				ID:             fileName,
				Type:           "cr:FileObject",
				Name:           fileName,
				ContentSize:    fmt.Sprintf("%d B", fileSize),
				ContentURL:     fileName,
				EncodingFormat: "text/csv",
				SHA256:         fileSHA256,
			},
		},
		RecordSets: []RecordSet{
			{
				ID:          "main",
				Type:        "cr:RecordSet",
				Name:        "main",
				Description: fmt.Sprintf("Records from %s", fileName),
				Fields:      fields,
			},
		},
	}

	// Write to file if output path is provided
	if outputPath != "" {
		// Marshal metadata to JSON-LD with proper indentation
		metadataJSON, err := json.MarshalIndent(metadata, "", "  ")
		if err != nil {
			return nil, CroissantError{Message: "failed to marshal JSON-LD: %w", Value: err}
		}

		// Validate that the generated JSON is valid JSON-LD
		processor := NewJSONLDProcessor()
		if err := processor.ValidateJSONLD(metadataJSON); err != nil {
			return nil, CroissantError{Message: "generated invalid JSON-LD", Value: err}
		}

		// Ensure directory exists
		if err := os.MkdirAll(filepath.Dir(outputPath), 0750); err != nil {
			return nil, CroissantError{Message: "failed to create directory", Value: err}
		}

		// Write metadata to file
		if err := os.WriteFile(outputPath, metadataJSON, 0600); err != nil {
			return nil, CroissantError{Message: "failed to write file", Value: err}
		}
	}

	// Create and validate metadata
	metadataWithValidation := &MetadataWithValidation{
		Metadata: metadata,
	}
	metadataWithValidation.Validate()

	return metadataWithValidation, nil
}

// cleanFieldName cleans field names to be valid identifiers
func cleanFieldName(name string) string {
	// Replace spaces and special characters with underscores
	reg := regexp.MustCompile(`[^a-zA-Z0-9_]`)
	cleaned := reg.ReplaceAllString(name, "_")

	// Remove leading/trailing underscores and multiple consecutive underscores
	cleaned = strings.Trim(cleaned, "_")
	reg2 := regexp.MustCompile(`_{2,}`)
	cleaned = reg2.ReplaceAllString(cleaned, "_")

	// Ensure it doesn't start with a number
	if len(cleaned) > 0 && cleaned[0] >= '0' && cleaned[0] <= '9' {
		cleaned = "field_" + cleaned
	}

	return cleaned
}
