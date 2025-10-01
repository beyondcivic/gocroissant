// File: pkg/croissant/validation_test.go
package croissant

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidateGoodCroissantFilesInDir(t *testing.T) {
	testDir := "testdata/1.0/good"
	files, err := os.ReadDir(testDir)
	if err != nil {
		t.Skip("Skipping directory test; " + testDir + " does not exist")

		return
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if filepath.Ext(file.Name()) == ".jsonld" {
			t.Run(file.Name(), func(t *testing.T) {
				issues, err := ValidateFile(filepath.Join(testDir, file.Name()))
				if err != nil {
					t.Fatalf("Failed to load croissant file %s: %v", file.Name(), err)
				}
				if issues.HasErrors() {
					t.Fatalf("Croissant file has errors: %s", issues.Report())
				}
			})
		}
	}
}

func TestValidateBadCroissantFilesInDir(t *testing.T) {
	testDir := "testdata/1.0/bad"
	files, err := os.ReadDir(testDir)
	if err != nil {
		t.Skip("Skipping directory test; " + testDir + " does not exist")

		return
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if filepath.Ext(file.Name()) == ".jsonld" {
			t.Run(file.Name(), func(t *testing.T) {
				issues, err := ValidateFile(filepath.Join(testDir, file.Name()))
				if err != nil {
					t.Fatalf("Failed to load croissant file %s: %v", file.Name(), err)
				}
				if !issues.HasErrors() {
					t.Fatalf("Croissant file has errors: %s", issues.Report())
				}
			})
		}
	}
}
