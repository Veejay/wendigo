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
)

type FileMatch struct {
	fileName    string
	lineNumber  int
	matchRegion []string
}

func searchTerm(filepath string, term string, matches chan<- FileMatch, done chan<- bool) {
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
      m := FileMatch{fileName: filepath, lineNumber: lineNumber, matchRegion: lines[lineNumber-2 : lineNumber+2]}
			matches <- m
		}
	}
  done <- true
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

func ProcessPath(path string, fi os.FileInfo, err error) error {
  // We don't process directories.
  // And I don't care for those silly symbolic links either :/
  if !(fi.IsDir() || ((fi.Mode() & os.ModeSymlink) > 0)) {
    fmt.Printf("Full Path: \t%s\nName: \t%s\n", path, fi.Name())
  }
  return nil
}

func main() {
	dirName := "/Users/bertrand/Programming/wendigo/test"
  filepath.Walk(dirName, ProcessPath)
}
