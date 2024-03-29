package database

import (
	"context"
	"encoding/json"

	"blockpropeller.dev/blockpropeller/infrastructure"
	"github.com/pkg/errors"
)

// DeploymentRepository is a databased backed implementation of a infrastructure.DeploymentRepository.
// @TODO: See if we can use GORM hooks for preparing and parsing the deployments.
type DeploymentRepository struct {
	db *DB
}

// NewDeploymentRepository returns a new DeploymentRepository instance.
func NewDeploymentRepository(db *DB) *DeploymentRepository {
	return &DeploymentRepository{db: db}
}

// Find a Deployment given a DeploymentID.
func (repo *DeploymentRepository) Find(ctx context.Context, id infrastructure.DeploymentID) (*infrastructure.Deployment, error) {
	var deployment infrastructure.Deployment
	err := repo.db.Model(ctx, &deployment).Where("id = ?", id).First(&deployment).Error
	if err != nil {
		return nil, errors.Wrap(err, "find deployment by ID")
	}

	return &deployment, nil
}

// FindByServer returns deployments on a given Server.
func (repo *DeploymentRepository) FindByServer(ctx context.Context, id infrastructure.ServerID) ([]*infrastructure.Deployment, error) {
	var deployments []*infrastructure.Deployment

	err := repo.db.Model(ctx, &deployments).
		Where("server_id = ?", id).
		Find(&deployments).
		Error
	if err != nil {
		return nil, errors.Wrap(err, "find deployments by server")
	}

	return deployments, nil
}

// Create a new Deployment.
func (repo *DeploymentRepository) Create(ctx context.Context, deployment *infrastructure.Deployment) error {
	err := repo.prepareDeployment(deployment)
	if err != nil {
		return errors.Wrap(err, "prepare deployment for creation")
	}

	err = repo.db.Model(ctx, deployment).Create(deployment).Error
	if err != nil {
		return errors.Wrap(err, "create deployment")
	}

	return nil
}

// Update an existing Deployment.
func (repo *DeploymentRepository) Update(ctx context.Context, deployment *infrastructure.Deployment) error {
	err := repo.prepareDeployment(deployment)
	if err != nil {
		return errors.Wrap(err, "prepare deployment for update")
	}

	err = repo.db.Model(ctx, deployment).Save(deployment).Error
	if err != nil {
		return errors.Wrap(err, "update deployment")
	}

	return nil
}

// DeleteForServer deletes all deployments associated with a given Server.
func (repo *DeploymentRepository) DeleteForServer(ctx context.Context, srv *infrastructure.Server) error {
	err := repo.db.Model(ctx, infrastructure.Deployment{}).
		Where("server_id = ?", srv.ID).
		Updates(infrastructure.Deployment{State: infrastructure.DeploymentStateDeleted}).
		Error
	if err != nil {
		return errors.Wrap(err, "set state as deleted")
	}

	err = repo.db.Model(ctx, (*infrastructure.Deployment)(nil)).
		Where("server_id = ?", srv.ID).
		Delete(&infrastructure.Deployment{}).
		Error
	if err != nil {
		return errors.Wrap(err, "delete deployments for server")
	}

	return nil
}

func (repo *DeploymentRepository) prepareDeployment(deployment *infrastructure.Deployment) error {
	data, err := json.Marshal(deployment.Configuration.MarshalMap())
	if err != nil {
		return errors.Wrap(err, "marshal config")
	}

	deployment.RawConfiguration = string(data)

	return nil
}
