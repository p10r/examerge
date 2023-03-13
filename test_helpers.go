package main

import (
	cp "github.com/otiai10/copy"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"testing"
)

func TestSetupTestEnvironment(t *testing.T) string {
	tmpDir := TestCreateTmpDir()

	err := cp.Copy("resources/example/", tmpDir)

	if err != nil {
		t.Fatalf("Could not create tmp dir, error: %s", err)
	}

	dirs := TestListSubDirTree(tmpDir, t)
	for _, dir := range dirs {
		t.Logf("Created %s/%s", tmpDir, dir)
	}

	return tmpDir
}

func TestListSubDirTree(tmpDir string, t *testing.T) []string {
	var dirs []string

	err := filepath.Walk(tmpDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		dirs = append(dirs, info.Name())
		return nil
	})
	if err != nil {
		t.Fatalf("Cannot list all sub directories %s", err)
	}

	return dirs
}

func TestRemove(path string, t *testing.T) {
	err := os.RemoveAll(path)
	if err != nil {
		t.Fatalf("Could not remove %s, error %s", path, err)
	}
}

func TestCreateTmpDir() string {
	dir, err := os.MkdirTemp("", "")
	if err != nil {
		log.Fatal(err)
	}

	return dir
}

// TestExists returns whether the given file or directory TestExists
func TestExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
