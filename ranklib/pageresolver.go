package ranklib

import (
  "os"
  "io"
  "log"
  "fmt"
  "sort"
  "bufio"
  "strings"
  "runtime"
  "encoding/json"
  "io/ioutil"
)

type CategoryResolver struct {
  CategoryName string
  PageResolver *PageResolver
}

type PageResolver struct {
  trie *Trie
  pages []RankedPage
  categoryResolvers []CategoryResolver
}

func reindexPages(pageMap map[uint64] *RankedPage, reindexMap map[uint64] int) {
  for _, page := range pageMap {
    id := page.Id
    page.Id = uint64(reindexMap[id])
    for j := range page.Influencers {
      page.Influencers[j].Id = uint64(reindexMap[page.Influencers[j].Id])
    }
  }
}

func CreatePageResolver(inputFile string, limit int) (*PageResolver, error) {
  n := ReadLength(inputFile)
  if n > limit {
    n = limit
  }


  rpchan := make(chan *RankedPage, 10000)
  go ReadRankedPages(inputFile, rpchan)

  pageById := make(map[uint64] *RankedPage, n)
  reindexMap := make(map[uint64] int, n)

  log.Printf("PageResolver: Reading ranked pages")
  i := 0
  for page := range rpchan {
    pageById[page.Id] = page
    if _, ok := reindexMap[page.Id]; ok {
      panic(fmt.Sprintf("Found existing page in reindex map %s", page))
    }
    reindexMap[page.Id] = i + 1
    i++
    if i >= n {
      break
    }
  }

  log.Printf("PageResolver: Reindexing pages")
  reindexPages(pageById, reindexMap)

  log.Printf("PageResolver: Creating Trie")
  trie := NewTrie()
  rawPageList := make([]RankedPage, n, n)
  i = 0
  for _, startPage := range pageById {
    insertionValue := TrieValue{Id: startPage.Id, Rank: startPage.Rank}
    trie.AddEntry(normalizeTitle(startPage.Title), insertionValue)
    for _, alias := range startPage.Aliases {
      trie.AddEntry(normalizeTitle(alias), insertionValue)
    }

    rawPageList[startPage.Id-1] = *startPage

    i++
    if i % 100000 == 0 && i > 0 {
      log.Printf("Inserted page #%d", i)
      sug, _ := trie.GetTopSuggestions("", 10)
      log.Printf("Suggestions 1: %q", sug)
    }
  }

  log.Printf("PageResolver: Garbage collecting")
  reindexMap = nil
  pageById = nil
  runtime.GC()
  log.Printf("PageResolver: All Done!")

  return &PageResolver{
    trie: trie,
    pages: rawPageList,
  }, nil
}

func pageResolverFromList(pages []RankedPage) (*PageResolver) {
  trie := NewTrie()
  for i, page := range pages {
    trie.AddEntry(normalizeTitle(page.Title), TrieValue{Id: uint64(i+1), Rank: page.Rank})
  }

  return &PageResolver {
    trie: trie,
    pages: pages,
  }
}

func (this *PageResolver) AddCategoryFromFile(categoryName string, inputFile string) (err error) {
  f, err := os.Open(inputFile)
  if err != nil {
    return
  }

  pages := make([]RankedPage, 0)
  reader := bufio.NewReader(f)
  seenSet := make(map[uint64] bool)
  for {
    rawLine, err := reader.ReadString('\n')
    if err == io.EOF {
      break
    } else if err != nil {
      return err
    }

    pageName := strings.TrimSpace(rawLine)

    if page, ok := this.PageByTitle(pageName); ok && !seenSet[page.Id] {
      if page.Title == "Case Closed" || page.Title == "1987 in film" || page.Title == "George A. Romero" {
        log.Printf("Found %s from %s", page.Title, pageName)
      }
      pages = append(pages, *page)
      seenSet[page.Id] = true
    } else {
      continue
    }
  }

  sort.Sort(RankedPageList(pages))

  resolver := pageResolverFromList(pages)

  this.categoryResolvers = append(this.categoryResolvers, CategoryResolver{
    CategoryName: categoryName,
    PageResolver: resolver,
  })

  return nil
}

func (this *PageResolver) GetCategories() []CategoryResolver {
  return this.categoryResolvers
}

func (this *PageResolver) OrderedPageRange(start int, end int) []RankedPage {
  if start > len(this.pages) {
    return make([]RankedPage, 0)
  } else if end > len(this.pages) {
    end = len(this.pages)
  }

  return this.pages[start:end]
}

func (this *PageResolver) PrefixSuggestions(prefix string, n int) []RankedPage {
  suggestions, _ := this.trie.GetTopSuggestions(normalizeTitle(prefix), n)
  ret := make([]RankedPage, len(suggestions))
  for i, suggestion := range suggestions {
    sp, _ := this.PageById(suggestion.Id)
    ret[i] = *sp
  }

  return ret
}

func (this *PageResolver) PageByTitle(title string) (*RankedPage, bool) {
  trieValue, ok := this.trie.GetEntry(normalizeTitle(title))
  if !ok { return nil, false }
  return this.PageById(trieValue.Id)
}

func (this *PageResolver) PageById(id uint64) (*RankedPage, bool) {
  if id -1 > uint64(len(this.pages)) {
    return nil, false
  }

  p := &this.pages[int(id) - 1]
  return p, true
}

func (this *PageResolver) DumpPageList(outputFile string) (err error) {
  var data []byte
  if data, err = json.MarshalIndent(this.pages, "", "\t"); err != nil {
    return
  }

  return ioutil.WriteFile(outputFile, data, 0755)
}

func normalizeTitle(title string) string {
  return strings.ToUpper(title)
}
