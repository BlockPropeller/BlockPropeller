package routes

import (
	"context"

	"blockpropeller.dev/blockpropeller/binance"
	"blockpropeller.dev/blockpropeller/httpserver/request"
	"blockpropeller.dev/blockpropeller/infrastructure"
	"blockpropeller.dev/blockpropeller/provision"
	"github.com/blang/semver"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

// GetJobResponse is a response to the get job request.
type GetJobResponse struct {
	Job *provision.Job `json:"job"`
}

// CreateJobRequest holds the request payload for the create job endpoint.
type CreateJobRequest struct {
	ProviderSettingsID infrastructure.ProviderSettingsID `json:"provider_id" form:"provider_id" validate:"required"`

	NodeNetwork binance.Network  `json:"node_network" form:"node_network" validate:"required,valid"`
	NodeType    binance.NodeType `json:"node_type" form:"node_type" validate:"required,valid"`
	NodeVersion string           `json:"node_version" form:"node_version" validate:"required"`
}

// CreateJobResponse is a response to the create job request.
type CreateJobResponse struct {
	Job *provision.Job `json:"job"`
}

// Provision routes define ways to provision infrastructure via BlockPropeller.
type Provision struct {
	jobScheduler *provision.JobScheduler

	jobRepo      provision.JobRepository
	settingsRepo infrastructure.ProviderSettingsRepository
}

// NewProvisionRoutes returns a new Provision routes instance.
func NewProvisionRoutes(jobScheduler *provision.JobScheduler, jobRepo provision.JobRepository, settingsRepo infrastructure.ProviderSettingsRepository) *Provision {
	return &Provision{jobScheduler: jobScheduler, jobRepo: jobRepo, settingsRepo: settingsRepo}
}

// LoadJob is a middleware for loading Jobs into request context
// as well as checking for correct permissions of an authenticated user.
func (p *Provision) LoadJob(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authAcc := request.AuthFromContext(c)
		if authAcc == nil {
			return echo.ErrForbidden.SetInternal(errors.New("missing authenticated account"))
		}

		jobID := provision.JobID(c.Param("job_id"))
		if jobID.String() == "" {
			return next(c)
		}

		job, err := p.jobRepo.Find(context.Background(), jobID)
		if err != nil {
			return echo.ErrInternalServerError.SetInternal(err)
		}

		if authAcc.ID.String() != job.AccountID.String() {
			return echo.ErrForbidden.
				SetInternal(errors.Errorf("unauthorized job access: authenticated %s, job %s",
					authAcc.ID, jobID))
		}

		request.WithJob(c, job)

		return next(c)
	}
}

// GetJob returns a requested job.
func (p *Provision) GetJob(c echo.Context) error {
	job := request.JobFromContext(c)
	if job == nil {
		return echo.ErrNotFound.SetInternal(errors.New("job not found in context"))
	}

	return c.JSON(200, &GetJobResponse{Job: job})
}

// CreateJob creates a new Job to be executed and returns it.
func (p *Provision) CreateJob(c echo.Context) error {
	var req CreateJobRequest
	if err := request.Parse(c, &req); err != nil {
		return err
	}

	acc := request.AuthFromContext(c)
	if acc == nil {
		return echo.ErrInternalServerError.SetInternal(errors.New("missing authenticated user"))
	}

	settings, err := p.settingsRepo.Find(context.Background(), req.ProviderSettingsID)
	if err != nil {
		return errors.Wrap(err, "find provider settings")
	}
	if settings.AccountID.String() != acc.ID.String() {
		return echo.ErrForbidden.
			SetInternal(errors.Errorf("unauthorized job access: authenticated %s, provider settings %s",
				acc.ID, req.ProviderSettingsID))
	}

	nodeVersion, err := semver.Parse(req.NodeVersion)
	if err != nil {
		return echo.ErrBadRequest.SetInternal(err)
	}

	deployment := binance.NewNodeDeployment(req.NodeNetwork, req.NodeType, nodeVersion)

	srv, err := infrastructure.NewServerBuilder(acc.ID).
		Provider(settings.Type).
		Build()
	if err != nil {
		return errors.Wrap(err, "build server")
	}

	job, err := provision.NewJobBuilder(acc.ID).
		Provider(settings).
		Server(srv).
		Deployment(deployment).
		Build()
	if err != nil {
		return errors.Wrap(err, "build job")
	}

	err = p.jobScheduler.Schedule(context.Background(), job)
	if err != nil {
		return errors.Wrap(err, "create job")
	}

	return c.JSON(201, &CreateJobResponse{Job: job})
}
