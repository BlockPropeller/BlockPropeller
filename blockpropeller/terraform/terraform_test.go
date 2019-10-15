package terraform_test

import (
	"testing"

	"blockpropeller.dev/blockpropeller/terraform"
	"blockpropeller.dev/lib/test"
	"github.com/blang/semver"
)

func TestTerraformIsExecutable(t *testing.T) {
	tf := terraform.New("/usr/local/bin/terraform")

	version, err := tf.Version()
	test.CheckErr(t, "get terraform version", err)

	_, err = semver.New(version)
	test.CheckErr(t, "invalid terraform version format", err)
}

func TestPlanSimpleResource(t *testing.T) {
	//@TODO: Test terraform planning.
}

func TestApplySimpleResource(t *testing.T) {
	//@TODO: Test terraform executing a plan.
}

func TestCleanupAfterSelf(t *testing.T) {
	//@TODO: Test terraform cleaning up after itself.
}
