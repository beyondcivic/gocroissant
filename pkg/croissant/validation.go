// validation.go
package croissant

import (
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strings"
)

// ValidationOptions represents options for validation
type ValidationOptions struct {
	StrictMode      bool
	CheckDataTypes  bool
	ValidateURLs    bool
	CheckFileExists bool
}

// DefaultValidationOptions returns default validation options
func DefaultValidationOptions() ValidationOptions {
	return ValidationOptions{
		StrictMode:      true,
		CheckDataTypes:  true,
		ValidateURLs:    false, // Don't validate URLs by default to avoid network calls
		CheckFileExists: false, // Don't check file existence by default
	}
}

// ValidateFile validates a Croissant metadata file and returns issues
func ValidateFile(filePath string) (*Issues, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, CroissantError{Message: "failed to read file", Value: err}
	}

	return ValidateJSON(data)
}

// ValidateJSON validates Croissant metadata in JSON-LD format and returns issues
func ValidateJSON(data []byte) (*Issues, error) {
	// Use JSON-LD processor for proper validation and parsing
	processor := NewJSONLDProcessor()

	// First, validate that it's valid JSON-LD
	if err := processor.ValidateJSONLD(data); err != nil {
		return nil, CroissantError{Message: "invalid JSON-LD document", Value: err}
	}

	// Parse the metadata using JSON-LD processor
	metadata, err := processor.ParseCroissantMetadata(data)
	if err != nil {
		return nil, CroissantError{Message: "failed to parse Croissant metadata", Value: err}
	}

	return ValidateMetadata(*metadata), nil
}

// ValidateJSONWithOptions validates Croissant metadata in JSON-LD format with options and returns issues
func ValidateJSONWithOptions(data []byte, options ValidationOptions) (*Issues, error) {
	// Use JSON-LD processor for proper validation and parsing
	processor := NewJSONLDProcessor()

	// First, validate that it's valid JSON-LD
	if err := processor.ValidateJSONLD(data); err != nil {
		return nil, CroissantError{Message: "invalid JSON-LD document", Value: err}
	}

	// Parse the metadata using JSON-LD processor
	metadata, err := processor.ParseCroissantMetadata(data)
	if err != nil {
		return nil, CroissantError{Message: "failed to parse Croissant metadata", Value: err}
	}

	return ValidateMetadataWithOptions(*metadata, options), nil
}

// ValidateMetadata validates a Metadata struct and returns issues
func ValidateMetadata(metadata Metadata) *Issues {
	return ValidateMetadataWithOptions(metadata, DefaultValidationOptions())
}

// ValidateMetadataWithOptions validates a Metadata struct with specific options
func ValidateMetadataWithOptions(metadata Metadata, options ValidationOptions) *Issues {
	node := FromMetadata(metadata)
	issues := NewIssues()

	// Run comprehensive validation
	ValidateMetadataNode(node, issues, options)

	return issues
}

// ValidateMetadataNode performs comprehensive validation of a metadata node
func ValidateMetadataNode(node *MetadataNode, issues *Issues, options ValidationOptions) {
	// Basic metadata validation
	if node.Name == "" {
		issues.AddError("Property \"https://schema.org/name\" is mandatory, but does not exist.", node)
	}

	if node.Type != "sc:Dataset" {
		issues.AddError("The current JSON-LD doesn't extend https://schema.org/Dataset.", node)
	}

	if node.ConformsTo == "" {
		issues.AddWarning("Property \"http://purl.org/dc/terms/conformsTo\" is recommended, but does not exist.", node)
	} else if !isValidConformsTo(node.ConformsTo) {
		issues.AddWarning(fmt.Sprintf("ConformsTo value \"%s\" is not a recognized Croissant version.", node.ConformsTo), node)
	}

	// Strict mode validations for metadata
	if options.StrictMode {
		if node.Description == "" {
			issues.AddWarning("Dataset description is recommended for better documentation.", node)
		}
		if node.Version == "" {
			issues.AddWarning("Dataset version is recommended for proper versioning.", node)
		}
		if node.DatePublished == "" {
			issues.AddWarning("Date published is recommended for dataset tracking.", node)
		}
	}

	// Validate distributions
	if len(node.Distributions) == 0 {
		issues.AddError("Dataset must have at least one distribution.", node)
	}

	for _, dist := range node.Distributions {
		dist.SetParent(node)
		ValidateDistributionNode(dist, issues, options)
	}

	// Validate record sets
	if len(node.RecordSets) == 0 {
		issues.AddError("Dataset must have at least one recordSet.", node)
	}

	for _, rs := range node.RecordSets {
		rs.SetParent(node)
		ValidateRecordSetNode(rs, issues, options)
	}

	// Cross-references validation
	ValidateCrossReferences(node, issues)
}

