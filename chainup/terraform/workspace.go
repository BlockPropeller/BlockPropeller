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

// WorkspaceSnapshot creates a serializable data structure form
// an active Workspace, so we can persist it and recreate the Workspace
// at a later date if necessary.
//
// This allows us to modify and delete infrastructure initially deployed with Terraform.
type WorkspaceSnapshot struct {
	WorkspacePath        string
	TerraformDefinitions string
	TerraformPlan        string
	TerraformState       string
}

// Workspace handles laying out and a set of `Resource`s
// in a format compatible with the Terraform command line
// interface.
//
// Workspace should be initialized for each new execution,
// and closed afterwards in order to cleanup the filesystem.
type Workspace struct {
	workDir string

	readOnly  bool
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

// RestoreWorkspace restores the serialized workspace inside a WorkspaceSnapshot.
//
// Restored Workspaces are readonly and will panic if add methods are called.
func RestoreWorkspace(snap *WorkspaceSnapshot) (*Workspace, error) {
	err := os.MkdirAll(snap.WorkspacePath, 0755)
	if err != nil {
		return nil, errors.Wrap(err, "restore workspace work dir")
	}

	err = ioutil.WriteFile(filepath.Join(snap.WorkspacePath, "main.tf"), []byte(snap.TerraformDefinitions), 0655)
	if err != nil {
		return nil, errors.Wrap(err, "restore terraform definitions")
	}

	err = ioutil.WriteFile(filepath.Join(snap.WorkspacePath, "tfplan"), []byte(snap.TerraformPlan), 0655)
	if err != nil {
		return nil, errors.Wrap(err, "restore terraform plan")
	}

	err = ioutil.WriteFile(filepath.Join(snap.WorkspacePath, "terraform.tfstate"), []byte(snap.TerraformState), 0655)
	if err != nil {
		return nil, errors.Wrap(err, "restore terraform state")
	}

	return &Workspace{
		workDir:  snap.WorkspacePath,
		readOnly: true,
	}, nil
}

// Add adds provided items to the list of items that should be flushed to the workspace.
//
// Add method does not automatically flush the items. A separate Flush method should be called
// before using the workspace.
func (w *Workspace) Add(items ...Renderer) {
	if w.readOnly {
		panic("workspace is readonly")
	}

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
	if w.readOnly {
		panic("workspace is readonly")
	}

	if len(resources) == 0 {
		return
	}

	w.flushed = false
	w.resources = append(w.resources, resources...)
}

// Flush persists all items in a Terraform file in order to be executed by Terraform.
func (w *Workspace) Flush() error {
	if w.readOnly {
		panic("workspace is readonly")
	}

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

// Snapshot converts a Workspace into a serializable structure.
func (w *Workspace) Snapshot() (*WorkspaceSnapshot, error) {
	snap := &WorkspaceSnapshot{
		WorkspacePath: w.WorkDir(),
	}

	definitions, err := ioutil.ReadFile(filepath.Join(w.workDir, "main.tf"))
	if err != nil && !os.IsNotExist(err) {
		return nil, errors.Wrap(err, "read terraform definitions")
	}

	snap.TerraformDefinitions = string(definitions)

	plan, err := ioutil.ReadFile(filepath.Join(w.workDir, "tfplan"))
	if err != nil && !os.IsNotExist(err) {
		return nil, errors.Wrap(err, "read terraform plan")
	}

	snap.TerraformPlan = string(plan)

	state, err := ioutil.ReadFile(filepath.Join(w.workDir, "terraform.tfstate"))
	if err != nil && !os.IsNotExist(err) {
		return nil, errors.Wrap(err, "read terraform state")
	}

	snap.TerraformState = string(state)

	return snap, nil
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
