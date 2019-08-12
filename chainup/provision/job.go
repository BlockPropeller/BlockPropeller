package provision

import (
	"sync"
	"time"

	"chainup.dev/chainup/infrastructure"
	"chainup.dev/chainup/statemachine"
	"chainup.dev/chainup/terraform"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
)

var (
	// ErrJobNotFound is returned when a JobRepository does not find a job to return.
	ErrJobNotFound = errors.New("job not found")
	// ErrJobAlreadyExists is returned when a Job creation is attempted with an existing JobID.
	ErrJobAlreadyExists = errors.New("job already exists")
)

// JobID is a unique server identifier.
type JobID string

// NewJobID returns a new unique JobID.
func NewJobID() JobID {
	return JobID(uuid.NewV4().String())
}

// String satisfies the Stringer interface.
func (id JobID) String() string {
	return string(id)
}

// Job represents a single provisioning request for the lifetime of the provisioning process.
//
// The provisioning job contains all the necessary information required for creating new infrastructure,
// as well as the specification for the servers and services needed to be running.
//
// Once complete, a Job serves only for record keeping, and is not concerned with any other
// actions on the created entities.
type Job struct {
	ID JobID `json:"id"`

	statemachine.Resource

	ProviderSettingsID infrastructure.ProviderSettingsID `json:"-"`
	ProviderSettings   *infrastructure.ProviderSettings  `json:"provider_settings"`

	ServerID infrastructure.ServerID `json:"-"`
	Server   *infrastructure.Server  `json:"server"`

	DeploymentID infrastructure.DeploymentID `json:"-"`
	Deployment   *infrastructure.Deployment  `json:"deployment"`

	WorkspaceSnapshot *terraform.WorkspaceSnapshot

	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	FinishedAt time.Time `json:"finished_at,omitempty"`
}

// NewJob returns a new Job instance.
func NewJob(provider *infrastructure.ProviderSettings, server *infrastructure.Server, deployment *infrastructure.Deployment) *Job {
	return &Job{
		ID: NewJobID(),

		Resource: statemachine.NewResource(StateCreated),

		ProviderSettingsID: provider.ID,
		ProviderSettings:   provider,

		ServerID: server.ID,
		Server:   server,

		DeploymentID: deployment.ID,
		Deployment:   deployment,

		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// JobBuilder allows for fluent job definition.
type JobBuilder struct {
	provider      *infrastructure.ProviderSettings
	server        *infrastructure.Server
	serverBuilder *infrastructure.ServerBuilder
	deployment    *infrastructure.Deployment
}

// NewJobBuilder returns a new JobBuilder instance.
func NewJobBuilder() *JobBuilder {
	return &JobBuilder{}
}

// Server specification that should be provisioned
func (b *JobBuilder) Server(server *infrastructure.Server) *JobBuilder {
	b.server = server

	return b
}

// Provider which is to be used to provision new infrastructure.
func (b *JobBuilder) Provider(provider *infrastructure.ProviderSettings) *JobBuilder {
	b.provider = provider

	return b
}

// Deployment which is to be provisioned on new infrastructure.
func (b *JobBuilder) Deployment(deployment *infrastructure.Deployment) *JobBuilder {
	b.deployment = deployment

	return b
}

// Build constructs a Job instance along with a Server specification.
func (b *JobBuilder) Build() (*Job, error) {
	if b.provider == nil {
		return nil, errors.New("missing provider configuration")
	}
	if b.server == nil {
		return nil, errors.New("missing server configuration")
	}
	if b.deployment == nil {
		return nil, errors.New("missing deployment configuration")
	}

	return NewJob(b.provider, b.server, b.deployment), nil
}

// JobRepository defines an interface for storing and retrieving provisioning jobs.
type JobRepository interface {
	// Find a Job given a JobID.
	Find(id JobID) (*Job, error)

	// Create a new Job.
	Create(job *Job) error

	// Update an existing Job.
	Update(job *Job) error
}

// InMemoryJobRepository holds the jobs inside an in-memory map.
//
// Jobs are not persisted on disk and won't survive program restarts.
type InMemoryJobRepository struct {
	jobs sync.Map
}

// NewInMemoryJobRepository returns a new InMemoryJobRepository instance.
func NewInMemoryJobRepository() *InMemoryJobRepository {
	return &InMemoryJobRepository{}
}

// Find a Job given a JobID.
func (repo *InMemoryJobRepository) Find(id JobID) (*Job, error) {
	req, ok := repo.jobs.Load(id)
	if !ok {
		return nil, ErrJobNotFound
	}

	return req.(*Job), nil
}

// Create a new Job.
func (repo *InMemoryJobRepository) Create(job *Job) error {
	_, loaded := repo.jobs.LoadOrStore(job.ID, job)
	if loaded {
		return ErrJobAlreadyExists
	}

	return nil
}

// Update an existing Job.
func (repo *InMemoryJobRepository) Update(job *Job) error {
	repo.jobs.Store(job.ID, job)

	return nil
}
