package resource

import (
	"fmt"
	"strings"
)

// Property represents a generic value of a Terraform resource definition.
//
// Configuration is done by applying one or more named properties to a specific resource.
type Property interface {
	Render() string
}

// RawProperty is used when we want to generate terraform functions
// that usually reference some other part of the terraform plan definition.
//
// For example: `digitalocean_ssh_key.default.id`
//     becomes: `digitalocean_ssh_key.default.id`
//   which generates a pointer to the ID of a terraform managed SSH key.
type RawProperty struct {
	value string
}

// NewRawProperty returns a new instance of a `RawProperty`.
func NewRawProperty(value string) *RawProperty {
	return &RawProperty{
		value: value,
	}
}

// Render prepares the property into a format appropriate for resource generation.
//
// Each property is responsible for rendering itself in a valid Terraform syntax.
func (prop RawProperty) Render() string {
	return prop.value
}

// StringProperty is is a wrapper for a string constant.
//
// The difference between the `StringProperty` and `RawProperty` is that
// the `StringProperty` is escaped with double quotation marks, while
// the `RawProperty` is left as is.
//
// For example: `demo.example.com`
//     becomes: `"demo.example.com"`
//   which serves as a string constant.
type StringProperty struct {
	value string
}

// NewStringProperty returns a new instance of a StringProperty.
func NewStringProperty(value string) *StringProperty {
	return &StringProperty{value: value}
}

// Render prepares the property into a format appropriate for resource generation.
//
// Each property is responsible for rendering itself in a valid Terraform syntax.
func (prop StringProperty) Render() string {
	return fmt.Sprintf("\"%s\"", prop.value)
}

// IntegerProperty is is a wrapper for an integer constant.
//
// In Terraform, integers are defined as is, without needing
// any additional syntax.
//
// For example: `1000`
// 	   becomes: `1000`
//   which serves as an integer constant.
type IntegerProperty struct {
	value int
}

// NewIntegerProperty returns a new instance of an IntegerProperty.
func NewIntegerProperty(value int) *IntegerProperty {
	return &IntegerProperty{
		value: value,
	}
}

// Render prepares the property into a format appropriate for resource generation.
//
// Each property is responsible for rendering itself in a valid Terraform syntax.
func (prop IntegerProperty) Render() string {
	return fmt.Sprintf("%d", prop.value)
}

// ArrayProperty is an aggregating type which is able to hold
// an arbitrary number of `Property` types.
//
// In Terraform, arrays are defined as a sequence of properties,
// surrounded by square brackets, and delimited with a comma.
//
// For example: `[]int{1, 2, 3}`
// 	   becomes: `[ 1, 2, 3 ]`
//
// Additionally, each individual element could be of a different type,
// as long as it implements the `Property` interface.
//
// For example: `[]interface{1, "2", res.name.id, []int{0, 1}}`
// 	   becomes: `[ 1, "2", res.name.id, [ 0, 1 ] ]`
type ArrayProperty struct {
	values []Property
}

// NewArrayProperty returns a new instance of an ArrayProperty.
func NewArrayProperty(values ...Property) *ArrayProperty {
	return &ArrayProperty{
		values: values,
	}
}

// Render prepares the property into a format appropriate for resource generation.
//
// Each property is responsible for rendering itself in a valid Terraform syntax.
func (prop ArrayProperty) Render() string {
	var values []string
	for _, value := range prop.values {
		values = append(values, value.Render())
	}

	return fmt.Sprintf("[ %s ]", strings.Join(values, ", "))
}
