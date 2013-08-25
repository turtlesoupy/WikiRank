package main

import (
  "log"
  "os"
  "math"
  "strconv"
  "runtime"
  "encoding/gob"
  "github.com/cosbynator/external_sort"
  "github.com/cosbynator/wikirank/ranklib"
  "github.com/cosbynator/wikirank/rankhttp"
)

func main() {
  if len(os.Args) <= 1 {
    log.Fatal("Not enough arguments")
    return
  }

  runtime.GOMAXPROCS(runtime.NumCPU()) // Mostly I/O bound, but why not
  gob.Register(ranklib.PageRankedArticle{})
  gob.Register(ranklib.PreprocessedPage{})

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
    err = pageResolver.AddCategoryFromFile("Universities", "/home/tdimson/go/src/github.com/cosbynator/wikirank/data/list_of_world_universities.txt")
    if err != nil { panic(err) }
    err = pageResolver.AddCategoryFromFile("Movies", "/home/tdimson/go/src/github.com/cosbynator/wikirank/data/list_of_movies.txt")
    if err != nil { panic(err) }
    err = pageResolver.AddCategoryFromFile("Countries", "/home/tdimson/go/src/github.com/cosbynator/wikirank/data/list_of_countries.txt")
    if err != nil { panic(err) }

    rankhttp.Serve(pageResolver, port)

  case "dumpcategory":
    if len(os.Args) <= 4 {
      log.Fatal("DumpCategory requires 'rankedPageFile', 'categoryFile' and 'outputName'")
      return
    }

    rankedPageFile := os.Args[2]
    categoryFile := os.Args[3]
    outputName := os.Args[4]

    err := ranklib.DumpCategory(rankedPageFile, categoryFile, outputName)
    if err != nil { panic(err) }
  case "pagerank":
    if len(os.Args) <= 3 {
      log.Fatal("PageRank: Required argument 'input.gob' / 'ranked_output.gob' missing")
      return
    }
    inputName := os.Args[2]
    outputName := os.Args[3]
    log.Printf("Page ranking from '%s' into '%s'", inputName, outputName)
    err := ranklib.RankAndWrite(inputName, outputName)
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

  case "preprocess_xml":
    if len(os.Args) <= 3 {
      log.Fatal("Preprocess: required argument 'input.xml' / 'output.gob' missing")
      return
    }

    xmlFileName := os.Args[2]
    outputName := os.Args[3]

    log.Printf("Preprocessing '%s' into '%s'", xmlFileName, outputName)
    err := ranklib.PreprocessXML(xmlFileName, outputName)
    if err != nil { panic(err) }
  case "rank_preprocessed":
    if len(os.Args) <= 3 {
      log.Fatal("Rank: required argument 'input.gob' / 'output.gob' missing")
      return
    }

    inputFile := os.Args[2]
    outputFile := os.Args[3]
    log.Printf("Ranking from '%s' into '%s'", inputFile, outputFile)
    err := ranklib.PageRankPreprocessedPages(inputFile, outputFile)
    if err != nil { panic(err) }
  case "sort":
    if len(os.Args) <= 3 {
      log.Fatal("Sort: required argument 'input.gob' / 'output.gob' missing")
      return
    }

    inputFile := os.Args[2]
    outputFile := os.Args[3]
    log.Printf("Sorting '%s' into '%s'", inputFile, outputFile)
    input := make(chan *ranklib.PageRankedArticle, 20000)
    inputSort := make(chan external_sort.ComparableItem, 20000)
    outputSort:= make(chan external_sort.ComparableItem, 20000)
    output := make(chan *ranklib.PageRankedArticle, 20000)
    writeDoneChan := make(chan bool)

    toSortAdapter := func() { for i := range input { inputSort <- external_sort.ComparableItem(i) }; close(inputSort) }
    fromSortAdapter:= func() { for i := range outputSort { output <- i.(*ranklib.PageRankedArticle) }; close(output) }

    go ranklib.ReadPageRankedArticles(inputFile, input)
    go toSortAdapter()
    go external_sort.ExternalSort(1000000, ranklib.PageRankedArticleGobHelper{}, inputSort, outputSort)
    go fromSortAdapter()
    go ranklib.WritePageRankedArticles(outputFile, output, writeDoneChan)

    <-writeDoneChan

  case "print":
    if len(os.Args) <= 2 {
      log.Fatal("Print requires input argument")
    }

    inputFile := os.Args[2]
    articles := make(chan *ranklib.PageRankedArticle, 20000)
    go ranklib.ReadPageRankedArticles(inputFile, articles)
    for a := range articles {
      log.Printf("%s: %f", a.Title, a.PageRank)
    }

  default:
    log.Fatalf("Unknown command '%s'", cmd)
    return
  }
  log.Printf("All done!")
}
