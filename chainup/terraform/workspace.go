package terraform

import (
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
)

// Workspace handles laying out and a set of `Resource`s
// in a format compatible with the Terraform command line
// interface.
//
// Workspace should be initialized for each new execution,
// and closed afterwards in order to cleanup the filesystem.
type Workspace struct {
	workDir string

	flushed bool
}

// NewWorkspace returns a new Workspace instance.
//
// A new temporary directory is allocated for each workspace.
func NewWorkspace() (*Workspace, error) {
	workDir, err := ioutil.TempDir(os.TempDir(), "tf-workspace-")
	if err != nil {
		return nil, errors.Wrap(err, "create temp dir")
	}

	return &Workspace{
		workDir: workDir,

		flushed: true,
	}, nil
}

// WorkDir returns the absolute filesystem path for the current Workspace.
func (w *Workspace) WorkDir() string {
	return w.workDir
}

// Close cleans up any files that were created
// for the lifetime of the Workspace.
func (w *Workspace) Close() error {
	err := os.RemoveAll(w.workDir)
	if err != nil {
		return errors.Wrap(err, "cleanup workspace directory")
	}

	return nil
}
