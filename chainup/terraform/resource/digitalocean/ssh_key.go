package digitalocean

import (
	"chainup.dev/chainup/terraform/resource"
)

// SSHKey is a managed DigitalOcean SSH Key.
//
// An SSHKey resource can then be referenced from
// a Droplet in order to gain SSH access to the
// provisioned machine.
//
// A caveat to this resource is that an SSH Key
// we want to add must not already be present on
// DigitalOcean.
type SSHKey struct {
	name   string
	pubKey string
}

// NewSSHKey returns a new SSHKey instance.
func NewSSHKey(name string, pubKey string) *SSHKey {
	return &SSHKey{
		name:   name,
		pubKey: pubKey,
	}
}

// Type of the resource.
func (k *SSHKey) Type() string {
	return "digitalocean_ssh_key"
}

// Name of the resource.
func (k *SSHKey) Name() string {
	return k.name
}

// Properties associated with the resource.
func (k *SSHKey) Properties() *resource.Properties {
	return resource.NewProperties().
		Prop("name", resource.NewStringProperty(k.name)).
		Prop("public_key", resource.NewStringProperty(k.pubKey))
}
