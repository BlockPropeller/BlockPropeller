package provision

import (
	"context"
	"time"

	"chainup.dev/chainup/statemachine"
	"chainup.dev/chainup/statemachine/middleware"
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
	txMiddleware *middleware.Transactional,
) *JobStateMachine {
	return &JobStateMachine{
		StateMachine: statemachine.Builder(ValidStates).
			Middleware(txMiddleware).
			Step(StateCreated, tfStep).
			Step(StateServerCreated, ansibleStep).
			Build(),
	}
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