// ValidateDistributionNode validates a distribution node
func ValidateDistributionNode(dist *DistributionNode, issues *Issues, options ValidationOptions) {
	// Required fields validation
	if dist.Name == "" {
		issues.AddError("Property \"https://schema.org/name\" is mandatory, but does not exist.", dist)
	}

	if dist.Type != "cr:FileObject" && dist.Type != "cr:FileSet" {
		issues.AddError(fmt.Sprintf("\"%s\" should have an attribute \"@type\": \"http://mlcommons.org/croissant/FileObject\" or \"@type\": \"http://mlcommons.org/croissant/FileSet\". Got %s instead.", dist.Name, dist.Type), dist)
	}

	if dist.ContentURL == "" {
		issues.AddError("Property \"https://schema.org/contentUrl\" is mandatory, but does not exist.", dist)
	} else if options.ValidateURLs && !isValidURL(dist.ContentURL) {
		issues.AddError(fmt.Sprintf("ContentURL \"%s\" is not a valid URL.", dist.ContentURL), dist)
	}

	if dist.EncodingFormat == "" {
		issues.AddError("Property \"https://schema.org/encodingFormat\" is mandatory, but does not exist.", dist)
	} else if !isValidEncodingFormat(dist.EncodingFormat) {
		issues.AddWarning(fmt.Sprintf("EncodingFormat \"%s\" is not a recognized MIME type.", dist.EncodingFormat), dist)
	}

	// Hash validation
	if dist.SHA256 != "" && !isValidSHA256(dist.SHA256) {
		issues.AddError(fmt.Sprintf("SHA256 hash \"%s\" is not a valid SHA-256 hash.", dist.SHA256), dist)
	} else if options.StrictMode && dist.SHA256 == "" {
		issues.AddWarning("SHA256 hash is recommended for file integrity verification.", dist)
	}

	if dist.MD5 != "" && !isValidMD5(dist.MD5) {
		issues.AddError(fmt.Sprintf("MD5 hash \"%s\" is not a valid MD5 hash.", dist.MD5), dist)
	}

	// File existence check
	if options.CheckFileExists && isLocalFile(dist.ContentURL) {
		if _, err := os.Stat(dist.ContentURL); os.IsNotExist(err) {
			issues.AddWarning(fmt.Sprintf("File \"%s\" does not exist.", dist.ContentURL), dist)
		}
	}
}

// ValidateRecordSetNode validates a record set node
func ValidateRecordSetNode(rs *RecordSetNode, issues *Issues, options ValidationOptions) {
	if rs.Name == "" {
		issues.AddError("Property \"https://schema.org/name\" is mandatory, but does not exist.", rs)
	}

	if rs.Type != "cr:RecordSet" {
		issues.AddError(fmt.Sprintf("\"%s\" should have an attribute \"@type\": \"http://mlcommons.org/croissant/RecordSet\". Got %s instead.", rs.Name, rs.Type), rs)
	}

	if len(rs.Fields) == 0 {
		issues.AddWarning("RecordSet has no fields defined.", rs)
	}

	// Validate key references if key is specified
	if rs.Key != nil {
		validateRecordSetKey(rs, issues)
	}

	for _, field := range rs.Fields {
		field.SetParent(rs)
		ValidateFieldNode(field, issues, options)
	}
}

// validateRecordSetKey validates that key references point to existing fields
func validateRecordSetKey(rs *RecordSetNode, issues *Issues) {
	if rs.Key == nil {
		return
	}

	keyIDs := rs.Key.GetKeyIDs()
	if len(keyIDs) == 0 {
		issues.AddError("Record set key is empty", rs)
		return
	}

	// Build a map of available field IDs for this record set
	fieldIDs := make(map[string]bool)
	for _, field := range rs.Fields {
		if field.ID != "" {
			fieldIDs[field.ID] = true
		}
		if field.Name != "" {
			fieldIDs[field.Name] = true
		}
	}

	// Check that all key IDs reference existing fields
	for _, keyID := range keyIDs {
		if !fieldIDs[keyID] {
			if rs.Key.IsComposite() {
				issues.AddError(fmt.Sprintf("Composite key references non-existent field \"%s\"", keyID), rs)
			} else {
				issues.AddError(fmt.Sprintf("Key references non-existent field \"%s\"", keyID), rs)
			}
		}
	}
}

