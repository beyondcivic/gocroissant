// main.go
package main

import (
	cmd "github.com/beyondcivic/gocroissant/cmd/gocroissant"
)

func main() {
	cmd.Init()
	cmd.Execute()
}
