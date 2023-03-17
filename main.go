package main

import (
	cp "github.com/otiai10/copy"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"log"
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

func Merge(exam Exam, rating Rating) (Exam, error) {
	inputPaths := []string{rating.file}

	err := api.MergeAppendFile(inputPaths, exam.file, nil)
	if err != nil {
		return Exam{}, err
	}

	return exam, nil
}
