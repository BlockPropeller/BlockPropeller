package routes

import (
	"context"

	"blockpropeller.dev/blockpropeller/httpserver/request"
	"blockpropeller.dev/blockpropeller/infrastructure"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

// ListDeploymentsResponse is a response to the list deployments request.
type ListDeploymentsResponse struct {
	Deployments []*infrastructure.Deployment `json:"deployments"`
}

// Deployment REST Resource for accessing deployment information.
type Deployment struct {
	deploymentRepo infrastructure.DeploymentRepository
}

// NewDeploymentRoutes returns a new Deployment routes instance.
func NewDeploymentRoutes(deploymentRepo infrastructure.DeploymentRepository) *Deployment {
	return &Deployment{deploymentRepo: deploymentRepo}
}

// List all Deployments for a Server.
func (s *Deployment) List(c echo.Context) error {
	srv := request.ServerFromContext(c)
	if srv == nil {
		return echo.ErrNotFound.SetInternal(errors.New("server not found in context"))
	}

	deployments, err := s.deploymentRepo.FindByServer(context.Background(), srv.ID)
	if err != nil {
		return errors.Wrap(err, "list deployments")
	}

	return c.JSON(200, &ListDeploymentsResponse{
		Deployments: deployments,
	})
}
