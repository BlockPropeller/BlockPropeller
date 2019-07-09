package terraform_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"chainup.dev/chainup/terraform"
	"chainup.dev/lib/test"
)

func TestWorkspaceAllocatesWorkingDir(t *testing.T) {
	w, err := terraform.NewWorkspace()

	test.CheckErr(t, "NewWorkspace()", err)
	defer test.Close(t, w)

	workDir := w.WorkDir()
	if workDir == "" {
		t.Errorf("Workspace.WorkDir(): expected non-empty work dir")
		return
	}

	err = ioutil.WriteFile(filepath.Join(workDir, "test.txt"), []byte("Hello World!"), 0644)
	test.CheckErr(t, "Workspace.WorkDir(): expected work dir to be writable", err)
}

func TestWorkspaceHasUniqueWorkingDir(t *testing.T) {
	w1, err := terraform.NewWorkspace()
	test.CheckErr(t, "w1 := NewWorkspace()", err)
	defer test.Close(t, w1)

	w2, err := terraform.NewWorkspace()
	test.CheckErr(t, "w2 := NewWorkspace()", err)
	defer test.Close(t, w2)

	if w1.WorkDir() == w2.WorkDir() {
		t.Errorf("expected each workspace to have a unique work dir: got %s", w1.WorkDir())
		return
	}
}

func TestWorkspaceCleanupOnClose(t *testing.T) {
	w, err := terraform.NewWorkspace()
	test.CheckErr(t, "NewWorkspace()", err)

	test.Close(t, w)

	stat, err := os.Stat(w.WorkDir())
	if !os.IsNotExist(err) {
		t.Errorf("expected workspace work dir to not exist, got stat %v, err %v", stat, err)
		return
	}
}