// metadata_node.go
package croissant

import "fmt"

// MetadataNode represents a Croissant metadata document
type MetadataNode struct {
	BaseNode
	Context       Context             `json:"@context"`
	Type          string              `json:"@type"`
	Description   string              `json:"description,omitempty"`
	ConformsTo    string              `json:"conformsTo,omitempty"`
	DatePublished string              `json:"datePublished,omitempty"`
	Version       string              `json:"version,omitempty"`
	Distributions []*DistributionNode `json:"distribution"`
	RecordSets    []*RecordSetNode    `json:"recordSet"`
	Issues        *Issues             `json:"-"` // Not serialized to JSON
}

// NewMetadataNode creates a new MetadataNode
func NewMetadataNode() *MetadataNode {
	return &MetadataNode{
		BaseNode: BaseNode{
			Name: "",
		},
		Context:       CreateDefaultContext(),
		Type:          "sc:Dataset",
		Distributions: make([]*DistributionNode, 0),
		RecordSets:    make([]*RecordSetNode, 0),
		Issues:        NewIssues(),
	}
}

// Validate validates the metadata node
func (m *MetadataNode) Validate(issues *Issues) {
	// Validate required fields
	if m.Name == "" {
		issues.AddError("Property \"https://schema.org/name\" is mandatory, but does not exist.", m)
	}

	// Validate type
	if m.Type != "sc:Dataset" {
		issues.AddError("The current JSON-LD doesn't extend https://schema.org/Dataset.", m)
	}

	// Validate distributions
	for _, dist := range m.Distributions {
		dist.SetParent(m)
		dist.Validate(issues)
	}

	// Validate record sets
	for _, rs := range m.RecordSets {
		rs.SetParent(m)
		rs.Validate(issues)
	}

	// Validate conformsTo is set
	if m.ConformsTo == "" {
		issues.AddWarning("Property \"http://purl.org/dc/terms/conformsTo\" is recommended, but does not exist.", m)
	}
}

// FromMetadata converts a Metadata struct to a MetadataNode
func FromMetadata(metadata Metadata) *MetadataNode {
	node := &MetadataNode{
		BaseNode: BaseNode{
			Name: metadata.Name,
		},
		Context:       metadata.Context,
		Type:          metadata.Type,
		Description:   metadata.Description,
		ConformsTo:    metadata.ConformsTo,
		DatePublished: metadata.DatePublished,
		Version:       metadata.Version,
		Issues:        NewIssues(),
	}

	// Convert distributions
	for _, dist := range metadata.Distributions {
		distNode := &DistributionNode{
			BaseNode: BaseNode{
				ID:   dist.ID,
				Name: dist.Name,
			},
			Type:           dist.Type,
			ContentSize:    dist.ContentSize,
			ContentURL:     dist.ContentURL,
			EncodingFormat: dist.EncodingFormat,
			SHA256:         dist.SHA256,
		}
		distNode.SetParent(node)
		node.Distributions = append(node.Distributions, distNode)
	}

	// Convert record sets
	for _, rs := range metadata.RecordSets {
		rsNode := &RecordSetNode{
			BaseNode: BaseNode{
				ID:   rs.ID,
				Name: rs.Name,
			},
			Type:        rs.Type,
			Description: rs.Description,
			Key:         rs.Key,
			Data:        rs.Data,
		}

		// Handle DataType safely - check for nil
		if rs.DataType != nil {
			rsNode.DataType = *rs.DataType
		} else {
			// Set empty DataType if nil
			rsNode.DataType = DataType{}
		}

		rsNode.SetParent(node)

		// Convert fields
		for _, field := range rs.Fields {
			fieldNode := convertFieldToNode(field, rsNode)
			rsNode.Fields = append(rsNode.Fields, fieldNode)
		}

		node.RecordSets = append(node.RecordSets, rsNode)
	}

	return node
}

// convertFieldToNode converts a Field to a FieldNode with proper nil handling
func convertFieldToNode(field Field, parent Node) *FieldNode {
	fieldNode := &FieldNode{
		BaseNode: BaseNode{
			ID:   field.ID,
			Name: field.Name,
		},
		Type:        field.Type,
		Description: field.Description,
		DataType:    field.DataType,
		Source: SourceNode{
			Extract: ExtractNode{
				Column:       field.Source.Extract.Column,
				JSONPath:     field.Source.Extract.JSONPath,
				Regex:        field.Source.Extract.Regex,
				FileProperty: field.Source.Extract.FileProperty,
			},
			FileObject: FileObjectRef{
				ID: field.Source.FileObject.ID,
			},
			FileSet: FileObjectRef{
				ID: field.Source.FileSet.ID,
			},
			Transform: field.Source.Transform,
			Format:    field.Source.Format,
		},
		Repeated:   field.Repeated,
		Examples:   field.Examples,
		References: field.References,
	}
	fieldNode.SetParent(parent)

	// Convert subfields if they exist
	for _, subField := range field.SubField {
		subFieldNode := convertFieldToNode(subField, fieldNode)
		fieldNode.SubField = append(fieldNode.SubField, subFieldNode)
	}

	return fieldNode
}

