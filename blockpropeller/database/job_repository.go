package database

import (
	"context"

	"blockpropeller.dev/blockpropeller/account"
	"blockpropeller.dev/blockpropeller/provision"
	"github.com/pkg/errors"
)

// JobRepository is a databased backed implementation of a provision.JobRepository.
type JobRepository struct {
	db *DB
}

// NewJobRepository returns a new JobRepository instance.
func NewJobRepository(db *DB) *JobRepository {
	return &JobRepository{db: db}
}

// Find a Job given a JobID.
func (repo *JobRepository) Find(ctx context.Context, id provision.JobID) (*provision.Job, error) {
	var job provision.Job
	err := repo.db.Model(ctx, &job).
		Preload("ProviderSettings").
		Where("id = ?", id).
		First(&job).
		Error
	if err != nil {
		return nil, errors.Wrap(err, "find job by ID")
	}

	return &job, nil
}

// List all jobs.
func (repo *JobRepository) List(ctx context.Context, accountID account.ID) ([]*provision.Job, error) {
	var jobs []*provision.Job
	err := repo.db.Model(ctx, &jobs).
		Where("account_id = ?", accountID).
		Find(&jobs).Error
	if err != nil {
		return nil, errors.Wrap(err, "find jobs")
	}

	return jobs, nil
}

// Create a new Job.
func (repo *JobRepository) Create(ctx context.Context, job *provision.Job) error {
	err := repo.db.Model(ctx, job).Create(job).Error
	if err != nil {
		return errors.Wrap(err, "create job")
	}

	return nil
}

// Update an existing Job.
func (repo *JobRepository) Update(ctx context.Context, job *provision.Job) error {
	err := repo.db.Model(ctx, job).Save(job).Error
	if err != nil {
		return errors.Wrap(err, "update job")
	}

	return nil
}
