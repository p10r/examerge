package main

import (
	"path/filepath"
	"testing"
)

func TestE2E(t *testing.T) {
	tmpDir := TestSetupTestEnvironment(t)
	Workflow(tmpDir, "example_rating")

	t.Run("moves all exams", func(t *testing.T) {
		firstExam := filepath.Join(tmpDir, "generated", "student1", "example_exam1.pdf")
		secondExam := filepath.Join(tmpDir, "generated", "student2", "example_exam2.pdf")

		AssertExists(t, firstExam)
		AssertExists(t, secondExam)
	})

	t.Run("merges all exams", func(t *testing.T) {
		firstExam := filepath.Join(tmpDir, "generated", "student1", "example_exam1.pdf")
		secondExam := filepath.Join(tmpDir, "generated", "student2", "example_exam2.pdf")

		AssertIsMerged(t, Exam{firstExam})
		AssertIsMerged(t, Exam{secondExam})
	})

	t.Run("removes ratings afterwards", func(t *testing.T) {
		firstRating := filepath.Join(tmpDir, "generated", "student1", "example_rating1.pdf")
		secondRating := filepath.Join(tmpDir, "generated", "student2", "example_rating2.pdf")

		AssertDoesntExist(t, firstRating)
		AssertDoesntExist(t, secondRating)
	})
}

func TestDirectories(t *testing.T) {
	t.Run("creates a dir for generated output", func(t *testing.T) {
		tmpDir := TestSetupTestEnvironment(t)
		defer TestTearDown(t, tmpDir)

		_, err := CreateOutputDirIn(tmpDir)
		AssertNoError(t, err)

		AssertExists(t, tmpDir+outputDirName)
	})

	t.Run("copies over existing files into /generated", func(t *testing.T) {
		tmpDir := TestSetupTestEnvironment(t)
		defer TestTearDown(t, tmpDir)

		outputDir, _ := CreateOutputDirIn(tmpDir)

		err := CopyExceptGenerated(tmpDir, outputDir)
		AssertNoError(t, err)

		AssertExists(t, outputDir+"/student1/example_exam1.pdf")
		AssertExists(t, outputDir+"/student1/example_rating1.pdf")
		AssertExists(t, outputDir+"/student2/example_exam2.pdf")
		AssertExists(t, outputDir+"/student2/example_exam2.pdf")
	})

	t.Run("merges files", func(t *testing.T) {
		tmpDir := TestSetupTestEnvironment(t)
		defer TestTearDown(t, tmpDir)

		exam := Exam{tmpDir + "/student1/example_exam1.pdf"}
		rating := Rating{tmpDir + "/student1/example_rating1.pdf"}

		merged, err := Merge(exam, rating)

		AssertNoError(t, err)
		AssertIsMerged(t, merged)
	})

	t.Run("deletes ratings", func(t *testing.T) {
		tmpDir := TestSetupTestEnvironment(t)
		defer TestTearDown(t, tmpDir)

		exam := Exam{tmpDir + "/student1/example_exam1.pdf"}
		rating := Rating{tmpDir + "/student1/example_rating1.pdf"}

		_, err := Merge(exam, rating)

		AssertNoError(t, err)
		TestExists(t, rating.file)
	})
}

func TestMergeAll(t *testing.T) {
	t.Run("merges files", func(t *testing.T) {
		tmpDir := TestSetupTestEnvironment(t)
		defer TestTearDown(t, tmpDir)

		exam := Exam{tmpDir + "/student1/example_exam1.pdf"}
		rating := Rating{tmpDir + "/student1/example_rating1.pdf"}

		merged, err := Merge(exam, rating)

		AssertNoError(t, err)
		AssertIsMerged(t, merged)
	})

	t.Run("merges all file in directory tree", func(t *testing.T) {
		tmpDir := TestSetupTestEnvironment(t)
		defer TestTearDown(t, tmpDir)

		MergeAll(tmpDir, "example_rating")
		AssertIsMerged(t, Exam{tmpDir + "/student1/example_exam1.pdf"})
		AssertIsMerged(t, Exam{tmpDir + "/student2/example_exam2.pdf"})
	})

	t.Run("removes ratings afterwards", func(t *testing.T) {
		tmpDir := TestSetupTestEnvironment(t)
		defer TestTearDown(t, tmpDir)

		MergeAll(tmpDir, "example_rating")

		deletedFile := tmpDir + "/student1/example_rating1.pdf"
		exists := TestExists(t, deletedFile)

		if exists {
			t.Errorf("Expected %q to be deleted but is not", deletedFile)
		}
	})
}

func TestExamAndRatingFrom(t *testing.T) {
	t.Run("returns exam and rating from dir", func(t *testing.T) {
		tmpDir := TestSetupTestEnvironment(t)
		defer TestTearDown(t, tmpDir)

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
		defer TestTearDown(t, tmpDir)

		TestRemove(t, tmpDir+"/student1/example_exam1.pdf")

		_, _, err := ExamAndRatingFrom(tmpDir+"/student1", "example_rating")

		AssertError(t, err, "exam file not found")
	})

	t.Run("returns error when rating cannot be found", func(t *testing.T) {
		tmpDir := TestSetupTestEnvironment(t)
		defer TestTearDown(t, tmpDir)

		TestRemove(t, tmpDir+"/student1/example_rating1.pdf")

		_, _, err := ExamAndRatingFrom(tmpDir+"/student1", "example_rating")

		AssertError(t, err, "rating file not found")
	})
}
