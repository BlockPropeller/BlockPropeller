package digitalocean

import (
	"blockpropeller.dev/blockpropeller/infrastructure"
	"blockpropeller.dev/blockpropeller/terraform"
	"blockpropeller.dev/blockpropeller/terraform/cloudprovider"
	"blockpropeller.dev/blockpropeller/terraform/resource"
	"blockpropeller.dev/blockpropeller/terraform/resource/digitalocean"
	"blockpropeller.dev/lib/log"
	"github.com/pkg/errors"
)

func init() {
	cloudprovider.RegisterProvider(infrastructure.ProviderDigitalOcean, &CloudProvider{})
}

var (
	image         = "ubuntu-18-04-x64"
	region        = "fra1"
	serverSizeMap = map[infrastructure.ServerSize]string{
		infrastructure.ServerSizeTest: "s-1vcpu-1gb",
		infrastructure.ServerSizeProd: "s-4vcpu-8gb",
	}
	volumeSizeMap = map[infrastructure.ServerSize]int{
		infrastructure.ServerSizeTest: 0,
		infrastructure.ServerSizeProd: 500,
	}
)

// CloudProvider is a terraform
type CloudProvider struct {
}

// Register satisfies the CloudProvider interface.
func (c *CloudProvider) Register(workspace *terraform.Workspace, settings *infrastructure.ProviderSettings) error {
	workspace.Add(digitalocean.NewProvider(settings.Credentials))

	return nil
}

// AddServer satisfies the CloudProvider interface.
func (c *CloudProvider) AddServer(workspace *terraform.Workspace, srv *infrastructure.Server) error {
	sshKey := srv.SSHKey

	doSSHKey := digitalocean.NewSSHKey(sshKey.Name, sshKey.EncodedPublicKey())
	log.Debug("using ssh key", log.Fields{
		"pub":  sshKey.EncodedPublicKey(),
		"priv": sshKey.EncodedPrivateKey(),
	})

	size, err := c.getDropletSize(srv.Size)
	if err != nil {
		return errors.Wrap(err, "get server size")
	}

	doDroplet := digitalocean.NewDroplet(srv.Name, image, region, size, []*digitalocean.SSHKey{doSSHKey})

	workspace.AddResource(doSSHKey, doDroplet)

	ipAddressOut := resource.NewOutput("ip-address", resource.ToPropSelector(doDroplet, "ipv4_address"))

	workspace.Add(ipAddressOut)

	// Add volume if necessary.
	volumeSize, err := c.getVolumeSize(srv.Size)
	if err != nil {
		return errors.Wrap(err, "get volume size")
	}

	if volumeSize > 0 {
		doVolume := digitalocean.NewVolume(srv.Name, region, volumeSize)
		doVolumeAttachment := digitalocean.NewVolumeAttachment(srv.Name, doDroplet, doVolume)

		workspace.AddResource(doVolume, doVolumeAttachment)
	}

	return nil
}

func (c *CloudProvider) getDropletSize(serverSize infrastructure.ServerSize) (string, error) {
	size, ok := serverSizeMap[serverSize]
	if !ok {
		return "", errors.Errorf("invalid server size: %s", serverSize)
	}

	return size, nil
}

func (c *CloudProvider) getVolumeSize(serverSize infrastructure.ServerSize) (int, error) {
	volumeSize, ok := volumeSizeMap[serverSize]
	if !ok {
		return 0, errors.Errorf("invalid server size: %s", serverSize)
	}

	return volumeSize, nil
}
