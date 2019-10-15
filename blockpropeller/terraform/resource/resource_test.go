package resource_test

import (
	"testing"

	"blockpropeller.dev/blockpropeller/terraform/resource"
	"blockpropeller.dev/lib/test"
)

type testResource struct {
}

func (testResource) Type() string {
	return "test_resource"
}

func (testResource) Name() string {
	return "test-0"
}

func (testResource) Properties() *resource.Properties {
	return resource.NewProperties().
		Prop("name", resource.NewStringProperty("test-0"))
}

func TestResourceToID(t *testing.T) {
	res := testResource{}

	got := resource.ToID(res).Render()
	want := "test_resource.test-0.id"

	test.AssertStringsEqual(t, "ToID()", got, want)
}

func TestResourceToIPAddress(t *testing.T) {
	res := testResource{}

	got := resource.ToPropSelector(res, "ipv4_address").Render()
	want := "test_resource.test-0.ipv4_address"

	test.AssertStringsEqual(t, "ToPropSelector()", got, want)
}

func TestResourceRendering(t *testing.T) {
	res := testResource{}

	got := resource.Render(res)
	want := `resource "test_resource" "test-0" {
  name="test-0"
}
`

	test.AssertStringsEqual(t, "Render()", got, want)
}

func TestPropertiesRendering(t *testing.T) {
	props := resource.NewProperties().
		Prop("foo", resource.NewStringProperty("foo")).
		Prop("bar", resource.NewStringProperty("baz"))

	got := props.Render()
	want := `foo="foo"
bar="baz"
`

	test.AssertStringsEqual(t, "props.Render()", got, want)
}

func TestPropertiesRenderingWithIndentation(t *testing.T) {
	props := resource.NewProperties().
		Prop("foo", resource.NewStringProperty("foo")).
		Prop("bar", resource.NewStringProperty("baz")).
		Indent(4)

	got := props.Render()
	want := `    foo="foo"
    bar="baz"
`

	test.AssertStringsEqual(t, "props.Render with indentation", got, want)
}
