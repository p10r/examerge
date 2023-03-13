package main

import (
	"testing"
)

func TestDirectories(t *testing.T) {
	t.Run("creates a dir for generated output", func(t *testing.T) {
		tmpDir := TestSetupTestEnvironment(t)
		defer TestRemove(tmpDir, t)

		Workflow(tmpDir)

		result, err := TestExists(tmpDir + outputDir)
		if err != nil || result == false {
			t.Errorf("Expected %s/generated, but doesn't exist", tmpDir)
		}
	})

	t.Run("copies over existing files into /generated", func(t *testing.T) {
		tmpDir := TestSetupTestEnvironment(t)
		defer TestRemove(tmpDir, t)

		CreateOutputDirIn(tmpDir)
		Copy(tmpDir)

		dirs := TestListSubDirTree(tmpDir+"/generated", t)
		for _, dir := range dirs {
			t.Logf("Created %s/%s", tmpDir, dir)
		}

	})
}

func Copy(path string) {

}
