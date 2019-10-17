// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package blockpropeller

import (
	"blockpropeller.dev/blockpropeller/account"
	"blockpropeller.dev/blockpropeller/ansible"
	"blockpropeller.dev/blockpropeller/database"
	"blockpropeller.dev/blockpropeller/database/transaction"
	"blockpropeller.dev/blockpropeller/httpserver"
	middleware2 "blockpropeller.dev/blockpropeller/httpserver/middleware"
	"blockpropeller.dev/blockpropeller/httpserver/routes"
	"blockpropeller.dev/blockpropeller/infrastructure"
	"blockpropeller.dev/blockpropeller/provision"
	"blockpropeller.dev/blockpropeller/statemachine/middleware"
	"blockpropeller.dev/blockpropeller/terraform"
	"blockpropeller.dev/lib/log"
	"blockpropeller.dev/lib/server"
	"github.com/google/wire"
	"testing"
)

// Injectors from inject_database.go:

func SetupDatabaseApp() (*App, func(), error) {
	provider := ProvideFileConfigProvider()
	config := ProvideConfig(provider)
	databaseConfig := config.Database
	logConfig := config.Log
	db, cleanup, err := database.ProvideDB(databaseConfig, logConfig)
	if err != nil {
		return nil, nil, err
	}
	accountRepository := database.NewAccountRepository(db)
	jwtConfig := config.JWT
	tokenService := account.ConfigureTokenService(jwtConfig)
	service := account.NewService(accountRepository, tokenService)
	providerSettingsRepository := database.NewProviderSettingsRepository(db)
	serverRepository := database.NewServerRepository(db)
	jobRepository := database.NewJobRepository(db)
	deploymentRepository := database.NewDeploymentRepository(db)
	jobScheduler := provision.NewJobScheduler(db, jobRepository, serverRepository, deploymentRepository)
	terraformConfig := config.Terraform
	terraformTerraform := terraform.ConfigureTerraform(terraformConfig)
	serverProvisioner := provision.NewServerProvisioner(terraformTerraform, serverRepository)
	stepProvisionServer := provision.NewStepProvisionServer(serverProvisioner, jobRepository)
	ansibleConfig := config.Ansible
	ansibleAnsible := ansible.ConfigureAnsible(ansibleConfig)
	deploymentProvisioner := provision.NewDeploymentProvisioner(ansibleAnsible, deploymentRepository)
	stepProvisionDeployment := provision.NewStepProvisionDeployment(deploymentProvisioner, jobRepository)
	failureMiddleware := provision.NewFailureMiddleware(jobRepository)
	transactional := middleware.NewTransactional(db)
	jobStateMachine := provision.ConfigureJobStateMachine(stepProvisionServer, stepProvisionDeployment, failureMiddleware, transactional)
	serverDestroyer := provision.NewServerDestroyer(terraformTerraform, db, serverRepository, deploymentRepository)
	provisioner := provision.NewProvisioner(jobStateMachine, jobScheduler, terraformTerraform, serverDestroyer)
	consoleLogger := log.NewConsoleLogger(logConfig)
	app := NewApp(config, accountRepository, service, providerSettingsRepository, serverRepository, jobRepository, jobScheduler, provisioner, consoleLogger)
	return app, func() {
		cleanup()
	}, nil
}

