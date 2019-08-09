package ansible_test

import (
	"testing"

	"chainup.dev/chainup/ansible"
	"chainup.dev/lib/test"
	"github.com/blang/semver"
)

func TestAnsibleIsExecutable(t *testing.T) {
	ans := ansible.New("/usr/local/bin/ansible-playbook")

	version, err := ans.Version()
	test.CheckErr(t, "get ansible version", err)

	_, err = semver.New(version)
	test.CheckErr(t, "invalid ansible version format", err)
}
