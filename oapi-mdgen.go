package main

import (
	"context"
	_ "embed"
	"fmt"
	"iter"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"text/template"

	"github.com/pb33f/libopenapi"
	"github.com/pb33f/libopenapi/datamodel/high/base"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/pb33f/libopenapi/orderedmap"
	"github.com/urfave/cli/v3"
)

func main() {
	app := &cli.Command{
		Name:  "oapi-mdgen",
		Usage: "Convert OpenAPI specifications to Markdown documentation",
		Description: "The command lets you convert OpenAPI specifications to markdown documentation.\n" +
			"More information on command can be found here: https://github.com/masagroup/oapi-mdgen",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "input",
				Aliases:  []string{"i"},
				Usage:    "input OpenAPI specification file",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "output",
				Aliases:  []string{"o"},
				Usage:    "output directory for the generated markdown files",
				Required: true,
			}},
		UseShortOptionHandling: true,
		Action: func(ctx context.Context, c *cli.Command) error {
			input := c.String("input")
			output := c.String("output")
			if input == "" || output == "" {
				return cli.Exit("input and output flags are required", 1)
			}
			return generateMarkdown(input, output)
		},
	}
	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

//go:embed oapi-mdgen.tmpl
var markDownTemplate string

type Model struct {
	Version        string
	Infos          *base.Info
	TagToEndpoints map[string][]*Endpoint
	RefToSchema    map[string]*base.Schema
}

type Endpoint struct {
	Method      string
	Path        string
	Description string
	Summary     string
	OperationId string
	Parameters  []*v3.Parameter
	RequestBody *v3.RequestBody
	Responses   *v3.Responses
}

func mapIterator[K comparable, V any](m *orderedmap.Map[K, V]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for pair := m.First(); pair != nil; pair = pair.Next() {
			if !yield(pair.Key(), pair.Value()) {
				return
			}
		}
	}
}

func orString(a, b string) string {
	if a != "" {
		return a
	}
	return b
}

func refName(ref string) string {
	return filepath.Base(ref)
}

func refAnchor(ref string) string {
	return strings.ReplaceAll(ref, "/", "")
}

func refID(ref string) string {
	if ref[0] == '#' {
		ref = ref[1:]
	}
	return strings.ReplaceAll(ref, "/", "")
}

func extractModel(document *v3.Document, model *Model) error {
	// endpoints are grouped by tags, so we need to iterate over all paths
	// and operations to extract the tags and their corresponding endpoints.
	for path, pathItem := range mapIterator(document.Paths.PathItems) {
		for _, param := range pathItem.Parameters {
			extractSchema(param.Schema, model)
		}
		for method, operation := range mapIterator(pathItem.GetOperations()) {
			// for each operation, we create an endpoint and add it to the corresponding tag in the model.
			endpoint := &Endpoint{
				Method:      strings.ToUpper(method),
				Path:        path,
				Summary:     orString(operation.Summary, pathItem.Summary),
				Description: orString(operation.Description, pathItem.Description),
				OperationId: operation.OperationId,
				Parameters:  append(pathItem.Parameters, operation.Parameters...),
				RequestBody: operation.RequestBody,
				Responses:   operation.Responses,
			}
			for _, tag := range operation.Tags {
				model.TagToEndpoints[tag] = append(model.TagToEndpoints[tag], endpoint)
			}
			// for each operation parameter, requestbody , response, we need to extract the schema and add it to the model.
			for _, param := range operation.Parameters {
				extractSchema(param.Schema, model)
			}
			if operation.RequestBody != nil {
				for _, content := range operation.RequestBody.Content.FromOldest() {
					extractSchema(content.Schema, model)
				}
			}
			for _, response := range operation.Responses.Codes.FromOldest() {
				for _, content := range response.Content.FromOldest() {
					extractSchema(content.Schema, model)
				}
				for _, header := range response.Headers.FromOldest() {
					extractSchema(header.Schema, model)
				}
			}
		}
	}
	return nil
}

func extractSchema(schemaProxy *base.SchemaProxy, model *Model) {
	if schemaProxy == nil {
		return
	}
	schema := schemaProxy.Schema()
	if schema == nil {
		return
	}
	if schemaProxy.IsReference() {
		// reference
		model.RefToSchema[schemaProxy.GetReference()] = schemaProxy.Schema()
	} else if items := schema.Items; items != nil && items.IsA() {
		// inline
		extractSchema(items.A, model)
	}
	if properties := schema.Properties; properties != nil {
		for _, prop := range properties.FromOldest() {
			extractSchema(prop, model)
		}
	}
}

func generateMarkdown(input, output string) (err error) {
	// load an OpenAPI 3 specification from bytes
	bytes, err := os.ReadFile(input)
	if err != nil {
		return fmt.Errorf("cannot read input file: %w", err)
	}

	// create a new document from specification bytes
	document, err := libopenapi.NewDocument(bytes)
	// if anything went wrong, an error is thrown
	if err != nil {
		panic(fmt.Sprintf("cannot create new document: %e", err))
	}

	// because we know this is a v3 spec, we can build a ready to go model from it.
	v3Model, err := document.BuildV3Model()
	if err != nil {
		return fmt.Errorf("cannot build v3 model: %w", err)
	}

	mdModel := &Model{
		Version:        v3Model.Model.Version,
		Infos:          v3Model.Model.Info,
		TagToEndpoints: map[string][]*Endpoint{},
		RefToSchema:    map[string]*base.Schema{},
	}

	if err := extractModel(&v3Model.Model, mdModel); err != nil {
		return fmt.Errorf("cannot extract model: %w", err)
	}

	funcMap := template.FuncMap{
		"refName":   refName,
		"refAnchor": refAnchor,
		"refID":     refID,
		"contains":  slices.Contains[[]string, string],
	}

	// we can now use the v3 model to generate markdown documentation.
	tmpl, err := template.New("markdown").Funcs(funcMap).Parse(markDownTemplate)
	if err != nil {
		return fmt.Errorf("cannot parse markdown template: %w", err)
	}

	// create output file
	f, err := os.Create(output)
	if err != nil {
		return fmt.Errorf("cannot create output file: %w", err)
	}
	defer func() {
		if err2 := f.Close(); err2 != nil && err == nil {
			err = fmt.Errorf("cannot close output file: %w", err2)
		}
	}()

	// execute template with v3 model as data
	err = tmpl.ExecuteTemplate(f, "main", mdModel)
	if err != nil {
		return fmt.Errorf("cannot execute markdown template: %w", err)
	}

	return nil
}
