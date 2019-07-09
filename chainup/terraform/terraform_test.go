package terraform_test

import (
	"testing"

	"chainup.dev/chainup/terraform"
)

func TestTerraformIsExecutable(t *testing.T) {
	tf := terraform.New("/usr/local/bin/terraform")

	_, err := tf.Version()
	if err != nil {
		t.Errorf("failed checking terraform version: %s", err)
		return
	}
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
