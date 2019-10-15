package terraform

// Config object for working with Terraform.
type Config struct {
	Path string `yaml:"path"`
}

// Validate satisfies the Config interface.
func (cfg *Config) Validate() error {
	if cfg.Path == "" {
		cfg.Path = "/usr/local/bin/terraform"
	}

	return nil
}
