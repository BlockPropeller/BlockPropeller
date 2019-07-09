package digitalocean

import (
	"bytes"

	"chainup.dev/chainup/terraform/resource"
	"github.com/pkg/errors"
)

// Config specifies all the configuration options for
// setting up new DigitalOcean Providers.
type Config struct {
	Key string `yaml:"key" json:"key"`
}

// Validate satisfies the log.Config interface.
//
// Validation requires that the caller provides a value for the
// DigitalOcean API key.
func (cfg *Config) Validate() error {
	if cfg.Key == "" {
		return errors.New("missing DigitalOcean access key")
	}

	return nil
}

// Provider configures Terraform to know how to authenticate
// DigitalOcean requests for provisioning resources.
type Provider struct {
	props *resource.Properties
}

// NewProvider returns a new Provider instance.
func NewProvider(key string) *Provider {
	return &Provider{
		props: resource.NewProperties().
			Prop("key", resource.NewStringProperty(key)),
	}
}

// Render satisfies the resource.Provider interface.
func (p *Provider) Render() string {
	var buf bytes.Buffer

	buf.WriteString("provider \"digitalocean\" {\n")
	buf.WriteString(p.props.Indent(2).Render())
	buf.WriteString("}\n")

	return buf.String()
}
