package ranklib

import (
  "os"
  "log"
  "bufio"
  "strings"
  "io/ioutil"
  "encoding/json"
)

type CategoryRankedPage struct {
  Title string
  PageRank float32
}

type DumpConfig struct {
  CaseSensitive bool
}

func DefaultDumpConfig() DumpConfig {
  return DumpConfig{
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

func DumpCategory(rankedPageFile string, categoryFile string, outputFile string) error {
  dumpConfig := DefaultDumpConfig()

  rpchan := make(chan *RankedPage, 1000)
  output := make([]CategoryRankedPage, 0, 1000)

  titles, err := readCategoryTitles(dumpConfig, categoryFile)
  if err != nil { return err }

  go ReadRankedPages(rankedPageFile, rpchan)
  for rankedPage := range rpchan {
    noisy := false
    if strings.HasPrefix(strings.ToUpper(rankedPage.Title), "IRON MAN") {
      noisy = true
        log.Printf("%s from %q", rankedPage.Title, rankedPage.Aliases)
    }

    nt := normalizedTitle(dumpConfig, rankedPage.Title)
    _, found := titles[nt]
    if found {
      titles[nt] = false
      if noisy {
        log.Printf("Invaliding %s from %s", nt, rankedPage.Title)
      }
    }

    for _, alias := range rankedPage.Aliases {
      nt = normalizedTitle(dumpConfig, alias)
      _, aliasFound := titles[nt]
      if aliasFound {
        titles[nt] = false
        found = true
        if noisy {
          log.Printf("Invaliding %s from %q", nt, rankedPage.Aliases)
        }
      }
    }

    if !found {
      continue
    }


    cp := CategoryRankedPage{Title: rankedPage.Title, PageRank: rankedPage.Rank}
    output = append(output, cp)
  }

  for title, needsMatch := range titles {
    if needsMatch {
      log.Printf("WARNING: Unable to match '%s' to an article", title)
    }
  }

  if data, err := json.MarshalIndent(output, "", "\t"); err != nil {
    return err
  } else {
    return ioutil.WriteFile(outputFile, data, 0755)
  }
}
