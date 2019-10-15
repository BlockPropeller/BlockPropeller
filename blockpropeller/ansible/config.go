package ansible

import "github.com/pkg/errors"

// Config object for working with Ansible.
type Config struct {
	Path string `yaml:"path"`

	PlaybooksDir string `yaml:"playbooks_dir"`
	KeysDir      string `yaml:"keys_dir"`
}

// Validate conforms to the config.Config interface.
func (cfg *Config) Validate() error {
	if cfg.Path == "" {
		cfg.Path = "/usr/local/bin/ansible-playbook"
	}

	if cfg.PlaybooksDir == "" {
		return errors.New("missing ansible playbooks dir")
	}

	if cfg.KeysDir == "" {
		cfg.KeysDir = "/blockpropeller/ansible/keys"
	}

	return nil
}
