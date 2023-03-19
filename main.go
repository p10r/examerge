package main

import (
	"errors"
	cp "github.com/otiai10/copy"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"golang.org/x/exp/slices"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const outputDirName = "/generated"
const logFileName = "examerge.log"

type Exam struct {
	file string
}

type Rating struct {
	file string
}

func main() {
	ratingPrefix := "rating_"

	workingDir, _ := os.Getwd()

	logFile := setupFileLogging(workingDir)
	defer logFile.Close()

	log.Printf("Target path is %s\n", workingDir)

	examCount := Workflow(workingDir, ratingPrefix)

	log.Printf("Processed %v exams", examCount)
}

func setupFileLogging(workingDir string) *os.File {
	logFile, err := os.OpenFile(
		filepath.Join(workingDir, logFileName),
		os.O_APPEND|os.O_RDWR|os.O_CREATE,
		0644,
	)
	if err != nil {
		log.Panic(err)
	}
	log.SetOutput(logFile)
	return logFile
}

func Workflow(path, ratingPrefix string) int {
	outputPath, err := CreateOutputDirIn(path)
	if err != nil {
		log.Fatalf("Could not create output directory %v", outputDirName)
	}

	err = CopyExceptGenerated(path, outputPath)
	if err != nil {
		log.Fatalf("Error when trying to copy exams into /generated %q", err)
	}

	return MergeAll(outputPath, ratingPrefix)
}

func CreateOutputDirIn(path string) (string, error) {
	destination := path + outputDirName
	err := os.Mkdir(destination, 0750)
	return destination, err
}

func CopyExceptGenerated(input, output string) error {
	options := cp.Options{
		Skip: func(srcInfo os.FileInfo, src, dest string) (bool, error) {
			dirName := "/" + srcInfo.Name()
			return dirName == outputDirName, nil
		},
	}
	return cp.Copy(input, output, options)
}

func MergeAll(parentDir, ratingPrefix string) int {
	dirs := findDirsIn(parentDir)
	var fileCount = 0
	for _, dir := range dirs {
		log.Printf("Current dir: %s\n", dir)
		exam, rating, err := ExamAndRatingFrom(dir, ratingPrefix)
		if err == nil {
			fileCount++
		}
		if err != nil {
			log.Fatalf("handle me") //TODO log
		}
		_, err = Merge(exam, rating)
		if err != nil {
			log.Fatalf("could not merge %s and %s, error: %s", exam.file, rating.file, err)
		}
	}
	return fileCount
}

func findDirsIn(parentDir string) []string {
	items, err := os.ReadDir(parentDir)
	if err != nil {
		log.Fatalf("Could not read items in %s, got %s", parentDir, err)
	}

	var dirs []string
	for _, item := range items {
		if item.IsDir() {
			dirs = append(dirs, filepath.Join(parentDir, item.Name()))
		}
	}
	return dirs
}

func Merge(exam Exam, rating Rating) (Exam, error) {
	inputPaths := []string{rating.file}

	err := api.MergeAppendFile(inputPaths, exam.file, nil)
	if err != nil {
		return Exam{}, err
	}
	err = os.Remove(rating.file)
	if err != nil {
		log.Fatalf("could not remove %s, error: %s", rating.file, err)
	}

	return exam, nil
}

func ExamAndRatingFrom(dir, ratingPrefix string) (Exam, Rating, error) {
	items, err := os.ReadDir(dir)
	if err != nil {
		log.Fatalf("Could not read items in %s, got %s", dir, err)
		return Exam{}, Rating{}, err
	}

	ratingIndex := slices.IndexFunc(items,
		func(item os.DirEntry) bool {
			return strings.HasPrefix(item.Name(), ratingPrefix)
		})

	if ratingIndex == -1 {
		return Exam{}, Rating{}, errors.New("rating file not found")
	}

	examIndex := slices.IndexFunc(items,
		func(item os.DirEntry) bool {
			return !strings.HasPrefix(item.Name(), ratingPrefix)
		})

	if examIndex == -1 {
		return Exam{}, Rating{}, errors.New("exam file not found")
	}

	return Exam{filepath.Join(dir, items[examIndex].Name())},
		Rating{filepath.Join(dir, items[ratingIndex].Name())},
		nil
}
