package files

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type directoryTree struct {
	// Name left empty on initialisation, filled when files are created
	Name string
	// if directory is empty, this is a file, if not a dir
	// contains
	DirectoryTree []directoryTree

	Result *os.File
}

func (d *directoryTree) IsFile() bool {
	if len(d.DirectoryTree) == 0 {
		return true
	}
	return false
}

func generateTestEnv(plan directoryTree) error {
	// Create top level dir
	path, err := os.MkdirTemp("", plan.Name)
	if err != nil {
		return err
	}

	plan.Name = path

	return generateTestEnvRec(plan)

}

func generateTestEnvRec(plan directoryTree) error {
	for idx := range plan.DirectoryTree {
		if plan.DirectoryTree[idx].IsFile() {
			file, err := os.CreateTemp(plan.Name, "")
			if err != nil {
				return err
			}

			plan.DirectoryTree[idx].Result = file
			// Name not set, inherent in file pointer
		} else {
			// is directory
			path, err := os.MkdirTemp(plan.Name, "")
			if err != nil {
				return err
			}

			plan.DirectoryTree[idx].Name = path
		}
	}
	return nil
}

func CleanupEnv(path directoryTree) {
	os.RemoveAll(path.Name)
}

func TestFileListing(t *testing.T) {

	// Empty model, just a directory
	testEnv0Model := directoryTree{}
	err := generateTestEnv(testEnv0Model)
	assert.NoError(t, err)
	defer CleanupEnv(testEnv0Model)

}
