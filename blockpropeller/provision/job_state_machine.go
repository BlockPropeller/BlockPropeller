package provision

import (
	"context"
	"time"

	"blockpropeller.dev/blockpropeller/statemachine"
	"blockpropeller.dev/blockpropeller/statemachine/middleware"
	"blockpropeller.dev/lib/log"
	"github.com/pkg/errors"
)

var (
	// StateCreated is the starting point for a provisioning job.
	StateCreated = statemachine.NewState("job_created")

	// StateServerCreated is the state after terraform successfully creates the requested server.
	StateServerCreated = statemachine.NewState("server_created")

	// StateCompleted is the terminating state representing a successful provisioning job.
	StateCompleted = statemachine.NewState("completed").Successful()

	// StateFailed is the terminating state representing provisioning server failure.
	// @TODO: Add failure message to job somewhere.
	StateFailed = statemachine.NewState("failed").Failure()

	// ValidStates of a provision.Job.
	ValidStates = []statemachine.State{StateCreated, StateServerCreated, StateCompleted, StateFailed}
)

// JobStateMachine defines the state machine for running provisioning jobs.
type JobStateMachine struct {
	*statemachine.StateMachine
}

// ConfigureJobStateMachine returns a preconfigured StateMachine
// for running provisioning jobs.
func ConfigureJobStateMachine(
	tfStep *StepProvisionServer,
	ansibleStep *StepProvisionDeployment,
	failureMiddleware *FailureMiddleware,
	txMiddleware *middleware.Transactional,
) *JobStateMachine {
	return &JobStateMachine{
		StateMachine: statemachine.Builder(ValidStates).
			Middleware(failureMiddleware).
			Middleware(txMiddleware).
			Step(StateCreated, tfStep).
			Step(StateServerCreated, ansibleStep).
			Build(),
	}
}

// FailureMiddleware transitions a Job into failed state if an error is returned
// from a regular step.
type FailureMiddleware struct {
	jobRepo JobRepository
}

// NewFailureMiddleware returns a new FailureMiddleware instance.
func NewFailureMiddleware(jobRepo JobRepository) *FailureMiddleware {
	return &FailureMiddleware{jobRepo: jobRepo}
}

// Wrap satisfies the Middleware interface.
func (f *FailureMiddleware) Wrap(step statemachine.Step) statemachine.Step {
	return statemachine.StepFn(func(ctx context.Context, res statemachine.StatefulResource) error {
		err := step.Step(ctx, res)
		if err == nil {
			// Resume execution on no error.
			return nil
		}

		// Mark the Job as failed.
		job, ok := res.(*Job)
		if !ok {
			panic("expected Job instance in FailureMiddleware")
		}

		log.ErrorErr(err, "failed running job state machine", log.Fields{
			"job_id":    job.ID,
			"last_step": job.GetState(),
		})

		job.SetState(StateFailed)

		now := time.Now()
		job.FinishedAt = &now

		err = f.jobRepo.Update(ctx, job)
		if err != nil {
			return errors.Wrap(err, "update job to failed state")
		}

		return nil
	})
}

// StepProvisionServer creates a plan for creating new infrastructure,
// executes it against the given cloud provider and waits for the
// provisioning to finish.
type StepProvisionServer struct {
	serverProvisioner *ServerProvisioner

	jobRepo JobRepository
}

// NewStepProvisionServer returns a new StepProvisionServer instance.
func NewStepProvisionServer(serverProvisioner *ServerProvisioner, jobRepo JobRepository) *StepProvisionServer {
	return &StepProvisionServer{serverProvisioner: serverProvisioner, jobRepo: jobRepo}
}

// Step satisfies the State Machine step interface.
func (step *StepProvisionServer) Step(ctx context.Context, res statemachine.StatefulResource) error {
	job := res.(*Job)

	err := step.serverProvisioner.Provision(ctx, job.ProviderSettings, job.Server)
	if err != nil {
		return errors.Wrap(err, "run server provisioning")
	}

	job.SetState(StateServerCreated)

	err = step.jobRepo.Update(ctx, job)
	if err != nil {
		return errors.Wrap(err, "update job")
	}

	return nil
}

// StepProvisionDeployment connects to a previously created server
// and runs an Ansible playbook for provisioning deployments on top of it.
type StepProvisionDeployment struct {
	deploymentProvisioner *DeploymentProvisioner

	jobRepo JobRepository
}

// NewStepProvisionDeployment returns a new StepProvisionDeployment instance.
func NewStepProvisionDeployment(deploymentProvisioner *DeploymentProvisioner, jobRepo JobRepository) *StepProvisionDeployment {
	return &StepProvisionDeployment{deploymentProvisioner: deploymentProvisioner, jobRepo: jobRepo}
}

// Step satisfies the Step interface.
func (step *StepProvisionDeployment) Step(ctx context.Context, res statemachine.StatefulResource) error {
	job := res.(*Job)

	err := step.deploymentProvisioner.Provision(ctx, job.Server, job.Deployment)
	if err != nil {
		return errors.Wrap(err, "failed provisioning deployment")
	}

	job.SetState(StateCompleted)

	finishedAt := time.Now()
	job.FinishedAt = &finishedAt

	err = step.jobRepo.Update(ctx, job)
	if err != nil {
		return errors.Wrap(err, "update job")
	}

	return nil
}
