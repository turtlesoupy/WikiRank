package ranklib

import (
  "math"
  "log"
)

type RankedPage struct {
  Page Page
  Rank float64
}

func pageRank(pages []RankedPage, teleportProbability float64, convergenceCriteron float64) {
  beta, epsilon := teleportProbability, convergenceCriteron
  log.Printf("Ranking with beta='%f', epsilon='%f'", beta, epsilon)
  n := len(pages)
  idRemap := make(map[uint64] uint32, n)
  lastRank := make([]float64, n)
  thisRank := make([]float64, n)

  for i := 0; i < n; i++ {
    idRemap[pages[i].Page.Id] = uint32(i)
  }

  for iteration, lastChange := 1, math.MaxFloat64; lastChange > epsilon; iteration++ {
    thisRank, lastRank = lastRank, thisRank
    if iteration > 1 {
      // Clear out old values
      for i:=0; i < n; i++ {
        thisRank[i] = 0.0
      }
    } else {
      // Base case: everything uniform
      for i:= 0; i < n; i++ {
        lastRank[i] = 1.0 / float64(n)
      }
    }

    // Single power iteration
    for i := 0; i < n; i++ {
      contribution := beta * lastRank[i] / float64(len(pages[i].Page.Links))
      for _, outboundId := range pages[i].Page.Links {
        thisRank[idRemap[outboundId]] += contribution
      }
    }

    // Reinsert leaked probability
    S := float64(0.0)
    for i := 0; i < n; i++ {
      S += thisRank[i]
    }
    leakedRank := (1.0 - S) / float64(n)
    lastChange = 0.0 // and calculate L1-difference too
    for i := 0; i < n; i++ {
      thisRank[i] += leakedRank
      lastChange += math.Abs(thisRank[i] - lastRank[i])
    }

    log.Printf("Pagerank iteration #%d delta=%f", iteration, lastChange)
  }

  for i := 0; i < n; i++ {
    pages[i].Rank = thisRank[i]
  }
}

func RankAndWrite(inputName string, outputName string) (err error) {
  n := ReadLength(inputName)

  rankList := make([]RankedPage, n)
  log.Printf("Write Ranked Results: %d pages", n)
  inputChan := make(chan *Page, 100)
  go ReadPages(inputName, inputChan)

  log.Printf("Reading pages...")
  i := 0
  for page := range inputChan {
    rankList[i].Page = *page
    rankList[i].Rank = 0.0
    i++
  }

  log.Printf("Computing PageRank...")
  pageRank(rankList, 0.85, 0.001)

  log.Printf("Writing ranked pages...")
  outputChan := make(chan *RankedPage, 100)
  defer close(outputChan)
  go WriteRankedPages(outputName, n, outputChan)
  for i := 0; i < n; i++ {
    outputChan <- &(rankList[i])
  }

  return
}

