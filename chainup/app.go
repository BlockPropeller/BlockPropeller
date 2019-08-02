package chainup

import (
	"chainup.dev/chainup/infrastructure"
	"chainup.dev/chainup/provision"
)

// App is a container that holds all necessary dependencies
// to provide ChainUP functionality.
//
// App can be used to run ChainUP tasks through API, CLI, test cases or similar entry points.
type App struct {
	Config *Config

	Provisioner *provision.Provisioner

	ServerRepository infrastructure.ServerRepository
}

// NewApp returns a new App instance.
func NewApp(config *Config, provisioner *provision.Provisioner, serverRepository infrastructure.ServerRepository) *App {
	return &App{Config: config, Provisioner: provisioner, ServerRepository: serverRepository}
}
