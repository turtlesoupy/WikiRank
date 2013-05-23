package main

import (
  "log"
  "os"
  "math"
  "strconv"
  "runtime"
  "github.com/cosbynator/wikirank/ranklib"
  "github.com/cosbynator/wikirank/rankhttp"
)

func main() {
  runtime.GOMAXPROCS(runtime.NumCPU()) // Mostly I/O bound, but why not

  if len(os.Args) <= 1 {
    log.Fatal("Not enough arguments")
    return
  }

  switch cmd := os.Args[1]; cmd {
  case "serve":
    if len(os.Args) < 4 {
      log.Fatal("Http server requires two arguments, trie location and port")
      return
    }

    limit := math.MaxInt32
    if len(os.Args) > 4 {
      var err error
      limit, err = strconv.Atoi(os.Args[4])
      if err != nil { panic(err) }
      log.Printf("Page limit is %d", limit)
    }

    trieLocation := os.Args[2]
    port, err := strconv.Atoi(os.Args[3])
    if err != nil { panic(err) }
    log.Printf("Building resolver from %s", trieLocation)
    pageResolver, err := ranklib.CreatePageResolver(trieLocation, limit)
    if err != nil { panic(err) }

    rankhttp.Serve(pageResolver, port)
  case "pagerank":
    if len(os.Args) <= 3 {
      log.Fatal("PageRank: Required argument 'input.gob' / 'ranked_output.gob' missing")
      return
    }
    inputName := os.Args[2]
    outputName := os.Args[3]
    log.Printf("Page ranking from '%s' into '%s'", inputName, outputName)
    err := ranklib.RankAndWrite(inputName, outputName, 10)
    if err != nil { panic(err) }

  case "extract_graph":
    if len(os.Args) <= 3 {
      log.Fatal("Extract: Required argument 'input.xml' / 'output.gob' missing")
      return
    }
    filename := os.Args[2]
    outputName := os.Args[3]

    log.Printf("Extracting from '%s' into '%s'", filename, outputName)
    err := ranklib.ReadFrom(filename, outputName)
    if err != nil { panic(err) }
  default:
    log.Fatalf("Unknown command '%s'", cmd)
    return
  }
  log.Printf("All done!")
}