func SetupDatabaseServer() (*AppServer, func(), error) {
	provider := ProvideFileConfigProvider()
	config := ProvideConfig(provider)
	serverConfig := config.Server
	databaseConfig := config.Database
	logConfig := config.Log
	db, cleanup, err := database.ProvideDB(databaseConfig, logConfig)
	if err != nil {
		return nil, nil, err
	}
	accountRepository := database.NewAccountRepository(db)
	jwtConfig := config.JWT
	tokenService := account.ConfigureTokenService(jwtConfig)
	service := account.NewService(accountRepository, tokenService)
	authenticationMiddleware := middleware2.NewAuthenticationMiddleware(service)
	authentication := routes.NewAuthenticationRoutes(service)
	routesAccount := routes.NewAccountRoutes(accountRepository)
	providerSettingsRepository := database.NewProviderSettingsRepository(db)
	providerSettings := routes.NewProviderSettingsRoutes(providerSettingsRepository)
	jobRepository := database.NewJobRepository(db)
	serverRepository := database.NewServerRepository(db)
	deploymentRepository := database.NewDeploymentRepository(db)
	jobScheduler := provision.NewJobScheduler(db, jobRepository, serverRepository, deploymentRepository)
	routesProvision := routes.NewProvisionRoutes(jobScheduler, jobRepository, providerSettingsRepository)
	router := &httpserver.Router{
		AuthenticatedMiddleware: authenticationMiddleware,
		AuthRoutes:              authentication,
		AccountRoutes:           routesAccount,
		ProviderSettingsRoutes:  providerSettings,
		ProvisionRoutes:         routesProvision,
	}
	consoleLogger := log.NewConsoleLogger(logConfig)
	serverServer, err := server.ProvideServer(serverConfig, router, consoleLogger)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	workerPoolConfig := config.WorkerPool
	terraformConfig := config.Terraform
	terraformTerraform := terraform.ConfigureTerraform(terraformConfig)
	serverProvisioner := provision.NewServerProvisioner(terraformTerraform, serverRepository)
	stepProvisionServer := provision.NewStepProvisionServer(serverProvisioner, jobRepository)
	ansibleConfig := config.Ansible
	ansibleAnsible := ansible.ConfigureAnsible(ansibleConfig)
	deploymentProvisioner := provision.NewDeploymentProvisioner(ansibleAnsible, deploymentRepository)
	stepProvisionDeployment := provision.NewStepProvisionDeployment(deploymentProvisioner, jobRepository)
	failureMiddleware := provision.NewFailureMiddleware(jobRepository)
	transactional := middleware.NewTransactional(db)
	jobStateMachine := provision.ConfigureJobStateMachine(stepProvisionServer, stepProvisionDeployment, failureMiddleware, transactional)
	serverDestroyer := provision.NewServerDestroyer(terraformTerraform, db, serverRepository, deploymentRepository)
	provisioner := provision.NewProvisioner(jobStateMachine, jobScheduler, terraformTerraform, serverDestroyer)
	workerPool := provision.NewWorkerPool(workerPoolConfig, jobRepository, provisioner)
	appServer := NewAppServer(serverServer, workerPool, consoleLogger)
	return appServer, func() {
		cleanup()
	}, nil
}

// Injectors from inject_memory.go:

func SetupInMemoryApp() *App {
	provider := ProvideFileConfigProvider()
	config := ProvideConfig(provider)
	inMemoryRepository := account.NewInMemoryRepository()
	jwtConfig := config.JWT
	tokenService := account.ConfigureTokenService(jwtConfig)
	service := account.NewService(inMemoryRepository, tokenService)
	inMemoryProviderSettingsRepository := infrastructure.NewInMemoryProviderSettingsRepository()
	inMemoryServerRepository := infrastructure.NewInMemoryServerRepository()
	inMemoryJobRepository := provision.NewInMemoryJobRepository()
	inMemoryTxContext := transaction.NewInMemoryTransactionContext()
	inMemoryDeploymentRepository := infrastructure.NewInMemoryDeploymentRepository()
	jobScheduler := provision.NewJobScheduler(inMemoryTxContext, inMemoryJobRepository, inMemoryServerRepository, inMemoryDeploymentRepository)
	terraformConfig := config.Terraform
	terraformTerraform := terraform.ConfigureTerraform(terraformConfig)
	serverProvisioner := provision.NewServerProvisioner(terraformTerraform, inMemoryServerRepository)
	stepProvisionServer := provision.NewStepProvisionServer(serverProvisioner, inMemoryJobRepository)
	ansibleConfig := config.Ansible
	ansibleAnsible := ansible.ConfigureAnsible(ansibleConfig)
	deploymentProvisioner := provision.NewDeploymentProvisioner(ansibleAnsible, inMemoryDeploymentRepository)
	stepProvisionDeployment := provision.NewStepProvisionDeployment(deploymentProvisioner, inMemoryJobRepository)
	failureMiddleware := provision.NewFailureMiddleware(inMemoryJobRepository)
	transactional := middleware.NewTransactional(inMemoryTxContext)
	jobStateMachine := provision.ConfigureJobStateMachine(stepProvisionServer, stepProvisionDeployment, failureMiddleware, transactional)
	serverDestroyer := provision.NewServerDestroyer(terraformTerraform, inMemoryTxContext, inMemoryServerRepository, inMemoryDeploymentRepository)
	provisioner := provision.NewProvisioner(jobStateMachine, jobScheduler, terraformTerraform, serverDestroyer)
	logConfig := config.Log
	consoleLogger := log.NewConsoleLogger(logConfig)
	app := NewApp(config, inMemoryRepository, service, inMemoryProviderSettingsRepository, inMemoryServerRepository, inMemoryJobRepository, jobScheduler, provisioner, consoleLogger)
	return app
}

