package infrastructure

import (
	"context"
	"sync"
	"time"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

var (
	// ErrDeploymentNotFound is returned when a DeploymentRepository does not find a deployment to return.
	ErrDeploymentNotFound = errors.New("deployment not found")
	// ErrDeploymentAlreadyExists is returned when a Deployment creation is attempted with an existing DeploymentID.
	ErrDeploymentAlreadyExists = errors.New("deployment already exists")
)

// DeploymentID is a unique server identifier.
type DeploymentID string

// NewDeploymentID returns a new unique DeploymentID.
func NewDeploymentID() DeploymentID {
	return DeploymentID(uuid.NewV4().String())
}

// String satisfies the Stringer interface.
func (id DeploymentID) String() string {
	return string(id)
}

// Deployment is used to define what service needs to be provisioned on a particular Server.
type Deployment struct {
	ID       DeploymentID `json:"id"`
	ServerID ServerID     `json:"-" sql:"type:varchar(255) REFERENCES servers(id)"`

	Type             DeploymentType   `json:"type"`
	Configuration    DeploymentConfig `json:"configuration" gorm:"-"`
	RawConfiguration string           `json:"-" gorm:"configuration"`

	State DeploymentState `json:"state"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"-"`
}

// NewDeployment returns a new Deployment instance.
func NewDeployment(typ DeploymentType, config DeploymentConfig) *Deployment {
	return &Deployment{
		ID: NewDeploymentID(),

		Type:          typ,
		Configuration: config,

		State: DeploymentStateRequested,

		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// DeploymentRepository defines an interface for storing and retrieving deployments.
//
// @TODO: Consider moving deployments under the Server repository.
type DeploymentRepository interface {
	// Find a Deployment given a DeploymentID.
	Find(ctx context.Context, id DeploymentID) (*Deployment, error)

	// Create a new Deployment.
	Create(ctx context.Context, deployment *Deployment) error

	// Update an existing Deployment.
	Update(ctx context.Context, deployment *Deployment) error

	// DeleteForServer deletes all deployments associated with a given Server.
	DeleteForServer(ctx context.Context, srv *Server) error
}

// InMemoryDeploymentRepository holds the deployments inside an in-memory map.
//
// Deployments are not persisted on disk and won't survive program restarts.
type InMemoryDeploymentRepository struct {
	deployments sync.Map
}

// NewInMemoryDeploymentRepository returns a new InMemoryDeploymentRepository instance.
func NewInMemoryDeploymentRepository() *InMemoryDeploymentRepository {
	return &InMemoryDeploymentRepository{}
}

// Find a Deployment given a DeploymentID.
func (repo *InMemoryDeploymentRepository) Find(ctx context.Context, id DeploymentID) (*Deployment, error) {
	req, ok := repo.deployments.Load(id)
	if !ok {
		return nil, ErrDeploymentNotFound
	}

	return req.(*Deployment), nil
}

// Create a new Deployment.
func (repo *InMemoryDeploymentRepository) Create(ctx context.Context, deployment *Deployment) error {
	_, loaded := repo.deployments.LoadOrStore(deployment.ID, deployment)
	if loaded {
		return ErrDeploymentAlreadyExists
	}

	return nil
}

// Update an existing Deployment.
func (repo *InMemoryDeploymentRepository) Update(ctx context.Context, deployment *Deployment) error {
	repo.deployments.Store(deployment.ID, deployment)

	return nil
}

// DeleteForServer deletes all deployments associated with a given Server.
func (repo *InMemoryDeploymentRepository) DeleteForServer(ctx context.Context, srv *Server) error {
	repo.deployments.Range(func(k, v interface{}) bool {
		deployment := v.(*Deployment)
		if deployment.ServerID.String() != srv.ID.String() {
			return true
		}

		repo.deployments.Delete(k)

		return true
	})

	return nil
}
