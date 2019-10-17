package provision

import (
	"context"
	"sync"
	"time"

	"blockpropeller.dev/lib/log"
)

// WorkerPoolConfig holds configuration for a provisioning worker pool.
type WorkerPoolConfig struct {
	WorkerCount int `yaml:"worker_count"`
}

// Validate satisfies the config.Config interface.
func (cfg *WorkerPoolConfig) Validate() error {
	if cfg.WorkerCount == 0 {
		cfg.WorkerCount = 20
	}

	return nil
}

// WorkerPool is responsible for concurrently processing
// provisioning jobs.
type WorkerPool struct {
	workerCount int

	jobCh      chan JobID
	activeJobs sync.Map

	jobRepo     JobRepository
	provisioner *Provisioner
}

// NewWorkerPool returns a new WorkerPool instance.
func NewWorkerPool(cfg *WorkerPoolConfig, jobRepo JobRepository, provisioner *Provisioner) *WorkerPool {
	return &WorkerPool{
		workerCount: cfg.WorkerCount,

		jobCh: make(chan JobID),

		jobRepo:     jobRepo,
		provisioner: provisioner,
	}
}

// Start the WorkerPool and all its workers, and wait for the Context to finish.
func (wp *WorkerPool) Start(ctx context.Context) {
	go wp.producerLoop(ctx)

	wp.startWorkers(ctx)
}

func (wp *WorkerPool) producerLoop(ctx context.Context) {
	for ctx.Err() == nil {
		jobs, err := wp.jobRepo.FindIncomplete(ctx, wp.getActiveJobs()...)
		if err != nil {
			log.ErrorErr(err, "failed finding incomplete jobs", log.Fields{
				"sleeping": 10,
			})
			wp.sleep(ctx, 10*time.Second)
			continue
		}
		if len(jobs) == 0 {
			// No jobs to schedule, sleeping.
			wp.sleep(ctx, 10*time.Second)
			continue
		}

		for _, job := range jobs {
			log.Info("scheduling job", log.Fields{
				"job_id": job.ID,
			})
			wp.addActiveJob(job.ID)
			wp.jobCh <- job.ID
		}
	}
}

func (wp *WorkerPool) startWorkers(ctx context.Context) {
	var wg sync.WaitGroup
	for i := 0; i < wp.workerCount; i++ {
		wg.Add(1)
		go wp.startWorker(ctx)
	}

	wg.Wait()
}

func (wp *WorkerPool) startWorker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case jobID := <-wp.jobCh:
			log.Info("starting job", log.Fields{
				"job_id": jobID,
			})
			err := wp.provisioner.Provision(ctx, jobID)
			if err != nil {
				// Provisioning process finished with an error.
				log.ErrorErr(err, "run provision job", log.Fields{
					"job_id": jobID,
				})
			}

			wp.removeActiveJob(jobID)
			log.Info("finished job", log.Fields{
				"job_id": jobID,
			})
		}
	}
}

func (wp *WorkerPool) getActiveJobs() []JobID {
	var jobIDs []JobID

	wp.activeJobs.Range(func(k, v interface{}) bool {
		jobIDs = append(jobIDs, k.(JobID))

		return true
	})

	if len(jobIDs) == 0 {
		return []JobID{}
	}

	return jobIDs
}

func (wp *WorkerPool) addActiveJob(id JobID) {
	wp.activeJobs.Store(id, true)
}

func (wp *WorkerPool) removeActiveJob(id JobID) {
	wp.activeJobs.Delete(id)
}

func (wp *WorkerPool) sleep(ctx context.Context, d time.Duration) {
	select {
	case <-ctx.Done():
	case <-time.After(d):
	}
}
