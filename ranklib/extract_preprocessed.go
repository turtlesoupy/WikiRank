package ranklib

import (
  "log"
)

func PageRankPreprocessedPages(inputName string, outputName string) (err error) {
  sequentialIdMap := make(map[string] uint32, 600000)
  redirectTitleMap := make(map[string] string, 6000000)
  preprocessedPageInputChannel := make(chan *PreprocessedPage, 20000)

  go ReadPreprocessedPages(inputName, preprocessedPageInputChannel)

  // Get ids
  log.Printf("Creating id map...")
  var sequentialId uint32
  sequentialId = 0
  for pp := range preprocessedPageInputChannel {
    if len(pp.RedirectTo) > 0 {
      redirectTitleMap[pp.Title] = pp.RedirectTo
    } else {
      sequentialIdMap[pp.Title] = sequentialId
      sequentialId++
      if sequentialId % 10000 == 0 {
        log.Printf("Page %d", sequentialId)
      }
    }
  }

  // Create graph
  log.Printf("Creating graph...")
  nodes := make([]GraphNode, sequentialId)
  preprocessedPageInputChannel = make(chan *PreprocessedPage, 20000)
  sequentialId = 0
  go ReadPreprocessedPages(inputName, preprocessedPageInputChannel)
  for pp := range preprocessedPageInputChannel {
    if len(pp.RedirectTo) > 0 {
      continue
    }


    for _, textLink := range pp.TextLinks {
      linkTitle := textLink.ArticleTitle
      var linkId uint32
      if redirectTo, ok := redirectTitleMap[linkTitle]; ok {
        linkId, ok = sequentialIdMap[redirectTo]
        if !ok {
          log.Printf("Unresolvable redirect in '%s': '%s'", pp.Title, linkTitle)
          continue
        }
      } else {
        linkId, ok = sequentialIdMap[linkTitle]
        if !ok {
          log.Printf("Bad link in '%s': '%s'", pp.Title, linkTitle)
          continue
        }
      }
      nodes[sequentialId].OutboundNeighbors = append(nodes[sequentialId].OutboundNeighbors, linkId)
    }
    sequentialId++
    if sequentialId % 10000 == 0 {
      log.Printf("Page %d", sequentialId)
    }
  }

  sequentialIdMap = nil

  // Output cool data format
  log.Printf("Page ranking...")
  g := Graph{Nodes: nodes}
  rankVector := pageRankGraph(g, 0.85, 0.0001)
  g = Graph{}
  nodes = nil

  log.Printf("Outputting data...")
  preprocessedPageInputChannel = make(chan *PreprocessedPage, 20000)
  sequentialId = 0
  go ReadPreprocessedPages(inputName, preprocessedPageInputChannel)
  for pp := range preprocessedPageInputChannel {
    if len(pp.RedirectTo) > 0 {
      continue
    }

    pageRank := rankVector[sequentialId]
    log.Printf("Page rank: %f", pageRank)

    sequentialId++
    if sequentialId % 10000 == 0 {
      log.Printf("Page %d", sequentialId)
    }
  }

  return
}
