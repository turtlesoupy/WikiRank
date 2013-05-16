package ranklib

import (
  "encoding/xml"
  "log"
  "fmt"
  "regexp"
  "os"
)

var titleFilter, _ = regexp.Compile("^(File|Talk:|Special|Wikipedia|Wiktionary|User|User Talk:)")

type redirect struct {
  Title string `xml:"title, attr"`
}

type pageElement struct {
  Title string `xml:"title"`
  Redirect redirect `xml:"redirect"`
  Text string `xml:"revision>text"`
  Id uint64 `xml:"id"`
}

func (p *pageElement) String() string {
  return fmt.Sprintf("pageElement[title=%s,id=%d]", p.Title, p.Id)
}

func ReadFrom(fileName string) (err error) {
  xmlFile, err := os.Open(fileName)
  if(err != nil) { return }
  defer xmlFile.Close()

  titleIdMap := make(map[string]uint64, 10000000)
  decoder := xml.NewDecoder(xmlFile)
  log.Printf("Starting parse")
  for {
    token, err := decoder.Token()
    if err != nil {
      return err
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

        titleIdMap[p.Title] = p.Id
        log.Printf("Articles: %d", len(titleIdMap))
      case "mediawiki":
      default:
        decoder.Skip()
      }
    default:
    }
  }

  return
}
