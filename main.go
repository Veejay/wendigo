package main

import (
  "fmt"
  "strings"
  "os"
  "io/ioutil"
  "runtime"
)

type Search struct {
  filename string
  matches chan<-FileMatch
}

type FileMatch struct {
  filename string
  linenumber int
  match string
}

// Number of workers we'll be spawning
var workers = runtime.NumCPU()

func main() {
  term := os.Args[1]
  findMatches(term, os.Args[2:])
}


func findMatches(term string, filenames []string) {
  searches := make(chan Search, workers)
  done := make(chan bool, workers)
  matches := make(chan FileMatch, len(filenames))

  go queueSearches(searches, matches, filenames) // Send all the searches on the searches channel 

  for i := 0; i < workers; i++ {
    go executeSearches(term, done, searches) // Get the searches from the channel and execute them
  }
  // Wait for all the workers to be totally done (i.e. they have consumed everything on the searches channel)
  go tallyCompleted(matches, done)

  // Will try to read from the matches channel until it's closed, which
  // will happen when each single worker is done with its workload
  for match := range matches {
    // TODO: Add a formatter 
    fmt.Printf("Result:\t%#v\n", match)
  }
  fmt.Println("DONE")
}

func queueSearches(searches chan<-Search, matches chan<-FileMatch, filenames []string) {
  for _, filename := range filenames {
    searches<-Search{filename: filename, matches: matches}
  }
  close(searches)
}

func executeSearches(term string, done chan<-bool, searches <-chan Search) {
  // Will keep on reading on the channel until a value arrives
  for search := range searches {
	  contents, err := ioutil.ReadFile(search.filename)
	  if err != nil {
      panic(err)
	  }
	  lines := strings.Split(string(contents), "\n")
	  if err != nil {
      panic(err)
	  }
	  for n, line := range lines {
		  if strings.Contains(line, term) {
			  search.matches <- FileMatch{filename: search.filename, linenumber: n, match: line}
		  }
	  }
  }
  // The channel is dry to the bone, let's bounce
  done <- true
}

func tallyCompleted(matches chan FileMatch, done <-chan bool) {
  // Once all the workers have sent their "done" signal, we can close the results channel and exit the program
  for w := 0; w < workers; w++ {
    <-done
  }
  close(matches)
}
