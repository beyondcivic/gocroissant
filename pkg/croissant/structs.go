package croissant

// Field represents a field in the Croissant metadata
type Field struct {
	ID          string      `json:"@id"`
	Type        string      `json:"@type"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	DataType    string      `json:"dataType"`
	Source      FieldSource `json:"source"`
}

// FieldSource represents the source information for a field
type FieldSource struct {
	Extract    Extract    `json:"extract"`
	FileObject FileObject `json:"fileObject"`
}

// Extract represents the extraction information for a field source
type Extract struct {
	Column string `json:"column"`
}

// FileObject represents a file object reference
type FileObject struct {
	ID string `json:"@id"`
}

// Distribution represents a file in the Croissant metadata
type Distribution struct {
	ID             string `json:"@id"`
	Type           string `json:"@type"`
	Name           string `json:"name"`
	ContentSize    string `json:"contentSize"`
	ContentURL     string `json:"contentUrl"`
	EncodingFormat string `json:"encodingFormat"`
	SHA256         string `json:"sha256"`
}

// RecordSet represents a record set in the Croissant metadata
type RecordSet struct {
	ID          string  `json:"@id"`
	Type        string  `json:"@type"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Fields      []Field `json:"field"`
}

// Context represents the JSON-LD context in the Croissant metadata
type Context struct {
	Language     string          `json:"@language"`
	Vocab        string          `json:"@vocab"`
	CiteAs       string          `json:"citeAs"`
	Column       string          `json:"column"`
	ConformsTo   string          `json:"conformsTo"`
	CR           string          `json:"cr"`
	DCT          string          `json:"dct"`
	Data         DataContext     `json:"data"`
	DataType     DataTypeContext `json:"dataType"`
	Extract      string          `json:"extract"`
	Field        string          `json:"field"`
	FileObject   string          `json:"fileObject"`
	FileProperty string          `json:"fileProperty"`
	SC           string          `json:"sc"`
	Source       string          `json:"source"`
}

// DataContext represents the data field in the context
type DataContext struct {
	ID   string `json:"@id"`
	Type string `json:"@type"`
}

// DataTypeContext represents the dataType field in the context
type DataTypeContext struct {
	ID   string `json:"@id"`
	Type string `json:"@type"`
}

// Metadata represents the complete Croissant metadata
type Metadata struct {
	Context       Context        `json:"@context"`
	Type          string         `json:"@type"`
	Name          string         `json:"name"`
	Description   string         `json:"description"`
	ConformsTo    string         `json:"conformsTo"`
	DatePublished string         `json:"datePublished"`
	Version       string         `json:"version"`
	Distributions []Distribution `json:"distribution"`
	RecordSets    []RecordSet    `json:"recordSet"`
}
