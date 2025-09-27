// utils.go
package croissant

import (
	"crypto/sha256"
	"encoding/csv"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// CalculateSHA256 calculates the SHA-256 hash of a file
func CalculateSHA256(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", CroissantError{Message: "failed to open file: %w", Value: err}
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", CroissantError{Message: "failed to calculate hash: %w", Value: err}
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// GetCSVColumns reads the column names and first row from a CSV file
func GetCSVColumns(csvPath string) ([]string, []string, error) {
	file, err := os.Open(csvPath)
	if err != nil {
		return nil, nil, CroissantError{Message: "failed to open CSV file: %w", Value: err}
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true

	// Read headers
	headers, err := reader.Read()
	if err != nil {
		return nil, nil, CroissantError{Message: "failed to read CSV headers: %w", Value: err}
	}

	// Clean headers
	for i, header := range headers {
		headers[i] = strings.TrimSpace(header)
	}

	// Read first row for data type inference
	firstRow, err := reader.Read()
	if err != nil && !errors.Is(err, io.EOF) {
		return nil, nil, CroissantError{Message: "failed to read first CSV row: %w", Value: err}
	}

	// If we hit EOF, there's no data row
	if errors.Is(err, io.EOF) {
		return headers, nil, nil
	}

	return headers, firstRow, nil
}

// GetCSVColumnsAndSampleRows reads column names and multiple sample rows for better type inference
func GetCSVColumnsAndSampleRows(csvPath string, maxRows int) ([]string, [][]string, error) {
	file, err := os.Open(csvPath)
	if err != nil {
		return nil, nil, CroissantError{Message: "failed to open CSV file: %w", Value: err}
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true

	// Handle common CSV format issues
	reader.LazyQuotes = true
	reader.FieldsPerRecord = -1 // Allow variable number of fields

	// Read headers
	headers, err := reader.Read()
	if err != nil {
		return nil, nil, CroissantError{Message: "failed to read CSV headers: %w", Value: err}
	}

	// Clean headers
	for i, header := range headers {
		headers[i] = strings.TrimSpace(header)
	}

	// Read sample rows
	var sampleRows [][]string
	rowCount := 0

	for rowCount < maxRows {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, nil, CroissantError{
				Message: fmt.Sprintf("failed to read CSV row %d", rowCount+1),
				Value:   err,
			}
		}

		sampleRows = append(sampleRows, row)
		rowCount++
	}

	return headers, sampleRows, nil
}

// ValidateOutputPath validates if the given path is a valid file path
func ValidateOutputPath(outputPath string) error {
	if outputPath == "" {
		return CroissantError{Message: "output path cannot be empty"}
	}

	// Check if the directory exists or can be created
	dir := filepath.Dir(outputPath)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return CroissantError{
				Message: fmt.Sprintf("cannot create directory %s", dir),
				Value:   err,
			}
		}
	}

	// Check if we can write to the file (create a temporary file to test)
	tempFile := outputPath + ".tmp"
	file, err := os.Create(tempFile)
	if err != nil {
		return CroissantError{
			Message: fmt.Sprintf("cannot write to path %s", outputPath),
			Value:   err,
		}
	}
	file.Close()
	return os.Remove(tempFile) // Clean up the temporary file
}

// DetectCSVDelimiter attempts to detect the CSV delimiter
func DetectCSVDelimiter(csvPath string) (rune, error) {
	file, err := os.Open(csvPath)
	if err != nil {
		return ',', CroissantError{Message: "failed to open CSV file: %w", Value: err}
	}
	defer file.Close()

	// Read first few lines to detect delimiter
	buffer := make([]byte, 1024)
	n, err := file.Read(buffer)
	if err != nil && !errors.Is(err, io.EOF) {
		return ',', CroissantError{Message: "failed to read file sample: %w", Value: err}
	}

	sample := string(buffer[:n])

	// Count occurrences of common delimiters
	delimiters := map[rune]int{
		',':  strings.Count(sample, ","),
		';':  strings.Count(sample, ";"),
		'\t': strings.Count(sample, "\t"),
		'|':  strings.Count(sample, "|"),
	}

	// Find the most common delimiter
	maxCount := 0
	bestDelimiter := ','
	for delimiter, count := range delimiters {
		if count > maxCount {
			maxCount = count
			bestDelimiter = delimiter
		}
	}

	return bestDelimiter, nil
}

// ParseCSVWithOptions parses a CSV file with custom options
func ParseCSVWithOptions(csvPath string, delimiter rune, hasHeader bool) ([]string, [][]string, error) {
	file, err := os.Open(csvPath)
	if err != nil {
		return nil, nil, CroissantError{Message: "failed to open CSV file: %w", Value: err}
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = delimiter
	reader.TrimLeadingSpace = true
	reader.LazyQuotes = true
	reader.FieldsPerRecord = -1

	var headers []string
	var rows [][]string

	// Read all records
	records, err := reader.ReadAll()
	if err != nil {
		return nil, nil, CroissantError{
			Message: "failed to read CSV records",
			Value:   err,
		}
	}

	if len(records) == 0 {
		return nil, nil, CroissantError{Message: "CSV file is empty"}
	}

	if hasHeader {
		headers = records[0]
		rows = records[1:]
	} else {
		// Generate default headers
		if len(records) > 0 {
			for i := 0; i < len(records[0]); i++ {
				headers = append(headers, fmt.Sprintf("column_%d", i+1))
			}
		}
		rows = records
	}

	// Clean headers
	for i, header := range headers {
		headers[i] = strings.TrimSpace(header)
		if headers[i] == "" {
			headers[i] = fmt.Sprintf("column_%d", i+1)
		}
	}

	return headers, rows, nil
}

// GetFileStats returns basic statistics about a file
func GetFileStats(filePath string) (map[string]interface{}, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, CroissantError{Message: "failed to get file stats: %w", Value: err}
	}

	stats := map[string]interface{}{
		"name":      fileInfo.Name(),
		"size":      fileInfo.Size(),
		"mode":      fileInfo.Mode(),
		"modTime":   fileInfo.ModTime(),
		"isDir":     fileInfo.IsDir(),
		"extension": filepath.Ext(filePath),
		"basename":  strings.TrimSuffix(fileInfo.Name(), filepath.Ext(fileInfo.Name())),
	}

	return stats, nil
}

// CountCSVRows counts the total number of rows in a CSV file (including header)
func CountCSVRows(csvPath string) (int, error) {
	file, err := os.Open(csvPath)
	if err != nil {
		return 0, CroissantError{
			Message: "failed to open CSV file",
			Value:   err,
		}
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.LazyQuotes = true
	reader.FieldsPerRecord = -1

	rowCount := 0
	for {
		_, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return 0, CroissantError{
				Message: fmt.Sprintf("failed to read CSV row %d", rowCount+1),
				Value:   err}
		}
		rowCount++
	}

	return rowCount, nil
}

// ValidateCSVStructure performs basic validation on CSV file structure
func ValidateCSVStructure(csvPath string) error {
	file, err := os.Open(csvPath)
	if err != nil {
		return CroissantError{
			Message: "failed to open CSV file",
			Value:   err,
		}
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.LazyQuotes = true
	reader.FieldsPerRecord = -1

	// Read first row (headers)
	headers, err := reader.Read()
	if err != nil {
		return CroissantError{
			Message: "failed to read CSV headers",
			Value:   err,
		}
	}

	if len(headers) == 0 {
		return CroissantError{
			Message: "CSV file has no columns",
			Value:   nil,
		}
	}

	// Check for duplicate headers
	headerMap := make(map[string]bool)
	for _, header := range headers {
		cleanHeader := strings.TrimSpace(header)
		if cleanHeader == "" {
			return CroissantError{
				Message: "CSV file has empty column header",
				Value:   nil,
			}
		}
		if headerMap[cleanHeader] {
			return CroissantError{
				Message: "CSV file has duplicate column header",
				Value:   cleanHeader,
			}
		}
		headerMap[cleanHeader] = true
	}

	return nil
}

// GetCSVColumnTypes analyzes a CSV file and returns inferred data types for each column
func GetCSVColumnTypes(csvPath string, sampleSize int) ([]string, []string, error) {
	headers, rows, err := GetCSVColumnsAndSampleRows(csvPath, sampleSize)
	if err != nil {
		return nil, nil, err
	}

	if len(rows) == 0 {
		// No data rows, default all to Text
		types := make([]string, len(headers))
		for i := range types {
			types[i] = "sc:Text"
		}
		return headers, types, nil
	}

	// Analyze each column
	columnTypes := make([]string, len(headers))
	for i := range headers {
		typeCounts := make(map[string]int)
		totalSamples := 0

		// Collect samples for this column
		for _, row := range rows {
			if i < len(row) && strings.TrimSpace(row[i]) != "" {
				dataType := InferDataType(row[i])
				typeCounts[dataType]++
				totalSamples++
			}
		}

		if totalSamples == 0 {
			columnTypes[i] = "sc:Text"
			continue
		}

		// Find the most common type
		maxCount := 0
		mostCommonType := "sc:Text"
		for dataType, count := range typeCounts {
			if count > maxCount {
				maxCount = count
				mostCommonType = dataType
			}
		}

		// If less than 70% of samples match the most common type, default to Text
		if float64(maxCount)/float64(totalSamples) < 0.7 {
			columnTypes[i] = "sc:Text"
		} else {
			columnTypes[i] = mostCommonType
		}
	}

	return headers, columnTypes, nil
}

// IsCSVFile checks if a file appears to be a CSV file based on extension
func IsCSVFile(filePath string) bool {
	ext := strings.ToLower(filepath.Ext(filePath))
	return ext == ".csv" || ext == ".tsv" || ext == ".txt"
}

// SanitizeFileName removes or replaces invalid characters in filenames
func SanitizeFileName(fileName string) string {
	// Replace invalid characters with underscores
	invalid := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|"}
	result := fileName
	for _, char := range invalid {
		result = strings.ReplaceAll(result, char, "_")
	}
	return result
}
