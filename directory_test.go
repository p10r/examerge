package main

import (
	"testing"
)

func TestDirectories(t *testing.T) {
	t.Run("creates a dir for generated output", func(t *testing.T) {
		tmpDir := TestSetupTestEnvironment(t)
		defer TestTearDown(tmpDir, t)

		Workflow(tmpDir)

		result, err := TestExists(tmpDir + outputDir)
		if err != nil || result == false {
			t.Errorf("Expected %s/generated, but doesn't exist", tmpDir)
		}
	})

	t.Run("copies over existing files into /generated", func(t *testing.T) {
		tmpDir := TestSetupTestEnvironment(t)
		defer TestTearDown(tmpDir, t)

		dir, _ := CreateOutputDirIn(tmpDir)

		CopyExceptGenerated(tmpDir, dir)
		TestExistsOrThrow(dir+"/student1/example_exam1.pdf", t)
		TestExistsOrThrow(dir+"/student1/example_rating1.pdf", t)
		TestExistsOrThrow(dir+"/student2/example_exam2.pdf", t)
		TestExistsOrThrow(dir+"/student2/example_exam2.pdf", t)
	})

	t.Run("merges files in a given dir", func(t *testing.T) {
		tmpDir := TestSetupTestEnvironment(t)
		exam := Exam{tmpDir + "/student1/example_exam1.pdf"}
		rating := Rating{tmpDir + "/student1/example_rating1.pdf"}

		merged, err := Merge(exam, rating)
		if err != nil {
			t.Errorf("Error! %q", err)
		}

		AssertIsMerged(t, merged)
	})
}
