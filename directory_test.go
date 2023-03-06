package main

import (
	"os"
	"reflect"
	"testing"
)

func TestDirectories(t *testing.T) {
	tmpDir := SetupTestEnvironment(t, 2)
	defer os.RemoveAll(tmpDir)

	t.Run("lists all sub-directories", func(t *testing.T) {
		want := []string{"student1", "student2"}

		got, err := ListDirectories(tmpDir)

		assertNoError(t, err)

		if !reflect.DeepEqual(fileNames(got), want) {
			t.Errorf("got %q want %q", fileNames(got), want)
		}
	})

	t.Run("creates a dir for generated output", func(t *testing.T) {
		Workflow(tmpDir)

		result, err := exists(tmpDir + "/generated")
		if err != nil || result == false {
			t.Errorf("Expected %s/generated, but doesn't exist", tmpDir)
		}
	})
}

func ListDirectories(path string) ([]os.DirEntry, error) {
	return os.ReadDir(path)
}

func assertNoError(t testing.TB, got error) {
	t.Helper()
	if got != nil {
		t.Fatal("got an error but didn't none")
	}
}

func fileNames(directories []os.DirEntry) []string {
	mapped := make([]string, len(directories))

	for i, e := range directories {
		mapped[i] = e.Name()
	}

	return mapped
}
