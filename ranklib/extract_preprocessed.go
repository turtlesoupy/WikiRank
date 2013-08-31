package ranklib

import (
  "log"
  "sort"
  "encoding/gob"
  "github.com/cosbynator/external_sort"
)

type PageRankedArticle struct {
  PageRank float64
  Aliases []string
  PreprocessedPage
}

func (self *PageRankedArticle) LessThan(other external_sort.ComparableItem) bool {
  return self.PageRank > other.(*PageRankedArticle).PageRank // desc
}

type PageRankedArticleGobHelper struct {}
func (PageRankedArticleGobHelper) EncodeComparable(g *gob.Encoder, item external_sort.ComparableItem) error {
  return g.Encode(item)
}
func (PageRankedArticleGobHelper) DecodeComparable(g *gob.Decoder) (external_sort.ComparableItem, error) {
  var tmp PageRankedArticle
  err := g.Decode(&tmp)
  return &tmp, err
}

func IntDedupe(s *[]int) {
  newArr := make([]int, 0, len(*s))
  sort.Sort(sort.IntSlice(*s))
  lastItem := 0
  for i, item := range (*s) {
    if i == 0 || item != lastItem {
      newArr = append(newArr, item)
    }
    lastItem = item
  }
  *s = newArr
}

func PageRankPreprocessedPages(inputName string, outputName string) (err error) {
  sequentialIdMap := make(map[string] int, 600000)
  redirectTitleMap := make(map[string] string, 6000000)
  preprocessedPageInputChannel := make(chan *PreprocessedPage, 20000)

  go ReadPreprocessedPages(inputName, preprocessedPageInputChannel)

  // Get ids
  log.Printf("Creating id map...")
  var sequentialId int
  sequentialId = 0
  for pp := range preprocessedPageInputChannel {
    if len(pp.RedirectTo) > 0 {
      redirectTitleMap[pp.Title] = pp.RedirectTo
    } else {
      sequentialIdMap[pp.Title] = sequentialId
      sequentialId++
      if sequentialId % 100000 == 0 {
        log.Printf("Id Map Page %d", sequentialId)
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
      var linkId int
      if redirectTo, ok := redirectTitleMap[linkTitle]; ok {
        linkId, ok = sequentialIdMap[redirectTo]
        if !ok {
          if !titleFilter.MatchString(redirectTo) {
            log.Printf("Unresolvable redirect in '%s': '%s'", pp.Title, redirectTo)
          }
          continue
        }
      } else {
        linkId, ok = sequentialIdMap[linkTitle]
        if !ok {
          //log.Printf("Bad link in '%s': '%s'", pp.Title, linkTitle)
          continue
        }
      }

      if linkId == sequentialId {
        continue // disallow self links
      }

      // Warn: this could introduce duplicates and must be filtered
      nodes[sequentialId].OutboundNeighbors = append(nodes[sequentialId].OutboundNeighbors, linkId)
    }
    sequentialId++
    if sequentialId % 100000 == 0 {
      log.Printf("Graph Page %d", sequentialId)
    }
  }

  log.Printf("Deduping...")
  for i := 0; i < len(nodes); i++ {
    IntDedupe(&nodes[i].OutboundNeighbors)
  }

  sequentialIdMap = nil
  aliasMap := make(map[string] []string, sequentialId)
  for redirectFrom, redirectTo := range redirectTitleMap {
    if aliases, ok := aliasMap[redirectTo]; ok {
      aliasMap[redirectTo] = append(aliases, redirectFrom)
    } else {
      l := make([]string, 0, 1)
      l = append(l, redirectFrom)
      aliasMap[redirectTo] = l
    }
  }
  redirectTitleMap = nil

  // Output cool data format
  log.Printf("Page ranking...")
  g := Graph{Nodes: nodes}
  rankVector := pageRankGraph(g, 0.85, 0.0001)
  g = Graph{}
  nodes = nil

  log.Printf("Outputting data into '%s'...", outputName)
  preprocessedPageInputChannel = make(chan *PreprocessedPage, 20000)
  pageOutputChan := make(chan *PageRankedArticle, 20000)
  writeDoneChan := make(chan bool)
  sequentialId = 0
  go WritePageRankedArticles(outputName, pageOutputChan, writeDoneChan)
  go ReadPreprocessedPages(inputName, preprocessedPageInputChannel)
  for pp := range preprocessedPageInputChannel {
    if len(pp.RedirectTo) > 0 {
      continue
    }

    aliases := aliasMap[pp.Title]
    sort.Strings(aliases)
    pageRank := rankVector[sequentialId]
    article := &PageRankedArticle{
      PreprocessedPage: *pp,
      PageRank: pageRank,
      Aliases: aliases,
    }
    pageOutputChan <- article

    sequentialId++
    if sequentialId % 100000 == 0 {
      log.Printf("Output page %d", sequentialId)
    }
  }
  close(pageOutputChan)
  <-writeDoneChan

  return
}
