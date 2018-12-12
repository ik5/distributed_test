package actions

import (
	"os"
	"path/filepath"
	"testing"
)

func TestProcExists(t *testing.T) {
	name := filepath.Base(os.Args[0])
	pid, err := ProcExists(name)
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	if !pid {
		t.Errorf("Unable to find '%s'", name)
	}
}

func TestIsFileContains(t *testing.T) {
	contains, err := IsFileContains("./actions_test.go", "actions")
	if err != nil {
		t.Errorf("Error: %s", err)
	}

	if !contains {
		t.Errorf("'actions' was not found")
	}
}

func TestIsFileExists(t *testing.T) {
	exists, err := IsFileExists("./actions_test.go")
	if err != nil {
		t.Errorf("Error: %s", err)
	}

	if !exists {
		t.Errorf("File was not found")
	}
}
