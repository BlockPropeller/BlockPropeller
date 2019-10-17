package provision

import (
	"context"
	"sync"
	"time"

	"blockpropeller.dev/blockpropeller/account"
	"blockpropeller.dev/blockpropeller/infrastructure"
	"blockpropeller.dev/blockpropeller/statemachine"
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
	ID        JobID      `json:"id" gorm:"type:varchar(36) not null"`
	AccountID account.ID `json:"-" gorm:"type:varchar(36) not null references accounts(id)"`

	statemachine.Resource `gorm:"embedded"`

	ProviderSettingsID infrastructure.ProviderSettingsID `json:"-" gorm:"type:varchar(255) references provider_settings(id)"`
	ProviderSettings   *infrastructure.ProviderSettings  `json:"provider_settings"`

	ServerID infrastructure.ServerID `json:"-" gorm:"type:varchar(36) references servers(id)"`
	Server   *infrastructure.Server  `json:"server"`

	DeploymentID infrastructure.DeploymentID `json:"-" gorm:"type:varchar(36) references deployments(id)"`
	Deployment   *infrastructure.Deployment  `json:"deployment"`

	CreatedAt  time.Time  `json:"created_at" gorm:"type:datetime not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt  time.Time  `json:"updated_at" gorm:"type:datetime not null;default:CURRENT_TIMESTAMP"`
	FinishedAt *time.Time `json:"finished_at,omitempty" gorm:"type:datetime"`
}

// NewJob returns a new Job instance.
func NewJob(
	accountID account.ID,
	provider *infrastructure.ProviderSettings,
	server *infrastructure.Server,
	deployment *infrastructure.Deployment,
) *Job {
	return &Job{
		ID:        NewJobID(),
		AccountID: accountID,

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
	accountID     account.ID
	provider      *infrastructure.ProviderSettings
	server        *infrastructure.Server
	serverBuilder *infrastructure.ServerBuilder
	deployment    *infrastructure.Deployment
}

// NewJobBuilder returns a new JobBuilder instance.
func NewJobBuilder(accountID account.ID) *JobBuilder {
	return &JobBuilder{
		accountID: accountID,
	}
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

	b.server.AddDeployment(b.deployment)

	return NewJob(b.accountID, b.provider, b.server, b.deployment), nil
}

// JobRepository defines an interface for storing and retrieving provisioning jobs.
type JobRepository interface {
	// FindIncomplete Jobs with the exclusion of provided JobIDs.
	FindIncomplete(ctx context.Context, excl ...JobID) ([]*Job, error)

	// Find a Job given a JobID.
	Find(ctx context.Context, id JobID) (*Job, error)

	// List all jobs.
	List(ctx context.Context, accountID account.ID) ([]*Job, error)

	// Create a new Job.
	Create(ctx context.Context, job *Job) error

	// Update an existing Job.
	Update(ctx context.Context, job *Job) error
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

// FindIncomplete Jobs with the exclusion of provided JobIDs.
func (repo *InMemoryJobRepository) FindIncomplete(ctx context.Context, excl ...JobID) ([]*Job, error) {
	var jobs []*Job

	isExcluded := func(job *Job) bool {
		for _, jobID := range excl {
			if job.ID == jobID {
				return true
			}
		}

		return false
	}

	repo.jobs.Range(func(key, v interface{}) bool {
		job := v.(*Job)

		if isExcluded(job) {
			return true
		}

		jobs = append(jobs, job)

		return true
	})

	return jobs, nil
}

// Find a Job given a JobID.
func (repo *InMemoryJobRepository) Find(ctx context.Context, id JobID) (*Job, error) {
	req, ok := repo.jobs.Load(id)
	if !ok {
		return nil, ErrJobNotFound
	}

	return req.(*Job), nil
}

// List all jobs.
func (repo *InMemoryJobRepository) List(ctx context.Context, accountID account.ID) ([]*Job, error) {
	var jobs []*Job

	repo.jobs.Range(func(key, v interface{}) bool {
		job := v.(*Job)
		if job.AccountID != accountID {
			return true
		}

		jobs = append(jobs, job)

		return true
	})

	return jobs, nil
}

// Create a new Job.
func (repo *InMemoryJobRepository) Create(ctx context.Context, job *Job) error {
	_, loaded := repo.jobs.LoadOrStore(job.ID, job)
	if loaded {
		return ErrJobAlreadyExists
	}

	return nil
}

// Update an existing Job.
func (repo *InMemoryJobRepository) Update(ctx context.Context, job *Job) error {
	repo.jobs.Store(job.ID, job)

	return nil
}
