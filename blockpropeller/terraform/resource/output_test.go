package resource_test

import (
	"testing"

	"blockpropeller.dev/blockpropeller/terraform/resource"
	"blockpropeller.dev/lib/test"
)

func TestOutputRendering(t *testing.T) {
	output := resource.NewOutput("test-ip", resource.NewStringProperty("192.168.1.1"))

	got := output.Render()
	want := `output "test-ip" {
  value="192.168.1.1"
}
`

	test.AssertStringsEqual(t, "Output.Render()", got, want)
}
