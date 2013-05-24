package ranklib

import (
  "log"
  "runtime"
  "strings"
)

const (
  maxRedirects = 100
)

type PageResolver struct {
  trie *Trie
  pages map[uint64] RankedPage
}

func CreatePageResolver(inputFile string, limit int) (*PageResolver, error) {
  n := ReadLength(inputFile)
  allowDanglingRedirects := false
  if n > limit {
    n = limit
    allowDanglingRedirects = true
  }

  pageList := make([]RankedPage, n)
  rpchan := make(chan *RankedPage, 10000)
  go ReadRankedPages(inputFile, rpchan)

  log.Printf("PageResolver: Reading ranked pages")
  countNonRedirects := 0
  i := 0
  for page := range rpchan {
    if !page.IsRedirect() {
      countNonRedirects++
    }
    pageList[i] = *page
    i++
    if i >= n {
      break
    }
  }
  log.Printf("PageResolver: Found %d/%d non-redirects", countNonRedirects, n)
  log.Printf("PageResolver: Creating trie")
  trie, err := createTrie(pageList, allowDanglingRedirects)
  if err != nil { return nil, err }
  runtime.GC()

  log.Printf("PageResolver: Creating non-redirect map")
  nonRedirects := make(map[uint64] RankedPage, countNonRedirects)
  for i := 0; i < n; i++ {
    page := pageList[i]
    if !page.IsRedirect() {
      nonRedirects[page.Id] = page
    }
  }
  log.Printf("PageResolver: All done...")

  pageList = nil // Force GC ?
  runtime.GC()

  return &PageResolver {
    trie: trie,
    pages: nonRedirects,
  }, nil
}

func (this *PageResolver) PrefixSuggestions(prefix string, n int) []RankedPage {
  suggestions, _ := this.trie.GetTopSuggestions(normalizeTitle(prefix), n)
  ret := make([]RankedPage, len(suggestions))
  for i, suggestion := range suggestions {
    ret[i] = this.pages[suggestion.Id]
  }

  return ret
}

func (this *PageResolver) PageByTitle(title string) (*RankedPage, bool) {
  trieValue, ok := this.trie.GetEntry(normalizeTitle(title))
  if !ok { return nil, false }
  page, ok := this.pages[trieValue.Id]
  if !ok { return nil, false }
  return &page, true
}

func (this *PageResolver) PageById(id uint64) (*RankedPage, bool) {
  p, ok := this.pages[id]
  return &p, ok
}

func normalizeTitle(title string) string {
  return strings.ToUpper(title)
}

func createTrie(pages []RankedPage, allowDanglingRedirects bool) (trie *Trie, err error) {
  trie = NewTrie()

  log.Printf("PageResolver: Trie: Building page map")
  pageById := make(map[uint64] *RankedPage, len(pages))
  for i := 0; i < len(pages); i++ {
    pageById[pages[i].Id] = &pages[i]
  }

  log.Printf("PageResolver: Trie: Inserting %d pages", len(pages))
  for i := 0; i < len(pages); i++ {
    startPage := &pages[i]
    insertionValue := TrieValue{Id: startPage.Id, Rank: startPage.Rank}
    if startPage.IsRedirect() {
      p := startPage
      redirectCount := 0
      for p.IsRedirect() && redirectCount < maxRedirects {
        var ok bool
        p, ok = pageById[p.RedirectToId]
        if !ok {
          if allowDanglingRedirects {
            break
          } else {
            panic("We got a non-okay redirect")
          }
        }
        redirectCount++
      }

      if p == nil {
        continue
      } else if redirectCount == maxRedirects {
        log.Printf("Infinite redirect loop for %s", startPage)
        continue;
      }

      insertionValue.Id = p.Id // Old rank, new id
    }

    trie.AddEntry(normalizeTitle(startPage.Title), insertionValue)

    if i % 100000 == 0 && i > 0 {
      log.Printf("Inserted page #%d", i)
      sug, _ := trie.GetTopSuggestions("", 10)
      log.Printf("Suggestions 1: %q", sug)
    }
  }

  return trie, nil
}

