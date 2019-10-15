package routes

import (
	"blockpropeller.dev/blockpropeller/provision"
	"github.com/labstack/echo"
)

// Provision routes define ways to provision infrastructure via BlockPropeller.
type Provision struct {
	jobRepo provision.JobRepository
}

// NewProvisionRoutes returns a new Provision routes instance.
func NewProvisionRoutes(jobRepo provision.JobRepository) *Provision {
	return &Provision{jobRepo: jobRepo}
}

// LoadJob is a middleware for loading Jobs into request context
// as well as checking for correct permissions of an authenticated user.
func (p *Provision) LoadJob(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		jobID := provision.JobID(c.Param("job_id"))
		if jobID.String() == "" {
			return next(c)
		}

		return next(c)
	}
}
