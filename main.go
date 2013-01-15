package main

import (
  "fmt"
  "net/http"
  "log"
  "io"
  "os"
  "bufio"
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

// Signature copied from fmt.Printf
// See http://golang.org/pkg/fmt/#Printf
func PrintFile(fileName string) (n int, err error) {
  file, err := os.Open("/Users/bertrand/diff_branches.txt")
  if err != nil {
    fmt.Printf("Error while opening a file: ", err)
  }
  defer file.Close()
  fileReader := bufio.NewReader(file)
  for {
    line, err := fileReader.ReadString('\n')
    if err != nil {
      if err != io.EOF {
        fmt.Println(err)
      }
      break
    }
    fmt.Print(line)
  }
  // Faking it for now, I still don't know how to "save" 
  // that piece of code in a function
  return 100, nil
}

func main() {

  http.HandleFunc("/", searchGetHandler)
  // Not clear exactly what should be used to discriminate between GETs and POSTs
  // Maybe http.Request.Method??
  http.HandleFunc("/search", searchPostHandler)
  err := http.ListenAndServe(":8080", nil)
  if err != nil {
    log.Fatal("ListenAndServe: ", err)
  }
}
