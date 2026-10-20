package driver

import (
	"slices"

	"github.com/pb33f/libopenapi/datamodel/high/base"
)

// Composition keywords, as returned by compositionKind.
const (
	kindAllOf = "allOf"
	kindOneOf = "oneOf"
	kindAnyOf = "anyOf"
)

// Variant is one alternative of a oneOf/anyOf schema. Value is the
// discriminator value selecting it: the explicit mapping key when there is
// one, otherwise the referenced schema name (OpenAPI implicit mapping).
// Value is empty for an inline alternative without mapping entry.
type Variant struct {
	Value       string
	SchemaProxy *base.SchemaProxy
}

// Property is a property of a schema, with the allOf composition flattened.
type Property struct {
	Name        string
	SchemaProxy *base.SchemaProxy
	Required    bool
}

// effectiveType returns the declared type of the schema, or infers "object"
// (properties or allOf present) and "array" (items present) when the type
// keyword is missing. It returns "" when nothing can be inferred.
func effectiveType(schema *base.Schema) string {
	if schema == nil {
		return ""
	}
	if len(schema.Type) > 0 {
		return schema.Type[0]
	}
	if (schema.Properties != nil && schema.Properties.Len() > 0) || len(schema.AllOf) > 0 {
		return "object"
	}
	if schema.Items != nil && schema.Items.IsA() {
		return "array"
	}
	return ""
}

// compositionKind returns which composition keyword the schema uses: "allOf",
// "oneOf", "anyOf" (in that precedence order) or "" for a plain schema.
func compositionKind(schema *base.Schema) string {
	switch {
	case schema == nil:
		return ""
	case len(schema.AllOf) > 0:
		return kindAllOf
	case len(schema.OneOf) > 0:
		return kindOneOf
	case len(schema.AnyOf) > 0:
		return kindAnyOf
	default:
		return ""
	}
}

// variants lists the alternatives of a oneOf (or anyOf) schema. Alternatives
// named in the discriminator mapping come first in mapping order; the others
// follow in declaration order with their implicit discriminator value. A
// mapping entry pointing outside the alternatives list is kept as a bare
// reference so that it is still documented.
func variants(schema *base.Schema) []Variant {
	if schema == nil {
		return nil
	}
	alternatives := schema.OneOf
	if len(alternatives) == 0 {
		alternatives = schema.AnyOf
	}
	remaining := make([]*base.SchemaProxy, 0, len(alternatives))
	byRef := make(map[string]*base.SchemaProxy, len(alternatives))
	for _, alternative := range alternatives {
		if alternative == nil {
			continue
		}
		if alternative.IsReference() {
			byRef[alternative.GetReference()] = alternative
		}
		remaining = append(remaining, alternative)
	}
	result := make([]Variant, 0, len(alternatives))
	if schema.Discriminator != nil {
		for value, ref := range schema.Discriminator.Mapping.FromOldest() {
			proxy, found := byRef[ref]
			if found {
				remaining = removeProxy(remaining, proxy)
			} else {
				proxy = base.CreateSchemaProxyRef(ref)
			}
			result = append(result, Variant{Value: value, SchemaProxy: proxy})
		}
	}
	for _, proxy := range remaining {
		value := ""
		if proxy.IsReference() {
			value = refName(proxy.GetReference())
		}
		result = append(result, Variant{Value: value, SchemaProxy: proxy})
	}
	return result
}

func removeProxy(proxies []*base.SchemaProxy, proxy *base.SchemaProxy) []*base.SchemaProxy {
	if i := slices.Index(proxies, proxy); i >= 0 {
		return slices.Delete(proxies, i, i+1)
	}
	return proxies
}

// properties returns the properties of the schema in declaration order. For an
// allOf composition, the properties of every part are merged (parts first, in
// order, then the schema's own properties): a property keeps the position of
// its first occurrence, the last definition wins and required flags are
// combined. Parts referencing an already visited schema are skipped, which
// keeps recursive compositions finite.
func properties(schema *base.Schema) []Property {
	collector := &propertyCollector{index: map[string]int{}, visited: map[string]bool{}}
	collector.collect(schema)
	return collector.result
}

type propertyCollector struct {
	result  []Property
	index   map[string]int
	visited map[string]bool
}

func (c *propertyCollector) collect(schema *base.Schema) {
	if schema == nil {
		return
	}
	for _, part := range schema.AllOf {
		if part == nil {
			continue
		}
		if part.IsReference() {
			ref := part.GetReference()
			if c.visited[ref] {
				continue
			}
			c.visited[ref] = true
		}
		c.collect(part.Schema())
	}
	for name, proxy := range schema.Properties.FromOldest() {
		c.add(Property{Name: name, SchemaProxy: proxy, Required: slices.Contains(schema.Required, name)})
	}
}

func (c *propertyCollector) add(property Property) {
	if i, found := c.index[property.Name]; found {
		property.Required = property.Required || c.result[i].Required
		c.result[i] = property
		return
	}
	c.index[property.Name] = len(c.result)
	c.result = append(c.result, property)
}
