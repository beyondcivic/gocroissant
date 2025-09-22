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
		}
		rsNode.SetParent(node)

		// Convert fields
		for _, field := range rs.Fields {
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
						Column: field.Source.Extract.Column,
					},
					FileObject: FileObjectRef{
						ID: field.Source.FileObject.ID,
					},
				},
			}
			fieldNode.SetParent(rsNode)
			rsNode.Fields = append(rsNode.Fields, fieldNode)
		}

		node.RecordSets = append(node.RecordSets, rsNode)
	}

	return node
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
	Type        string        `json:"@type"`
	Description string        `json:"description,omitempty"`
	Fields      []*FieldNode  `json:"field"`
	Key         *RecordSetKey `json:"key,omitempty"`
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
	Type        string     `json:"@type"`
	Description string     `json:"description,omitempty"`
	DataType    DataType   `json:"dataType,omitempty"`
	Source      SourceNode `json:"source"`
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
	}

	// Validate source
	if !f.Source.ValidateSource() {
		issues.AddError(fmt.Sprintf("Node \"%s\" is a field and has no source. Please, use http://mlcommons.org/croissant/source to specify the source.", f.ID), f)
	}
}

// SourceNode represents a source
type SourceNode struct {
	Extract    ExtractNode   `json:"extract"`
	FileObject FileObjectRef `json:"fileObject"`
}

// ValidateSource validates the source node
func (s *SourceNode) ValidateSource() bool {
	// Check if both extract and file object references are valid
	return s.Extract.Column != "" && s.FileObject.ID != ""
}

// ExtractNode represents extraction details
type ExtractNode struct {
	Column string `json:"column,omitempty"`
}

// FileObjectRef represents a reference to a file object
type FileObjectRef struct {
	ID string `json:"@id"`
}
