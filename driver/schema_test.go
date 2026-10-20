package driver

import (
	"testing"

	"github.com/pb33f/libopenapi"
	"github.com/pb33f/libopenapi/datamodel/high/base"
	"github.com/pb33f/testify/require"
)

// loadSchemas parses the given component schemas (yaml, indented under
// components.schemas) and returns them by name.
func loadSchemas(t *testing.T, schemas string) map[string]*base.Schema {
	t.Helper()
	spec := "openapi: '3.2.0'\ninfo: {title: t, version: '1'}\npaths: {}\ncomponents:\n  schemas:\n" + schemas
	document, err := libopenapi.NewDocument([]byte(spec))
	require.NoError(t, err)
	model, err := document.BuildV3Model()
	require.NoError(t, err)
	result := map[string]*base.Schema{}
	for name, proxy := range model.Model.Components.Schemas.FromOldest() {
		result[name] = proxy.Schema()
		require.NotNil(t, result[name], name)
	}
	return result
}

func variantValues(variants []Variant) (values []string, refs []string) {
	for _, variant := range variants {
		values = append(values, variant.Value)
		refs = append(refs, variant.SchemaProxy.GetReference())
	}
	return values, refs
}

func propertyNames(properties []Property) (names []string, required []bool) {
	for _, property := range properties {
		names = append(names, property.Name)
		required = append(required, property.Required)
	}
	return names, required
}

func TestEffectiveType(t *testing.T) {
	schemas := loadSchemas(t, `
    Typed:
      type: string
    Props:
      properties:
        a: {type: string}
    Composed:
      allOf:
        - $ref: '#/components/schemas/Typed'
    List:
      items: {type: string}
    Empty:
      description: nothing
`)
	require.Equal(t, "string", effectiveType(schemas["Typed"]))
	require.Equal(t, "object", effectiveType(schemas["Props"]))
	require.Equal(t, "object", effectiveType(schemas["Composed"]))
	require.Equal(t, "array", effectiveType(schemas["List"]))
	require.Equal(t, "", effectiveType(schemas["Empty"]))
	require.Equal(t, "", effectiveType(nil))
}

func TestCompositionKind(t *testing.T) {
	schemas := loadSchemas(t, `
    Plain:
      type: object
    All:
      allOf: [{type: object}]
    One:
      oneOf: [{type: string}, {type: integer}]
    Any:
      anyOf: [{type: string}, {type: integer}]
    Both:
      allOf: [{type: object}]
      oneOf: [{type: string}]
`)
	require.Equal(t, "", compositionKind(schemas["Plain"]))
	require.Equal(t, kindAllOf, compositionKind(schemas["All"]))
	require.Equal(t, kindOneOf, compositionKind(schemas["One"]))
	require.Equal(t, kindAnyOf, compositionKind(schemas["Any"]))
	require.Equal(t, kindAllOf, compositionKind(schemas["Both"]))
	require.Equal(t, "", compositionKind(nil))
}

func TestVariants(t *testing.T) {
	schemas := loadSchemas(t, `
    A: {type: object}
    B: {type: object}
    C: {type: object}
    Plain:
      type: string
    NoDiscriminator:
      oneOf:
        - $ref: '#/components/schemas/A'
        - type: object
        - $ref: '#/components/schemas/B'
    Mapped:
      oneOf:
        - $ref: '#/components/schemas/A'
        - $ref: '#/components/schemas/B'
        - $ref: '#/components/schemas/C'
      discriminator:
        propertyName: kind
        mapping:
          "null": '#/components/schemas/C'
          a: '#/components/schemas/A'
          external: '#/components/schemas/Plain'
    Any:
      anyOf:
        - $ref: '#/components/schemas/A'
`)
	require.Nil(t, variants(nil))
	require.Empty(t, variants(schemas["Plain"]))

	values, refs := variantValues(variants(schemas["NoDiscriminator"]))
	require.Equal(t, []string{"A", "", "B"}, values)
	require.Equal(t, []string{"#/components/schemas/A", "", "#/components/schemas/B"}, refs)

	// mapping entries first in mapping order, then unmapped alternatives
	// with their implicit value, mapping to an unknown schema kept as a ref.
	values, refs = variantValues(variants(schemas["Mapped"]))
	require.Equal(t, []string{"null", "a", "external", "B"}, values)
	require.Equal(t, []string{
		"#/components/schemas/C",
		"#/components/schemas/A",
		"#/components/schemas/Plain",
		"#/components/schemas/B",
	}, refs)

	values, _ = variantValues(variants(schemas["Any"]))
	require.Equal(t, []string{"A"}, values)
}

func TestProperties(t *testing.T) {
	schemas := loadSchemas(t, `
    Plain:
      type: object
      properties:
        b: {type: string}
        a: {type: integer}
      required: [a]
    Base:
      allOf:
        - $ref: '#/components/schemas/Plain'
      properties:
        c: {type: boolean}
    Derived:
      allOf:
        - $ref: '#/components/schemas/Base'
        - type: object
          properties:
            a: {type: string, description: override}
            d: {type: string}
          required: [d]
      properties:
        b: {type: number}
      required: [b]
    Recursive:
      allOf:
        - $ref: '#/components/schemas/Recursive'
        - $ref: '#/components/schemas/Plain'
    Empty:
      type: object
`)
	require.Empty(t, properties(nil))
	require.Empty(t, properties(schemas["Empty"]))

	names, required := propertyNames(properties(schemas["Plain"]))
	require.Equal(t, []string{"b", "a"}, names)
	require.Equal(t, []bool{false, true}, required)

	// nested allOf through a reference, parts before own properties,
	// first position kept, last definition wins, required flags combined.
	result := properties(schemas["Derived"])
	names, required = propertyNames(result)
	require.Equal(t, []string{"b", "a", "c", "d"}, names)
	require.Equal(t, []bool{true, true, false, true}, required)
	require.Equal(t, "override", result[1].SchemaProxy.Schema().Description)
	require.Equal(t, []string{"number"}, result[0].SchemaProxy.Schema().Type)

	names, _ = propertyNames(properties(schemas["Recursive"]))
	require.Equal(t, []string{"b", "a"}, names)
}
