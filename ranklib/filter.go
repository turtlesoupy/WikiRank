package ranklib

import (
  "os"
  "bufio"
)

type FilterOptions struct {}

func normalizePageTitle(title string) string {
  return title
}

func FilterPageRankedArticles(filterName string, filterOptions FilterOptions,
                              input chan *PageRankedArticle,
                              output chan *PageRankedArticle) {
  defer close(output)
  titleMap := make(map[string] bool, 10000)
  cf, err := os.Open(filterName)
  if err != nil { panic(err) }
  defer cf.Close()

  scanner := bufio.NewScanner(cf)
  for scanner.Scan() {
    titleMap[normalizeTitle(scanner.Text())] = true
  }

  if err := scanner.Err(); err != nil { panic(err)}

  for i := range input {
    if _, ok := titleMap[normalizeTitle(i.Title)]; ok {
      output <- i
      continue
    } else {
      for _, alias := range i.Aliases {
        if _, ok := titleMap[normalizeTitle(alias)]; ok {
          output <- i
          break
        }
      }
    }
  }
}
