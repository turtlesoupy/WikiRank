package ranklib

import (
  "encoding/xml"
  "log"
  "fmt"
  "regexp"
  "os"
  "encoding/gob"
)

var titleFilter = regexp.MustCompile("^(File|Talk:|Special|Wikipedia|Wiktionary|User|User Talk:)")

type redirect struct {
  Title string `xml:"title, attr"`
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
  Links []uint64
}

func (p *pageElement) String() string {
  return fmt.Sprintf("pageElement[title=%s,id=%d]", p.Title, p.Id)
}

func yieldPages(fileName string, cp chan *pageElement) {
  xmlFile, err := os.Open(fileName)
  if(err != nil) { panic(err) }

  defer xmlFile.Close()

  log.Printf("Starting parse")
  decoder := xml.NewDecoder(xmlFile)
  for {
    token, err := decoder.Token()
    if err != nil {
      panic(err)
    } else if token == nil {
      break
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

  close(cp)
}


var linkRegex = regexp.MustCompile(`\[\[(?:([^|\]]*)\|)?([^\]]+)\]\]`)
var cleanSectionRegex = regexp.MustCompile(`^[^#]*`)
func newPage(pe *pageElement, titleIdMap map[string]uint64) *Page {
  p := Page{Title: pe.Title, Id: pe.Id}
  submatches := linkRegex.FindAllStringSubmatch(pe.Text, -1)
  p.Links = make([]uint64, 0, len(submatches))
  for _, submatch := range submatches {
    var dirtyLinkName string
    if len(submatch[1]) == 0 {
      dirtyLinkName = submatch[2]
    } else {
      dirtyLinkName = cleanSectionRegex.FindString(submatch[1])
    }
    if linkId, ok := titleIdMap[dirtyLinkName]; ok {
      p.Links = append(p.Links, linkId)
    }
  }

  return &p
}


func ReadFrom(fileName string, outputName string) (err error) {
  outputFile, err := os.OpenFile(outputName, os.O_WRONLY|os.O_CREATE, 0600)
  if err != nil { panic(err) }
  defer outputFile.Close()

  pageChan := make(chan *pageElement, 1000)
  go yieldPages(fileName, pageChan)

  titleIdMap := make(map[string]uint64, 10000000)
  numPages := 0
  log.Printf("Starting title pass")
  // First pass: fill title id map
  for page := range pageChan {
    titleIdMap[page.Title] = page.Id
    numPages++

    if numPages % 10000 == 0 {
      log.Printf("Page #%d", numPages)
    }
  }

  log.Printf("Done title pass, starting write pass")

  pageChan = make(chan *pageElement, 1000)
  go yieldPages(fileName, pageChan)
  gobEncoder := gob.NewEncoder(outputFile)
  gobEncoder.Encode(numPages)

  i := 0
  for page := range pageChan {
    p := newPage(page, titleIdMap)
    gobEncoder.Encode(p)
    i++
    if i % 10000 == 0 {
      log.Printf("Page #%d", numPages)
    }
  }

  log.Printf("Done write pass")

  return
}
