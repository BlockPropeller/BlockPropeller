package chainup

import (
	"chainup.dev/chainup/infrastructure"
	"chainup.dev/chainup/provision"
	"chainup.dev/lib/log"
)

// App is a container that holds all necessary dependencies
// to provide ChainUP functionality.
//
// App can be used to run ChainUP tasks through API, CLI, test cases or similar entry points.
type App struct {
	Config *Config

	ProviderSettingsRepository infrastructure.ProviderSettingsRepository

	JobScheduler *provision.JobScheduler
	Provisioner  *provision.Provisioner

	Logger log.Logger
}

// NewApp returns a new App instance.
func NewApp(
	config *Config,
	providerSettingsRepo infrastructure.ProviderSettingsRepository,
	jobScheduler *provision.JobScheduler,
	provisioner *provision.Provisioner,
	logger log.Logger,
) *App {
	log.SetGlobal(logger)

	return &App{
		Config:                     config,
		ProviderSettingsRepository: providerSettingsRepo,
		JobScheduler:               jobScheduler,
		Provisioner:                provisioner,
		Logger:                     logger,
	}
}
