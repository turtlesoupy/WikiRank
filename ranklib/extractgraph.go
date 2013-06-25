package ranklib

import (
  "os"
  "io"
  "log"
  "fmt"
  "bufio"
  "regexp"
  "strings"
  "encoding/xml"
  "compress/bzip2"
  "compress/gzip"
)

const (
  readerBufferSize = 32 * 1024 * 1024
)

var titleFilter = regexp.MustCompile("^(File|Talk|Special|Wikipedia|Wiktionary|User|User Talk|Category|Portal|Template|Mediawiki):")

type redirect struct {
  Title string `xml:"title,attr"`
}

type pageElement struct {
  Title string `xml:"title"`
  Redirect redirect `xml:"redirect"`
  Text string `xml:"revision>text"`
  Id uint64 `xml:"id"`
}

type Page struct {
  Title string
  Id uint64
  Coordinate Coordinate
  RedirectToId uint64
  Links []uint64
}

func (page *Page) IsRedirect() bool {
  return page.RedirectToId != 0
}

func (p *pageElement) String() string {
  return fmt.Sprintf("pageElement[title=%s,id=%d]", p.Title, p.Id)
}

func yieldPageElements(fileName string, cp chan *pageElement) {
  xmlFile, err := os.Open(fileName)
  if(err != nil) { panic(err) }

  var xmlReader io.Reader = bufio.NewReaderSize(xmlFile, readerBufferSize)
  if strings.HasSuffix(fileName, ".bz2") {
    log.Printf("Assuming bzip2 compressed dump")
    xmlReader = bzip2.NewReader(xmlReader)
  } else if strings.HasSuffix(fileName, ".gz") {
    log.Printf("Assuming gzip compressed dump")
    xmlReader, err = gzip.NewReader(xmlReader)
  } else {
    log.Printf("Assuming uncompressed dump")
  }


  if(err != nil) { panic(err) }

  defer xmlFile.Close()
  defer close(cp)

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
        cp <- &p
      case "mediawiki":
      default:
        decoder.Skip()
      }
    default:
    }
  }
}

func newPage(pe *pageElement, titleIdMap map[string]uint64) *Page {
  p := Page{Title: pe.Title, Id: pe.Id, RedirectToId: 0}
  if len(pe.Redirect.Title) > 0 {
    if redirectId, ok := titleIdMap[cleanSectionRegex.FindString(pe.Redirect.Title)]; ok {
      p.RedirectToId = redirectId
    }
  }

  dedupeLinks := make(map[uint64] bool)
  submatches := linkRegex.FindAllStringSubmatch(pe.Text, -1)
  p.Links = make([]uint64, 0, len(submatches))
  for _, submatch := range submatches {
    var dirtyLinkName string
    if len(submatch[1]) == 0 {
      dirtyLinkName = submatch[2]
    } else {
      dirtyLinkName = cleanSectionRegex.FindString(submatch[1])
    }
    if linkId, ok := titleIdMap[dirtyLinkName]; ok && linkId != p.Id {
      if _, alreadySeen := dedupeLinks[linkId]; !alreadySeen {
        dedupeLinks[linkId] = true
        p.Links = append(p.Links, linkId)
      }
    }
  }

  if c, ok := coordinatesFromWikiText(pe); ok {
    p.Coordinate = c
  }

  return &p
}


func ReadFrom(fileName string, outputName string) (err error) {
  pageInputChan := make(chan *pageElement, 1000)
  go yieldPageElements(fileName, pageInputChan)

  titleIdMap := make(map[string]uint64, 12000000)
  numPages := 0
  log.Printf("Starting title pass")
  // First pass: fill title id map
  for page := range pageInputChan {
    titleIdMap[page.Title] = page.Id
    numPages++
    newPage(page, titleIdMap)

    if numPages % 10000 == 0 {
      log.Printf("Page #%d", numPages)
    }
  }

  log.Printf("Done title pass, starting write pass")

  pageInputChan = make(chan *pageElement, 1000)
  pageOutputChan := make(chan *Page, 1000)
  writeDoneChan := make(chan bool)
  go yieldPageElements(fileName, pageInputChan)
  go WritePages(outputName, numPages, pageOutputChan, writeDoneChan)
  i := 0
  for page := range pageInputChan {
    pageOutputChan <- newPage(page, titleIdMap)
    i++
    if i % 10000 == 0 {
      log.Printf("Page #%d", i)
    }
  }
  close(pageOutputChan)
  <-writeDoneChan

  log.Printf("Done write pass")

  return
}
