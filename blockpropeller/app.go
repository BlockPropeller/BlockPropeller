package blockpropeller

import (
	"context"

	"blockpropeller.dev/blockpropeller/account"
	"blockpropeller.dev/blockpropeller/encryption"
	"blockpropeller.dev/blockpropeller/infrastructure"
	"blockpropeller.dev/blockpropeller/provision"
	"blockpropeller.dev/lib/log"
	"blockpropeller.dev/lib/server"
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

// InitGlobal configures the global dependencies of the App.
func (app *App) InitGlobal() {
	log.SetGlobal(app.Logger)
	encryption.Init(app.Config.Encryption.Secret)
}

// AppServer is a wrapper around an App that also serves traffic and processes provisioning jobs through a worker pool.
type AppServer struct {
	App        *App
	srv        *server.Server
	workerPool *provision.WorkerPool
}

// NewAppServer returns a new AppServer instance.
func NewAppServer(app *App, srv *server.Server, workerPool *provision.WorkerPool) *AppServer {
	return &AppServer{App: app, srv: srv, workerPool: workerPool}
}

// Start runs the worker pool in the background and a HTTP server in the foreground.
func (app *AppServer) Start(ctx context.Context) error {
	go app.workerPool.Start(ctx)

	return app.srv.Start()
}
