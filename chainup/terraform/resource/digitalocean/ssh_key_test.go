package digitalocean_test

import (
	"testing"

	"chainup.dev/chainup/terraform/resource"
	"chainup.dev/chainup/terraform/resource/digitalocean"
)

func TestSSHKeyRendering(t *testing.T) {
	sshKey := digitalocean.NewSSHKey("example", "ssh-rsa example@example.com")

	want := `resource "digitalocean_ssh_key" "example" {
  "name" = "example"
  "public_key" = "ssh-rsa example@example.com"
}
`

	got := resource.Render(sshKey)
	if got != want {
		t.Errorf("SSH key rendering missmatch:\ngot:\n'%s'\nwant:\n'%s'", got, want)
	}
}
