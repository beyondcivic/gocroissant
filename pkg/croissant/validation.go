// validation.go
package croissant

import (
	"encoding/json"
	"io/ioutil"
)

// ValidateFile validates a Croissant metadata file and returns issues
func ValidateFile(filePath string) (*Issues, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	return ValidateJSON(data)
}

// ValidateJSON validates Croissant metadata in JSON format and returns issues
func ValidateJSON(data []byte) (*Issues, error) {
	var metadata Metadata
	if err := json.Unmarshal(data, &metadata); err != nil {
		return nil, err
	}

	return ValidateMetadata(metadata), nil
}

// ValidateMetadata validates a Metadata struct and returns issues
func ValidateMetadata(metadata Metadata) *Issues {
	node := FromMetadata(metadata)
	issues := NewIssues()

	// Run validation
	node.Validate(issues)

	return issues
}

// AddValidationToMetadata adds validation functionality to the Metadata struct
type MetadataWithValidation struct {
	Metadata
	issues *Issues
}

// Validate runs validation on the metadata
func (m *MetadataWithValidation) Validate() {
	m.issues = ValidateMetadata(m.Metadata)
}

// Report returns a string report of validation issues
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
