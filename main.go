package main

import (
  "fmt"
  "net/http"
  "strings"
  "log"
  "io/ioutil"
  "path/filepath"
  "os"
)

// Handles the initial GET on /
func searchGetHandler(w http.ResponseWriter, r *http.Request) {
  // FIXME: Use templates for all this
  searchForm := `
    <html>
      <head>
      </head>
      <body>
        <form action="/search" method="POST">
          <div>
            <input name="term" size="60" />
          </div>
          <div>
            <input name="directory" size="60" />
          </div>
          <div>
            <input type="submit" value="Search">
          </div>
        </form>
      </body>
    </html>
  `
  fmt.Fprintf(w, searchForm)
}

// TODO: This should be hooked to the main chunk of the (i.e. once the directory
// and the search term have been obtained, we should launch the crawling)

// Handles POST requests to /search
func searchPostHandler(w http.ResponseWriter, r *http.Request) {
  directoryName, searchTerm := r.FormValue("directory"), r.FormValue("term")

  fmt.Println(directoryName)
  fmt.Println(searchTerm)

  dirHandle, err := os.Open(directoryName)

  // Ensure that the file descriptor gets properly closed
  defer dirHandle.Close()

  if err != nil {
    log.Fatal("Error while trying to open the directory: ", err)
    os.Exit(1)
  }
  // Read all the contents of the directory at once 
  // (See http://golang.org/pkg/os/#File.Readdir)
  files, err := dirHandle.Readdir(0)
  if err != nil {
    log.Fatal("Error while reading the contents of the directory: ", err)
  }
  for _, f := range files {
    // Sorry son, you got a bad case of ugly
    if f.IsDir() {
      fmt.Fprintf(w, "%s is a directory\n", f.Name())
    } else {
      fmt.Fprintf(w, "%s is a file\n", f.Name())
    }
  }
}

func CountFiles(fi []os.FileInfo) (count int) {
  count = 0
  for _, elem := range fi {
    if !elem.IsDir() {
      count++
    }
  }
  return
}

// Greps all the file in a given directory and returns the results 
// as an array of strings
func DirGrep(directoryName string, searchTerm string) (matches []string) {

  // Acquire a file handle
  dir, err := os.Open(directoryName)
  defer dir.Close()
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
  // Gets the contents of the directory (files and subdirectories)
  contents, err := dir.Readdir(0)
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
  leftToProcess := CountFiles(contents)
  // Create a channel so that we can grep all the files concurrently and synchronize
  // (Whatever that means)
  ch := make(chan string)

  for _, f := range contents {
    // If the entry is a file, start a new goroutine
    // and write on the channel when a match is found
    // If it's a directory we want to 
    // recursively call the DirGrep function on it
    if f.IsDir() {
      fmt.Printf("\n\nDIRECTORY: %s\n\n", filepath.Join("/Users", "bertrand", "Programming", "wendigo", "test", f.Name()))
      continue
    } else {
      fmt.Printf("\n\nFILE: %s\n\n", filepath.Join("/Users", "bertrand", "Programming", "wendigo", "test", f.Name()))
      go func(fileName string) {
        stringContents, err := ioutil.ReadFile(filepath.Join("/Users", "bertrand", "Programming", "wendigo", "test", fileName))
        if err != nil {
          fmt.Println(err)
          os.Exit(1)
        }
        lines := strings.Split(string(stringContents), "\n")
        for _, line := range lines {
          if strings.Contains(line, searchTerm) {
            ch <- fmt.Sprintf("Found the search term in file @@ %s @@\n\n", fileName)
          }
        }
        leftToProcess--
      }(f.Name())
    }
  }
  for {
    select {
      case m, ok := <-ch:
        if !ok {
          fmt.Println("Channel is closed")
        } else {
          matches = append(matches, m)
          if leftToProcess == 0 {
            goto END
          }
        }
    }
  }
END:
  return matches
}

// TODO: Pass a string channel as a parameter and stuff 
// the results in there for the caller to consume
// Signature copied from fmt.Printf
// See http://golang.org/pkg/fmt/#Printf
/* func PrintFile(fileName string) (n int, err error) { */
/*   file, err := os.Open(fileName) */
/*   if err != nil { */
/*     fmt.Printf("Error while opening a file: ", err) */
/*   } */
/*   defer file.Close() */
/*   fileReader := bufio.NewReader(file) */
/*   lineNumber := 1 */
/*   for { */
/*     line, err := fileReader.ReadString('\n') */
/*     if err != nil { */
/*       if err != io.EOF { */
/*         fmt.Println(err) */
/*       } */
/*       break */
/*     } */
/*     if strings.Contains(line, "work_week") { */
/*       fmt.Printf("Term \"%s\" found at line %d\n", "work_week", lineNumber) */
/*     } */
/*     lineNumber++ */
/*   } */
/*   // Faking it for now, I still don't know how to "save" */
/*   // that piece of code in a function */
/*   return 100, nil */
/* } */

func main() {
  found := DirGrep("/Users/bertrand/Programming/wendigo/test", "FOOBAR")

  for _, match := range found {
    fmt.Printf("\n\nMATCH: %s\n\n", match)
  }
  /* http.HandleFunc("/", searchGetHandler) */
  /* // Not clear exactly what should be used to discriminate between GETs and POSTs */
  /* // Maybe http.Request.Method?? */
  /* http.HandleFunc("/search", searchPostHandler) */
  /* err := http.ListenAndServe(":8080", nil) */
  /* if err != nil { */
  /*   log.Fatal("ListenAndServe: ", err) */
  /* } */
}
