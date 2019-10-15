package blockpropeller

import (
	"blockpropeller.dev/blockpropeller/account"
	"blockpropeller.dev/blockpropeller/infrastructure"
	"blockpropeller.dev/blockpropeller/provision"
	"blockpropeller.dev/lib/log"
)

// App is a container that holds all necessary dependencies
// to provide BlockPropeller functionality.
//
// App can be used to run BlockPropeller tasks through API, CLI, test cases or similar entry points.
type App struct {
	Config *Config

	AccountRepository account.Repository
	AccountService    *account.Service

	ProviderSettingsRepository infrastructure.ProviderSettingsRepository
	ServerRepository           infrastructure.ServerRepository
	JobRepository              provision.JobRepository

	JobScheduler *provision.JobScheduler
	Provisioner  *provision.Provisioner

	Logger log.Logger
}

// NewApp returns a new App instance.
func NewApp(
	config *Config,
	accRepo account.Repository,
	accSvc *account.Service,
	providerSettingsRepo infrastructure.ProviderSettingsRepository,
	serverRepo infrastructure.ServerRepository,
	jobRepo provision.JobRepository,
	jobScheduler *provision.JobScheduler,
	provisioner *provision.Provisioner,
	logger log.Logger,
) *App {
	log.SetGlobal(logger)

	return &App{
		Config:                     config,
		AccountRepository:          accRepo,
		AccountService:             accSvc,
		ProviderSettingsRepository: providerSettingsRepo,
		ServerRepository:           serverRepo,
		JobRepository:              jobRepo,
		JobScheduler:               jobScheduler,
		Provisioner:                provisioner,
		Logger:                     logger,
	}
}
