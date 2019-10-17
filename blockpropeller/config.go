package blockpropeller

import (
	"blockpropeller.dev/blockpropeller/account"
	"blockpropeller.dev/blockpropeller/ansible"
	"blockpropeller.dev/blockpropeller/database"
	"blockpropeller.dev/blockpropeller/provision"
	"blockpropeller.dev/blockpropeller/terraform"
	"blockpropeller.dev/lib/log"
	"blockpropeller.dev/lib/server"
	"github.com/pkg/errors"
)

// Config is the root config structure declaring all possible configuration parameters.
type Config struct {
	Log        *log.Config                 `yaml:"log"`
	Server     *server.Config              `yaml:"server"`
	WorkerPool *provision.WorkerPoolConfig `yaml:"worker_pool"`

	Database *database.Config   `yaml:"database"`
	JWT      *account.JWTConfig `yaml:"jwt"`

	DigitalOcean *DigitalOceanConfig `yaml:"digital_ocean"`

	Terraform *terraform.Config `yaml:"terraform"`
	Ansible   *ansible.Config   `yaml:"ansible"`
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
