package provision

import (
	"context"

	"blockpropeller.dev/blockpropeller/database/transaction"
	"blockpropeller.dev/blockpropeller/infrastructure"
	"blockpropeller.dev/blockpropeller/terraform"
	"blockpropeller.dev/lib/log"
	"github.com/pkg/errors"
)

// ServerDestroyer is responsible for, given a server entity, destroying the
// infrastructure associated with it via Terraform.
type ServerDestroyer struct {
	tf *terraform.Terraform

	txContext      transaction.TxContext
	srvRepo        infrastructure.ServerRepository
	deploymentRepo infrastructure.DeploymentRepository
}

// NewServerDestroyer returns a new ServerDestroyer instance.
func NewServerDestroyer(tf *terraform.Terraform, txContext transaction.TxContext, srvRepo infrastructure.ServerRepository, deploymentRepo infrastructure.DeploymentRepository) *ServerDestroyer {
	return &ServerDestroyer{tf: tf, txContext: txContext, srvRepo: srvRepo, deploymentRepo: deploymentRepo}
}

// Destroy runs the destruction of resources associated with the Server entity.
func (sd *ServerDestroyer) Destroy(ctx context.Context, srv *infrastructure.Server) error {
	if srv.WorkspaceSnapshot == nil {
		return errors.New("missing workspace snapshot")
	}

	workspace, err := terraform.RestoreWorkspace(srv.WorkspaceSnapshot)
	if err != nil {
		return errors.Wrap(err, "restore workspace")
	}
	defer func() {
		log.Debug("cleaning up Terraform workspace")
		log.Closer(workspace)
	}()

	err = sd.tf.Init(workspace)
	if err != nil {
		return errors.Wrap(err, "init workspace")
	}

	err = sd.tf.Destroy(workspace)
	if err != nil {
		return errors.Wrap(err, "destroy workspace")
	}

	srv.WorkspaceSnapshot, err = workspace.Snapshot()
	if err != nil {
		return errors.Wrap(err, "snapshot workspace")
	}

	return sd.txContext.RunInTransaction(ctx, func(ctx context.Context) error {
		err = sd.deploymentRepo.DeleteForServer(ctx, srv)
		if err != nil {
			return errors.Wrap(err, "delete deployments")
		}

		err = sd.srvRepo.Delete(ctx, srv)
		if err != nil {
			return errors.Wrap(err, "delete server")
		}

		return nil
	})
}
