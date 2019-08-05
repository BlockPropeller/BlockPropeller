package terraform

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"

	"chainup.dev/chainup/terraform/resource"
	"github.com/pkg/errors"
)

// Renderer is any struct that can render itself into a Terraform compatible syntax.
//
// Renderers are used as individual units that need to be deployed inside a single workspace.
type Renderer interface {
	Render() string
}

// Workspace handles laying out and a set of `Resource`s
// in a format compatible with the Terraform command line
// interface.
//
// Workspace should be initialized for each new execution,
// and closed afterwards in order to cleanup the filesystem.
type Workspace struct {
	workDir string

	flushed   bool
	items     []Renderer
	resources []resource.Resource
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

// Add adds provided items to the list of items that should be flushed to the workspace.
//
// Add method does not automatically flush the items. A separate Flush method should be called
// before using the workspace.
func (w *Workspace) Add(items ...Renderer) {
	if len(items) == 0 {
		return
	}

	w.flushed = false
	w.items = append(w.items, items...)
}

// AddResource acts in the same way as the Add method, only difference being that it
// accepts a variadic number of Resources instead of items.
//
// @TODO: Refactor Workspace and related interfaces to require only one method of adding contents to a workspace.
func (w *Workspace) AddResource(resources ...resource.Resource) {
	if len(resources) == 0 {
		return
	}

	w.flushed = false
	w.resources = append(w.resources, resources...)
}

// Flush persists all items in a Terraform file in order to be executed by Terraform.
func (w *Workspace) Flush() error {
	if w.flushed {
		return nil
	}

	var buf bytes.Buffer

	for _, item := range w.items {
		buf.WriteString(item.Render())
		buf.WriteRune('\n')
	}

	for _, res := range w.resources {
		buf.WriteString(resource.Render(res))
		buf.WriteRune('\n')
	}

	err := ioutil.WriteFile(filepath.Join(w.workDir, "main.tf"), buf.Bytes(), 0644)
	if err != nil {
		return errors.Wrap(err, "write items to disk")
	}

	return nil
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
