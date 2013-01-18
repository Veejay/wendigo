package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
  "io/ioutil"
  "strings"
)

type FileMatch struct {
  fileName string
  lineNumber int
  matchRegion []string
}

// TODO: We might want to add a colorize function to format the matches for the terminal
// Ou alors tout simplement créer un objet HTMLFormatter/TerminalFormatter
func (match FileMatch) String() string {
  return fmt.Sprintf("\nMatch found in %s\n%d:\t%s",
    match.fileName,
    match.lineNumber,
    strings.Replace(strings.Join(match.matchRegion, "\n"), "work_week", "[31;43mwork_week[0m", -1))
}

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
			walk(path, filepaths)
		} else {
			filepaths <- path
		}
	}
}

func recursiveWalk(path string, filepaths chan<- string) {
	walk(path, filepaths)
	close(filepaths)
}

func searchTerm(filepath string, term string) (matches []FileMatch) {
  // Read the file
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
      matches = append(matches, FileMatch{fileName: filepath, lineNumber: lineNumber, matchRegion: lines[lineNumber-2:lineNumber+2]})
    }
  }
  return matches
}

func main() {
	dirName := "/Users/bertrand/Programming/wendigo/test"
	filepaths := make(chan string)

  go recursiveWalk(dirName, filepaths)

	for path := range filepaths {
    for _, match := range searchTerm(path, "work_week") {
      fmt.Println(match.String())
    }
	}
}
