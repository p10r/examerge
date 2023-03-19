package main

import (
	"errors"
	"fmt"
	cp "github.com/otiai10/copy"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"golang.org/x/exp/slices"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const outputDir = "/generated"

type Exam struct {
	file string
}

type Rating struct {
	file string
}

func main() {

}

func Workflow(path string) {
	path, err := CreateOutputDirIn(path)
	if err != nil {
		log.Fatalf("Could not create output directory %v", outputDir)
	}
	workingDir, _ := os.Getwd()

	CopyExceptGenerated(workingDir, path)
}

func CreateOutputDirIn(path string) (string, error) {
	destination := path + outputDir
	err := os.Mkdir(destination, 0750)
	return destination, err
}

func CopyExceptGenerated(input, output string) error {
	options := cp.Options{
		Skip: func(srcInfo os.FileInfo, src, dest string) (bool, error) {
			dirName := "/" + srcInfo.Name()
			return dirName == outputDir, nil
		},
	}
	return cp.Copy(input, output, options)
}

func MergeAll(parentDir, ratingPrefix string) {
	dirs := findDirsIn(parentDir)

	for _, dir := range dirs {
		fmt.Println("Current dir: " + dir) //TODO log
		exam, rating, err := ExamAndRatingFrom(dir, ratingPrefix)
		if err != nil {
			log.Fatalf("handle me") //TODO log
		}
		Merge(exam, rating) //TODO handle
		os.Remove(rating.file)
	}
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
