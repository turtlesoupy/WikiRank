package ranklib

import (
  "os"
  "io"
  "log"
  "bytes"
  "regexp"
  "encoding/xml"
  "compress/gzip"
  "compress/bzip2"
  "strings"
)

type PreprocessedPage struct {
  Title string
  Id uint64
  RedirectTo string
  Infobox *ParsedInfobox
  LanguageTitles map[string] string
  TextLinks []ParsedTextLink
}

type ParsedInfobox struct {
  InfoboxType string
  Attributes map[string] string
}

type ParsedTextLink struct {
  ArticleTitle string
  Count uint32
}

type wikiParser struct {
  pe *pageElement
  cleanedText string
}

func PreprocessXML(inputFile string, outputFile string) error {
  pageInputChan := make(chan *pageElement, 20000)
  pageOutputChan := make(chan *PreprocessedPage, 20000)
  writeDoneChan := make(chan bool)

  go yieldPageElementsFromFile(inputFile, pageInputChan)
  go WritePreprocessedPages(outputFile, pageOutputChan, writeDoneChan)
  i := 0
  for pe := range pageInputChan {
    pageOutputChan <- parseXMLPageElement(pe)
    i++
    if i % 10000 == 0 {
      log.Printf("Reached page %d", i)
    }
  }
  close(pageOutputChan)
  <-writeDoneChan

  return nil
}


var removeComments = regexp.MustCompile(`<!--(.*?)-->`)
func parseXMLPageElement(pe *pageElement) *PreprocessedPage {
  parser := wikiParser {
    pe: pe,
    cleanedText: removeComments.ReplaceAllLiteralString(pe.Text, ""),
  }

  ret := &PreprocessedPage{
    Title: pe.Title,
    Id: pe.Id,
    RedirectTo: pe.Redirect.Title,
    Infobox: parser.parseInfobox(),
    LanguageTitles: parser.parseInterlanguagePageTitles(),
    TextLinks: parser.parseTextLinks(),
  }
  return ret
}

var interlanguageR = regexp.MustCompile(`\[\[(en|nl|de|fr|sv|it|es|ru):(.*?)\]\]`)
func (this *wikiParser) parseInterlanguagePageTitles() map[string] string {
  ret := make(map[string] string)
  submatches := interlanguageR.FindAllStringSubmatch(this.cleanedText, -1)
  for _, submatch := range submatches {
    language := submatch[1]
    title := submatch[2]
    ret[language] = title
  }

  return ret
}

var linkRegex = regexp.MustCompile(`\[\[(?:([^|\]]*)\|)?([^\]]+)\]\]`)
func (this *wikiParser) parseTextLinks() []ParsedTextLink {
  linkCounts := make(map[string] uint32, 100)
  submatches := linkRegex.FindAllStringSubmatch(this.cleanedText, -1)
  for _, submatch := range submatches {
    var realLinkName string
    if len(submatch[1]) == 0 {
      realLinkName = submatch[2]
    } else {
      realLinkName = cleanSectionRegex.FindString(submatch[1])
    }
    linkCounts[realLinkName]++
  }

  links := make([]ParsedTextLink, 0, len(linkCounts))
  for k,v := range linkCounts {
    links = append(links, ParsedTextLink{ArticleTitle: k, Count: v})
  }

  return links
}

type scanState uint8
const (
  _ = iota
  NONE = 1 << iota
  IN_KEY
  IN_VALUE
)


var infoboxStartR = regexp.MustCompile(`(?i){{ *Infobox *([^|}=]*)`)
func (this *wikiParser) parseInfobox() *ParsedInfobox {
  infoboxMatches := infoboxStartR.FindStringSubmatchIndex(this.cleanedText)
  if len(infoboxMatches) == 0 {
    return nil
  }

  subject := strings.TrimSpace(this.cleanedText[infoboxMatches[2]:infoboxMatches[3]])
  ret := &ParsedInfobox{InfoboxType: subject, Attributes: make(map[string] string)}

  // Don't try this at home
  var state scanState = NONE
  var buffer bytes.Buffer
  var lastKey string
  squiggleDepth := 2
  chevronDepth := 0
  squareDepth := 0

  flushScanState := func() {
    switch state {
    case IN_KEY:
      key := buffer.String()
      if len(key) > 0 && key[len(key) - 1] == '=' {
        key = key[0:len(key) - 1]
      }
      key = strings.TrimSpace(key)

      lastKey = key
    case IN_VALUE:
      value := buffer.String()
      if len(value) > 0 && value[len(value) -1] == '|' {
        value = value[:len(value) - 1]
      }
      if squiggleDepth == 0 && value[len(value)-1] == '}' && value[len(value) - 2] == '}' {
        value = value[:len(value) - 2]
      }
      value = strings.TrimSpace(value)

      if len(lastKey) > 0 && len(value) > 0 {
        ret.Attributes[lastKey] = value
      }

    }
    buffer.Reset()
  }

  for _, c := range this.cleanedText[infoboxMatches[3]:] {
    if state == IN_KEY || state == IN_VALUE {
      buffer.WriteRune(c)
    }

    if c == '{' {
      squiggleDepth++
    } else if c == '}' {
      squiggleDepth--
      if squiggleDepth == 0 {
        flushScanState()
        break
      }
    } else if c == '<' {
      chevronDepth++
    } else if c == '>' {
      chevronDepth--
    } else if c == '[' {
      squareDepth++
    } else if c == ']' {
      squareDepth--
    }

    if squiggleDepth == 2 && chevronDepth == 0 && squareDepth == 0 {
      if c == '|' {
        flushScanState()
        state = IN_KEY
      } else if c == '=' {
        flushScanState()
        state = IN_VALUE
      }
    }
  }

  return ret
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
