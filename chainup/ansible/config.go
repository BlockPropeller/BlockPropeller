package ansible

// Config object for working with Ansible.
type Config struct {
	Path string `yaml:"path"`
}

// Validate conforms to the config.Config interface.
func (cfg *Config) Validate() error {
	if cfg.Path == "" {
		cfg.Path = "/usr/local/bin/ansible-playbook"
	}

	return nil
}
