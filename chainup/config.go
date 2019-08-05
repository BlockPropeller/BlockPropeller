package chainup

import (
	"chainup.dev/chainup/terraform"
	"github.com/pkg/errors"
)

// Config is the root config structure declaring all possible configuration parameters.
type Config struct {
	DigitalOcean *DigitalOceanConfig `yaml:"digital_ocean"`

	Terraform *terraform.Config `yaml:"terraform"`
}

// Validate satisfies the config.Config interface.
func (cfg *Config) Validate() error {
	return nil
}

// DigitalOceanConfig specifies all configuration parameters for DigitalOcean provider.
type DigitalOceanConfig struct {
	AccessToken string `yaml:"access_token"`
}

// Validate satisfies the config.Config interface.
func (cfg *DigitalOceanConfig) Validate() error {
	if cfg.AccessToken == "" {
		return errors.New("missing DigitalOcean access token")
	}

	return nil
}
