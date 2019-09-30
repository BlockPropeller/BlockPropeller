package provision

import (
	"context"
	"time"

	"chainup.dev/chainup/ansible"
	"chainup.dev/chainup/infrastructure"
	"chainup.dev/lib/log"
	"github.com/pkg/errors"
)

var (
	// ErrServerNotReadyForDeployments is returned for Servers that are not ready to accept new deployments.
	ErrServerNotReadyForDeployments = errors.New("server not ready for deployments")
	// ErrDeploymentNotInRequestedState is returned for Deployments that are not ready to be provisioned on a Server.
	ErrDeploymentNotInRequestedState = errors.New("deployment not in requested state")
)

// DeploymentProvisioner is responsible for configuring Deployments on a target
// Server via Ansible playbooks.
type DeploymentProvisioner struct {
	ans *ansible.Ansible

	deploymentRepo infrastructure.DeploymentRepository
}

// NewDeploymentProvisioner returns a new DeploymentProvisioner instance.
func NewDeploymentProvisioner(ans *ansible.Ansible, deploymentRepo infrastructure.DeploymentRepository) *DeploymentProvisioner {
	return &DeploymentProvisioner{ans: ans, deploymentRepo: deploymentRepo}
}

// Provision configures the specified Deployment on a target Server.
func (dp *DeploymentProvisioner) Provision(ctx context.Context, srv *infrastructure.Server, deployment *infrastructure.Deployment) error {
	if srv.State != infrastructure.ServerStateOk {
		return ErrServerNotReadyForDeployments
	}

	if deployment.State != infrastructure.DeploymentStateRequested {
		return ErrDeploymentNotInRequestedState
	}

	version, err := dp.ans.Version()
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

		err = dp.ans.RunPlaybook(srv, deployment)
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

	deployment.State = infrastructure.DeploymentStateOk

	err = dp.deploymentRepo.Update(ctx, deployment)
	if err != nil {
		return errors.Wrap(err, "update deployment")
	}

	return nil
}
