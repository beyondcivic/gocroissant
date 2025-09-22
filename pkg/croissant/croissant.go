// croissant.go
package croissant

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// InferDataType infers the schema.org data type from a value
func InferDataType(value string) string {
	// Trim whitespace
	value = strings.TrimSpace(value)
	if value == "" {
		return "sc:Text"
	}

	// Try to parse as boolean
	lowerVal := strings.ToLower(value)
	if lowerVal == "true" || lowerVal == "false" {
		return "sc:Boolean"
	}

	// Try to parse as integer
	if _, err := strconv.ParseInt(value, 10, 64); err == nil {
		return "sc:Integer"
	}

	// Try to parse as float
	if _, err := strconv.ParseFloat(value, 64); err == nil {
		return "sc:Number"
	}

	// Try to parse as date (various formats)
	dateFormats := []string{
		"2006-01-02",
		"01/02/2006",
		"2006/01/02",
		"2006-01-02T15:04:05",
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05-07:00",
	}
	for _, format := range dateFormats {
		if _, err := time.Parse(format, value); err == nil {
			return "sc:DateTime"
		}
	}

	// Try to parse as URL
	if _, err := url.ParseRequestURI(value); err == nil && (strings.HasPrefix(value, "http://") || strings.HasPrefix(value, "https://")) {
		return "sc:URL"
	}

	// Try to detect email
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if emailRegex.MatchString(value) {
		return "sc:Text" // Email is still text but we could add a custom type if needed
	}

	// Default to Text
	return "sc:Text"
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
		return "", fmt.Errorf("validation failed: %s", metadata.Report())
	}

	return outputPath, nil
}

// GenerateMetadataWithValidation generates Croissant metadata with validation from a CSV file
func GenerateMetadataWithValidation(csvPath string, outputPath string) (*MetadataWithValidation, error) {
	// Get file information
	fileName := filepath.Base(csvPath)
	fileInfo, err := os.Stat(csvPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}
	fileSize := fileInfo.Size()

	// Calculate SHA-256 hash
	fileSHA256, err := CalculateSHA256(csvPath)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate SHA-256: %w", err)
	}

	// Get column information
	headers, firstRow, err := GetCSVColumns(csvPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV: %w", err)
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
			return nil, fmt.Errorf("failed to marshal JSON-LD: %w", err)
		}

		// Validate that the generated JSON is valid JSON-LD
		processor := NewJSONLDProcessor()
		if err := processor.ValidateJSONLD(metadataJSON); err != nil {
			return nil, fmt.Errorf("generated invalid JSON-LD: %w", err)
		}

		// Ensure directory exists
		if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
			return nil, fmt.Errorf("failed to create directory: %w", err)
		}

		// Write metadata to file
		if err := os.WriteFile(outputPath, metadataJSON, 0644); err != nil {
			return nil, fmt.Errorf("failed to write file: %w", err)
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
