package blockpropeller

import (
	"blockpropeller.dev/blockpropeller/account"
	"blockpropeller.dev/blockpropeller/ansible"
	"blockpropeller.dev/blockpropeller/provision"
	"blockpropeller.dev/blockpropeller/statemachine/middleware"
	"blockpropeller.dev/blockpropeller/terraform"
	"blockpropeller.dev/lib/config"
	"github.com/google/wire"
	"github.com/pkg/errors"
)

// AppSet keeps a set of all app level dependencies.
var AppSet = wire.NewSet(
	account.Set,
	terraform.Set,
	ansible.Set,
	provision.Set,
	middleware.Set,

	ProvideConfig,
	wire.FieldsOf(new(*Config),
		"Log", "Server", "WorkerPool", "Database", "JWT", "Terraform", "Ansible"),
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
		config.WithOverride("test"),
		config.WithName("config"),
		config.WithPath("config"),
		config.WithPath("."),
		config.SearchForPath(),
	})

	return config.NewFileProvider(opts...)
}
