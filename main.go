package main

import "os"

func main() {

}

func Workflow(path string) {
	os.Chdir(path)
	os.Mkdir("generated", 0750)
}
