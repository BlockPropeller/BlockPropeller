package ansible

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"blockpropeller.dev/blockpropeller/infrastructure"
	"blockpropeller.dev/lib/log"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

var (
	// ErrServerUnreachable indicates that the server we are trying to access is not reachable.
	ErrServerUnreachable = errors.New("server unreachable")
)

// Ansible is a wrapper around an ansible-playbook command line utility
// exposing the ability to provision arbitrary servers.
type Ansible struct {
	path string

	playbooksDir string
	keysDir      string
}

// ConfigureAnsible returns a configured Terraform instance.
func ConfigureAnsible(cfg *Config) *Ansible {
	return New(cfg.Path, cfg.PlaybooksDir, cfg.KeysDir)
}

// New returns a new Terraform instance.
func New(path string, playbooksDir string, keysDir string) *Ansible {
	return &Ansible{
		path:         path,
		playbooksDir: playbooksDir,
		keysDir:      keysDir,
	}
}

// RunPlaybook executes the playbook on a specified Server
// and applying the provided deployment configuration.
func (ans *Ansible) RunPlaybook(srv *infrastructure.Server, deployment *infrastructure.Deployment) error {
	keyPath, err := ans.setupSSHKey(srv.SSHKey)
	if err != nil {
		return errors.Wrap(err, "setup ssh key")
	}
	defer ans.cleanupSSHKey(keyPath)

	var extraVars []string
	for key, value := range deployment.Configuration.MarshalMap() {
		extraVars = append(extraVars, key+"="+value)
	}

	out, err := ans.exec(
		ans.playbooksDir,
		"--inventory", srv.IPAddress.String()+",",
		"--key-file", keyPath,
		"--extra-vars", strings.Join(extraVars, " "),
		"site.yaml",
	)
	log.Debug("run ansible-playbook", log.Fields{
		"stdout": string(out),
	})
	if err != nil {
		return errors.Wrap(err, "execute ansible playbook")
	}

	return nil
}

// Version returns the version of the underlying binary.
//
// This method can be used as a health check whether the
// binary is correctly configured.
func (ans *Ansible) Version() (string, error) {
	out, err := ans.exec(ans.playbooksDir, "--version")
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
		if execErr.ExitCode() == 4 {
			return output, ErrServerUnreachable
		}

		return output, errors.Wrapf(execErr,
			"execution error for [ansible-playbook %s]: %s",
			strings.Join(args, " "),
			string(execErr.Stderr),
		)
	}
	if err != nil {
		return output, errors.Wrapf(err, "exec command for [ansible-playbook %s]", strings.Join(args, " "))
	}

	return output, nil
}

func (ans *Ansible) setupSSHKey(sshKey *infrastructure.SSHKey) (string, error) {
	keyFile := filepath.Join(ans.keysDir, uuid.NewV4().String())
	err := os.MkdirAll(filepath.Dir(keyFile), 0755)
	if err != nil {
		return "", errors.Wrap(err, "create ansible keys dir")
	}

	err = ioutil.WriteFile(keyFile, []byte(sshKey.EncodedPrivateKey()), 0400)
	if err != nil {
		log.ErrorErr(err, "failed writing private key", log.Fields{
			"path": keyFile,
		})
		return "", errors.Wrap(err, "write private key")
	}

	return keyFile, nil
}

func (ans *Ansible) cleanupSSHKey(keyFile string) {
	err := os.Remove(keyFile)
	if err != nil && !os.IsNotExist(err) {
		log.ErrorErr(err, "failed cleaning up ansible ssh key", log.Fields{
			"path": keyFile,
		})
	}
}
