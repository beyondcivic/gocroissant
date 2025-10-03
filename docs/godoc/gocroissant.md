# croissant

```go
import "github.com/beyondcivic/gocroissant/pkg/croissant"
```

Package croissant provides functionality for working with the ML Commons Croissant metadata format \- a standardized way to describe machine learning datasets using JSON\-LD.

This package simplifies the creation of Croissant\-compatible metadata from CSV data sources by automatically inferring schema types from dataset content, generating complete valid JSON\-LD metadata, providing validation tools to ensure compatibility, and supporting the full Croissant specification.

### Basic Usage

Generate metadata from a CSV file:

```
outputPath, err := croissant.GenerateMetadata("data.csv", "dataset.jsonld")
if err != nil {
	log.Fatalf("Error generating metadata: %v", err)
}
fmt.Printf("Metadata saved to: %s\n", outputPath)
```

### Advanced Generation with Validation

Generate metadata and get the parsed structure for further processing:

```
metadata, err := croissant.GenerateMetadataWithValidation("data.csv", "dataset.jsonld")
if err != nil {
	log.Fatalf("Error generating metadata: %v", err)
}

// Validate the generated metadata
options := croissant.DefaultValidationOptions()
options.StrictMode = true
validationResult := metadata.ValidateWithOptions(options)

if validationResult.HasErrors() {
	fmt.Println("Validation issues found:")
	fmt.Println(validationResult.Report())
}
```

### Data Type Inference

The package automatically infers schema.org data types from CSV content:

