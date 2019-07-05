package resource_test

import (
	"testing"

	"chainup.dev/chainup/terraform/resource"
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

	if got != want {
		t.Errorf("resource.ToID: got %s, want %s", got, want)
		return
	}
}

func TestResourceRendering(t *testing.T) {
	res := testResource{}

	got := resource.Render(res)
	want := `resource "test_resource" "test-0" {
  "name" = "test-0"
}
`
	if got != want {
		t.Errorf("resource.Render: got '%s', want '%s'", got, want)
		return
	}
}

func TestPropertiesRendering(t *testing.T) {
	props := resource.NewProperties().
		Prop("foo", resource.NewStringProperty("foo")).
		Prop("bar", resource.NewStringProperty("baz"))

	got := props.Render()
	want := `"foo" = "foo"
"bar" = "baz"
`

	if got != want {
		t.Errorf("props.Render: got '%s', want '%s'", got, want)
		return
	}
}

func TestPropertiesRenderingWithIndentation(t *testing.T) {
	props := resource.NewProperties().
		Prop("foo", resource.NewStringProperty("foo")).
		Prop("bar", resource.NewStringProperty("baz")).
		Indent(4)

	got := props.Render()
	want := `    "foo" = "foo"
    "bar" = "baz"
`

	if got != want {
		t.Errorf("props.Render with indentation: got '%s', want '%s'", got, want)
		return
	}
}
