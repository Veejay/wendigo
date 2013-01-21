package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

type FileMatch struct {
	fileName    string
	lineNumber  int
	matchRegion []string
}

func (match FileMatch) String() string {
	return strings.Join(match.matchRegion, "\n")
}

func searchTerm(filepath string, term string) (matches []FileMatch) {
	contents, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Println(err)
		return
	}
	lines := strings.Split(string(contents), "\n")
	if err != nil {
		log.Println(err)
		return
	}

	for lineNumber, line := range lines {
		if strings.Contains(line, term) {
			matches = append(matches, FileMatch{fileName: filepath, lineNumber: lineNumber, matchRegion: lines[lineNumber-2 : lineNumber+2]})
		}
	}
	return matches
}

func Pygmentize(snippet string) (pygmentedSnippet string) {

	var out bytes.Buffer
	pygmentize := exec.Command("pygmentize", "-f", "html", "-l", "ruby")

	pygmentize.Stdin = strings.NewReader(snippet)
	pygmentize.Stdout = &out

	pygmentize.Start()

	pygmentize.Wait()

	return string(out.Bytes())
}

func main() {
	directoryPath := "/Users/bertrand/Programming/wendigo/test"

	// Empty slice of strings, will eventually contain the paths to the files
	// gathered by the recursive walk of the directory
	paths := []string{}
	var wg sync.WaitGroup

	filepath.Walk(directoryPath, func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			panic(err)
		}
		// TODO: Take care of the symbolic links as well 
		// See os.ModeSymlink
		if !fi.IsDir() {
			paths = append(paths, path)
		}
		return nil
	})

	for _, path := range paths {
		wg.Add(1)
		// Calls Done() on the waiting to indicate that the task has been completed
		// but since the work is not being done sequentially, the goroutine itself needs to send that 
		// signal
		go func(path string) {
			matches := searchTerm(path, "assignments")
			for _, match := range matches {
				fmt.Println(Pygmentize(match.String()))
			}
			wg.Done()
		}(path)
	}
	// Since all the jobs are launched in goroutines, nothing will block
	// needs to wait for everything to finish before exiting the program 
	wg.Wait()
	fmt.Println("Done processing all the files")
}