// ValidateFieldNode validates a field node
func ValidateFieldNode(field *FieldNode, issues *Issues, options ValidationOptions) {
	if field == nil {
		return
	}

	if field.Name == "" {
		issues.AddError("Property \"https://schema.org/name\" is mandatory, but does not exist.", field)
	}

	if field.Type != "cr:Field" {
		issues.AddError(fmt.Sprintf("\"%s\" should have an attribute \"@type\": \"http://mlcommons.org/croissant/Field\". Got %s instead.", field.Name, field.Type), field)
	}

	// More explicit dataType validation
	if field.DataType.GetFirstType() == "" {
		issues.AddError(fmt.Sprintf("Field \"%s\" is missing required \"dataType\" property.", field.Name), field)
	} else if options.CheckDataTypes {
		if isValid, invalidTypes := validateDataTypes(field.DataType); !isValid {
			for _, invalidType := range invalidTypes {
				issues.AddError(fmt.Sprintf("Field \"%s\" has invalid dataType \"%s\". Valid types include: sc:Text, sc:Number, sc:Boolean, sc:DateTime, sc:URL, sc:GeoCoordinates, sc:ImageObject, cr:BoundingBox, etc.", field.Name, invalidType), field)
			}
		}
	}

	// Strict mode validations for fields
	if options.StrictMode {
		if field.Description == "" {
			issues.AddWarning(fmt.Sprintf("Field \"%s\" is missing recommended \"description\" property.", field.Name), field)
		}
	}

	// Check if field has subfields - use len check that's safe even if SubFields is nil
	hasSubFields := field.SubField != nil && len(field.SubField) > 0

	// Only validate source for leaf fields (fields without subfields)
	if !hasSubFields {
		// Check if source is properly configured
		if !hasValidFieldSource(field) {
			issues.AddError(fmt.Sprintf("Field \"%s\" has invalid or missing source configuration.", field.Name), field)
		}
	}

	// Validate subfields recursively
	if field.SubField != nil {
		for _, subField := range field.SubField {
			if subField != nil {
				subField.SetParent(field)
				ValidateFieldNode(subField, issues, options)
			}
		}
	}
}

// hasValidFieldSource checks if a field node has valid source configuration
func hasValidFieldSource(field *FieldNode) bool {
	if field == nil {
		return false
	}

	// Check if the source has file object reference and extraction method
	hasFileObject := field.Source.FileObject.ID != ""
	hasExtract := field.Source.Extract.Column != "" ||
		field.Source.Extract.JSONPath != "" ||
		field.Source.Extract.FileProperty != "" ||
		field.Source.Extract.Regex != ""
	hasFormat := field.Source.Format != ""

	// A valid source needs either a file object reference with extraction info, or format
	return hasFileObject && (hasExtract || hasFormat)
}

// ValidateCrossReferences validates that all references are valid
func ValidateCrossReferences(node *MetadataNode, issues *Issues) {
	// Build a map of all available IDs
	availableIDs := make(map[string]bool)

	// Add distribution IDs
	for _, dist := range node.Distributions {
		if dist.ID != "" {
			availableIDs[dist.ID] = true
		}
		if dist.Name != "" {
			availableIDs[dist.Name] = true
		}
	}

	// Add record set IDs
	for _, rs := range node.RecordSets {
		if rs.ID != "" {
			availableIDs[rs.ID] = true
		}
		if rs.Name != "" {
			availableIDs[rs.Name] = true
		}

		// Add field IDs
		for _, field := range rs.Fields {
			if field.ID != "" {
				availableIDs[field.ID] = true
			}
			fieldPath := fmt.Sprintf("%s/%s", rs.Name, field.Name)
			availableIDs[fieldPath] = true
		}
	}

	// Check field sources reference valid file objects
	for _, rs := range node.RecordSets {
		for _, field := range rs.Fields {
			if field.Source.FileObject.ID != "" {
				if !availableIDs[field.Source.FileObject.ID] {
					issues.AddError(fmt.Sprintf("Field \"%s\" references non-existent file object \"%s\".", field.Name, field.Source.FileObject.ID), field)
				}
			}
		}
	}
}

