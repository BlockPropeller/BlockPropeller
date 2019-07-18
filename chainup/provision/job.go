package provision

import (
	"sync"
	"time"

	"chainup.dev/chainup/statemachine"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

var (
	// ErrJobNotFound is returned when a JobRepository does not find a job to return.
	ErrJobNotFound = errors.New("job not found")
	// ErrJobAlreadyExists is returned when a Job creation is attempted with an existing JobID.
	ErrJobAlreadyExists = errors.New("job already exists")
)

// JobID is a unique job identifier.
type JobID string

// NewJobID returns a new unique JobID.
func NewJobID() JobID {
	return JobID(uuid.NewV4().String())
}

// String satisfies the Stringer interface.
func (id JobID) String() string {
	return string(id)
}

// JobBuilder helps construct provisioning jobs by providing
// a fluent interface for configuring job details.
type JobBuilder struct {
}

// NewJobBuilder starts the process of building a job.
func NewJobBuilder() *JobBuilder {
	return &JobBuilder{}
}

// Build assembles all the job configuration into a single job object.
func (b *JobBuilder) Build() *Job {
	return NewJob()
}

// Job holds all the configuration values for a single provisioning job.
type Job struct {
	ID JobID `json:"id"`

	statemachine.Resource

	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	FinishedAt time.Time `json:"finished_at,omitempty"`
}

// NewJob allows you to construct a provision job in a single line.
//
// If you need a fluent interface for constructing the Job, you can use the JobBuilder.
func NewJob() *Job {
	return &Job{
		ID: NewJobID(),

		Resource: statemachine.NewResource(StateCreated),

		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
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
