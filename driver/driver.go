package driver

import (
	_ "embed"
	"errors"
	"fmt"
	"html/template"
	"io"
	"iter"
	"path/filepath"
	"slices"
	"strings"
	"unicode"

	"github.com/pb33f/libopenapi"
	"github.com/pb33f/libopenapi/datamodel/high/base"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/pb33f/libopenapi/orderedmap"
)

//go:embed driver.tmpl
var markDownTemplate string

type Model struct {
	Version        string
	Infos          *base.Info
	Tags           map[string]*base.Tag
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

func pathToCamelCase(path string) (string, error) {
	// Split by non-alphanumeric characters (slashes, hyphens, braces, etc.)
	words := strings.FieldsFunc(path, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	})
	if len(words) == 0 {
		return "", nil
	}
	var builder strings.Builder
	for i, word := range words {
		if i == 0 {
			// First word: force the first letter to lowercase
			if _, err := builder.WriteString(strings.ToLower(word[:1]) + word[1:]); err != nil {
				return "", err
			}
		} else {
			// Subsequent words: force the first letter to uppercase
			if _, err := builder.WriteString(strings.ToUpper(word[:1]) + word[1:]); err != nil {
				return "", err
			}
		}
	}
	return builder.String(), nil
}

func refName(ref string) string {
	return filepath.Base(ref)
}

func refAnchor(ref string) (string, error) {
	return pathToCamelCase(ref)
}

func refID(ref string) (string, error) {
	if ref[0] == '#' {
		ref = ref[1:]
	}
	return pathToCamelCase(ref)
}

func add(i, j int) int {
	return i + j
}

func dict(values ...any) (map[string]any, error) {
	if len(values)%2 != 0 {
		return nil, errors.New("invalid dict call")
	}
	dict := make(map[string]any, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return nil, errors.New("dict keys must be strings")
		}
		dict[key] = values[i+1]
	}
	return dict, nil
}

func until(end int) iter.Seq[int] {
	return func(yield func(int) bool) {
		for i := range end {
			if !yield(i) {
				return
			}
		}
	}
}

func null() any {
	return nil
}

const defaultTag = "default"

func extractModel(document *v3.Document, model *Model) error {
	// tags
	for _, tag := range document.Tags {
		model.Tags[tag.Name] = tag
	}
	// endpoints are grouped by tags, so we need to iterate over all paths
	// and operations to extract the tags and their corresponding endpoints.
	for path, pathItem := range mapIterator(document.Paths.PathItems) {
		for _, param := range pathItem.Parameters {
			extractSchema(param.Schema, model)
		}
		for method, operation := range mapIterator(pathItem.GetOperations()) {
			operationMethod := strings.ToUpper(method)
			operationID, err := pathToCamelCase(method + "-" + path)
			if err != nil {
				return err
			}
			// for each operation, we create an endpoint and add it to the corresponding tag in the model.
			endpoint := &Endpoint{
				Method:      operationMethod,
				Path:        path,
				Summary:     orString(operation.Summary, pathItem.Summary),
				Description: orString(operation.Description, pathItem.Description),
				OperationId: orString(operation.OperationId, operationID),
				Parameters:  append(pathItem.Parameters, operation.Parameters...),
				RequestBody: operation.RequestBody,
				Responses:   operation.Responses,
			}
			if len(operation.Tags) == 0 {
				model.TagToEndpoints[defaultTag] = append(model.TagToEndpoints[defaultTag], endpoint)
			} else {
				for _, tag := range operation.Tags {
					model.TagToEndpoints[tag] = append(model.TagToEndpoints[tag], endpoint)
				}
			}
			// for each operation parameter, requestbody , response, we need to extract the schema and add it to the model.
			for _, param := range operation.Parameters {
				extractSchema(param.Schema, model)
			}
			// request body
			if operation.RequestBody != nil {
				for _, content := range operation.RequestBody.Content.FromOldest() {
					extractSchema(content.Schema, model)
				}
			}
			// responses
			if operation.Responses != nil {
				for _, response := range operation.Responses.Codes.FromOldest() {
					for _, content := range response.Content.FromOldest() {
						extractSchema(content.Schema, model)
					}
					// for _, header := range response.Headers.FromOldest() {
					// 	extractSchema(header.Schema, model)
					// }
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
		model.RefToSchema[schemaProxy.GetReference()] = schema
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

func GenerateMarkdown(input io.Reader, output io.Writer) (err error) {
	bytes, err := io.ReadAll(input)
	if err != nil {
		return err
	}
	// create a new document from specification bytes
	document, err := libopenapi.NewDocument(bytes)
	// if anything went wrong, an error is thrown
	if err != nil {
		return fmt.Errorf("cannot create new document: %w", err)
	}

	// because we know this is a v3 spec, we can build a ready to go model from it.
	v3Model, err := document.BuildV3Model()
	if err != nil {
		return fmt.Errorf("cannot build v3 model: %w", err)
	}

	mdModel := &Model{
		Version:        v3Model.Model.Version,
		Infos:          v3Model.Model.Info,
		Tags:           map[string]*base.Tag{},
		TagToEndpoints: map[string][]*Endpoint{},
		RefToSchema:    map[string]*base.Schema{},
	}

	if err := extractModel(&v3Model.Model, mdModel); err != nil {
		return fmt.Errorf("cannot extract model: %w", err)
	}

	funcMap := template.FuncMap{
		"dict":      dict,
		"add":       add,
		"refName":   refName,
		"refAnchor": refAnchor,
		"refID":     refID,
		"contains":  slices.Contains[[]string, string],
		"until":     until,
		"null":      null,
	}

	// we can now use the v3 model to generate markdown documentation.
	tmpl, err := template.New("markdown").Funcs(funcMap).Parse(markDownTemplate)
	if err != nil {
		return fmt.Errorf("cannot parse markdown template: %w", err)
	}

	// execute template with v3 model as data
	err = tmpl.ExecuteTemplate(output, "main", mdModel)
	if err != nil {
		return fmt.Errorf("cannot execute markdown template: %w", err)
	}

	return nil
}
