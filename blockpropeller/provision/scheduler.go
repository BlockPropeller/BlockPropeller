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
