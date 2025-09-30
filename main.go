// Gocroissant is a command-line tool and Go library for working with the ML Commons
// Croissant metadata format.
//
// Croissant is a standardized way to describe machine learning datasets using JSON-LD.
// This tool simplifies the creation of Croissant-compatible metadata from CSV data sources.
//
// # Installation
//
// Install the latest version:
//
//	go install github.com/beyondcivic/gocroissant@latest
//
// # Usage
//
// Generate metadata from a CSV file:
//
//	gocroissant generate data.csv -o metadata.jsonld
//
// Validate existing metadata:
//
//	gocroissant validate metadata.jsonld
//
// For detailed usage information, run:
//
//	gocroissant --help
package main

import (
	cmd "github.com/beyondcivic/gocroissant/cmd/gocroissant"
)

func main() {
	cmd.Init()
	cmd.Execute()
}