func SetupInMemoryServer() (*AppServer, func(), error) {
	provider := ProvideFileConfigProvider()
	config := ProvideConfig(provider)
	serverConfig := config.Server
	inMemoryRepository := account.NewInMemoryRepository()
	jwtConfig := config.JWT
	tokenService := account.ConfigureTokenService(jwtConfig)
	service := account.NewService(inMemoryRepository, tokenService)
	authenticationMiddleware := middleware2.NewAuthenticationMiddleware(service)
	authentication := routes.NewAuthenticationRoutes(service)
	routesAccount := routes.NewAccountRoutes(inMemoryRepository)
	inMemoryProviderSettingsRepository := infrastructure.NewInMemoryProviderSettingsRepository()
	providerSettings := routes.NewProviderSettingsRoutes(inMemoryProviderSettingsRepository)
	inMemoryTxContext := transaction.NewInMemoryTransactionContext()
	inMemoryJobRepository := provision.NewInMemoryJobRepository()
	inMemoryServerRepository := infrastructure.NewInMemoryServerRepository()
	inMemoryDeploymentRepository := infrastructure.NewInMemoryDeploymentRepository()
	jobScheduler := provision.NewJobScheduler(inMemoryTxContext, inMemoryJobRepository, inMemoryServerRepository, inMemoryDeploymentRepository)
	routesProvision := routes.NewProvisionRoutes(jobScheduler, inMemoryJobRepository, inMemoryProviderSettingsRepository)
	router := &httpserver.Router{
		AuthenticatedMiddleware: authenticationMiddleware,
		AuthRoutes:              authentication,
		AccountRoutes:           routesAccount,
		ProviderSettingsRoutes:  providerSettings,
		ProvisionRoutes:         routesProvision,
	}
	logConfig := config.Log
	consoleLogger := log.NewConsoleLogger(logConfig)
	serverServer, err := server.ProvideServer(serverConfig, router, consoleLogger)
	if err != nil {
		return nil, nil, err
	}
	workerPoolConfig := config.WorkerPool
	terraformConfig := config.Terraform
	terraformTerraform := terraform.ConfigureTerraform(terraformConfig)
	serverProvisioner := provision.NewServerProvisioner(terraformTerraform, inMemoryServerRepository)
	stepProvisionServer := provision.NewStepProvisionServer(serverProvisioner, inMemoryJobRepository)
	ansibleConfig := config.Ansible
	ansibleAnsible := ansible.ConfigureAnsible(ansibleConfig)
	deploymentProvisioner := provision.NewDeploymentProvisioner(ansibleAnsible, inMemoryDeploymentRepository)
	stepProvisionDeployment := provision.NewStepProvisionDeployment(deploymentProvisioner, inMemoryJobRepository)
	failureMiddleware := provision.NewFailureMiddleware(inMemoryJobRepository)
	transactional := middleware.NewTransactional(inMemoryTxContext)
	jobStateMachine := provision.ConfigureJobStateMachine(stepProvisionServer, stepProvisionDeployment, failureMiddleware, transactional)
	serverDestroyer := provision.NewServerDestroyer(terraformTerraform, inMemoryTxContext, inMemoryServerRepository, inMemoryDeploymentRepository)
	provisioner := provision.NewProvisioner(jobStateMachine, jobScheduler, terraformTerraform, serverDestroyer)
	workerPool := provision.NewWorkerPool(workerPoolConfig, inMemoryJobRepository, provisioner)
	appServer := NewAppServer(serverServer, workerPool, consoleLogger)
	return appServer, func() {
	}, nil
}

// Injectors from inject_testing.go:

