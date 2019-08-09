package ansible

import (
	"os/exec"
	"strings"

	"github.com/pkg/errors"
)

// Ansible is a wrapper around an ansible-playbook command line utility
// exposing the ability to provision arbitrary servers.
type Ansible struct {
	path string
}

// ConfigureAnsible returns a configured Terraform instance.
func ConfigureAnsible(cfg *Config) *Ansible {
	return New(cfg.Path)
}

// New returns a new Terraform instance.
func New(path string) *Ansible {
	return &Ansible{
		path: path,
	}
}

// Version returns the version of the underlying binary.
//
// This method can be used as a health check whether the
// binary is correctly configured.
func (ans *Ansible) Version() (string, error) {
	out, err := ans.exec("", "--version")
	if err != nil {
		return "", errors.Wrap(err, "get ansible version")
	}

	outLines := strings.Split(string(out), "\n")
	versionLine := outLines[0]
	versionParts := strings.Split(versionLine, " ")
	if len(versionParts) != 2 {
		return versionLine, nil
	}

	return versionParts[1], nil
}

// exec wraps the interaction with the underlying binary.
func (ans *Ansible) exec(dir string, args ...string) ([]byte, error) {
	cmd := exec.Command(ans.path, args...)
	if dir != "" {
		cmd.Dir = dir
	}

	output, err := cmd.Output()
	if execErr, ok := err.(*exec.ExitError); ok {
		return nil, errors.Wrapf(execErr,
			"execution error for [ansible-playbook %s]: %s",
			strings.Join(args, " "),
			string(execErr.Stderr),
		)
	}
	if err != nil {
		return nil, errors.Wrapf(err, "exec command for [ansible-playbook %s]", strings.Join(args, " "))
	}

	return output, nil
}
