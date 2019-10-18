package digitalocean_test

import (
	"testing"

	"blockpropeller.dev/blockpropeller/terraform/resource"
	"blockpropeller.dev/blockpropeller/terraform/resource/digitalocean"
	"blockpropeller.dev/lib/test"
)

func TestVolumeRendering(t *testing.T) {
	volume := digitalocean.NewVolume("example", "fra1", 500)

	want := `resource "digitalocean_volume" "example" {
  name="volume"
  region="fra1"
  size=500
  initial_filesystem_type="ext4"
}
`

	got := resource.Render(volume)

	test.AssertStringsEqual(t, "Volume.Render()", got, want)
}

func TestVolumeAttachmentRendering(t *testing.T) {
	droplet := digitalocean.NewDroplet("example-0", "", "", "", nil)
	volume := digitalocean.NewVolume("volume-500", "fra1", 500)
	attachment := digitalocean.NewVolumeAttachment("example-volume-att", droplet, volume)

	want := `resource "digitalocean_volume_attachment" "example-volume-att" {
  droplet_id=digitalocean_droplet.example-0.id
  volume_id=digitalocean_volume.volume-500.id
}
`

	got := resource.Render(attachment)

	test.AssertStringsEqual(t, "VolumeAttachment.Render()", got, want)
}
