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
