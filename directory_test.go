package main

import (
	"os"
	"testing"
)

func TestDirectories(t *testing.T) {
	t.Run("creates a dir for generated output", func(t *testing.T) {
		tmpDir := TestSetupTestEnvironment(t)
		defer TestTearDown(tmpDir, t)

		_, err := CreateOutputDirIn(tmpDir)
		AssertNoError(t, err)

		result, err := TestExists(t, tmpDir+outputDir)
		if err != nil || result == false {
			t.Errorf("Expected %s/generated, but doesn't exist", tmpDir)
		}
	})

	t.Run("copies over existing files into /generated", func(t *testing.T) {
		tmpDir := TestSetupTestEnvironment(t)
		defer TestTearDown(tmpDir, t)

		dir, _ := CreateOutputDirIn(tmpDir)

		err := CopyExceptGenerated(tmpDir, dir)
		AssertNoError(t, err)

		TestExistsOrThrow(dir+"/student1/example_exam1.pdf", t)
		TestExistsOrThrow(dir+"/student1/example_rating1.pdf", t)
		TestExistsOrThrow(dir+"/student2/example_exam2.pdf", t)
		TestExistsOrThrow(dir+"/student2/example_exam2.pdf", t)
	})

	t.Run("merges files", func(t *testing.T) {
		tmpDir := TestSetupTestEnvironment(t)
		exam := Exam{tmpDir + "/student1/example_exam1.pdf"}
		rating := Rating{tmpDir + "/student1/example_rating1.pdf"}

		merged, err := Merge(exam, rating)

		AssertNoError(t, err)
		AssertIsMerged(t, merged)
	})

	//t.Run("merges all file in directory tree", func(t *testing.T) {
	//	tmpDir := TestSetupTestEnvironment(t)
	//	MergeAll(tmpDir, "example_rating")
	//	AssertIsMerged(t, Exam{tmpDir + "/student1/example_exam1.pdf"})
	//	AssertIsMerged(t, Exam{tmpDir + "/student2/example_exam2.pdf"})
	//})

	//t.Run("deletes ratings", func(t *testing.T) {
	//
	//})

	//t.Run("exits if directory has more than two files", func(t *testing.T) {
	//
	//})
}

func TestMergeAll(t *testing.T) {
	t.Run("merges files", func(t *testing.T) {
		tmpDir := TestSetupTestEnvironment(t)
		exam := Exam{tmpDir + "/student1/example_exam1.pdf"}
		rating := Rating{tmpDir + "/student1/example_rating1.pdf"}

		merged, err := Merge(exam, rating)

		AssertNoError(t, err)
		AssertIsMerged(t, merged)
	})

	t.Run("merges all file in directory tree", func(t *testing.T) {
		tmpDir := TestSetupTestEnvironment(t)
		MergeAll(tmpDir, "example_rating")
		AssertIsMerged(t, Exam{tmpDir + "/student1/example_exam1.pdf"})
		AssertIsMerged(t, Exam{tmpDir + "/student2/example_exam2.pdf"})
	})
}

func TestExamAndRatingFrom(t *testing.T) {
	t.Run("returns exam and rating from dir", func(t *testing.T) {
		tmpDir := TestSetupTestEnvironment(t)
		dir := tmpDir + "/student1"

		exam, rating, err := ExamAndRatingFrom(dir, "example_rating")
		AssertNoError(t, err)

		expectedExamPath := dir + "/example_exam1.pdf"
		if exam.file != expectedExamPath {
			t.Fatalf("Expected %q, got %q", expectedExamPath, exam.file)
		}

		expectedRatingPath := dir + "/example_rating1.pdf"
		if rating.file != expectedRatingPath {
			t.Fatalf("Expected %q, got %q", expectedRatingPath, rating.file)
		}
	})

	t.Run("returns error when exam cannot be found", func(t *testing.T) {
		tmpDir := TestSetupTestEnvironment(t)
		os.Remove(tmpDir + "/student1/example_exam1.pdf")

		_, _, err := ExamAndRatingFrom(tmpDir+"/student1", "example_rating")

		if err.Error() != "exam file not found" {
			t.Fatalf("Didn't match expected error, was %q", err)
		}
	})

	t.Run("returns error when rating cannot be found", func(t *testing.T) {
		tmpDir := TestSetupTestEnvironment(t)
		os.Remove(tmpDir + "/student1/example_rating1.pdf")

		_, _, err := ExamAndRatingFrom(tmpDir+"/student1", "example_rating")

		if err.Error() != "rating file not found" {
			t.Fatalf("Didn't match expected error, was %q", err)
		}
	})
}
