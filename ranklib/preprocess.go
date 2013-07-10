package ranklib

import (
  "os"
  "io"
  "log"
  "encoding/xml"
  "compress/gzip"
  "compress/bzip2"
  "strings"
)


func PreprocessXML(inputFile string, outputFile string) error {
  pageInputChan := make(chan *pageElement, 20000)
  go yieldPageElementsFromFile(inputFile, pageInputChan)
  for pe := range pageInputChan {
    wikiParser := NewWikiParser(pe)
    wikiParser.ParseOut()
  }

  return nil
}

func yieldPageElementsFromFile(fileName string, cp chan *pageElement) {
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
