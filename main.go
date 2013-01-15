package main

import "fmt"
import "os"
import "strings"
import "io/ioutil"

func main() {
  var dirName, searchTerm string
  fmt.Scanf("%s", &dirName)
  fmt.Scanf("%s", &searchTerm)


  fmt.Printf("The name of the directory is %s\n", dirName)
  fmt.Printf("The search term is %s\n", searchTerm)

  contents, err := ioutil.ReadFile("/Users/Bertrand/diff_branches.txt")
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
  for i, line := range strings.Split(string(contents), "\n") {
    if strings.Contains(string(line), "%input") {
      fmt.Printf("Found the string %s on line %d\n", "%input", i)
    }
  }
}
