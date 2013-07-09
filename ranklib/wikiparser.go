package ranklib

import (
  "log"
  "bytes"
  "regexp"
  "strings"
)

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

var infoboxStart = regexp.MustCompile(`(?i){{ *Infobox *([^|}=]*)`)
var removeComments = regexp.MustCompile(`<!--(.*?)-->`)
func parseInfobox(pe *pageElement) *parsedInfobox {
  text := removeComments.ReplaceAllLiteralString(pe.Text, "")

  infoboxMatches := infoboxStart.FindStringSubmatchIndex(text)
  if len(infoboxMatches) == 0 {
    return nil
  }

  subject := strings.TrimSpace(text[infoboxMatches[2]:infoboxMatches[3]])
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

  for _, c := range text[infoboxMatches[3]:] {
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

  log.Printf("%s %s", pe.Title, ret)
  return ret
}