// DistributionNode represents a file distribution
type DistributionNode struct {
	BaseNode
	Type           string `json:"@type"`
	ContentSize    string `json:"contentSize,omitempty"`
	ContentURL     string `json:"contentUrl,omitempty"`
	EncodingFormat string `json:"encodingFormat,omitempty"`
	SHA256         string `json:"sha256,omitempty"`
	MD5            string `json:"md5,omitempty"`
}

// Validate validates the distribution node
func (d *DistributionNode) Validate(issues *Issues) {
	// Validate required fields
	if d.Name == "" {
		issues.AddError("Property \"https://schema.org/name\" is mandatory, but does not exist.", d)
	}

	// Validate type
	if d.Type != "cr:FileObject" && d.Type != "cr:FileSet" {
		issues.AddError(fmt.Sprintf("\"%s\" should have an attribute \"@type\": \"http://mlcommons.org/croissant/FileObject\" or \"@type\": \"http://mlcommons.org/croissant/FileSet\". Got %s instead.", d.Name, d.Type), d)
	}

	// Validate content URL
	if d.ContentURL == "" {
		issues.AddError("Property \"https://schema.org/contentUrl\" is mandatory, but does not exist.", d)
	}

	// Validate encoding format
	if d.EncodingFormat == "" {
		issues.AddError("Property \"https://schema.org/encodingFormat\" is mandatory, but does not exist.", d)
	}
}

// RecordSetNode represents a record set
type RecordSetNode struct {
	BaseNode
	Type        string                   `json:"@type"`
	Description string                   `json:"description,omitempty"`
	DataType    DataType                 `json:"dataType,omitempty"`
	Fields      []*FieldNode             `json:"field"`
	Key         *RecordSetKey            `json:"key,omitempty"`
	Data        []map[string]interface{} `json:"data,omitempty"`
}

// Validate validates the record set node
func (r *RecordSetNode) Validate(issues *Issues) {
	// Validate required fields
	if r.Name == "" {
		issues.AddError("Property \"https://schema.org/name\" is mandatory, but does not exist.", r)
	}

	// Validate type
	if r.Type != "cr:RecordSet" {
		issues.AddError(fmt.Sprintf("\"%s\" should have an attribute \"@type\": \"http://mlcommons.org/croissant/RecordSet\". Got %s instead.", r.Name, r.Type), r)
	}

	// Validate dataType if specified
	if r.DataType.GetFirstType() != "" {
		r.validateDataType(issues)
	}

	// Validate key references if key is specified
	if r.Key != nil {
		r.validateKey(issues)
	}

	// Validate fields
	for _, field := range r.Fields {
		field.SetParent(r)
		field.Validate(issues)
	}
}

// validateDataType validates RecordSet dataType and associated requirements
func (r *RecordSetNode) validateDataType(issues *Issues) {
	firstType := r.DataType.GetFirstType()

	// Validate enumeration RecordSets
	if firstType == "sc:Enumeration" {
		r.validateEnumeration(issues)
	}

	// Validate split RecordSets
	if firstType == "cr:Split" {
		r.validateSplit(issues)
	}
}

// validateEnumeration validates enumeration-specific requirements
func (r *RecordSetNode) validateEnumeration(issues *Issues) {
	// Enumeration RecordSets must have a key
	if r.Key == nil {
		issues.AddError("Enumeration RecordSet must specify a key", r)
		return
	}

	// Should have a name field
	hasNameField := false
	for _, field := range r.Fields {
		if field.Name == fmt.Sprintf("%s/name", r.Name) || field.Name == "name" {
			hasNameField = true
			break
		}
	}
	if !hasNameField {
		issues.AddWarning("Enumeration RecordSet should have a 'name' field", r)
	}
}

// validateSplit validates split-specific requirements
func (r *RecordSetNode) validateSplit(issues *Issues) {
	// Split RecordSets should have name and url fields
	hasNameField := false
	hasUrlField := false

	for _, field := range r.Fields {
		fieldName := field.Name
		if fieldName == fmt.Sprintf("%s/name", r.Name) || fieldName == "name" {
			hasNameField = true
		}
		if fieldName == fmt.Sprintf("%s/url", r.Name) || fieldName == "url" {
			hasUrlField = true
		}
	}

	if !hasNameField {
		issues.AddWarning("Split RecordSet should have a 'name' field", r)
	}
	if !hasUrlField {
		issues.AddWarning("Split RecordSet should have a 'url' field", r)
	}
}

