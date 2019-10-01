package account

import "github.com/pkg/errors"

// JWTConfig defines the configuration parameters for JWT services.
type JWTConfig struct {
	Secret string `yaml:"secret"`
}

// Validate satisfies the config.Config interface.
func (cfg *JWTConfig) Validate() error {
	if cfg.Secret == "" {
		return errors.New("missing JWT secret")
	}

	return nil
}
