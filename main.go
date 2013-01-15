package main

import "fmt"
import "net/http"
import "log"
import "os"

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
    fmt.Println(f.Name())
  }
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
