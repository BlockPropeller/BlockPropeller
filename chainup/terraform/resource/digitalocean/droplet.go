package digitalocean

import (
	"chainup.dev/chainup/terraform/resource"
)

// Droplet is a DigitalOcean definition of a server.
//
// In order to configure a droplet, properties such as `image`, `region` and `size` must be set.
type Droplet struct {
	name    string
	image   string
	region  string
	size    string
	sshKeys []*SSHKey
}

// NewDroplet returns a new instance of a Droplet.
func NewDroplet(name string, image string, region string, size string, sshKeys []*SSHKey) *Droplet {
	return &Droplet{name: name, image: image, region: region, size: size, sshKeys: sshKeys}
}

// Type of the resource.
func (d *Droplet) Type() string {
	return "digitalocean_droplet"
}

// Name of the resource.
func (d *Droplet) Name() string {
	return d.name
}

// Properties associated with the resource.
func (d *Droplet) Properties() *resource.Properties {
	var sshKeys []resource.Property
	for _, sshKey := range d.sshKeys {
		sshKeys = append(sshKeys, resource.ToID(sshKey))
	}

	return resource.NewProperties().
		Prop("name", resource.NewStringProperty(d.name)).
		Prop("image", resource.NewStringProperty(d.image)).
		Prop("region", resource.NewStringProperty(d.region)).
		Prop("size", resource.NewStringProperty(d.size)).
		Prop("ssh_keys", resource.NewArrayProperty(sshKeys...))
}
