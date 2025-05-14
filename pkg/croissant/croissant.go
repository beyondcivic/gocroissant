// croissant.go
package croissant

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// InferDataType infers the schema.org data type from a value
func InferDataType(value string) string {
	// Try to parse as integer
	if _, err := strconv.ParseInt(value, 10, 64); err == nil {
		return "sc:Integer"
	}

	// Try to parse as float
	if _, err := strconv.ParseFloat(value, 64); err == nil {
		return "sc:Float"
	}

	// Try to parse as date (YYYY-MM-DD)
	if _, err := time.Parse("2006-01-02", value); err == nil {
		return "sc:Date"
	}

	// Default to Text
	return "sc:Text"
}

// CreateDefaultContext creates the default context for Croissant metadata
func CreateDefaultContext() Context {
	return Context{
		Language:   "en",
		Vocab:      "https://schema.org/",
		CiteAs:     "cr:citeAs",
		Column:     "cr:column",
		ConformsTo: "dct:conformsTo",
		CR:         "http://mlcommons.org/croissant/",
		DCT:        "http://purl.org/dc/terms/",
		Data: DataContext{
			ID:   "cr:data",
			Type: "@json",
		},
		DataType: DataTypeContext{
			ID:   "cr:dataType",
			Type: "@vocab",
		},
		Extract:      "cr:extract",
		Field:        "cr:field",
		FileObject:   "cr:fileObject",
		FileProperty: "cr:fileProperty",
		SC:           "https://schema.org/",
		Source:       "cr:source",
	}
}

// GenerateMetadataWithValidation generates Croissant metadata with validation from a CSV file
func GenerateMetadataWithValidation(csvPath string, outputPath string) (*MetadataWithValidation, error) {
	// Get file information
	fileName := filepath.Base(csvPath)
	fileInfo, err := os.Stat(csvPath)
	if err != nil {
		return nil, err
	}
	fileSize := fileInfo.Size()

	// Calculate SHA-256 hash
	fileSHA256, err := CalculateSHA256(csvPath)
	if err != nil {
		return nil, err
	}

	// Get column information
	headers, firstRow, err := GetCSVColumns(csvPath)
	if err != nil {
		return nil, err
	}

	// Create fields based on CSV columns
	fields := make([]Field, 0, len(headers))
	for i, header := range headers {
		fieldID := fmt.Sprintf("main/%s", header)
		dataType := "sc:Text" // Default

		// Try to infer data type from first row if available
		if firstRow != nil && i < len(firstRow) {
			dataType = InferDataType(firstRow[i])
		}

		field := Field{
			ID:          fieldID,
			Type:        "cr:Field",
			Name:        header,
			Description: fmt.Sprintf("Field for %s", header),
			DataType:    dataType,
			Source: FieldSource{
				Extract: Extract{
					Column: header,
				},
				FileObject: FileObject{
					ID: fileName,
				},
			},
		}

		fields = append(fields, field)
	}

	// Create metadata structure
	metadata := Metadata{
		Context:       CreateDefaultContext(),
		Type:          "sc:Dataset",
		Name:          fmt.Sprintf("%s_dataset", strings.Split(fileName, ".")[0]),
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

	// Set default output path if not provided
	if outputPath != "" {
		// Marshal metadata to JSON
		metadataJSON, err := json.MarshalIndent(metadata, "", "  ")
		if err != nil {
			return nil, err
		}

		// Write metadata to file
		if err := os.WriteFile(outputPath, metadataJSON, 0644); err != nil {
			return nil, err
		}
	}

	// Create and validate metadata
	metadataWithValidation := &MetadataWithValidation{
		Metadata: metadata,
	}
	metadataWithValidation.Validate()

	return metadataWithValidation, nil
}
