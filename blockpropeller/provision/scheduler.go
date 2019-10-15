package provision

import (
	"context"

	"blockpropeller.dev/blockpropeller/database/transaction"
	"blockpropeller.dev/blockpropeller/infrastructure"
	"github.com/pkg/errors"
)

// JobScheduler is responsible for taking a Job request, persisting it
// and queuing the job in order to be executed by the provisioner.
type JobScheduler struct {
	txContext transaction.TxContext

	jobRepo        JobRepository
	serverRepo     infrastructure.ServerRepository
	deploymentRepo infrastructure.DeploymentRepository
}

// NewJobScheduler returns a new JobScheduler instance.
func NewJobScheduler(
	txContext transaction.TxContext,
	jobRepo JobRepository,
	serverRepo infrastructure.ServerRepository,
	deploymentRepo infrastructure.DeploymentRepository,
) *JobScheduler {
	return &JobScheduler{
		txContext:      txContext,
		jobRepo:        jobRepo,
		serverRepo:     serverRepo,
		deploymentRepo: deploymentRepo,
	}
}

// Schedule a new Job by saving it to the repositories.
//
// The job will be picked up by the Provider by polling the JobRepository.
func (js *JobScheduler) Schedule(ctx context.Context, job *Job) error {
	err := js.txContext.RunInTransaction(ctx, func(ctx context.Context) error {
		err := js.serverRepo.Create(ctx, job.Server)
		if err != nil {
			return errors.Wrap(err, "create server request")
		}

		err = js.deploymentRepo.Create(ctx, job.Deployment)
		if err != nil {
			return errors.Wrap(err, "create deployment request")
		}

		err = js.jobRepo.Create(ctx, job)
		if err != nil {
			return errors.Wrap(err, "create job request")
		}

		return nil
	})
	if err != nil {
		return errors.Wrap(err, "failed scheduling job")
	}

	return nil
}

// FindScheduled fetches a job to be provisioned.
func (js *JobScheduler) FindScheduled(ctx context.Context, jobID JobID) (*Job, error) {
	job, err := js.jobRepo.Find(ctx, jobID)
	if err != nil {
		return nil, errors.Wrap(err, "find job")
	}

	job.Server, err = js.serverRepo.Find(ctx, job.ServerID)
	if err != nil {
		return nil, errors.Wrap(err, "find job server")
	}

	job.Deployment, err = js.deploymentRepo.Find(ctx, job.DeploymentID)
	if err != nil {
		return nil, errors.Wrap(err, "find job deployment")
	}

	return job, nil
}