func SetupTestApp(t *testing.T) *App {
	provider := ProvideTestConfigProvider()
	config := ProvideConfig(provider)
	inMemoryRepository := account.NewInMemoryRepository()
	jwtConfig := config.JWT
	tokenService := account.ConfigureTokenService(jwtConfig)
	service := account.NewService(inMemoryRepository, tokenService)
	inMemoryProviderSettingsRepository := infrastructure.NewInMemoryProviderSettingsRepository()
	inMemoryServerRepository := infrastructure.NewInMemoryServerRepository()
	inMemoryJobRepository := provision.NewInMemoryJobRepository()
	inMemoryTxContext := transaction.NewInMemoryTransactionContext()
	inMemoryDeploymentRepository := infrastructure.NewInMemoryDeploymentRepository()
	jobScheduler := provision.NewJobScheduler(inMemoryTxContext, inMemoryJobRepository, inMemoryServerRepository, inMemoryDeploymentRepository)
	terraformConfig := config.Terraform
	terraformTerraform := terraform.ConfigureTerraform(terraformConfig)
	serverProvisioner := provision.NewServerProvisioner(terraformTerraform, inMemoryServerRepository)
	stepProvisionServer := provision.NewStepProvisionServer(serverProvisioner, inMemoryJobRepository)
	ansibleConfig := config.Ansible
	ansibleAnsible := ansible.ConfigureAnsible(ansibleConfig)
	deploymentProvisioner := provision.NewDeploymentProvisioner(ansibleAnsible, inMemoryDeploymentRepository)
	stepProvisionDeployment := provision.NewStepProvisionDeployment(deploymentProvisioner, inMemoryJobRepository)
	failureMiddleware := provision.NewFailureMiddleware(inMemoryJobRepository)
	transactional := middleware.NewTransactional(inMemoryTxContext)
	jobStateMachine := provision.ConfigureJobStateMachine(stepProvisionServer, stepProvisionDeployment, failureMiddleware, transactional)
	serverDestroyer := provision.NewServerDestroyer(terraformTerraform, inMemoryTxContext, inMemoryServerRepository, inMemoryDeploymentRepository)
	provisioner := provision.NewProvisioner(jobStateMachine, jobScheduler, terraformTerraform, serverDestroyer)
	testingLogger := log.NewTestingLogger(t)
	app := NewApp(config, inMemoryRepository, service, inMemoryProviderSettingsRepository, inMemoryServerRepository, inMemoryJobRepository, jobScheduler, provisioner, testingLogger)
	return app
}

func SetupTestServer(t *testing.T) (*AppServer, func(), error) {
	provider := ProvideTestConfigProvider()
	config := ProvideConfig(provider)
	serverConfig := config.Server
	inMemoryRepository := account.NewInMemoryRepository()
	jwtConfig := config.JWT
	tokenService := account.ConfigureTokenService(jwtConfig)
	service := account.NewService(inMemoryRepository, tokenService)
	authenticationMiddleware := middleware2.NewAuthenticationMiddleware(service)
	authentication := routes.NewAuthenticationRoutes(service)
	routesAccount := routes.NewAccountRoutes(inMemoryRepository)
	inMemoryProviderSettingsRepository := infrastructure.NewInMemoryProviderSettingsRepository()
	providerSettings := routes.NewProviderSettingsRoutes(inMemoryProviderSettingsRepository)
	inMemoryTxContext := transaction.NewInMemoryTransactionContext()
	inMemoryJobRepository := provision.NewInMemoryJobRepository()
	inMemoryServerRepository := infrastructure.NewInMemoryServerRepository()
	inMemoryDeploymentRepository := infrastructure.NewInMemoryDeploymentRepository()
	jobScheduler := provision.NewJobScheduler(inMemoryTxContext, inMemoryJobRepository, inMemoryServerRepository, inMemoryDeploymentRepository)
	routesProvision := routes.NewProvisionRoutes(jobScheduler, inMemoryJobRepository, inMemoryProviderSettingsRepository)
	router := &httpserver.Router{
		AuthenticatedMiddleware: authenticationMiddleware,
		AuthRoutes:              authentication,
		AccountRoutes:           routesAccount,
		ProviderSettingsRoutes:  providerSettings,
		ProvisionRoutes:         routesProvision,
	}
	testingLogger := log.NewTestingLogger(t)
	serverServer, err := server.ProvideServer(serverConfig, router, testingLogger)
	if err != nil {
		return nil, nil, err
	}
	workerPoolConfig := config.WorkerPool
	terraformConfig := config.Terraform
	terraformTerraform := terraform.ConfigureTerraform(terraformConfig)
	serverProvisioner := provision.NewServerProvisioner(terraformTerraform, inMemoryServerRepository)
	stepProvisionServer := provision.NewStepProvisionServer(serverProvisioner, inMemoryJobRepository)
	ansibleConfig := config.Ansible
	ansibleAnsible := ansible.ConfigureAnsible(ansibleConfig)
	deploymentProvisioner := provision.NewDeploymentProvisioner(ansibleAnsible, inMemoryDeploymentRepository)
	stepProvisionDeployment := provision.NewStepProvisionDeployment(deploymentProvisioner, inMemoryJobRepository)
	failureMiddleware := provision.NewFailureMiddleware(inMemoryJobRepository)
	transactional := middleware.NewTransactional(inMemoryTxContext)
	jobStateMachine := provision.ConfigureJobStateMachine(stepProvisionServer, stepProvisionDeployment, failureMiddleware, transactional)
	serverDestroyer := provision.NewServerDestroyer(terraformTerraform, inMemoryTxContext, inMemoryServerRepository, inMemoryDeploymentRepository)
	provisioner := provision.NewProvisioner(jobStateMachine, jobScheduler, terraformTerraform, serverDestroyer)
	workerPool := provision.NewWorkerPool(workerPoolConfig, inMemoryJobRepository, provisioner)
	appServer := NewAppServer(serverServer, workerPool, testingLogger)
	return appServer, func() {
	}, nil
}

