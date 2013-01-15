package main

import "fmt"
import "net/http"
import "log"

// Handles the initial GET on /search 
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
// This should be hooked to the main chunk of the (i.e. once the directory
// and the search term have been obtained, we should launch the crawling)
func searchPostHandler(w http.ResponseWriter, r *http.Request) {
  fmt.Printf("The search term provided by the user is: %s\n", r.FormValue("term"))
  fmt.Printf("The name of the directory provided by the user is: %s\n", r.FormValue("directory"))
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
