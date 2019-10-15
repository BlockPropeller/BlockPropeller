package digitalocean_test

import (
	"testing"

	"blockpropeller.dev/blockpropeller/terraform/resource"
	"blockpropeller.dev/blockpropeller/terraform/resource/digitalocean"
	"blockpropeller.dev/lib/test"
)

func TestDropletRendering(t *testing.T) {
	droplet := digitalocean.NewDroplet(
		"example-0",
		"ubuntu-18-04-x64",
		"fra1",
		"s-4vcpu-8gb",
		nil,
	)

	want := `resource "digitalocean_droplet" "example-0" {
  name="example-0"
  image="ubuntu-18-04-x64"
  region="fra1"
  size="s-4vcpu-8gb"
  ssh_keys=[]
}
`

	got := resource.Render(droplet)
	test.AssertStringsEqual(t, "Droplet.Render()", got, want)
}

func TestDropletWithSSHKey(t *testing.T) {
	sshKey := digitalocean.NewSSHKey("default", "ssh-rsa example@example.com")

	droplet := digitalocean.NewDroplet(
		"example-0",
		"ubuntu-18-04-x64",
		"fra1",
		"s-4vcpu-8gb",
		[]*digitalocean.SSHKey{sshKey},
	)

	want := `resource "digitalocean_droplet" "example-0" {
  name="example-0"
  image="ubuntu-18-04-x64"
  region="fra1"
  size="s-4vcpu-8gb"
  ssh_keys=[digitalocean_ssh_key.default.id]
}
`

	got := resource.Render(droplet)
	test.AssertStringsEqual(t, "Droplet.Render()", got, want)
}

func TestDropletWithMultipleSSHKeys(t *testing.T) {
	sshKeys := []*digitalocean.SSHKey{
		digitalocean.NewSSHKey("example", "ssh-rsa example@example.com"),
		digitalocean.NewSSHKey("foo", "ssh-rsa foo@bar.com"),
	}

	droplet := digitalocean.NewDroplet(
		"example-0",
		"ubuntu-18-04-x64",
		"fra1",
		"s-4vcpu-8gb",
		sshKeys,
	)

	want := `resource "digitalocean_droplet" "example-0" {
  name="example-0"
  image="ubuntu-18-04-x64"
  region="fra1"
  size="s-4vcpu-8gb"
  ssh_keys=[digitalocean_ssh_key.example.id, digitalocean_ssh_key.foo.id]
}
`

	got := resource.Render(droplet)
	test.AssertStringsEqual(t, "Droplet.Render()", got, want)
}
