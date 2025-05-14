package croissant

import (
	"crypto/sha256"
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

// CalculateSHA256 calculates the SHA-256 hash of a file
func CalculateSHA256(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

// GetCSVColumns reads the column names and first row from a CSV file
func GetCSVColumns(csvPath string) ([]string, []string, error) {
	file, err := os.Open(csvPath)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Read headers
	headers, err := reader.Read()
	if err != nil {
		return nil, nil, err
	}

	// Read first row for data type inference
	firstRow, err := reader.Read()
	if err != nil && err != io.EOF {
		return nil, nil, err
	}

	// If we hit EOF, there's no data row
	if err == io.EOF {
		return headers, nil, nil
	}

	return headers, firstRow, nil
}
