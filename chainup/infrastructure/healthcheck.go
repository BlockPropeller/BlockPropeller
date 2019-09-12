package infrastructure

import (
	"net/http"

	"chainup.dev/lib/log"
	"github.com/pkg/errors"
)

// HealthCheck is used to keep track of the health of individual deployments,
// signaling whether the deployment is responding correctly.
type HealthCheck interface {
	Health() error
}

// CheckHealth checks a Deployment health based on a configured HealthCheck.
func CheckHealth(srv *Server, deployment *Deployment) error {
	spec, err := getDeploymentSpec(deployment.Type)
	if err != nil {
		return err
	}

	check, err := spec.HealthCheck(srv, deployment)
	if err != nil {
		return errors.Wrap(err, "get deployment health check")
	}

	err = check.Health()
	if err != nil {
		return errors.Wrap(err, "check health")
	}

	return nil
}

// HTTPHealthCheck sends a simple HTTP request to determine the health of the deployment.
type HTTPHealthCheck struct {
	client *http.Client

	Method string
	URL    string

	ExpectedStatusCode int
}

// NewHTTPHealthCheck returns a new HTTPHealthCheck instance.
func NewHTTPHealthCheck(method string, url string, status int) *HTTPHealthCheck {
	return &HTTPHealthCheck{
		client: &http.Client{},

		Method:             method,
		URL:                url,
		ExpectedStatusCode: status,
	}
}

// Health conforms to the HealthCheck interface.
func (hc *HTTPHealthCheck) Health() error {
	req, err := http.NewRequest(hc.Method, hc.URL, nil)
	if err != nil {
		return errors.Wrap(err, "build health check request")
	}

	resp, err := hc.client.Do(req)
	if err != nil {
		return errors.Wrap(err, "send health check request")
	}
	defer log.Closer(resp.Body)

	if resp.StatusCode != hc.ExpectedStatusCode {
		return errors.Errorf("unexpected status: %d", resp.StatusCode)
	}

	return nil
}
