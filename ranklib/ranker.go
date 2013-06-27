package ranklib

import (
  "fmt"
  "log"
  "sort"
)

const (
  PageRankWalkProbability = 0.85
  PageRankConvergence = 0.0001
)

type RankedPage struct {
  Page
  Order uint32
  Rank float32
  OutboundCount uint32
}

func (this *RankedPage) String() string {
  return fmt.Sprintf("RankedPage[title=%s, id=%d, rank=%f]", this.Title, this.Id, this.Rank)
}

type RankedPageList []RankedPage

func (s RankedPageList) Len() int {
  return len(s)
}

func (s RankedPageList) Less(i, j int) bool {
  return s[i].Rank > s[j].Rank
}

func (s RankedPageList) Swap(i, j int) {
    s[i], s[j] = s[j], s[i]
}

func RankAndWrite(inputName string, outputName string) (err error) {
  n := ReadLength(inputName)
  pageList := make([]Page, n)
  log.Printf("Write Ranked Results: %d pages", n)
  inputChan := make(chan *Page, 100)
  go ReadPages(inputName, inputChan)

  log.Printf("Reading pages...")
  i := 0
  for page := range inputChan {
    pageList[i] = *page
    if(pageList[i].Title == "The Wizard of Oz (1939 film)") {
      log.Printf("%s: %q", pageList[i].Title, pageList[i].Links)
    }
    i++
  }

  log.Printf("Computing PageRank...")
  rankVector := pageRank(pageList, PageRankWalkProbability, PageRankConvergence)

  log.Printf("Converting to ranked pages...")
  rankList := make([]RankedPage, n)
  for i := 0; i < n; i++ {
    page := &pageList[i];
    rankList[i] = RankedPage{
      Page: *page,
      Rank: float32(rankVector[i]),
      OutboundCount: uint32(len(page.Links)),
    }
    if(rankList[i].Title == "The Wizard of Oz (1939 film)") {
      log.Printf("%s: %q", rankList[i], rankList[i].Links)
    }
  }

  log.Printf("Calculating percentiles...")
  sort.Sort(RankedPageList(rankList))
  for i := 0; i < n; i++ {
    rankList[i].Order = uint32(i+1)
  }

  log.Printf("Writing ranked pages...")
  outputChan := make(chan *RankedPage, 100)
  writeDoneChan := make(chan bool, 1)
  go WriteRankedPages(outputName, n, outputChan, writeDoneChan)
  for i := 0; i < n; i++ {
    outputChan <- &(rankList[i])
  }
  close(outputChan)
  <-writeDoneChan

  return
}
