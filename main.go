package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func walk(path string, filepaths chan<- string) {
	dir, err := os.Open(path)
	if err != nil {
		log.Println(err)
		return
	}
	defer dir.Close()
	contents, err := dir.Readdir(0)
	if err != nil {
		log.Println(err)
		return
	}
	for _, f := range contents {
		path := filepath.Join(path, f.Name())
		if f.IsDir() {
			WalkPathRecursively(path, filepaths)
		} else {
			filepaths <- path
		}
	}
}

func recursiveWalk(path string, filepaths chan<- string) {
	walk(path, filepaths)
	close(filepaths)
}

func main() {
	dirName := "/Users/bertrand/Programming/wendigo/test"
	filepaths := make(chan string)
	go recursiveWalk(dirName, filepaths)

	for path := range filepaths {
		fmt.Println(path)
	}
}