// inject_database.go:

var dbAppSet = wire.NewSet(
	ProvideFileConfigProvider, log.NewConsoleLogger, wire.Bind(new(log.Logger), new(*log.ConsoleLogger)), database.Set, database.NewAccountRepository, wire.Bind(new(account.Repository), new(*database.AccountRepository)), database.NewJobRepository, wire.Bind(new(provision.JobRepository), new(*database.JobRepository)), database.NewServerRepository, wire.Bind(new(infrastructure.ServerRepository), new(*database.ServerRepository)), database.NewDeploymentRepository, wire.Bind(new(infrastructure.DeploymentRepository), new(*database.DeploymentRepository)), database.NewProviderSettingsRepository, wire.Bind(new(infrastructure.ProviderSettingsRepository), new(*database.ProviderSettingsRepository)), AppSet,
)

// inject_memory.go:

var inMemAppSet = wire.NewSet(
	ProvideFileConfigProvider, log.NewConsoleLogger, wire.Bind(new(log.Logger), new(*log.ConsoleLogger)), transaction.NewInMemoryTransactionContext, wire.Bind(new(transaction.TxContext), new(*transaction.InMemoryTxContext)), account.NewInMemoryRepository, wire.Bind(new(account.Repository), new(*account.InMemoryRepository)), provision.NewInMemoryJobRepository, wire.Bind(new(provision.JobRepository), new(*provision.InMemoryJobRepository)), infrastructure.NewInMemoryServerRepository, wire.Bind(new(infrastructure.ServerRepository), new(*infrastructure.InMemoryServerRepository)), infrastructure.NewInMemoryDeploymentRepository, wire.Bind(new(infrastructure.DeploymentRepository), new(*infrastructure.InMemoryDeploymentRepository)), infrastructure.NewInMemoryProviderSettingsRepository, wire.Bind(new(infrastructure.ProviderSettingsRepository), new(*infrastructure.InMemoryProviderSettingsRepository)), AppSet,
)

// inject_testing.go:

var testAppSet = wire.NewSet(
	ProvideTestConfigProvider, log.NewTestingLogger, wire.Bind(new(log.Logger), new(*log.TestingLogger)), transaction.NewInMemoryTransactionContext, wire.Bind(new(transaction.TxContext), new(*transaction.InMemoryTxContext)), account.NewInMemoryRepository, wire.Bind(new(account.Repository), new(*account.InMemoryRepository)), provision.NewInMemoryJobRepository, wire.Bind(new(provision.JobRepository), new(*provision.InMemoryJobRepository)), infrastructure.NewInMemoryServerRepository, wire.Bind(new(infrastructure.ServerRepository), new(*infrastructure.InMemoryServerRepository)), infrastructure.NewInMemoryDeploymentRepository, wire.Bind(new(infrastructure.DeploymentRepository), new(*infrastructure.InMemoryDeploymentRepository)), infrastructure.NewInMemoryProviderSettingsRepository, wire.Bind(new(infrastructure.ProviderSettingsRepository), new(*infrastructure.InMemoryProviderSettingsRepository)), AppSet,
)