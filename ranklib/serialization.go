package ranklib

import (
  "encoding/gob"
  "os"
)

func ReadLength(fileName string) int {
  f, err := os.Open(fileName)
  if err != nil { panic(err) }
  defer f.Close()
  gobDecoder := gob.NewDecoder(f)
  var length int
  gobDecoder.Decode(&length)
  return length
}

func ReadPages(fileName string, cp chan *Page) {
  defer close(cp)
  f, err := os.Open(fileName)
  if err != nil { panic(err) }
  defer f.Close()

  gobDecoder := gob.NewDecoder(f)
  var length int
  gobDecoder.Decode(&length)

  for i := 0; i < length; i++ {
    var page Page
    gobDecoder.Decode(&page)
    cp <- &page
  }
}

func WritePages(fileName string, numPages int, cp chan *Page) {
  outputFile, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0600)
  if err != nil { panic(err) }
  defer outputFile.Close()

  gobEncoder := gob.NewEncoder(outputFile)
  gobEncoder.Encode(numPages)

  for page := range cp {
    gobEncoder.Encode(page)
  }
}

func WriteRankedPages(fileName string, numPages int, cp chan *RankedPage) {
  outputFile, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0600)
  if err != nil { panic(err) }
  defer outputFile.Close()

  gobEncoder := gob.NewEncoder(outputFile)
  gobEncoder.Encode(numPages)

  for rp := range cp {
    gobEncoder.Encode(rp)
  }
}
