package provision

import (
	"context"
	"net"

	"blockpropeller.dev/blockpropeller/infrastructure"
	"blockpropeller.dev/blockpropeller/terraform"
	"blockpropeller.dev/blockpropeller/terraform/cloudprovider"
	"blockpropeller.dev/lib/log"
	"github.com/pkg/errors"
)

var (
	// ErrServerNotReadyForProvisioning is returned for Servers that are not in an appropriate state for provisioning.
	ErrServerNotReadyForProvisioning = errors.New("server not ready for provisioning")
)

// ServerProvisioner takes a Server and provisions desired infrastructure
// to the defined provider using Terraform.
type ServerProvisioner struct {
	tf *terraform.Terraform

	srvRepo infrastructure.ServerRepository
}

// NewServerProvisioner returns a new ServerProvisioner instance.
func NewServerProvisioner(tf *terraform.Terraform, srvRepo infrastructure.ServerRepository) *ServerProvisioner {
	return &ServerProvisioner{tf: tf, srvRepo: srvRepo}
}

// Provision runs the provisioning via Terraform for the provided provider and server spec.
func (sp *ServerProvisioner) Provision(ctx context.Context, provider *infrastructure.ProviderSettings, srv *infrastructure.Server) error {
	if srv.State != infrastructure.ServerStateRequested {
		return ErrServerNotReadyForProvisioning
	}

	var workspace *terraform.Workspace
	var err error
	if srv.WorkspaceSnapshot == nil {
		workspace, err = sp.setupWorkspace(provider, srv)
	} else {
		//@TODO: Test out this code path.
		workspace, err = terraform.RestoreWorkspace(srv.WorkspaceSnapshot)
	}
	if err != nil {
		return errors.Wrap(err, "failed setting up workspace")
	}
	defer func() {
		log.Debug("cleaning up Terraform workspace")
		log.Closer(workspace)
	}()

	log.Debug("running terraform init...")

	err = sp.tf.Init(workspace)
	if err != nil {
		return errors.Wrap(err, "init workspace")
	}

	log.Debug("running terraform plan...")

	err = sp.tf.Plan(workspace)
	if err != nil {
		return errors.Wrap(err, "prepare execution plan")
	}

	log.Debug("running terraform apply...")

	err = sp.tf.Apply(workspace)
	if err != nil {
		return errors.Wrap(err, "apply execution plan")
	}

	log.Debug("running terraform output...")

	rawIP, err := sp.tf.Output(workspace, "ip-address")
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

	srv.IPAddress = ip.String()
	srv.State = infrastructure.ServerStateOk

	snap, err := workspace.Snapshot()
	if err != nil {
		return errors.Wrap(err, "take workspace snapshot")
	}

	srv.WorkspaceSnapshot = snap

	err = sp.srvRepo.Update(ctx, srv)
	if err != nil {
		return errors.Wrap(err, "update server")
	}

	return nil
}

func (sp *ServerProvisioner) setupWorkspace(provider *infrastructure.ProviderSettings, srv *infrastructure.Server) (*terraform.Workspace, error) {
	// Prepare workspace in which to execute Terraform plan.
	workspace, err := terraform.NewWorkspace()
	if err != nil {
		return nil, errors.Wrap(err, "create new workspace")
	}

	log.Debug("created Terraform workspace", log.Fields{
		"dir": workspace.WorkDir(),
	})

	cloudProvider, err := cloudprovider.GetProvider(provider.Type)
	if err != nil {
		return nil, errors.Wrap(err, "get cloud provider")
	}

	err = cloudProvider.Register(workspace, provider)
	if err != nil {
		return nil, errors.Wrap(err, "register cloud provider in workspace")
	}

	log.Debug("using provider", log.Fields{
		"type":        provider.Type,
		"credentials": provider.Credentials,
	})

	err = cloudProvider.AddServer(workspace, srv)
	if err != nil {
		return nil, errors.Wrap(err, "add server to workspace")
	}

	err = workspace.Flush()
	if err != nil {
		return nil, errors.Wrap(err, "flush workspace")
	}

	return workspace, nil
}
