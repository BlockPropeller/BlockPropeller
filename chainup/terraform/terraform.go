package terraform

import (
	"os/exec"
	"strings"

	"github.com/pkg/errors"
)

// Terraform is a wrapper around a terraform command line utility
// exposing the ability to plan and provision Terraform resources.
type Terraform struct {
	path string
}

// New returns a new Terraform instance.
func New(path string) *Terraform {
	return &Terraform{
		path: path,
	}
}

// Version returns the version of the underlying binary.
//
// This method can be used as a health check whether the
// binary is correctly configured.
func (tf *Terraform) Version() (string, error) {
	ver, err := tf.exec("version")
	if err != nil {
		return "", errors.Wrap(err, "get terraform version")
	}

	return string(ver), nil
}

// exec wraps the interaction with the underlying binary.
func (tf *Terraform) exec(args ...string) ([]byte, error) {
	output, err := exec.Command(tf.path, args...).Output()
	if execErr, ok := err.(*exec.ExitError); ok {
		return nil, errors.Wrapf(execErr,
			"execution error for [terraform %s]: %s",
			strings.Join(args, " "),
			string(execErr.Stderr),
		)
	}
	if err != nil {
		return nil, errors.Wrapf(err, "exec command for [terraform %s]", strings.Join(args, " "))
	}

	return output, nil
}
