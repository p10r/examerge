package main

import (
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"os"
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
	CreateOutputDirIn(path)
}

func CreateOutputDirIn(path string) error {
	return os.Mkdir(path+outputDir, 0750)
}

func Merge(exam Exam, rating Rating, outputDestination string) error {
	inputPaths := []string{exam.file, rating.file}

	return api.MergeCreateFile(inputPaths, outputDestination, nil)
}
