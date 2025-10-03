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
// # Advanced Generation with Validation
//
// Generate metadata and get the parsed structure for further processing:
//
//	metadata, err := croissant.GenerateMetadataWithValidation("data.csv", "dataset.jsonld")
//	if err != nil {
//		log.Fatalf("Error generating metadata: %v", err)
//	}
//
//	// Validate the generated metadata
//	options := croissant.DefaultValidationOptions()
//	options.StrictMode = true
//	validationResult := metadata.ValidateWithOptions(options)
//
//	if validationResult.HasErrors() {
//		fmt.Println("Validation issues found:")
//		fmt.Println(validationResult.Report())
//	}
//
// # Data Type Inference
//
// The package automatically infers schema.org data types from CSV content:
//   - Boolean values (true/false, 1/0) → sc:Boolean
//   - Integer numbers (123, -456) → sc:Integer
//   - Floating-point numbers (3.14, 2.5e10) → sc:Float
//   - Dates in various formats (2023-01-01, 01/15/2023) → sc:Date
//   - URLs (https://example.com) → sc:URL
//   - Default to Text for other content → sc:Text
//
// # Validation
//
// Validate existing Croissant metadata:
//
//	issues, err := croissant.ValidateFile("metadata.jsonld")
//	if err != nil {
//		log.Fatalf("Validation error: %v", err)
//	}
//	if !issues.HasErrors() {
//		fmt.Println("Validation passed")
//	} else {
//		fmt.Println("Validation issues:")
//		fmt.Println(issues.Report())
//	}
//
// # Schema Compatibility Checking
//
// Compare two metadata files for schema compatibility:
//
//	reference, err := croissant.LoadMetadataFromFile("reference.jsonld")
//	if err != nil {
//		log.Fatalf("Error loading reference: %v", err)
//	}
//
//	candidate, err := croissant.LoadMetadataFromFile("candidate.jsonld")
//	if err != nil {
//		log.Fatalf("Error loading candidate: %v", err)
//	}
//
//	result := croissant.MatchMetadata(*reference, *candidate)
//	if result.IsMatch {
//		fmt.Printf("Compatible! %d fields matched\n", len(result.MatchedFields))
//	} else {
//		fmt.Printf("Incompatible: %d missing, %d type mismatches\n",
//			len(result.MissingFields), len(result.TypeMismatches))
//	}
//
// # JSON-LD Processing
//
// Work directly with JSON-LD data:
//
//	data, err := os.ReadFile("metadata.jsonld")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	issues, err := croissant.ValidateJSON(data)
//	if err != nil {
//		log.Fatalf("Validation error: %v", err)
//	}
//
//	fmt.Printf("Validation completed with %d errors and %d warnings\n",
//		len(issues.Errors), len(issues.Warnings))
//
// # Validation Options
//
// Customize validation behavior:
//
//	options := croissant.ValidationOptions{
//		StrictMode:      true,  // Enable additional warnings
//		CheckDataTypes:  true,  // Validate data type specifications
//		ValidateURLs:    false, // Skip network calls for URL validation
//		CheckFileExists: true,  // Verify referenced files exist
//	}
//
//	issues, err := croissant.ValidateJSONWithOptions(data, options)
//	if err != nil {
//		log.Fatal(err)
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

// CreateEnumerationRecordSet creates a RecordSet for categorical/enumeration data.
func CreateEnumerationRecordSet(id, name string, values []string, urls []string) RecordSet {
	fields := []Field{
		{
			ID:       fmt.Sprintf("%s/name", id),
			Type:     "cr:Field",
			Name:     fmt.Sprintf("%s/name", id),
			DataType: NewSingleDataType(VT_scText),
		},
	}

	// Add URL field if URLs are provided
	if len(urls) > 0 {
		urlField := Field{
			ID:       fmt.Sprintf("%s/url", id),
			Type:     "cr:Field",
			Name:     fmt.Sprintf("%s/url", id),
			DataType: NewSingleDataType(VT_scURL),
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
		DataType:    NewNullableSingleDataType(VT_scEnum),
		Fields:      fields,
		Key:         NewRecordSetKey(fmt.Sprintf("%s/name", id)),
		Data:        data,
	}

	return recordSet
}

// CreateSplitRecordSet creates a standard ML split RecordSet.
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

// CreateDefaultContext creates the ML Commons Croissant 1.0 compliant context.
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

// GenerateMetadata generates Croissant metadata from a CSV file (simple API).
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

// GenerateMetadataWithValidation generates Croissant metadata with validation from a CSV file.
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
		dataType := VT_scText // Default

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
			return nil, CroissantError{Message: "failed to marshal JSON-LD", Value: err}
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

// cleanFieldName cleans field names to be valid identifiers.
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