- Boolean values \(true/false, 1/0\) → sc:Boolean
- Integer numbers \(123, \-456\) → sc:Integer
- Floating\-point numbers \(3.14, 2.5e10\) → sc:Float
- Dates in various formats \(2023\-01\-01, 01/15/2023\) → sc:Date
- URLs \(https://example.com\) → sc:URL
- Default to Text for other content → sc:Text

### Validation

Validate existing Croissant metadata:

```
issues, err := croissant.ValidateFile("metadata.jsonld")
if err != nil {
	log.Fatalf("Validation error: %v", err)
}
if !issues.HasErrors() {
	fmt.Println("Validation passed")
} else {
	fmt.Println("Validation issues:")
	fmt.Println(issues.Report())
}
```

### Schema Compatibility Checking

Compare two metadata files for schema compatibility:

```
reference, err := croissant.LoadMetadataFromFile("reference.jsonld")
if err != nil {
	log.Fatalf("Error loading reference: %v", err)
}

candidate, err := croissant.LoadMetadataFromFile("candidate.jsonld")
if err != nil {
	log.Fatalf("Error loading candidate: %v", err)
}

result := croissant.MatchMetadata(*reference, *candidate)
if result.IsMatch {
	fmt.Printf("Compatible! %d fields matched\n", len(result.MatchedFields))
} else {
	fmt.Printf("Incompatible: %d missing, %d type mismatches\n",
		len(result.MissingFields), len(result.TypeMismatches))
}
```

### JSON\\\-LD Processing

Work directly with JSON\-LD data:

```
data, err := os.ReadFile("metadata.jsonld")
if err != nil {
	log.Fatal(err)
}

issues, err := croissant.ValidateJSON(data)
if err != nil {
	log.Fatalf("Validation error: %v", err)
}

fmt.Printf("Validation completed with %d errors and %d warnings\n",
	len(issues.Errors), len(issues.Warnings))
```

### Validation Options

Customize validation behavior:

```
options := croissant.ValidationOptions{
	StrictMode:      true,  // Enable additional warnings
	CheckDataTypes:  true,  // Validate data type specifications
	ValidateURLs:    false, // Skip network calls for URL validation
	CheckFileExists: true,  // Verify referenced files exist
}

issues, err := croissant.ValidateJSONWithOptions(data, options)
if err != nil {
	log.Fatal(err)
}
```

Package croissant provides comprehensive functionality for working with the ML Commons Croissant metadata format \- a standardized way to describe machine learning datasets using JSON\-LD.

### Overview

The Croissant metadata format is an open standard designed to improve dataset documentation, searchability, and usage in machine learning workflows. This package simplifies working with Croissant metadata by providing:

- Automatic metadata generation from CSV files with intelligent type inference
- Comprehensive validation tools with detailed error reporting
- Schema compatibility checking for dataset evolution
- Full JSON\-LD processing and validation support
- Extensible architecture supporting the complete Croissant specification

### Quick Start

Generate metadata from a CSV file:

```
outputPath, err := croissant.GenerateMetadata("data.csv", "metadata.jsonld")
if err != nil {
	log.Fatalf("Error: %v", err)
}
fmt.Printf("Metadata generated: %s\n", outputPath)
```

Validate existing metadata:

```
issues, err := croissant.ValidateFile("metadata.jsonld")
if err != nil {
	log.Fatalf("Validation error: %v", err)
}

if issues.HasErrors() {
	fmt.Println("Validation failed:")
	fmt.Println(issues.Report())
} else {
	fmt.Println("✓ Validation passed!")
}
```

Compare metadata files for compatibility:

```
ref, _ := croissant.LoadMetadataFromFile("reference.jsonld")
cand, _ := croissant.LoadMetadataFromFile("candidate.jsonld")

result := croissant.MatchMetadata(*ref, *cand)
if result.IsMatch {
	fmt.Printf("✓ Compatible schemas with %d matched fields\n", len(result.MatchedFields))
} else {
	fmt.Printf("✗ Incompatible: %d missing, %d type mismatches\n",
		len(result.MissingFields), len(result.TypeMismatches))
}
```

### Features

\#\# Metadata Generation

The package automatically generates Croissant\-compliant metadata from CSV files:

- Intelligent data type inference \(Boolean, Integer, Float, Date, URL, Text\)
- SHA\-256 hash calculation for file integrity verification
- Configurable output paths and validation options
- Support for environment variable configuration

\#\# Validation System

Comprehensive validation with multiple modes:

- JSON\-LD structure validation using the json\-gold library
- Croissant specification compliance checking
- Configurable validation modes \(standard, strict\)
- Optional file existence and URL accessibility verification
- Detailed error reporting with contextual information

\#\# Schema Compatibility

Advanced schema comparison for dataset evolution:

- Field\-by\-field compatibility analysis
- Intelligent type compatibility \(numeric type flexibility\)
- Support for schema evolution \(additional fields allowed\)
- Detailed reporting of matches, mismatches, and missing fields

### Data Type Inference

The package automatically maps CSV content to appropriate schema.org types:

```
Input Pattern              → Detected Type → Schema.org Type
===========================================================
true, false, 1, 0         → Boolean       → sc:Boolean
123, -456                 → Integer       → sc:Integer
3.14, 2.5e10             → Float         → sc:Float
2023-01-01, 01/15/2023   → Date          → sc:Date
https://example.com       → URL           → sc:URL
Everything else           → Text          → sc:Text
```

### Validation Options

Customize validation behavior using ValidationOptions:

```
options := croissant.ValidationOptions{
	StrictMode:      true,  // Enable additional warnings
	CheckDataTypes:  true,  // Validate data type specifications
	ValidateURLs:    false, // Skip network calls for URL validation
	CheckFileExists: true,  // Verify referenced files exist
}

issues, err := croissant.ValidateJSONWithOptions(data, options)
```

### Schema Compatibility Rules

When comparing metadata files, the following rules apply:

- All fields in the reference must exist in the candidate
- Field data types must be compatible \(exact or compatible numeric types\)
- Additional fields in the candidate are allowed \(backward compatibility\)
- Compatible numeric types: sc:Number accepts sc:Float and sc:Integer

### Error Handling

All functions follow Go error handling conventions. Common error scenarios:

- File I/O errors \(file not found, permission denied\)
- JSON parsing errors \(invalid JSON syntax\)
- JSON\-LD validation errors \(invalid JSON\-LD structure\)
- Croissant validation errors \(specification non\-compliance\)
- CSV parsing errors \(invalid structure or encoding\)

### Performance Considerations

- Metadata objects can be cached to avoid repeated file parsing
- Large CSV files are processed incrementally for memory efficiency
- URL validation is optional to avoid network latency
- File existence checks can be disabled for performance

### Examples

See the examples directory for comprehensive usage examples:

- Basic metadata generation and validation
- Advanced validation with custom options
- Schema compatibility checking
- Error handling patterns

### Related Tools

This package includes a command\-line tool \(gocroissant\) that provides:

- generate: Convert CSV files to Croissant metadata
- validate: Validate existing metadata files
- match: Compare metadata files for compatibility
- info: Analyze CSV file structure
- version: Display version information

### Specification Compliance

This implementation supports:

- Croissant specification version 1.0
- JSON\-LD 1.1 processing
- Schema.org vocabulary
- Full Croissant metadata structure

### License

MIT License \- see LICENSE file for details.

### Related Projects

- ML Commons Croissant: https://github.com/mlcommons/croissant
- Croissant Editor: Web\-based metadata editor
- Python Croissant: Python implementation

issues.go

jsonld.go

metadata\_node.go

node.go

structs.go

utils.go

validation.go

## Index

- [croissant](#croissant)
    - [Basic Usage](#basic-usage)
    - [Advanced Generation with Validation](#advanced-generation-with-validation)
    - [Data Type Inference](#data-type-inference)
    - [Validation](#validation)
    - [Schema Compatibility Checking](#schema-compatibility-checking)
    - [JSON\\-LD Processing](#json-ld-processing)
    - [Validation Options](#validation-options)
    - [Overview](#overview)
    - [Quick Start](#quick-start)
    - [Features](#features)
    - [Data Type Inference](#data-type-inference-1)
    - [Validation Options](#validation-options-1)
    - [Schema Compatibility Rules](#schema-compatibility-rules)
    - [Error Handling](#error-handling)
    - [Performance Considerations](#performance-considerations)
    - [Examples](#examples)
    - [Related Tools](#related-tools)
    - [Specification Compliance](#specification-compliance)
    - [License](#license)
    - [Related Projects](#related-projects)
  - [Index](#index)
  - [Constants](#constants)
  - [func CalculateSHA256](#func-calculatesha256)
  - [func CountCSVRows](#func-countcsvrows)
  - [func DetectCSVDelimiter](#func-detectcsvdelimiter)
  - [func ExtractCroissantProperties](#func-extractcroissantproperties)
  - [func GenerateMetadata](#func-generatemetadata)
  - [func GetCSVColumnTypes](#func-getcsvcolumntypes)
  - [func GetCSVColumns](#func-getcsvcolumns)
  - [func GetCSVColumnsAndSampleRows](#func-getcsvcolumnsandsamplerows)
  - [func GetExpandedProperty](#func-getexpandedproperty)
  - [func GetFileStats](#func-getfilestats)
  - [func GetPropertyArray](#func-getpropertyarray)
  - [func GetPropertyValue](#func-getpropertyvalue)
  - [func InferDataType](#func-inferdatatype)
  - [func InferSemanticDataType](#func-infersemanticdatatype)
  - [func IsCSVFile](#func-iscsvfile)
  - [func IsValidDataType](#func-isvaliddatatype)
  - [func ParseCSVWithOptions](#func-parsecsvwithoptions)
  - [func SanitizeFileName](#func-sanitizefilename)
  - [func ValidateCSVStructure](#func-validatecsvstructure)
  - [func ValidateCrossReferences](#func-validatecrossreferences)
  - [func ValidateDistributionNode](#func-validatedistributionnode)
  - [func ValidateFieldNode](#func-validatefieldnode)
  - [func ValidateMetadataNode](#func-validatemetadatanode)
  - [func ValidateOutputPath](#func-validateoutputpath)
  - [func ValidateRecordSetNode](#func-validaterecordsetnode)
  - [type BaseNode](#type-basenode)
    - [func (\*BaseNode) GetID](#func-basenode-getid)
    - [func (\*BaseNode) GetName](#func-basenode-getname)
    - [func (\*BaseNode) GetParent](#func-basenode-getparent)
    - [func (\*BaseNode) SetParent](#func-basenode-setparent)
  - [type Context](#type-context)
    - [func CreateDefaultContext](#func-createdefaultcontext)
  - [type CroissantError](#type-croissanterror)
    - [func (CroissantError) Error](#func-croissanterror-error)
  - [type DataContext](#type-datacontext)
  - [type DataType](#type-datatype)
    - [func NewArrayDataType](#func-newarraydatatype)
    - [func NewNullableSingleDataType](#func-newnullablesingledatatype)
    - [func NewSingleDataType](#func-newsingledatatype)
    - [func (DataType) GetFirstType](#func-datatype-getfirsttype)
    - [func (DataType) GetTypes](#func-datatype-gettypes)
    - [func (DataType) IsArray](#func-datatype-isarray)
    - [func (DataType) MarshalJSON](#func-datatype-marshaljson)
    - [func (\*DataType) UnmarshalJSON](#func-datatype-unmarshaljson)
  - [type DataTypeContext](#type-datatypecontext)
  - [type Distribution](#type-distribution)
  - [type DistributionNode](#type-distributionnode)
    - [func (\*DistributionNode) Validate](#func-distributionnode-validate)
  - [type Extract](#type-extract)
  - [type ExtractNode](#type-extractnode)
  - [type Field](#type-field)
  - [type FieldMismatch](#type-fieldmismatch)
  - [type FieldNode](#type-fieldnode)
    - [func (\*FieldNode) Validate](#func-fieldnode-validate)
  - [type FieldRef](#type-fieldref)
  - [type FieldRefSlice](#type-fieldrefslice)
    - [func (FieldRefSlice) MarshalJSON](#func-fieldrefslice-marshaljson)
    - [func (\*FieldRefSlice) UnmarshalJSON](#func-fieldrefslice-unmarshaljson)
  - [type FieldSource](#type-fieldsource)
    - [func (FieldRefSlice) MarshalJSON](#func-fieldrefslice-marshaljson-1)
    - [func (\*FieldRefSlice) UnmarshalJSON](#func-fieldrefslice-unmarshaljson-1)
  - [type FieldSource](#type-fieldsource-1)
    - [func (FieldSource) ValidateSource](#func-fieldsource-validatesource)
  - [type FileObject](#type-fileobject)
  - [type FileObjectRef](#type-fileobjectref)
  - [type Issue](#type-issue)
  - [type IssueType](#type-issuetype)
  - [type Issues](#type-issues)
    - [func NewIssues](#func-newissues)
    - [func ValidateFile](#func-validatefile)
    - [func ValidateJSON](#func-validatejson)
    - [func ValidateJSONWithOptions](#func-validatejsonwithoptions)
    - [func ValidateMetadata](#func-validatemetadata)
    - [func ValidateMetadataWithOptions](#func-validatemetadatawithoptions)
    - [func (\*Issues) AddError](#func-issues-adderror)
    - [func (\*Issues) AddWarning](#func-issues-addwarning)
    - [func (\*Issues) ErrorCount](#func-issues-errorcount)
    - [func (\*Issues) HasErrors](#func-issues-haserrors)
    - [func (\*Issues) HasWarnings](#func-issues-haswarnings)
    - [func (\*Issues) Report](#func-issues-report)
    - [func (\*Issues) WarningCount](#func-issues-warningcount)
  - [type JSONLDProcessor](#type-jsonldprocessor)
    - [func NewJSONLDProcessor](#func-newjsonldprocessor)
    - [func (\*JSONLDProcessor) CompactJSONLD](#func-jsonldprocessor-compactjsonld)
    - [func (\*JSONLDProcessor) ParseCroissantMetadata](#func-jsonldprocessor-parsecroissantmetadata)
    - [func (\*JSONLDProcessor) ParseJSONLD](#func-jsonldprocessor-parsejsonld)
    - [func (\*JSONLDProcessor) ValidateJSONLD](#func-jsonldprocessor-validatejsonld)
  - [type KeyRef](#type-keyref)
  - [type MatchResult](#type-matchresult)
    - [func MatchMetadata](#func-matchmetadata)
  - [type Metadata](#type-metadata)
    - [func LoadMetadataFromFile](#func-loadmetadatafromfile)
  - [type MetadataNode](#type-metadatanode)
    - [func FromMetadata](#func-frommetadata)
    - [func NewMetadataNode](#func-newmetadatanode)
    - [func (\*MetadataNode) Validate](#func-metadatanode-validate)
  - [type MetadataWithValidation](#type-metadatawithvalidation)
    - [func GenerateMetadataWithValidation](#func-generatemetadatawithvalidation)
    - [func NewMetadataWithValidation](#func-newmetadatawithvalidation)
    - [func (\*MetadataWithValidation) GetIssues](#func-metadatawithvalidation-getissues)
    - [func (\*MetadataWithValidation) HasErrors](#func-metadatawithvalidation-haserrors)
    - [func (\*MetadataWithValidation) HasWarnings](#func-metadatawithvalidation-haswarnings)
    - [func (\*MetadataWithValidation) Report](#func-metadatawithvalidation-report)
    - [func (\*MetadataWithValidation) Validate](#func-metadatawithvalidation-validate)
    - [func (\*MetadataWithValidation) ValidateWithOptions](#func-metadatawithvalidation-validatewithoptions)
  - [type Node](#type-node)
  - [type RecordSet](#type-recordset)
    - [func CreateEnumerationRecordSet](#func-createenumerationrecordset)
    - [func CreateSplitRecordSet](#func-createsplitrecordset)
  - [type RecordSetKey](#type-recordsetkey)
    - [func NewCompositeKey](#func-newcompositekey)
    - [func NewRecordSetKey](#func-newrecordsetkey)
    - [func (RecordSetKey) GetKeyIDs](#func-recordsetkey-getkeyids)
    - [func (RecordSetKey) IsComposite](#func-recordsetkey-iscomposite)
    - [func (RecordSetKey) MarshalJSON](#func-recordsetkey-marshaljson)
    - [func (\*RecordSetKey) UnmarshalJSON](#func-recordsetkey-unmarshaljson)
  - [type RecordSetNode](#type-recordsetnode)
    - [func (\*RecordSetNode) Validate](#func-recordsetnode-validate)
  - [type Source](#type-source)
  - [type SourceNode](#type-sourcenode)
    - [func (\*SourceNode) ValidateSource](#func-sourcenode-validatesource)
  - [type Transform](#type-transform)
  - [type ValidationOptions](#type-validationoptions)
    - [func DefaultValidationOptions](#func-defaultvalidationoptions)


## Constants

<a name="VT_crBBox"></a>

```go
const VT_crBBox string = "cr:BoundingBox"
```

<a name="VT_crLabel"></a>Croissant\-specific types.

```go
const VT_crLabel string = "cr:Label"
```

<a name="VT_crSegMask"></a>

```go
const VT_crSegMask string = "cr:SegmentationMask"
```

<a name="VT_crSplit"></a>

```go
const VT_crSplit string = "cr:Split"
```

<a name="VT_crSplitTest"></a>

```go
const VT_crSplitTest string = "cr:TestSplit"
```

<a name="VT_crSplitTrain"></a>Croissant Split types.

```go
const VT_crSplitTrain string = "cr:TrainingSplit"
```

<a name="VT_crSplitVal"></a>

```go
const VT_crSplitVal string = "cr:ValidationSplit"
```

<a name="VT_scBool"></a>

```go
const VT_scBool string = "sc:Boolean"
```

<a name="VT_scDateT"></a>

```go
const VT_scDateT string = "sc:DateTime"
```

<a name="VT_scEnum"></a>

```go
const VT_scEnum string = "sc:Enumeration"
```

<a name="VT_scGeoCoord"></a>

```go
const VT_scGeoCoord string = "sc:GeoCoordinates"
```

<a name="VT_scGeoShape"></a>

```go
const VT_scGeoShape string = "sc:GeoShape"
```

<a name="VT_scImage"></a>

```go
const VT_scImage string = "sc:ImageObject"
```

<a name="VT_scInt"></a>

```go
const VT_scInt string = "sc:Integer"
```

<a name="VT_scNum"></a>

```go
const VT_scNum string = "sc:Number"
```

<a name="VT_scText"></a>Schema.org data types.

```go
const VT_scText string = "sc:Text"
```

<a name="VT_scURL"></a>

```go
const VT_scURL string = "sc:URL"
```

<a name="VT_scVideo"></a>

```go
const VT_scVideo string = "sc:VideoObject"
```

<a name="VT_wdPrefix"></a>Wikidata entities \(wd:Q...\).

```go
const VT_wdPrefix string = "wd:Q"
```

<a name="CalculateSHA256"></a>
## func [CalculateSHA256](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/utils.go#L17>)

```go
func CalculateSHA256(filePath string) (string, error)
```

CalculateSHA256 calculates the SHA\-256 hash of a file.

<a name="CountCSVRows"></a>
## func [CountCSVRows](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/utils.go#L260>)

```go
func CountCSVRows(csvPath string) (int, error)
```

CountCSVRows counts the total number of rows in a CSV file \(including header\).

<a name="DetectCSVDelimiter"></a>
## func [DetectCSVDelimiter](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/utils.go#L148>)

```go
func DetectCSVDelimiter(csvPath string) (rune, error)
```

DetectCSVDelimiter attempts to detect the CSV delimiter.

<a name="ExtractCroissantProperties"></a>
## func [ExtractCroissantProperties](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/jsonld.go#L169>)

```go
func ExtractCroissantProperties(expanded map[string]interface{}) map[string]interface{}
```

ExtractCroissantProperties extracts common Croissant properties from expanded JSON\-LD.

<a name="GenerateMetadata"></a>
## func [GenerateMetadata](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/croissant.go#L380>)

```go
func GenerateMetadata(csvPath string, outputPath string) (string, error)
```

GenerateMetadata generates Croissant metadata from a CSV file \(simple API\).

<a name="GetCSVColumnTypes"></a>
## func [GetCSVColumnTypes](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/utils.go#L345>)

```go
func GetCSVColumnTypes(csvPath string, sampleSize int) ([]string, []string, error)
```

GetCSVColumnTypes analyzes a CSV file and returns inferred data types for each column.

<a name="GetCSVColumns"></a>
## func [GetCSVColumns](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/utils.go#L33>)

```go
func GetCSVColumns(csvPath string) ([]string, []string, error)
```

GetCSVColumns reads the column names and first row from a CSV file.

<a name="GetCSVColumnsAndSampleRows"></a>
## func [GetCSVColumnsAndSampleRows](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/utils.go#L69>)

```go
func GetCSVColumnsAndSampleRows(csvPath string, maxRows int) ([]string, [][]string, error)
```

GetCSVColumnsAndSampleRows reads column names and multiple sample rows for better type inference.

<a name="GetExpandedProperty"></a>
## func [GetExpandedProperty](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/jsonld.go#L97>)

```go
func GetExpandedProperty(expanded map[string]interface{}, property string) interface{}
```

GetExpandedProperty retrieves a property from expanded JSON\-LD using its full IRI.

<a name="GetFileStats"></a>
## func [GetFileStats](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/utils.go#L240>)

```go
func GetFileStats(filePath string) (map[string]interface{}, error)
```

GetFileStats returns basic statistics about a file.

<a name="GetPropertyArray"></a>
## func [GetPropertyArray](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/jsonld.go#L155>)

```go
func GetPropertyArray(property interface{}) []interface{}
```

GetPropertyArray extracts an array of values from a JSON\-LD property.

<a name="GetPropertyValue"></a>
## func [GetPropertyValue](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/jsonld.go#L121>)

```go
func GetPropertyValue(property interface{}) string
```

GetPropertyValue extracts a simple string value from a JSON\-LD property.

<a name="InferDataType"></a>
## func [InferDataType](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/croissant.go#L132>)

```go
func InferDataType(value string) string
```

InferDataType infers the schema.org data type from a value.

<a name="InferSemanticDataType"></a>
## func [InferSemanticDataType](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/croissant.go#L222>)

```go
func InferSemanticDataType(fieldName, value string, context map[string]interface{}) []string
```

InferSemanticDataType attempts to infer semantic data types for ML datasets.

<a name="IsCSVFile"></a>
## func [IsCSVFile](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/utils.go#L402>)

```go
func IsCSVFile(filePath string) bool
```

IsCSVFile checks if a file appears to be a CSV file based on extension.

<a name="IsValidDataType"></a>
## func [IsValidDataType](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/croissant.go#L186>)

```go
func IsValidDataType(dataType string) bool
```

IsValidDataType checks if a dataType is valid according to Croissant specification.

<a name="ParseCSVWithOptions"></a>
## func [ParseCSVWithOptions](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/utils.go#L186>)

```go
func ParseCSVWithOptions(csvPath string, delimiter rune, hasHeader bool) ([]string, [][]string, error)
```

ParseCSVWithOptions parses a CSV file with custom options.

<a name="SanitizeFileName"></a>
## func [SanitizeFileName](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/utils.go#L408>)

```go
func SanitizeFileName(fileName string) string
```

SanitizeFileName removes or replaces invalid characters in filenames.

<a name="ValidateCSVStructure"></a>
## func [ValidateCSVStructure](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/utils.go#L292>)

```go
func ValidateCSVStructure(csvPath string) error
```

ValidateCSVStructure performs basic validation on CSV file structure.

<a name="ValidateCrossReferences"></a>
## func [ValidateCrossReferences](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/validation.go#L324>)

```go
func ValidateCrossReferences(node *MetadataNode, issues *Issues)
```

ValidateCrossReferences validates that all references are valid.

<a name="ValidateDistributionNode"></a>
## func [ValidateDistributionNode](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/validation.go#L150>)

```go
func ValidateDistributionNode(dist *DistributionNode, issues *Issues, options ValidationOptions)
```

ValidateDistributionNode validates a distribution node.

<a name="ValidateFieldNode"></a>
## func [ValidateFieldNode](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/validation.go#L252>)

```go
func ValidateFieldNode(field *FieldNode, issues *Issues, options ValidationOptions)
```

ValidateFieldNode validates a field node.

<a name="ValidateMetadataNode"></a>
## func [ValidateMetadataNode](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/validation.go#L96>)

```go
func ValidateMetadataNode(node *MetadataNode, issues *Issues, options ValidationOptions)
```

ValidateMetadataNode performs comprehensive validation of a metadata node.

<a name="ValidateOutputPath"></a>
## func [ValidateOutputPath](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/utils.go#L118>)

```go
func ValidateOutputPath(outputPath string) error
```

ValidateOutputPath validates if the given path is a valid file path.

<a name="ValidateRecordSetNode"></a>
## func [ValidateRecordSetNode](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/validation.go#L192>)

```go
func ValidateRecordSetNode(rs *RecordSetNode, issues *Issues, options ValidationOptions)
```

ValidateRecordSetNode validates a record set node.

<a name="BaseNode"></a>
## type [BaseNode](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/node.go#L14-L18>)

BaseNode implements common functionality for all nodes.

```go
type BaseNode struct {
    ID   string `json:"@id,omitempty"`
    Name string `json:"name,omitempty"`
    // contains filtered or unexported fields
}
```

<a name="BaseNode.GetID"></a>
### func \(\*BaseNode\) [GetID](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/node.go#L24>)

```go
func (n *BaseNode) GetID() string
```



<a name="BaseNode.GetName"></a>
### func \(\*BaseNode\) [GetName](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/node.go#L20>)

```go
func (n *BaseNode) GetName() string
```



<a name="BaseNode.GetParent"></a>
### func \(\*BaseNode\) [GetParent](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/node.go#L28>)

```go
func (n *BaseNode) GetParent() Node
```



<a name="BaseNode.SetParent"></a>
### func \(\*BaseNode\) [SetParent](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/node.go#L32>)

```go
func (n *BaseNode) SetParent(parent Node)
```



<a name="Context"></a>
## type [Context](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/structs.go#L247-L283>)

Context represents the complete JSON\-LD context for Croissant 1.0.

```go
type Context struct {
    Language      string          `json:"@language"`
    Vocab         string          `json:"@vocab"`
    CiteAs        string          `json:"citeAs"`
    Column        string          `json:"column"`
    ConformsTo    string          `json:"conformsTo"`
    CR            string          `json:"cr"`
    DCT           string          `json:"dct"`
    RAI           string          `json:"rai,omitempty"`
    WD            string          `json:"wd,omitempty"`
    Data          DataContext     `json:"data"`
    DataType      DataTypeContext `json:"dataType"`
    Examples      DataContext     `json:"examples"`
    Extract       string          `json:"extract"`
    Field         string          `json:"field"`
    FileObject    string          `json:"fileObject"`
    FileProperty  string          `json:"fileProperty"`
    FileSet       string          `json:"fileSet"`
    Format        string          `json:"format"`
    Includes      string          `json:"includes"`
    IsLiveDataset string          `json:"isLiveDataset"`
    JSONPath      string          `json:"jsonPath"`
    Key           string          `json:"key"`
    MD5           string          `json:"md5"`
    ParentField   string          `json:"parentField"`
    Path          string          `json:"path"`
    RecordSet     string          `json:"recordSet"`
    References    string          `json:"references"`
    Regex         string          `json:"regex"`
    Repeated      string          `json:"repeated"`
    Replace       string          `json:"replace"`
    SC            string          `json:"sc"`
    Separator     string          `json:"separator"`
    Source        string          `json:"source"`
    SubField      string          `json:"subField"`
    Transform     string          `json:"transform"`
}
```

<a name="CreateDefaultContext"></a>
### func [CreateDefaultContext](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/croissant.go#L331>)

```go
func CreateDefaultContext() Context
```

CreateDefaultContext creates the ML Commons Croissant 1.0 compliant context.

<a name="CroissantError"></a>
## type [CroissantError](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/error.go#L5-L10>)



```go
type CroissantError struct {
    // Message to show the user.
    Message string
    // Value to include with message
    Value any
}
```

<a name="CroissantError.Error"></a>
### func \(CroissantError\) [Error](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/error.go#L12>)

```go
func (e CroissantError) Error() string
```



<a name="DataContext"></a>
## type [DataContext](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/structs.go#L286-L289>)

DataContext represents the data field in the context.

```go
type DataContext struct {
    ID   string `json:"@id"`
    Type string `json:"@type"`
}
```

<a name="DataType"></a>
## type [DataType](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/structs.go#L112>)

DataType represents a data type that can be either a single string or an array of strings. It is represented internally as a list.

```go
type DataType []string
```

<a name="NewArrayDataType"></a>
### func [NewArrayDataType](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/structs.go#L362>)

```go
func NewArrayDataType(dataTypes ...string) DataType
```

NewArrayDataType creates a DataType with multiple types.

<a name="NewNullableSingleDataType"></a>
### func [NewNullableSingleDataType](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/structs.go#L352>)

```go
func NewNullableSingleDataType(dataType string) *DataType
```

NewSingleDataType creates a DataType with a single type.

<a name="NewSingleDataType"></a>
### func [NewSingleDataType](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/structs.go#L357>)

```go
func NewSingleDataType(dataType string) DataType
```

NewSingleDataType creates a DataType with a single type.

<a name="DataType.GetFirstType"></a>
### func \(DataType\) [GetFirstType](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/structs.go#L211>)

```go
func (d DataType) GetFirstType() string
```

GetFirstType returns the first data type \(useful for backward compatibility\).

<a name="DataType.GetTypes"></a>
### func \(DataType\) [GetTypes](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/structs.go#L206>)

```go
func (d DataType) GetTypes() []string
```

GetTypes returns all data types \(single or array\).

<a name="DataType.IsArray"></a>
### func \(DataType\) [IsArray](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/structs.go#L201>)

```go
func (d DataType) IsArray() bool
```

IsArray returns true if this is an array of data types.

<a name="DataType.MarshalJSON"></a>
### func \(DataType\) [MarshalJSON](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/structs.go#L167>)

```go
func (d DataType) MarshalJSON() ([]byte, error)
```

MarshalJSON implements custom JSON marshaling for DataType.

<a name="DataType.UnmarshalJSON"></a>
### func \(\*DataType\) [UnmarshalJSON](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/structs.go#L179>)

```go
func (d *DataType) UnmarshalJSON(data []byte) error
```

UnmarshalJSON implements custom JSON unmarshaling for DataType.

<a name="DataTypeContext"></a>
## type [DataTypeContext](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/structs.go#L292-L295>)

DataTypeContext represents the dataType field in the context.

```go
type DataTypeContext struct {
    ID   string `json:"@id"`
    Type string `json:"@type"`
}
```

<a name="Distribution"></a>
## type [Distribution](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/structs.go#L220-L232>)

Distribution represents a file in the Croissant metadata.

```go
type Distribution struct {
    ID   string `json:"@id"`
    Type string `json:"@type"`
    // The name of the file.
    Name string `json:"name"`
    // Description of the file.
    Description string `json:"description,omitempty"`
    // File size in kb, mb, gb etc...
    // Defaults to bytes if unit not specified.
    ContentSize string `json:"contentSize,omitempty"`
    // URL to the actual bytes of the file object.
    ContentURL string `json:"contentUrl,omitempty"`
    // Format of the file, given as a MIME type.
    EncodingFormat string `json:"encodingFormat"`
    // SHA256 checksum of the file contents.
    SHA256 string `json:"sha256,omitempty"`
    // MD5 checksum of the file contents.
    MD5 string `json:"md5,omitempty"`
    // A FileObject or FileSet this resource is contained in.
    ContainedIn *FileObjectRef `json:"containedIn,omitempty"`
    // A glob pattern of the files to include (FileSet).
    Includes string `json:"includes,omitempty"`
    // A glob pattern of the files to exclude (FileSet).
    Excludes string `json:"excludes,omitempty"`
}
```

<a name="DistributionNode"></a>
## type [DistributionNode](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/metadata_node.go#L173-L181>)

DistributionNode represents a file distribution.

```go
type DistributionNode struct {
    BaseNode
    Type           string `json:"@type"`
    ContentSize    string `json:"contentSize,omitempty"`
    ContentURL     string `json:"contentUrl,omitempty"`
    EncodingFormat string `json:"encodingFormat,omitempty"`
    SHA256         string `json:"sha256,omitempty"`
    MD5            string `json:"md5,omitempty"`
}
```

<a name="DistributionNode.Validate"></a>
### func \(\*DistributionNode\) [Validate](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/metadata_node.go#L184>)

```go
func (d *DistributionNode) Validate(issues *Issues)
```

Validate validates the distribution node.

<a name="Extract"></a>
## type [Extract](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/structs.go#L33-L39>)

Extract represents the extraction information for a field source.

```go
type Extract struct {
    // Extraction method
    FileProperty string `json:"fileProperty,omitempty"`
    // Name of the column (field) that contains values.
    Column string `json:"column,omitempty"`
    // A JSONPATH expression that extracts values.
    JSONPath  string `json:"jsonPath,omitempty"`
    Regex     string `json:"regex,omitempty"`
    Separator string `json:"separator,omitempty"`
}
```

<a name="ExtractNode"></a>
## type [ExtractNode](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/metadata_node.go#L417-L422>)

ExtractNode represents extraction details.

```go
type ExtractNode struct {
    Regex        string `json:"regex,omitempty"`
    Column       string `json:"column,omitempty"`
    JSONPath     string `json:"jsonPath,omitempty"`
    FileProperty string `json:"fileProperty,omitempty"`
}
```

<a name="Field"></a>
## type [Field](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/structs.go#L9-L21>)

Field represents a field in the Croissant metadata.

```go
type Field struct {
    ID          string        `json:"@id"`
    Type        string        `json:"@type"`
    Name        string        `json:"name"`
    Description string        `json:"description,omitempty"`
    DataType    DataType      `json:"dataType"`
    Source      FieldSource   `json:"source,omitempty"`
    Repeated    bool          `json:"repeated,omitempty"`
    Examples    interface{}   `json:"examples,omitempty"`
    SubField    []Field       `json:"subField,omitempty"`
    ParentField FieldRefSlice `json:"parentField,omitempty"`
    References  FieldRefSlice `json:"references,omitempty"`
}
```

<a name="FieldMismatch"></a>
## type [FieldMismatch](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/match.go#L54-L63>)

FieldMismatch represents a field that exists in both metadata files but has incompatible data types. This indicates a schema compatibility issue that prevents the candidate from being used as a drop\-in replacement for the reference.

```go
type FieldMismatch struct {
    // FieldName is the name of the field that has a type mismatch.
    FieldName string

    // ReferenceType is the data type expected by the reference metadata.
    ReferenceType string

    // CandidateType is the data type found in the candidate metadata.
    CandidateType string
}
```

<a name="FieldNode"></a>
## type [FieldNode](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/metadata_node.go#L342-L352>)

FieldNode represents a field.

```go
type FieldNode struct {
    BaseNode
    Type        string        `json:"@type"`
    Description string        `json:"description,omitempty"`
    DataType    DataType      `json:"dataType,omitempty"`
    Source      SourceNode    `json:"source,omitempty"`
    Repeated    bool          `json:"repeated,omitempty"`
    Examples    interface{}   `json:"examples,omitempty"`
    SubField    []*FieldNode  `json:"subField,omitempty"`
    References  FieldRefSlice `json:"references,omitempty"`
}
```

<a name="FieldNode.Validate"></a>
### func \(\*FieldNode\) [Validate](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/metadata_node.go#L355>)

```go
func (f *FieldNode) Validate(issues *Issues)
```

Validate validates the field node.

<a name="FieldRef"></a>
## type [FieldRef](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/structs.go#L52-L55>)

FieldRef represents a reference to another field.

```go
type FieldRef struct {
    ID    string  `json:"@id,omitempty"`
    Field *KeyRef `json:"field,omitempty"`
}
```

<a name="FieldRefSlice"></a>
## type [FieldRefSlice](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/structs.go#L58>)

Parses ONE or MANY FieldRefs.

```go
type FieldRefSlice []FieldRef
```

<a name="FieldRefSlice.MarshalJSON"></a>
### func \(FieldRefSlice\) [MarshalJSON](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/structs.go#L99>)

```go
func (ref FieldRefSlice) MarshalJSON() ([]byte, error)
```



<a name="FieldRefSlice.UnmarshalJSON"></a>
### func \(\*FieldRefSlice\) [UnmarshalJSON](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/structs.go#L67>)

```go
func (ref *FieldRefSlice) UnmarshalJSON(data []byte) error
```

In some test files, references are nested under a "field" property. In cases of reformatting, the property will be omitted.

Accepts:

- "references": \{ "@id": "..." \}
- "references": \{ \[\{"@id": "..."\}, \{"@id": "..."\}...\] \}
- "references": \{ field: \{"@id": "..."\} \}

<a name="FieldSource"></a>
## type [FieldSource](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/structs.go#L24-L30>)

Parses ONE or MANY FieldRefs.

```go
type FieldRefSlice []FieldRef
```

<a name="FieldRefSlice.MarshalJSON"></a>
### func \(FieldRefSlice\) [MarshalJSON](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/structs.go#L96>)

```go
func (ref FieldRefSlice) MarshalJSON() ([]byte, error)
```



<a name="FieldRefSlice.UnmarshalJSON"></a>
### func \(\*FieldRefSlice\) [UnmarshalJSON](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/structs.go#L74>)

```go
func (ref *FieldRefSlice) UnmarshalJSON(data []byte) error
```



<a name="FieldSource"></a>
## type [FieldSource](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/structs.go#L35-L41>)

FieldSource represents the source information for a field.

```go
type FieldSource struct {
    Extract    *Extract    `json:"extract,omitempty"`
    FileObject *FileObject `json:"fileObject,omitempty"`
    FileSet    *FileObject `json:"fileSet,omitempty"`
    Transform  *Transform  `json:"transform,omitempty"`
    Format     string      `json:"format,omitempty"`
}
```

<a name="FieldSource.ValidateSource"></a>
### func \(FieldSource\) [ValidateSource](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/structs.go#L367>)

```go
func (fs FieldSource) ValidateSource() bool
```

ValidateSource validates the source configuration.

<a name="FileObject"></a>
## type [FileObject](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/structs.go#L42-L44>)

FileObject represents a file object reference.

```go
type FileObject struct {
    ID string `json:"@id"`
}
```

<a name="FileObjectRef"></a>
## type [FileObjectRef](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/metadata_node.go#L425-L427>)

FileObjectRef represents a reference to a file object.

```go
type FileObjectRef struct {
    ID string `json:"@id"`
}
```

<a name="Issue"></a>
## type [Issue](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/issues.go#L19-L23>)

Issue represents a single validation issue.

```go
type Issue struct {
    Type    IssueType
    Message string
    Context string // For context like "Metadata(mydataset) > FileObject(a-csv-table)"
}
```

<a name="IssueType"></a>
## type [IssueType](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/issues.go#L11>)

IssueType represents the type of issue \(error or warning\).

```go
type IssueType int
```

<a name="ErrorIssue"></a>

```go
const (
    ErrorIssue IssueType = iota
    WarningIssue
)
```

<a name="Issues"></a>
## type [Issues](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/issues.go#L26-L29>)

Issues represents a collection of validation issues.

```go
type Issues struct {
    // contains filtered or unexported fields
}
```

<a name="NewIssues"></a>
### func [NewIssues](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/issues.go#L32>)

```go
func NewIssues() *Issues
```

NewIssues creates a new Issues instance.

<a name="ValidateFile"></a>
### func [ValidateFile](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/validation.go#L32>)

```go
func ValidateFile(filePath string) (*Issues, error)
```

ValidateFile validates a Croissant metadata file and returns issues.

<a name="ValidateJSON"></a>
### func [ValidateJSON](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/validation.go#L42>)

```go
func ValidateJSON(data []byte) (*Issues, error)
```

ValidateJSON validates Croissant metadata in JSON\-LD format and returns issues.

<a name="ValidateJSONWithOptions"></a>
### func [ValidateJSONWithOptions](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/validation.go#L61>)

```go
func ValidateJSONWithOptions(data []byte, options ValidationOptions) (*Issues, error)
```

ValidateJSONWithOptions validates Croissant metadata in JSON\-LD format with options and returns issues.

<a name="ValidateMetadata"></a>
### func [ValidateMetadata](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/validation.go#L80>)

```go
func ValidateMetadata(metadata Metadata) *Issues
```

ValidateMetadata validates a Metadata struct and returns issues.

<a name="ValidateMetadataWithOptions"></a>
### func [ValidateMetadataWithOptions](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/validation.go#L85>)

```go
func ValidateMetadataWithOptions(metadata Metadata, options ValidationOptions) *Issues
```

ValidateMetadataWithOptions validates a Metadata struct with specific options.

<a name="Issues.AddError"></a>
### func \(\*Issues\) [AddError](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/issues.go#L40>)

```go
func (i *Issues) AddError(message string, node ...Node)
```

AddError adds a new error to the issues collection.

<a name="Issues.AddWarning"></a>
### func \(\*Issues\) [AddWarning](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/issues.go#L52>)

```go
func (i *Issues) AddWarning(message string, node ...Node)
```

AddWarning adds a new warning to the issues collection.

<a name="Issues.ErrorCount"></a>
### func \(\*Issues\) [ErrorCount](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/issues.go#L74>)

```go
func (i *Issues) ErrorCount() int
```

ErrorCount returns the number of errors.

<a name="Issues.HasErrors"></a>
### func \(\*Issues\) [HasErrors](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/issues.go#L64>)

```go
func (i *Issues) HasErrors() bool
```

HasErrors returns true if there are any errors.

<a name="Issues.HasWarnings"></a>
### func \(\*Issues\) [HasWarnings](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/issues.go#L69>)

```go
func (i *Issues) HasWarnings() bool
```

HasWarnings returns true if there are any warnings.

<a name="Issues.Report"></a>
### func \(\*Issues\) [Report](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/issues.go#L84>)

```go
func (i *Issues) Report() string
```

Report generates a human\-readable report of all issues.

<a name="Issues.WarningCount"></a>
### func \(\*Issues\) [WarningCount](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/issues.go#L79>)

```go
func (i *Issues) WarningCount() int
```

WarningCount returns the number of warnings.

<a name="JSONLDProcessor"></a>
## type [JSONLDProcessor](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/jsonld.go#L12-L15>)

JSONLDProcessor handles JSON\-LD processing using json\-gold library.

```go
type JSONLDProcessor struct {
    // contains filtered or unexported fields
}
```

<a name="NewJSONLDProcessor"></a>
### func [NewJSONLDProcessor](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/jsonld.go#L18>)

```go
func NewJSONLDProcessor() *JSONLDProcessor
```

NewJSONLDProcessor creates a new JSON\-LD processor.

<a name="JSONLDProcessor.CompactJSONLD"></a>
### func \(\*JSONLDProcessor\) [CompactJSONLD](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/jsonld.go#L53>)

```go
func (j *JSONLDProcessor) CompactJSONLD(expanded interface{}, context map[string]interface{}) (map[string]interface{}, error)
```

CompactJSONLD compacts an expanded JSON\-LD document with the given context.

<a name="JSONLDProcessor.ParseCroissantMetadata"></a>
### func \(\*JSONLDProcessor\) [ParseCroissantMetadata](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/jsonld.go#L80>)

```go
func (j *JSONLDProcessor) ParseCroissantMetadata(data []byte) (*Metadata, error)
```

ParseCroissantMetadata parses Croissant JSON\-LD metadata and converts it to our Metadata struct.

<a name="JSONLDProcessor.ParseJSONLD"></a>
### func \(\*JSONLDProcessor\) [ParseJSONLD](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/jsonld.go#L26>)

```go
func (j *JSONLDProcessor) ParseJSONLD(data []byte) (map[string]interface{}, error)
```

ParseJSONLD parses and expands JSON\-LD document to a normalized form.

<a name="JSONLDProcessor.ValidateJSONLD"></a>
### func \(\*JSONLDProcessor\) [ValidateJSONLD](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/jsonld.go#L64>)

```go
func (j *JSONLDProcessor) ValidateJSONLD(data []byte) error
```

ValidateJSONLD validates that the document is valid JSON\-LD.

<a name="KeyRef"></a>
## type [KeyRef](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/structs.go#L47-L49>)

KeyRef represents a key reference in a composite key.

```go
type KeyRef struct {
    ID string `json:"@id"`
}
```

<a name="MatchResult"></a>
## type [MatchResult](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/match.go#L29-L49>)

MatchResult represents the result of comparing two Croissant metadata files for schema compatibility. It provides detailed information about field matches, mismatches, and additional fields.

The comparison checks whether a candidate metadata file is compatible with a reference metadata file. Compatibility means:

- All fields from the reference must exist in the candidate
- Field data types must be compatible \(exact match or compatible numeric types\)
- The candidate may have additional fields \(this doesn't affect compatibility\)

Example usage:

```
ref, _ := croissant.LoadMetadataFromFile("reference.jsonld")
cand, _ := croissant.LoadMetadataFromFile("candidate.jsonld")
result := croissant.MatchMetadata(*ref, *cand)

if result.IsMatch {
	fmt.Printf("Compatible! %d fields matched\n", len(result.MatchedFields))
} else {
	fmt.Printf("Issues: %d missing, %d type mismatches\n",
		len(result.MissingFields), len(result.TypeMismatches))
}
```

```go
type MatchResult struct {
    // IsMatch indicates whether the candidate is compatible with the reference.
    // True if all reference fields exist in candidate with compatible types.
    IsMatch bool

    // MissingFields lists field names that exist in reference but not in candidate.
    // These represent compatibility violations.
    MissingFields []string

    // TypeMismatches lists fields that exist in both files but have incompatible data types.
    // These represent compatibility violations.
    TypeMismatches []FieldMismatch

    // ExtraFields lists field names that exist in candidate but not in reference.
    // These do not affect compatibility but may be useful for information.
    ExtraFields []string

    // MatchedFields lists field names that exist in both files with compatible types.
    // These represent successful matches.
    MatchedFields []string
}
```

<a name="MatchMetadata"></a>
### func [MatchMetadata](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/match.go#L114>)

```go
func MatchMetadata(reference Metadata, candidate Metadata) *MatchResult
```

MatchMetadata compares two Croissant metadata objects to check if the candidate is compatible with the reference. The candidate can have additional fields, but all reference fields must exist in the candidate with matching data types.

Compatibility Rules:

- All fields in the reference must exist in the candidate
- Field data types must be compatible \(see type compatibility below\)
- Additional fields in the candidate are allowed and don't affect compatibility

Type Compatibility:

- Exact type matches \(sc:Text = sc:Text\)
- Numeric type compatibility \(sc:Number accepts sc:Float, sc:Integer\)
- Schema.org prefix normalization \(sc:Text = https://schema.org/Text\)

The function returns a MatchResult containing detailed information about:

- Whether the schemas are compatible \(IsMatch\)
- Successfully matched fields \(MatchedFields\)
- Missing required fields \(MissingFields\)
- Type mismatches \(TypeMismatches\)
- Additional fields in candidate \(ExtraFields\)

Example:

```
reference, err := croissant.LoadMetadataFromFile("reference.jsonld")
if err != nil {
	log.Fatal(err)
}

candidate, err := croissant.LoadMetadataFromFile("candidate.jsonld")
if err != nil {
	log.Fatal(err)
}

result := croissant.MatchMetadata(*reference, *candidate)

if result.IsMatch {
	fmt.Printf("✓ Compatible schemas with %d matched fields\n", len(result.MatchedFields))
	if len(result.ExtraFields) > 0 {
		fmt.Printf("  Candidate has %d additional fields\n", len(result.ExtraFields))
	}
} else {
	fmt.Printf("✗ Incompatible schemas:\n")
	if len(result.MissingFields) > 0 {
		fmt.Printf("  Missing %d required fields\n", len(result.MissingFields))
	}
	if len(result.TypeMismatches) > 0 {
		fmt.Printf("  %d type mismatches found\n", len(result.TypeMismatches))
	}
}
```

<a name="Metadata"></a>
## type [Metadata](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/structs.go#L298-L315>)

Metadata represents the complete Croissant metadata for a dataset.

```go
type Metadata struct {
    Context Context `json:"@context"`
    // Dataset Type.  Must by `schema.org/Dataset`
    Type string `json:"@type"`
    // Name of the dataset.
    Name string `json:"name"`
    // Description of the dataset.
    Description string `json:"description,omitempty"`
    // Versioned schema the croissant metadata conforms to.
    ConformsTo string `json:"conformsTo"`
    // Date the dataset was published.
    DatePublished string `json:"datePublished,omitempty"`
    // Version of the dataset.
    // Either an single int, or a MAJOR.MINOR.PATCH sematic version string.
    // [Semantic Versioning 2.0.0](https://semver.org/spec/v2.0.0.html)
    Version string `json:"version,omitempty"`
    // Url of the dataset, usually a webpage.
    URL string `json:"url,omitempty"`
    // Licenses of the dataset.
    // Spec suggests using references from https://spdx.org/licenses/.
    License string `json:"license,omitempty"`
    // A citation to the dataset itself, or a citation for a publication that describes the dataset.
    // Ideally, citations should be expressed using the bibtex format.
    // Note that this is different from schema.org/citation, which is used to make a citation to another publication from this dataset.
    CiteAs string `json:"citeAs,omitempty"`
    // Creator(s) of the dataset.
    Creator interface{} `json:"creator,omitempty"`
    // Publisher(s) of the dataset.
    Publisher interface{} `json:"publisher,omitempty"`
    // A set of keywords associated with the dataset, either as free text, or a DefinedTerm with a formal definition.
    Keywords []string `json:"keywords,omitempty"`
    // Set of FileObject and FileSet definitions that describe the raw files of the dataset.
    Distributions []Distribution `json:"distribution"`
    // Set of RecordSet definitions that describe the content of the dataset.
    RecordSets []RecordSet `json:"recordSet"`
    // If true, dataset is non-static and may change over time.
    // Distribution resources may not contain a checksum if they are expected to change.
    IsLiveDataset bool `json:"isLiveDataset,omitempty"`
}
```

<a name="LoadMetadataFromFile"></a>
### func [LoadMetadataFromFile](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/match.go#L271>)

```go
func LoadMetadataFromFile(filePath string) (*Metadata, error)
```

LoadMetadataFromFile loads and parses a Croissant metadata file from disk. It validates the JSON\-LD structure and parses it into a Metadata object.

The function performs the following steps:

1. Reads the file from the specified path
2. Validates that the content is valid JSON\-LD using the json\-gold library
3. Parses the JSON\-LD into a Croissant Metadata structure
4. Returns the parsed metadata or an error if any step fails

Supported file formats:

- JSON\-LD files \(.jsonld, .json\)
- Must conform to Croissant metadata specification
- Must be valid JSON\-LD documents

Example:

```
metadata, err := croissant.LoadMetadataFromFile("dataset.jsonld")
if err != nil {
	log.Fatalf("Failed to load metadata: %v", err)
}

fmt.Printf("Loaded dataset: %s\n", metadata.Name)
fmt.Printf("Record sets: %d\n", len(metadata.RecordSets))
fmt.Printf("Distributions: %d\n", len(metadata.Distributions))
```

Common errors:

- File not found or permission denied
- Invalid JSON syntax
- Invalid JSON\-LD structure
- Non\-compliant Croissant metadata format

<a name="MetadataNode"></a>
## type [MetadataNode](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/metadata_node.go#L7-L18>)

MetadataNode represents a Croissant metadata document.

```go
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
```

<a name="FromMetadata"></a>
### func [FromMetadata](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/metadata_node.go#L65>)

```go
func FromMetadata(metadata Metadata) *MetadataNode
```

FromMetadata converts a Metadata struct to a MetadataNode.

<a name="NewMetadataNode"></a>
### func [NewMetadataNode](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/metadata_node.go#L21>)

```go
func NewMetadataNode() *MetadataNode
```

NewMetadataNode creates a new MetadataNode.

<a name="MetadataNode.Validate"></a>
### func \(\*MetadataNode\) [Validate](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/metadata_node.go#L35>)

```go
func (m *MetadataNode) Validate(issues *Issues)
```

Validate validates the metadata node.

<a name="MetadataWithValidation"></a>
## type [MetadataWithValidation](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/validation.go#L370-L374>)

AddValidationToMetadata adds validation functionality to the Metadata struct.

```go
type MetadataWithValidation struct {
    Metadata
    // contains filtered or unexported fields
}
```

<a name="GenerateMetadataWithValidation"></a>
### func [GenerateMetadataWithValidation](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/croissant.go#L395>)

```go
func GenerateMetadataWithValidation(csvPath string, outputPath string) (*MetadataWithValidation, error)
```

GenerateMetadataWithValidation generates Croissant metadata with validation from a CSV file.

<a name="NewMetadataWithValidation"></a>
### func [NewMetadataWithValidation](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/validation.go#L377>)

```go
func NewMetadataWithValidation(metadata Metadata) *MetadataWithValidation
```

NewMetadataWithValidation creates a new MetadataWithValidation instance.

<a name="MetadataWithValidation.GetIssues"></a>
### func \(\*MetadataWithValidation\) [GetIssues](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/validation.go#L420>)

```go
func (m *MetadataWithValidation) GetIssues() *Issues
```

GetIssues returns the validation issues.

<a name="MetadataWithValidation.HasErrors"></a>
### func \(\*MetadataWithValidation\) [HasErrors](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/validation.go#L404>)

```go
func (m *MetadataWithValidation) HasErrors() bool
```

HasErrors returns true if there are validation errors.

<a name="MetadataWithValidation.HasWarnings"></a>
### func \(\*MetadataWithValidation\) [HasWarnings](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/validation.go#L412>)

```go
func (m *MetadataWithValidation) HasWarnings() bool
```

HasWarnings returns true if there are validation warnings.

<a name="MetadataWithValidation.Report"></a>
### func \(\*MetadataWithValidation\) [Report](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/validation.go#L396>)

```go
func (m *MetadataWithValidation) Report() string
```

Report returns a string report of validation issues.

<a name="MetadataWithValidation.Validate"></a>
### func \(\*MetadataWithValidation\) [Validate](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/validation.go#L385>)

```go
func (m *MetadataWithValidation) Validate()
```

Validate runs validation on the metadata.

<a name="MetadataWithValidation.ValidateWithOptions"></a>
### func \(\*MetadataWithValidation\) [ValidateWithOptions](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/validation.go#L390>)

```go
func (m *MetadataWithValidation) ValidateWithOptions(options ValidationOptions)
```

ValidateWithOptions runs validation with specific options.

<a name="Node"></a>
## type [Node](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/node.go#L5-L11>)

Node represents a node in the Croissant metadata structure.

```go
type Node interface {
    GetName() string
    GetID() string
    GetParent() Node
    SetParent(Node)
    Validate(*Issues)
}
```

<a name="RecordSet"></a>
## type [RecordSet](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/structs.go#L235-L244>)

RecordSet represents a record set in the Croissant metadata.

```go
type RecordSet struct {
    ID          string                   `json:"@id"`
    Type        string                   `json:"@type"`
    Name        string                   `json:"name"`
    Description string                   `json:"description,omitempty"`
    DataType    *DataType                `json:"dataType,omitempty"`
    Fields      []Field                  `json:"field"`
    Key         *RecordSetKey            `json:"key,omitempty"`
    Data        []map[string]interface{} `json:"data,omitempty"`
}
```

<a name="CreateEnumerationRecordSet"></a>
### func [CreateEnumerationRecordSet](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/croissant.go#L268>)

```go
func CreateEnumerationRecordSet(id, name string, values []string, urls []string) RecordSet
```

CreateEnumerationRecordSet creates a RecordSet for categorical/enumeration data.

<a name="CreateSplitRecordSet"></a>
### func [CreateSplitRecordSet](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/croissant.go#L316>)

```go
func CreateSplitRecordSet() RecordSet
```

CreateSplitRecordSet creates a standard ML split RecordSet.

<a name="RecordSetKey"></a>
## type [RecordSetKey](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/structs.go#L115>)

RecordSetKey represents a record set key that can be either a single key or composite key.

```go
type RecordSetKey []KeyRef
```

<a name="NewCompositeKey"></a>
### func [NewCompositeKey](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/structs.go#L343>)

```go
func NewCompositeKey(keyIDs ...string) *RecordSetKey
```

NewCompositeKey creates a RecordSetKey with multiple key references.

<a name="NewRecordSetKey"></a>
### func [NewRecordSetKey](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/structs.go#L336>)

```go
func NewRecordSetKey(keyID string) *RecordSetKey
```

NewRecordSetKey creates a RecordSetKey with a single key reference

<a name="RecordSetKey.GetKeyIDs"></a>
### func \(RecordSetKey\) [GetKeyIDs](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/structs.go#L155>)

```go
func (k RecordSetKey) GetKeyIDs() []string
```

GetKeyIDs returns all key IDs \(single or composite\).

<a name="RecordSetKey.IsComposite"></a>
### func \(RecordSetKey\) [IsComposite](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/structs.go#L150>)

```go
func (k RecordSetKey) IsComposite() bool
```

IsComposite returns true if this is a composite key.

<a name="RecordSetKey.MarshalJSON"></a>
### func \(RecordSetKey\) [MarshalJSON](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/structs.go#L118>)

```go
func (key RecordSetKey) MarshalJSON() ([]byte, error)
```

MarshalJSON implements custom JSON marshaling for RecordSetKey.

<a name="RecordSetKey.UnmarshalJSON"></a>
### func \(\*RecordSetKey\) [UnmarshalJSON](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/structs.go#L130>)

```go
func (key *RecordSetKey) UnmarshalJSON(data []byte) error
```

UnmarshalJSON implements custom JSON unmarshaling for RecordSetKey.

<a name="RecordSetNode"></a>
## type [RecordSetNode](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/metadata_node.go#L207-L215>)

RecordSetNode represents a record set.

```go
type RecordSetNode struct {
    BaseNode
    Type        string                   `json:"@type"`
    Description string                   `json:"description,omitempty"`
    DataType    DataType                 `json:"dataType,omitempty"`
    Fields      []*FieldNode             `json:"field"`
    Key         *RecordSetKey            `json:"key,omitempty"`
    Data        []map[string]interface{} `json:"data,omitempty"`
}
```

<a name="RecordSetNode.Validate"></a>
### func \(\*RecordSetNode\) [Validate](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/metadata_node.go#L218>)

```go
func (r *RecordSetNode) Validate(issues *Issues)
```

Validate validates the record set node.

<a name="Source"></a>
## type [Source](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/structs.go#L328-L333>)

Source represents a more complete source definition.

```go
type Source struct {
    Extract    *Extract    `json:"extract,omitempty"`
    FileObject *FileObject `json:"fileObject,omitempty"`
    Field      string      `json:"field,omitempty"`
    Transform  []Transform `json:"transform,omitempty"`
}
```

<a name="SourceNode"></a>
## type [SourceNode](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/metadata_node.go#L399-L405>)

SourceNode represents a source.

```go
type SourceNode struct {
    Extract    *ExtractNode   `json:"extract,omitempty"`
    FileObject *FileObjectRef `json:"fileObject,omitempty"`
    FileSet    *FileObjectRef `json:"fileSet,omitempty"`
    Transform  *Transform     `json:"transform,omitempty"`
    Format     string         `json:"format,omitempty"`
}
```

<a name="SourceNode.ValidateSource"></a>
### func \(\*SourceNode\) [ValidateSource](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/metadata_node.go#L408>)

```go
func (s *SourceNode) ValidateSource() bool
```

ValidateSource validates the source node.

<a name="Transform"></a>
## type [Transform](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/structs.go#L318-L325>)

Transform represents a data transformation.

```go
type Transform struct {
    Type      string `json:"@type"`
    Regex     string `json:"regex,omitempty"`
    Replace   string `json:"replace,omitempty"`
    Format    string `json:"format,omitempty"`
    JSONPath  string `json:"jsonPath,omitempty"`
    Separator string `json:"separator,omitempty"`
}
```

<a name="ValidationOptions"></a>
## type [ValidationOptions](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/validation.go#L14-L19>)

ValidationOptions represents options for validation.

```go
type ValidationOptions struct {
    StrictMode      bool
    CheckDataTypes  bool
    ValidateURLs    bool
    CheckFileExists bool
}
```

<a name="DefaultValidationOptions"></a>
### func [DefaultValidationOptions](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/validation.go#L22>)

```go
func DefaultValidationOptions() ValidationOptions
```

DefaultValidationOptions returns default validation options.