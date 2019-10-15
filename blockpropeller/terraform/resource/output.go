package resource

import (
	"bytes"
	"fmt"
)

// Output defines a named output to be used for extracting information
// about the provisioned architecture.
type Output struct {
	Name  string
	Value Property
}

// NewOutput returns a new Output instance.
func NewOutput(name string, value Property) *Output {
	return &Output{
		Name:  name,
		Value: value,
	}
}

// Render the output into Terraform syntax.
func (out *Output) Render() string {
	var buf bytes.Buffer

	props := NewProperties().
		Prop("value", out.Value)

	buf.WriteString(fmt.Sprintf("output \"%s\" {\n", FormatName(out.Name)))
	buf.WriteString(props.Indent(2).Render())
	buf.WriteString("}\n")

	return buf.String()
}
