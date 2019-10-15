package terraform

import (
	"os"
	"os/exec"
	"strings"

	"blockpropeller.dev/lib/log"
	"github.com/pkg/errors"
)

// Terraform is a wrapper around a terraform command line utility
// exposing the ability to plan and provision Terraform resources.
type Terraform struct {
	path string
}

// ConfigureTerraform returns a configured Terraform instance.
func ConfigureTerraform(cfg *Config) *Terraform {
	return New(cfg.Path)
}

// New returns a new Terraform instance.
func New(path string) *Terraform {
	return &Terraform{
		path: path,
	}
}

// Init executes terraform init in the provided workspace.
//
// terraform init must be called before terraform plan or apply.
func (tf *Terraform) Init(workspace *Workspace) error {
	out, err := tf.exec(workspace.WorkDir(), "init", "-no-color", "-input=false")
	log.Debug("terraform init", log.Fields{
		"stdout": string(out),
	})
	if err != nil {
		return errors.Wrap(err, "execute terraform init")
	}

	return nil
}

// Plan connects to the configured provider and creates a plan
// for infrastructure that needs to be provisioned on the provider.
func (tf *Terraform) Plan(workspace *Workspace) error {
	out, err := tf.exec(workspace.WorkDir(), "plan", "-out=tfplan", "-no-color", "-input=false")
	log.Debug("terraform plan", log.Fields{
		"stdout": string(out),
	})
	if err != nil {
		return errors.Wrap(err, "execute terraform plan")
	}

	return nil
}

// Apply executes the plan previously created by the Plan method.
//
// Plan method *must* be called before apply, otherwise apply will fail.
func (tf *Terraform) Apply(workspace *Workspace) error {
	out, err := tf.exec(workspace.WorkDir(), "apply", "-no-color", "-input=false", "tfplan")
	log.Debug("terraform apply", log.Fields{
		"stdout": string(out),
	})
	if err != nil {
		return errors.Wrap(err, "execute terraform apply")
	}

	return nil
}

// Output returns the value of a defined output by a given name.
//
// Terraform apply must have been called beforehand in order for the output command to work.
func (tf *Terraform) Output(workspace *Workspace, name string) (string, error) {
	out, err := tf.exec(workspace.WorkDir(), "output", "-no-color", name)
	log.Debug("terraform output", log.Fields{
		"stdout": string(out),
	})
	if err != nil {
		return "", errors.Wrap(err, "execute terraform output")
	}

	return strings.Trim(string(out), " \n"), nil
}

// Destroy destroys all resources provisioned on a configured provider.
func (tf *Terraform) Destroy(workspace *Workspace) error {
	out, err := tf.exec(workspace.WorkDir(), "destroy", "-no-color", "-auto-approve")
	log.Debug("terraform destroy", log.Fields{
		"stdout": string(out),
	})
	if err != nil {
		return errors.Wrap(err, "execute terraform destroy")
	}

	return nil
}

// Version returns the version of the underlying binary.
//
// This method can be used as a health check whether the
// binary is correctly configured.
func (tf *Terraform) Version() (string, error) {
	out, err := tf.exec("", "version")
	if err != nil {
		return "", errors.Wrap(err, "get terraform version")
	}

	outLines := strings.Split(string(out), "\n")
	versionLine := outLines[0]
	versionParts := strings.Split(versionLine, " ")
	if len(versionParts) != 2 {
		return versionLine, nil
	}

	return strings.TrimPrefix(versionParts[1], "v"), nil
}

// exec wraps the interaction with the underlying binary.
func (tf *Terraform) exec(dir string, args ...string) ([]byte, error) {
	cmd := exec.Command(tf.path, args...)
	cmd.Env = append(os.Environ(), "TF_IN_AUTOMATION=true")
	if dir != "" {
		cmd.Dir = dir
	}

	output, err := cmd.Output()
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
