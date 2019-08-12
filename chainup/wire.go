package chainup

import (
	"chainup.dev/chainup/ansible"
	"chainup.dev/chainup/infrastructure"
	"chainup.dev/chainup/provision"
	"chainup.dev/chainup/terraform"
	"chainup.dev/lib/config"
	"github.com/google/wire"
	"github.com/pkg/errors"
)

// AppSet keeps a set of all app level dependencies.
var AppSet = wire.NewSet(
	terraform.Set,
	ansible.Set,
	infrastructure.Set,
	provision.Set,

	ProvideConfig,
	wire.FieldsOf(new(*Config), "Log", "Terraform", "Ansible"),
	NewApp,
)

// ProvideConfig initializes and returns a new Config instance.
//
// ProvideConfig panics on failed configuration load.
func ProvideConfig(provider config.Provider) *Config {
	var cfg Config

	_, err := provider.Load(&cfg)
	if err != nil {
		panic(errors.Errorf("could not load config: %s", err))
	}

	return &cfg
}

// ProvideFileConfigProvider provides a configures file config provider to be used for
// configuration loading.
func ProvideFileConfigProvider() config.Provider {
	opts := append([]config.FileProviderOpt{
		config.WithName("config"),
		config.WithPath("config"),
		config.WithPath("."),
	})

	return config.NewFileProvider(opts...)
}

// ProvideTestConfigProvider provides a configures file config provider to be used for
// configuration loading.
//
// ProvideTestConfigProvider searches for configuration in any of the parent folders
// of the current working directory.
func ProvideTestConfigProvider() config.Provider {
	opts := append([]config.FileProviderOpt{
		config.WithName("config"),
		config.WithPath("config"),
		config.WithPath("."),
		config.SearchForPath(),
	})

	return config.NewFileProvider(opts...)
}
