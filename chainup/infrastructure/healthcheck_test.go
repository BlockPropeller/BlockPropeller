package infrastructure_test

import (
	"testing"

	"chainup.dev/chainup/infrastructure"
	"chainup.dev/lib/test"
	"github.com/pkg/errors"
)

func TestOkHTTPHealthCheck(t *testing.T) {
	test.Integration(t)

	okCheck := infrastructure.NewHTTPHealthCheck("GET", "https://example.com", 200)

	err := okCheck.Health()
	test.CheckErr(t, "expected health check to pass", err)
}

func TestFailingHTTPHealthCheck(t *testing.T) {
	test.Integration(t)

	failingCheck := infrastructure.NewHTTPHealthCheck("GET", "https://example.com/404", 200)

	err := failingCheck.Health()
	test.CheckErrExists(t, "expected health check to fail", err)
	test.AssertStringsEqual(t, "expected correct message",
		errors.Cause(err).Error(), "unexpected status: 404")
}
