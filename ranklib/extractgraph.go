package ranklib

import (
  "os"
  "io"
  "log"
  "fmt"
  "sort"
  "regexp"
  "strings"
  "strconv"
  "encoding/xml"
  "compress/bzip2"
  "compress/gzip"
)

const (
  readerBufferSize = 32 * 1024 * 1024
)


var titleFilter = regexp.MustCompile("^(File|Talk|Special|Wikipedia|Wiktionary|User|User Talk|Category|Portal|Template|Mediawiki|Help):")

type redirect struct {
  Title string `xml:"title,attr"`
}

const (
  hasCoordinate = iota << 1
)

type Page struct {
  Title string
  Id uint64
  Coordinate Coordinate
  Aliases []string
  Links []Link
  Flags uint32
  ReleaseYear int32
}

type Link struct {
  PageId uint64
  Count uint32
}

func (page *Page) HasCoordinate() bool {
  return page.Flags & hasCoordinate > 0
}

func (p *pageElement) String() string {
  return fmt.Sprintf("pageElement[title=%s,id=%d]", p.Title, p.Id)
}

func yieldPageElements(fileName string, cp chan *pageElement) {
  xmlFile, err := os.Open(fileName)
  if(err != nil) { panic(err) }

  var xmlReader io.Reader
  if strings.HasSuffix(fileName, ".bz2") {
    log.Printf("Assuming bzip2 compressed dump")
    xmlReader = bzip2.NewReader(xmlFile)
  } else if strings.HasSuffix(fileName, ".gz") {
    log.Printf("Assuming gzip compressed dump")
    xmlReader, err = gzip.NewReader(xmlFile)
  } else {
    xmlReader = xmlFile
    log.Printf("Assuming uncompressed dump")
  }

  if err != nil { panic(err) }

  defer xmlFile.Close()
  defer close(cp)

  pageCounter := 0
  log.Printf("Starting parse")
  decoder := xml.NewDecoder(xmlReader)
  for {
    token, err := decoder.Token()
    if err == io.EOF || token == nil {
      log.Printf("EOF")
      break
    } else if err != nil {
      panic(err)
    }


    switch e := token.(type) {
    case xml.StartElement:
      switch e.Name.Local {
      case "page":
        var p pageElement
        decoder.DecodeElement(&p, &e)
        if titleFilter.MatchString(p.Title) {
          continue
        }
        pageCounter++
        if pageCounter % 10000 == 0 {
          log.Printf("Reached page %d", pageCounter)
        }
        cp <- &p
      case "mediawiki":
      default:
        decoder.Skip()
      }
    default:
    }
  }
}

var filmDateInfoboxR = regexp.MustCompile(fmt.Sprintf(`(?i)[|] *(released) *= (.*?)`))
var yearR = regexp.MustCompile(`\d\d\d\d`)
func releaseYearFromWikiText(pe *pageElement) (int32, bool) {
  releaseString := extractFromInfobox(pe.Text, filmDateInfoboxR)
  if releaseString == "" {
    return -1, false
  }

  releaseYearString := yearR.FindString(releaseString)
  releaseYear, err := strconv.Atoi(releaseYearString)
  if err != nil {
    log.Printf("Error converting release year for '%s' in '%q'", pe.Title, releaseYear)
    return -1, false
  }

  return int32(releaseYear), true
}

func ReadFrom(fileName string, outputName string) (err error) {
  pages := make([]*Page, 0, 5000000)
  pageTitleMap := make(map[string] *Page, 12000000)

  log.Printf("Starting pass 1: pages")
  pageInputChan := make(chan *pageElement, 20000)
  go yieldPageElements(fileName, pageInputChan)
  for pe := range pageInputChan {
    if len(pe.Redirect.Title) == 0 {
      p := &Page {Title: pe.Title, Id: pe.Id, Links: make([]Link, 0, 10)}

      if c, ok := coordinatesFromWikiText(pe); ok {
        p.Coordinate = c
        p.Flags |= hasCoordinate
      }

      if c, ok := releaseYearFromWikiText(pe); ok {
        p.ReleaseYear = c
      } else {
        p.ReleaseYear = -1
      }

      pages = append(pages, p)
      pageTitleMap[pe.Title] = p
    }
  }

  log.Printf("Starting pass 2: redirects")
  pageInputChan = make(chan *pageElement, 1000)
  go yieldPageElements(fileName, pageInputChan)
  for pe := range pageInputChan {
    if len(pe.Redirect.Title) > 0 {
      redirectedTitle := cleanSectionRegex.FindString(pe.Redirect.Title)
      if redirectPage, ok := pageTitleMap[redirectedTitle]; ok {
        redirectPage.Aliases = append(redirectPage.Aliases, pe.Title)
        pageTitleMap[pe.Title] = redirectPage
      } else if !titleFilter.MatchString(redirectedTitle) {
        log.Printf("Unresolvable redirect: '%s' -> '%s' (cleaned '%s')", pe.Title, pe.Redirect.Title, redirectedTitle)
      }
    }
  }

  log.Printf("Starting pass 3: links")
  pageInputChan = make(chan *pageElement, 1000)
  go yieldPageElements(fileName, pageInputChan)
  for pe := range pageInputChan {
    fromPage, ok := pageTitleMap[pe.Title]
    if !ok {
      if pe.Redirect.Title != "" {
        log.Printf("Warning: page '%s' not in title map", pe.Title) // Errors only really matter in non-redirects
      }
      continue
    }

    linkCounts := make(map[uint64] uint32)
    submatches := linkRegex.FindAllStringSubmatch(pe.Text, -1)
    for _, submatch := range submatches {
      var dirtyLinkName string
      if len(submatch[1]) == 0 {
        dirtyLinkName = submatch[2]
      } else {
        dirtyLinkName = cleanSectionRegex.FindString(submatch[1])
      }

      if toPage, ok := pageTitleMap[dirtyLinkName]; ok && toPage.Id != fromPage.Id {
        linkCounts[toPage.Id]++
      }
    }

    for linkedId, count := range linkCounts {
      fromPage.Links = append(fromPage.Links, Link{PageId: linkedId, Count: count})
    }
  }

  log.Printf("Starting writing...")
  pageOutputChan := make(chan *Page, 1000)
  writeDoneChan := make(chan bool)
  go WritePages(outputName, len(pages), pageOutputChan, writeDoneChan)
  for _, p := range pages {
    sort.Strings(p.Aliases)
    pageOutputChan <- p
  }
  close(pageOutputChan)
  <-writeDoneChan

  log.Printf("Done write pass")

  return
}
