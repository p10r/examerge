package main

import (
	cp "github.com/otiai10/copy"
	"github.com/pdfcpu/pdfcpu/pkg/api"
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

func TestTearDown(t *testing.T, path string) {
	t.Helper()
	err := os.RemoveAll(path)
	if err != nil {
		t.Fatalf("Could not remove %s, error %s", path, err)
	}
}

func TestRemove(t *testing.T, path string) {
	t.Helper()
	err := os.Remove(path)
	if err != nil {
		t.Errorf("Got error when trying to remove %q, %q", path, err.Error())
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
func TestExists(t *testing.T, path string) bool {
	t.Helper()
	_, err := os.Stat(path)

	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func TestExistsOrThrow(path string, t *testing.T) {
	t.Helper()
	exists := TestExists(t, path)
	if !exists {
		t.Fatalf("%s does not exist", path)
	}
}

func AssertIsMerged(t *testing.T, exam Exam) {
	t.Helper()
	pageCount, err := api.PageCountFile(exam.file)
	if err != nil {
		t.Fatalf("Got error, expected none %q", err)
	}
	if pageCount != 4 {
		t.Fatalf("Expected page count to be 4, but was %d", pageCount)
	}
}

func AssertNoError(t testing.TB, got error) {
	t.Helper()
	if got != nil {
		t.Fatal("got an error but didn't expect one")
	}
}

func AssertError(t testing.TB, got error, message string) {
	t.Helper()
	if got.Error() != message {
		t.Fatalf("Expected error message to be %q, was %q", message, got.Error())
	}
	if got == nil {
		t.Fatal("didn't get an error but expected one")
	}
}

func AssertExists(t *testing.T, path string) {
	t.Helper()
	if !TestExists(t, path) {
		t.Errorf("Expected %s to exist, but doesn't", path)
	}
}

func AssertDoesntExist(t *testing.T, path string) {
	t.Helper()
	if TestExists(t, path) {
		t.Errorf("Expected %s to be removed, but isn't", path)
	}
}