// AddValidationToMetadata adds validation functionality to the Metadata struct
type MetadataWithValidation struct {
	Metadata
	issues  *Issues
	options ValidationOptions
}

// NewMetadataWithValidation creates a new MetadataWithValidation instance
func NewMetadataWithValidation(metadata Metadata) *MetadataWithValidation {
	return &MetadataWithValidation{
		Metadata: metadata,
		options:  DefaultValidationOptions(),
	}
}

// Validate runs validation on the metadata
func (m *MetadataWithValidation) Validate() {
	m.issues = ValidateMetadataWithOptions(m.Metadata, m.options)
}

// ValidateWithOptions runs validation with specific options
func (m *MetadataWithValidation) ValidateWithOptions(options ValidationOptions) {
	m.options = options
	m.issues = ValidateMetadataWithOptions(m.Metadata, options)
}

// Report returns a string report of validation issues
// Returns an empty string is there are no issues or warnings.
func (m *MetadataWithValidation) Report() string {
	if m.issues == nil {
		m.Validate()
	}
	return m.issues.Report()
}

// HasErrors returns true if there are validation errors
func (m *MetadataWithValidation) HasErrors() bool {
	if m.issues == nil {
		m.Validate()
	}
	return m.issues.HasErrors()
}

// HasWarnings returns true if there are validation warnings
func (m *MetadataWithValidation) HasWarnings() bool {
	if m.issues == nil {
		m.Validate()
	}
	return m.issues.HasWarnings()
}

// GetIssues returns the validation issues
func (m *MetadataWithValidation) GetIssues() *Issues {
	if m.issues == nil {
		m.Validate()
	}
	return m.issues
}

// Validation helper functions

func isValidConformsTo(conformsTo string) bool {
	validVersions := []string{
		"http://mlcommons.org/croissant/1.0",
		"http://mlcommons.org/croissant/0.8",
	}
	for _, version := range validVersions {
		if conformsTo == version {
			return true
		}
	}
	return false
}

func isValidURL(urlStr string) bool {
	if urlStr == "" {
		return false
	}

	// Allow relative URLs and file paths
	if !strings.Contains(urlStr, "://") {
		return true
	}

	_, err := url.ParseRequestURI(urlStr)
	return err == nil
}

func isLocalFile(urlStr string) bool {
	return !strings.HasPrefix(urlStr, "http://") && !strings.HasPrefix(urlStr, "https://") && !strings.HasPrefix(urlStr, "ftp://")
}

func isValidEncodingFormat(format string) bool {
	// Common MIME types for datasets
	validFormats := map[string]bool{
		"text/csv":                  true,
		"application/json":          true,
		"application/jsonl":         true,
		"text/plain":                true,
		"application/xml":           true,
		"text/xml":                  true,
		"application/parquet":       true,
		"text/tab-separated-values": true,
		"application/zip":           true,
		"application/gzip":          true,
		"application/x-tar":         true,
		"image/jpeg":                true,
		"image/png":                 true,
		"image/tiff":                true,
		"audio/wav":                 true,
		"audio/mpeg":                true,
		"video/mp4":                 true,
		"application/pdf":           true,
	}

	return validFormats[format] || strings.HasPrefix(format, "text/") || strings.HasPrefix(format, "application/") || strings.HasPrefix(format, "image/") || strings.HasPrefix(format, "audio/") || strings.HasPrefix(format, "video/")
}

// validateDataTypes validates all data types in a DataType (single or array)
func validateDataTypes(dt DataType) (bool, []string) {
	types := dt.GetTypes()
	if len(types) == 0 {
		return false, []string{}
	}

	var invalidTypes []string
	allValid := true

	for _, dataType := range types {
		if !IsValidDataType(dataType) {
			allValid = false
			invalidTypes = append(invalidTypes, dataType)
		}
	}

	return allValid, invalidTypes
}

func isValidSHA256(hash string) bool {
	matched, _ := regexp.MatchString("^[a-fA-F0-9]{64}$", hash)
	return matched
}

func isValidMD5(hash string) bool {
	matched, _ := regexp.MatchString("^[a-fA-F0-9]{32}$", hash)
	return matched
}
