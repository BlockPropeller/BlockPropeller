package infrastructure

import (
	"time"

	uuid "github.com/satori/go.uuid"
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
	ServerID ServerID     `json:"-"`

	Type          DeploymentType    `json:"type"`
	Configuration map[string]string `json:"configuration"`

	State DeploymentState `json:"state"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewDeployment returns a new Deployment instance.
func NewDeployment(typ DeploymentType, config map[string]string) *Deployment {
	return &Deployment{
		ID: NewDeploymentID(),

		Type:          typ,
		Configuration: config,

		State: DeploymentStateRequested,

		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
