package provision

import (
	"context"
	"net"
	"time"

	"chainup.dev/chainup/ansible"
	"chainup.dev/chainup/infrastructure"
	"chainup.dev/chainup/statemachine"
	"chainup.dev/chainup/statemachine/middleware"
	"chainup.dev/chainup/terraform"
	"chainup.dev/chainup/terraform/resource"
	"chainup.dev/chainup/terraform/resource/digitalocean"
	"chainup.dev/lib/log"
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
	tfStep *TerraformStep,
	ansibleStep *AnsibleStep,
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

// TerraformStep creates a plan for creating new infrastructure,
// executes it against the given cloud provider and waits for the
// provisioning to finish.
type TerraformStep struct {
	tf *terraform.Terraform

	srvRepo infrastructure.ServerRepository
	jobRepo JobRepository
}

// NewTerraformStep returns a new TerraformStep instance.
func NewTerraformStep(tf *terraform.Terraform, srvRepo infrastructure.ServerRepository, jobRepo JobRepository) *TerraformStep {
	return &TerraformStep{tf: tf, srvRepo: srvRepo, jobRepo: jobRepo}
}

// Step satisfies the State Machine step interface.
func (step *TerraformStep) Step(ctx context.Context, res statemachine.StatefulResource) error {
	job := res.(*Job)

	providerSettings := job.ProviderSettings

	server := job.Server
	sshKey := server.SSHKey

	// Prepare workspace in which to execute Terraform plan.
	workspace, err := terraform.NewWorkspace()
	if err != nil {
		return errors.Wrap(err, "create new workspace")
	}

	defer func() {
		log.Debug("cleaning up Terraform workspace")
		log.Closer(workspace)
	}()

	log.Debug("created Terraform workspace", log.Fields{
		"dir": workspace.WorkDir(),
	})

	workspace.Add(digitalocean.NewProvider(providerSettings.Credentials))

	log.Debug("using provider", log.Fields{
		"type":        providerSettings.Type,
		"credentials": providerSettings.Credentials,
	})

	doSSHKey := digitalocean.NewSSHKey(sshKey.Name, sshKey.EncodedPublicKey())
	log.Debug("using ssh key", log.Fields{
		"pub":  sshKey.EncodedPublicKey(),
		"priv": sshKey.EncodedPrivateKey(),
	})

	doDroplet := digitalocean.NewDroplet(
		server.Name,
		"ubuntu-18-04-x64",
		"fra1",
		"s-1vcpu-1gb",
		[]*digitalocean.SSHKey{doSSHKey},
	)

	workspace.AddResource(doSSHKey, doDroplet)

	ipAddressOut := resource.NewOutput("ip-address", resource.ToPropSelector(doDroplet, "ipv4_address"))

	workspace.Add(ipAddressOut)

	err = workspace.Flush()
	if err != nil {
		return errors.Wrap(err, "flush workspace")
	}

	log.Debug("running terraform init...")

	err = step.tf.Init(workspace)
	if err != nil {
		return errors.Wrap(err, "init workspace")
	}

	log.Debug("running terraform plan...")

	err = step.tf.Plan(workspace)
	if err != nil {
		return errors.Wrap(err, "prepare execution plan")
	}

	log.Debug("running terraform apply...")

	err = step.tf.Apply(workspace)
	if err != nil {
		return errors.Wrap(err, "apply execution plan")
	}

	log.Debug("running terraform output...")

	rawIP, err := step.tf.Output(workspace, "ip-address")
	if err != nil {
		return errors.Wrap(err, "get ip address of provisioned server")
	}

	ip := net.ParseIP(rawIP)
	if ip == nil {
		return errors.Errorf("invalid server IP: %s", rawIP)
	}

	log.Debug("server provisioned", log.Fields{
		"ip": ip.String(),
	})

	server.IPAddress = ip
	server.State = infrastructure.ServerStateRunning

	snap, err := workspace.Snapshot()
	if err != nil {
		return errors.Wrap(err, "take workspace snapshot")
	}

	job.WorkspaceSnapshot = snap
	job.SetState(StateServerCreated)

	err = step.srvRepo.Update(ctx, server)
	if err != nil {
		return errors.Wrap(err, "update server")
	}

	err = step.jobRepo.Update(ctx, job)
	if err != nil {
		return errors.Wrap(err, "update job")
	}

	return nil
}

// AnsibleStep connects to a previously created server
// and runs an Ansible playbook for provisioning deployments on top of it.
type AnsibleStep struct {
	ans *ansible.Ansible

	deploymentRepo infrastructure.DeploymentRepository
	jobRepo        JobRepository
}

// NewAnsibleStep returns a new AnsibleStep instance.
func NewAnsibleStep(ans *ansible.Ansible, deploymentRepo infrastructure.DeploymentRepository, jobRepo JobRepository) *AnsibleStep {
	return &AnsibleStep{ans: ans, deploymentRepo: deploymentRepo, jobRepo: jobRepo}
}

// Step satisfies the Step interface.
func (step *AnsibleStep) Step(ctx context.Context, res statemachine.StatefulResource) error {
	job := res.(*Job)
	srv := job.Server
	deployment := job.Deployment

	version, err := step.ans.Version()
	if err != nil {
		return errors.Wrap(err, "check ansible version")
	}

	log.Debug("using ansible", log.Fields{
		"version": version,
	})

	log.Debug("running playbook...")

	for tries := 5; tries > 0; tries-- {
		log.Debug("waiting for server to become available", log.Fields{
			"seconds": 5,
		})
		time.Sleep(5 * time.Second)

		err = step.ans.RunPlaybook(srv, deployment)
		if err != nil {
			log.ErrorErr(err, "failed running playbook on server", log.Fields{
				"tries": tries,
			})
			continue
		}

		break
	}
	if err != nil {
		return errors.Wrap(err, "failed running playbook on server")
	}

	deployment.State = infrastructure.DeploymentStateRunning

	job.SetState(StateCompleted)

	err = step.deploymentRepo.Update(ctx, deployment)
	if err != nil {
		return errors.Wrap(err, "update deployment")
	}

	err = step.jobRepo.Update(ctx, job)
	if err != nil {
		return errors.Wrap(err, "update job")
	}

	return nil
}
