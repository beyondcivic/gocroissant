// Generates docs for golang code and commands.
//
//go:generate go run ./docs-gen.go
package main

import (
	"errors"
	"go/build"
	"log"
	"os"

	cmd "github.com/beyondcivic/gocroissant/cmd/gocroissant"
	"github.com/beyondcivic/gocroissant/pkg/croissant"
	"github.com/invopop/jsonschema"
	"github.com/princjef/gomarkdoc"
	"github.com/princjef/gomarkdoc/lang"
	"github.com/princjef/gomarkdoc/logger"
	"github.com/spf13/cobra/doc"
)

// Generate tool docs from cobra commands
func GenerateCobraDocs() error {
	err := os.Mkdir("cmd", 0750)
	if err != nil && !errors.Is(err, os.ErrExist) {
		return err
	}
	// Setup cobra config
	cmd.Init()
	err = doc.GenMarkdownTree(cmd.RootCmd, "./cmd")
	if err != nil {
		return err
	}

	return nil
}

func GenerateReflector(pkg string, path string) (*jsonschema.Reflector, error) {
	reflector := new(jsonschema.Reflector)
	err := reflector.AddGoComments(pkg, path)

	return reflector, err
}

func WriteSchema(filename string, schema *jsonschema.Schema) error {
	outFile, err := schema.MarshalJSON()
	if err != nil {
		return err
	}

	return os.WriteFile(filename, outFile, 0600)
}

// Generates jsonschema files for reflected resources.
// In this case, reflects croissant.Metadata{}.
//
//nolint:exhaustruct
func GenerateTypeSchemas() error {
	err := os.Mkdir("schemas", 0750)
	if err != nil && !errors.Is(err, os.ErrExist) {
		return err
	}

	reflector, err := GenerateReflector(
		"github.com/beyondcivic/gocroissant/pkg/croissant", "../pkg/croissant",
	)
	if err != nil {
		return err
	}
	reflector.AllowAdditionalProperties = true
	// Specify struct to 'schema-fy'
	schema := reflector.Reflect(&croissant.Metadata{})

	if err := WriteSchema("./schemas/generated-schema.json", schema); err != nil {
		return err
	}

	return nil
}

// Generate markdown docs for packages
func GoMarkDoc() error {
	docRenderer, err := gomarkdoc.NewRenderer()
	if err != nil {
		return err
	}

	repo := lang.Repo{
		Remote:        "https://github.com:beyondcivic/gocroissant",
		DefaultBranch: "main",
		PathFromRoot:  "",
	}

	err = os.Mkdir("godoc", 0750)
	if err != nil && !errors.Is(err, os.ErrExist) {
		return err
	}

	docFiles := map[string]string{
		"../pkg/croissant": "gocroissant.md",
	}

	for pkgPath, pkgDoc := range docFiles {
		buildPkg, err := build.ImportDir(pkgPath, build.ImportComment)
		if err != nil {
			return err
		}

		logger := logger.New(logger.DebugLevel)
		pkg, err := lang.NewPackageFromBuild(logger, buildPkg, lang.PackageWithRepositoryOverrides(&repo))
		if err != nil {
			return err
		}
		output, err := docRenderer.Package(pkg)
		if err != nil {
			return err
		}
		err = os.WriteFile("godoc/"+pkgDoc, []byte(output), 0600)
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	if err := GenerateCobraDocs(); err != nil {
		log.Fatal(err)
	}
	if err := GoMarkDoc(); err != nil {
		log.Fatal(err)
	}
	if err := GenerateTypeSchemas(); err != nil {
		log.Fatal(err)
	}
}
