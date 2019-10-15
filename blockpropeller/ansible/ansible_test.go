package ansible_test

import (
	"testing"

	"blockpropeller.dev/blockpropeller/ansible"
	"blockpropeller.dev/lib/test"
	"github.com/blang/semver"
)

func TestAnsibleIsExecutable(t *testing.T) {
	ans := ansible.New(
		"/usr/local/bin/ansible-playbook",
		"../../playbooks",
		"/tmp/blockpropeller/ansible/keys",
	)

	version, err := ans.Version()
	test.CheckErr(t, "get ansible version", err)

	_, err = semver.New(version)
	test.CheckErr(t, "invalid ansible version format", err)
}
