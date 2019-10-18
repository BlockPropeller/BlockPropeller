package digitalocean

import (
	"blockpropeller.dev/blockpropeller/terraform/resource"
)

// Volume is a DigitalOcean volume that can be attached to a Droplet.
type Volume struct {
	name   string
	region string
	size   int
}

// NewVolume returns a new Volume instance.
func NewVolume(name string, region string, size int) *Volume {
	return &Volume{
		name:   name,
		region: region,
		size:   size,
	}
}

// Type of the resource.
func (v *Volume) Type() string {
	return "digitalocean_volume"
}

// Name of the resource.
func (v *Volume) Name() string {
	return v.name
}

// Properties associated with the resource.
func (v *Volume) Properties() *resource.Properties {
	return resource.NewProperties().
		Prop("name", resource.NewStringProperty("volume")).
		Prop("region", resource.NewStringProperty(v.region)).
		Prop("size", resource.NewIntegerProperty(v.size)).
		Prop("initial_filesystem_type", resource.NewStringProperty("ext4"))
}

// VolumeAttachment connects a Volume with a DigitalOcean Droplet.
type VolumeAttachment struct {
	name    string
	droplet *Droplet
	volume  *Volume
}

// NewVolumeAttachment returns a new VolumeAttachment instance.
func NewVolumeAttachment(name string, droplet *Droplet, volume *Volume) *VolumeAttachment {
	return &VolumeAttachment{
		name:    name,
		droplet: droplet,
		volume:  volume,
	}
}

// Type of the resource.
func (a *VolumeAttachment) Type() string {
	return "digitalocean_volume_attachment"
}

// Name of the resource.
func (a *VolumeAttachment) Name() string {
	return a.name
}

// Properties associated with the resource.
func (a *VolumeAttachment) Properties() *resource.Properties {
	return resource.NewProperties().
		Prop("droplet_id", resource.ToID(a.droplet)).
		Prop("volume_id", resource.ToID(a.volume))
}
