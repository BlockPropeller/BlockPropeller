package digitalocean_test

import (
	"testing"

	"chainup.dev/chainup/terraform/resource"
	"chainup.dev/chainup/terraform/resource/digitalocean"
	"chainup.dev/lib/test"
)

func TestSSHKeyRendering(t *testing.T) {
	sshKey := digitalocean.NewSSHKey("example", "ssh-rsa example@example.com")

	want := `resource "digitalocean_ssh_key" "example" {
  name="example"
  public_key="ssh-rsa example@example.com"
}
`

	got := resource.Render(sshKey)

	test.AssertStringsEqual(t, "SSHKey.Render()", got, want)
}
