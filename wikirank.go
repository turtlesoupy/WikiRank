package main

import (
  "log"
  "os"
  "github.com/cosbynator/wikirank/ranklib"
)

func main() {
  if len(os.Args) <= 1 {
    log.Fatal("Not enough arguments")
    return
  }

  switch cmd := os.Args[1]; cmd {
  case "extract_graph":
    if len(os.Args) <= 2 {
      log.Fatal("Extract: Required argument 'input.xml' missing")
      return
    }
    filename := os.Args[2]
    log.Printf("Extracting from %s", filename)
    err := ranklib.ReadFrom(filename)
    if err != nil { panic(err) }
    log.Printf("All done!")
  default:
    log.Fatalf("Unknown command '%s'", cmd)
    return
  }
}
