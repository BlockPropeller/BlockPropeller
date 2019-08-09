package resource

import (
	"bytes"
	"fmt"
	"strings"
)

// Provider represents a destination where Terraform
// resources should be deployed to.
//
// Each family of resources that is about to be deployed
// should have a provider defined for them.
type Provider interface {
	// Render a provider into Terraform syntax.
	Render() string
}

// Resource is an abstraction over different kinds of resources
// Terraform is able to provision.
//
// Resources are what is passed to Terraform in order to provision
// infrastructure.
//
// The actual implementations of the `Resource` interface are located
// in child packages specific to the Terraform Provider that supports them.
// It is up to these packages to provide syntactically correct definitions
// in order to be executed by Terraform.
type Resource interface {
	// Type of the resource.
	Type() string
	// Name of the resource.
	Name() string
	// Properties associated with the resource.
	Properties() *Properties
}

// ToID returns a `Property` containing the pointer to the provided resource.
func ToID(res Resource) Property {
	return NewRawProperty(fmt.Sprintf("%s.%s.id", res.Type(), FormatName(res.Name())))
}

// ToPropSelector returns a `Property` containing the pointer for a specific property of a resource.
func ToPropSelector(res Resource, name string) Property {
	return NewRawProperty(fmt.Sprintf("%s.%s.%s", res.Type(), FormatName(res.Name()), name))
}

// FormatName converts the resource name into a format suitable for use in Terraform resource names.
func FormatName(name string) string {
	return strings.ReplaceAll(name, " ", "")
}

// Render transforms a provided Resource into a textual Terraform resource definition.
//
// That result can then be passed on the the Terraform executable in order to be applied.
func Render(res Resource) string {
	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf("resource \"%s\" \"%s\" {\n",
		res.Type(), FormatName(res.Name())))
	buf.WriteString(res.Properties().Indent(2).Render())
	buf.WriteString("}\n")

	return buf.String()
}

// namedProperty binds a name to a particular Resource.
//
// An array of namedProperties is later used as the body of a `Resource`.
type namedProperty struct {
	name string
	prop Property
}

// Properties represents the body of a `Resource`.
//
// The fluent interface of the Properties structure allows for quick
// definition of resource body to Terraform Providers.
type Properties struct {
	indent int
	props  []namedProperty
}

// NewProperties creates a new instance of Properties,
// ready to be chained into a complete Properties definition.
func NewProperties() *Properties {
	return &Properties{}
}

// Prop appends a new named `Property` to an existing set of `Properties`.
//
// Prop returns a fluent interface to provide better developer experience
// while instantiating `Properties` structures.
func (p *Properties) Prop(name string, prop Property) *Properties {
	p.props = append(p.props, namedProperty{name, prop})

	return p
}

// Indent allow the called to define the number of spaces that the list
// of properties will be indented with. This is mainly an effort to improve
// the readability of generated Terraform files.
func (p *Properties) Indent(indent int) *Properties {
	p.indent = indent
	return p
}

// Render assembles a list of named properties into a
// string that can be embedded as the body of a `Resource`.
//
// Unlike traditional Golang maps, which render key value pairs
// randomly, Properties's Render method preserves the order of
// added properties in order to make testing of Properties easier.
func (p Properties) Render() string {
	var buf bytes.Buffer

	for _, nprop := range p.props {
		buf.WriteString(fmt.Sprintf(
			"%s%s=%s\n",
			strings.Repeat(" ", p.indent),
			nprop.name,
			nprop.prop.Render(),
		))
	}

	return buf.String()
}