// validateKey validates that key references point to existing fields
func (r *RecordSetNode) validateKey(issues *Issues) {
	if r.Key == nil {
		return
	}

	keyIDs := r.Key.GetKeyIDs()
	if len(keyIDs) == 0 {
		issues.AddError("Record set key is empty", r)
		return
	}

	// Build a map of available field IDs for this record set
	fieldIDs := make(map[string]bool)
	for _, field := range r.Fields {
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
			if r.Key.IsComposite() {
				issues.AddError(fmt.Sprintf("Composite key references non-existent field \"%s\"", keyID), r)
			} else {
				issues.AddError(fmt.Sprintf("Key references non-existent field \"%s\"", keyID), r)
			}
		}
	}
}

// FieldNode represents a field
type FieldNode struct {
	BaseNode
	Type        string          `json:"@type"`
	Description string          `json:"description,omitempty"`
	DataType    DataType        `json:"dataType,omitempty"`
	Source      SourceNode      `json:"source,omitempty"`
	Repeated    bool            `json:"repeated,omitempty"`
	Examples    interface{}     `json:"examples,omitempty"`
	SubField    []*FieldNode    `json:"subField,omitempty"`
	References  []FieldRefSlice `json:"references,omitempty"`
}

// Validate validates the field node
func (f *FieldNode) Validate(issues *Issues) {
	// Validate required fields
	if f.Name == "" {
		issues.AddError("Property \"https://schema.org/name\" is mandatory, but does not exist.", f)
	}

	// Validate type
	if f.Type != "cr:Field" {
		issues.AddError(fmt.Sprintf("\"%s\" should have an attribute \"@type\": \"http://mlcommons.org/croissant/Field\". Got %s instead.", f.Name, f.Type), f)
	}

	// Validate data type
	if f.DataType.GetFirstType() == "" {
		issues.AddError("The field does not specify a valid http://mlcommons.org/croissant/dataType, neither does any of its predecessor.", f)
	} else {
		// Validate that all dataTypes are valid
		for _, dataType := range f.DataType.GetTypes() {
			if !IsValidDataType(dataType) {
				issues.AddWarning(fmt.Sprintf("DataType \"%s\" may not be recognized", dataType), f)
			}
		}
	}

	// Validate source - but skip validation for fields in enumeration RecordSets with inline data
	if !f.Source.ValidateSource() {
		// Check if this field belongs to an enumeration RecordSet with inline data
		if parent := f.GetParent(); parent != nil {
			if recordSet, ok := parent.(*RecordSetNode); ok {
				// If it's an enumeration with inline data, skip source validation
				if recordSet.DataType.GetFirstType() == "sc:Enumeration" && len(recordSet.Data) > 0 {
					// Skip validation for enumeration fields with inline data
				} else {
					issues.AddError(fmt.Sprintf("Node \"%s\" is a field and has no source. Please, use http://mlcommons.org/croissant/source to specify the source.", f.ID), f)
				}
			} else {
				issues.AddError(fmt.Sprintf("Node \"%s\" is a field and has no source. Please, use http://mlcommons.org/croissant/source to specify the source.", f.ID), f)
			}
		} else {
			issues.AddError(fmt.Sprintf("Node \"%s\" is a field and has no source. Please, use http://mlcommons.org/croissant/source to specify the source.", f.ID), f)
		}
	}
}

// SourceNode represents a source
type SourceNode struct {
	Extract    ExtractNode   `json:"extract,omitempty"`
	FileObject FileObjectRef `json:"fileObject,omitempty"`
	FileSet    FileObjectRef `json:"fileSet,omitempty"`
	Transform  *Transform    `json:"transform,omitempty"`
	Format     string        `json:"format,omitempty"`
}

// ValidateSource validates the source node
func (s *SourceNode) ValidateSource() bool {
	// Check if either extract has content or file object/set references are valid
	hasExtract := s.Extract.Column != "" || s.Extract.JSONPath != "" || s.Extract.FileProperty != ""
	hasFileRef := s.FileObject.ID != "" || s.FileSet.ID != ""

	return hasExtract && hasFileRef
}

// ExtractNode represents extraction details
type ExtractNode struct {
	Regex        string `json:"regex,omitempty"`
	Column       string `json:"column,omitempty"`
	JSONPath     string `json:"jsonPath,omitempty"`
	FileProperty string `json:"fileProperty,omitempty"`
}

// FileObjectRef represents a reference to a file object
type FileObjectRef struct {
	ID string `json:"@id"`
}
