package main

import (
	"fmt"
	"log"
	"os"
	"testing"
)

func SetupTestEnvironment(t *testing.T, studentCount int) string {
	tmpDir := CreateTmpDir()
	t.Logf("tmpDir: %s", tmpDir)

	for i := 1; i <= studentCount; i++ {
		fullNewDir := fmt.Sprintf("%s/student%v", tmpDir, i)
		t.Logf("fullNewDir: %s", fullNewDir)

		mkDirErr := os.Mkdir(fullNewDir, 0750)
		if mkDirErr != nil {
			t.Fatalf("Could not create dir %q", mkDirErr)
		}

		exampleFile, _ := os.ReadFile("example_exam.pdf")

		examFileName := fmt.Sprintf("%s/exam_%v", fullNewDir, i)

		err := os.WriteFile(examFileName, exampleFile, 0750)
		if err != nil {
			t.Fatalf("Could not create %s", examFileName)
		}
		t.Logf("Created %s/%s", fullNewDir, examFileName)
	}

	return tmpDir
}

func CreateTmpDir() string {
	dir, err := os.MkdirTemp("", "")
	if err != nil {
		log.Fatal(err)
	}

	return dir
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
