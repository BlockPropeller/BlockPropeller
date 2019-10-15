package digitalocean_test

import (
	"testing"

	"blockpropeller.dev/blockpropeller/terraform/resource/digitalocean"
	"blockpropeller.dev/lib/test"
)

func TestProviderRendering(t *testing.T) {
	provider := digitalocean.NewProvider("foobar")

	want := `provider "digitalocean" {
  token="foobar"
}
`

	got := provider.Render()

	test.AssertStringsEqual(t, "Provider.Render()", got, want)
}
