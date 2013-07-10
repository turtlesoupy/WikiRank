package ranklib

import (
  "log"
  "bytes"
  "regexp"
  "strings"
)

type WikiParser struct {
  pageTitle string
  cleanedText string
}

var removeComments = regexp.MustCompile(`<!--(.*?)-->`)
func NewWikiParser(pe *pageElement) *WikiParser {
  p := &WikiParser{}
  p.pageTitle = pe.Title
  p.cleanedText = removeComments.ReplaceAllLiteralString(pe.Text, "")
  return p
}

var linkRegex = regexp.MustCompile(`\[\[(?:([^|\]]*)\|)?([^\]]+)\]\]`)
func (this *WikiParser) ParseTextLinks() []string {
  links := make([]string, 100)
  submatches := linkRegex.FindAllStringSubmatch(this.cleanedText, -1)
  for _, submatch := range submatches {
    var realLinkName string
    if len(submatch[1]) == 0 {
      realLinkName = submatch[2]
    } else {
      realLinkName = cleanSectionRegex.FindString(submatch[1])
    }
    links = append(links, realLinkName)
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

type parsedInfobox struct {
  InfoboxType string
  Attributes map[string] string
}

var infoboxStartR = regexp.MustCompile(`(?i){{ *Infobox *([^|}=]*)`)
func (this *WikiParser) ParseInfobox() *parsedInfobox {
  infoboxMatches := infoboxStartR.FindStringSubmatchIndex(this.cleanedText)
  if len(infoboxMatches) == 0 {
    return nil
  }

  subject := strings.TrimSpace(this.cleanedText[infoboxMatches[2]:infoboxMatches[3]])
  ret := &parsedInfobox{InfoboxType: subject, Attributes: make(map[string] string)}

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

  log.Printf("%s !! %s", this.pageTitle, ret)
  return ret
}
