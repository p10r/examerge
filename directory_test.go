package main

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"testing"
)

func TestDirectories(t *testing.T) {
	t.Run("lists all sub-directories", func(t *testing.T) {
		want := []string{"student1", "student2"}
		tmpDir := setupTestEnvironment(t, 2)

		got, err := ListDirectories(tmpDir)

		assertNoError(t, err)

		if !reflect.DeepEqual(fileNames(got), want) {
			t.Errorf("got %q want %q", fileNames(got), want)
		}
	})

	t.Run("creates a dir for generated output", func(t *testing.T) {
		tmpDir := setupTestEnvironment(t, 2)

		Workflow(tmpDir)

		result, err := exists(tmpDir + "/generated")
		if err != nil || result == false {
			t.Errorf("Expected %s/generated, but doesn't exist", tmpDir)
		}
	})
}

func setupTestEnvironment(t *testing.T, studentCount int) string {
	tmpDir := CreateTmpDir()
	t.Logf("tmpDir: %s", tmpDir)

	for i := 1; i <= studentCount; i++ {
		newDir := fmt.Sprintf("%s/student%v", tmpDir, i)
		t.Logf("newDir: %s", newDir)

		mkDirErr := os.Mkdir(newDir, 0750)
		if mkDirErr != nil {
			t.Fatalf("Could not create dir %q", mkDirErr)
		}

		exampleFile, _ := os.ReadFile("example.pdf")

		examFileName := fmt.Sprintf("%s/exam_%v", newDir, i)

		err := os.WriteFile(examFileName, exampleFile, 0750)
		if err != nil {
			t.Fatalf("Could not create %s", examFileName)
		}
		t.Logf("Created %s/%s", newDir, examFileName)
	}
	//defer os.RemoveAll(tmpDir)
	return tmpDir
}

// exists returns whether the given file or directory exists
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func CreateTmpDir() string {
	dir, err := os.MkdirTemp("", "")
	if err != nil {
		log.Fatal(err)
	}

	return dir
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
