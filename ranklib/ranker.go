package ranklib

import (
  "log"
  "sort"
  "math"
)

const (
  pageRankWalkProbability = 0.85
  pageRankConvergence = 0.001
)

type Influencer struct {
  Id uint64
  Influence float32
}

type InfluencerList []Influencer

func (s InfluencerList) Len() int {
  return len(s)
}

func (s InfluencerList) Less(i, j int) bool {
  return s[i].Influence > s[j].Influence
}

func (s InfluencerList) Swap(i, j int) {
    s[i], s[j] = s[j], s[i]
}

type RankedPage struct {
  Title string
  Id uint64
  Order uint32
  Rank float32
  OutboundCount uint32
  Influencers []Influencer
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

func computeInfluencers(numInfluencers int, pageList []Page, rankList []RankedPage, idRemap map[uint64] uint32) {
  n := len(pageList)
  for i := 0; i < n; i++ {
    page := &pageList[i]
    rankedPage := &rankList[i]
    perLinkContribution := rankedPage.Rank * pageRankWalkProbability  / float32(len(page.Links))
    for _, linkedId := range page.Links {
      linkedRank := &rankList[idRemap[linkedId]]
      if len(linkedRank.Influencers) < numInfluencers {
        linkedRank.Influencers = append(linkedRank.Influencers, Influencer{
          Id: page.Id,
          Influence: perLinkContribution,
        })
      } else {
        minVal := float32(math.MaxFloat32)
        minIndex := -1
        for i, influencer := range linkedRank.Influencers {
          if influencer.Influence < minVal {
            minVal = influencer.Influence
            minIndex = i
          }
        }
        if perLinkContribution > minVal {
          linkedRank.Influencers[minIndex] = Influencer{
            Id: page.Id,
            Influence: perLinkContribution,
          }
        }
      }
    }
  }

  // Re-sort after, since it is weird otherwise
  for i := 0; i < n; i++ {
    rankedPage := &rankList[i]
    sort.Sort(InfluencerList(rankedPage.Influencers))
  }
}

func RankAndWrite(inputName string, outputName string, numInfluencers int) (err error) {
  n := ReadLength(inputName)
  pageList := make([]Page, n)
  log.Printf("Write Ranked Results: %d pages", n)
  inputChan := make(chan *Page, 100)
  go ReadPages(inputName, inputChan)

  log.Printf("Reading pages...")
  i := 0
  for page := range inputChan {
    pageList[i] = *page
    i++
  }

  log.Printf("Computing PageRank...")
  rankVector, idRemap := pageRank(pageList, pageRankWalkProbability, pageRankConvergence)

  log.Printf("Converting to ranked pages...")
  rankList := make([]RankedPage, n)
  for i := 0; i < n; i++ {
    page := &pageList[i];
    rankList[i] = RankedPage{
      Title: page.Title,
      Id: page.Id,
      Rank: float32(rankVector[i]),
      OutboundCount: uint32(len(page.Links)),
      Influencers: make([]Influencer, 0, numInfluencers),
    }
  }

  log.Printf("Calculating largest influencers...")
  computeInfluencers(numInfluencers, pageList, rankList, idRemap)

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
