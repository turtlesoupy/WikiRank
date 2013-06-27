package ranklib

import (
  "os"
  "log"
  "sort"
  "bufio"
  "strings"
  "io/ioutil"
  "encoding/json"
)

type CategoryPage struct {
  RankedPage
  CategoryInfluencers CategoryInfluencers
}

type CategoryInfluencer struct {
  Title string
  Influence float32
}
type CategoryInfluencers []CategoryInfluencer
func (s CategoryInfluencers) Len() int { return len(s) }
func (s CategoryInfluencers) Less(i, j int) bool { return s[i].Influence > s[j].Influence }
func (s CategoryInfluencers) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

type CategoryPageList []*CategoryPage
func (s CategoryPageList) Len() int { return len(s) }
func (s CategoryPageList) Less(i, j int) bool { return s[i].Rank > s[j].Rank }
func (s CategoryPageList) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

type DumpConfig struct {
  NumInfluencers int
  NumItems int
  CaseSensitive bool
}

func DefaultDumpConfig() DumpConfig {
  return DumpConfig{
    NumInfluencers: -1,
    NumItems: -1,
    CaseSensitive: true,
  }
}

func normalizedTitle(dumpConfig DumpConfig, t string) string {
  if dumpConfig.CaseSensitive {
    return t
  } else {
    return strings.ToUpper(t)
  }
}

func DumpCategory(rankedPageFile string, categoryFile string, outputFile string) error {
  dumpConfig := DefaultDumpConfig()
  dumpConfig.NumItems = 25
  dumpConfig.NumInfluencers = 10

  rpchan := make(chan *RankedPage, 1000)
  categoryPages := make(map[uint64] *CategoryPage, 1000)


  titles, err := readCategoryTitles(dumpConfig, categoryFile)
  if err != nil { return err }

  go ReadRankedPages(rankedPageFile, rpchan)
  for rankedPage := range rpchan {
    if invalidateTitles(dumpConfig, rankedPage, titles) {
      categoryPages[rankedPage.Id] = &CategoryPage{
        RankedPage: *rankedPage,
        CategoryInfluencers: make(CategoryInfluencers, 0),
      }
    }
  }

  computeCategoryPageMetrics(dumpConfig, categoryPages)
  err = sortAndWrite(dumpConfig, categoryPages, outputFile)

  if err != nil { return err }

  for title, needsMatch := range titles {
    if needsMatch {
      log.Printf("WARNING: Unable to match '%s' to an article", title)
    }
  }
  return nil
}

func sortAndWrite(dumpConfig DumpConfig, categoryPages map[uint64] *CategoryPage, outputFile string) error {
  pageList := make(CategoryPageList, 0, len(categoryPages))
  totalRank := 0.0
  for _, v := range categoryPages {
    totalRank += float64(v.Rank)
    sort.Sort(v.CategoryInfluencers)
    pageList = append(pageList, v)
  }
  sort.Sort(pageList)

  jsonList := make([](map[string] interface{}), 0, len(pageList))
  for i, categoryPage := range pageList {
    if dumpConfig.NumItems > 0 && i >= dumpConfig.NumItems {
      break
    }

    influencersJson := make([]map[string] interface{}, 0, len(categoryPage.CategoryInfluencers))
    for j, categoryInfluencer := range categoryPage.CategoryInfluencers {
      if dumpConfig.NumInfluencers > 0 && j >= dumpConfig.NumInfluencers {
        break
      }

      influencersJson = append(influencersJson, map[string] interface{} {
        "title": categoryInfluencer.Title,
        "influence": categoryInfluencer.Influence,
      })
    }

    jsonList = append(jsonList, map[string] interface{} {
      "title": categoryPage.Title,
      "pageRank": float64(categoryPage.Rank) / totalRank, // conditional pageRank
      "position": i+1,
      "influencers": influencersJson,
    })
  }

  if data, err := json.MarshalIndent(jsonList, "", "\t"); err != nil {
    return err
  } else {
    return ioutil.WriteFile(outputFile, data, 0755)
  }
}

func computeCategoryPageMetrics(dumpConfig DumpConfig, categoryPages map[uint64] *CategoryPage) {
  for _, categoryPage := range categoryPages {
    for _, link := range categoryPage.Links {
      if linkedPage, ok := categoryPages[link.PageId]; ok {
        linkedPage.CategoryInfluencers = append(linkedPage.CategoryInfluencers, CategoryInfluencer {
          Title: categoryPage.Title,
          Influence: PageRankWalkProbability * (categoryPage.Rank / float32(len(categoryPage.Links))) / linkedPage.Rank,
        })
      }
    }
  }
}

func readCategoryTitles(dumpConfig DumpConfig, categoryFile string) (map[string] bool, error) {
  titleMap := make(map[string] bool, 10000)

  cf, err := os.Open(categoryFile)
  if err != nil { return titleMap, err }

  defer cf.Close()
  scanner := bufio.NewScanner(cf)
  for scanner.Scan() {
    nt := normalizedTitle(dumpConfig, scanner.Text())
    titleMap[nt] = true
  }

  if err := scanner.Err(); err != nil {
    return titleMap, err
  }

  return titleMap, nil
}

func computeCategoryPages(dumpConfig DumpConfig, pages map[uint64] *CategoryPage) {
  for _, v := range pages {
    for _, link := range v.Links {
      if linkedPage, ok := pages[link.PageId]; ok {
        linkedPage.CategoryInfluencers = append(linkedPage.CategoryInfluencers, CategoryInfluencer{
          Title: v.Title,
          Influence: v.Rank / float32(len(v.Links)),
        })
      }
    }
  }
}

func invalidateTitles(dumpConfig DumpConfig, page *RankedPage, titles map[string] bool) bool {
  noisy := false
  if page.Title == "Krazy Kat" {
    noisy = true
  }

  nt := normalizedTitle(dumpConfig, page.Title)
  _, found := titles[nt]
  if found {
    titles[nt] = false
    if noisy {
      log.Printf("Invaliding %s from %s", nt, page.Title)
    }
  }

  for _, alias := range page.Aliases {
    nt = normalizedTitle(dumpConfig, alias)
    _, aliasFound := titles[nt]
    if aliasFound {
      titles[nt] = false
      found = true
      if noisy {
        log.Printf("Invaliding %s from %q", nt, page.Aliases)
      }
    }
  }

  return found
}
