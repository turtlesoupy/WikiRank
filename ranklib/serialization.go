package ranklib

import (
  "encoding/gob"
  "os"
  "log"
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
    err := gobDecoder.Decode(&page)
    if err != nil {
      panic(err)
    }
    cp <- &page
  }
}

func ReadRankedPages(fileName string, cp chan *RankedPage) {
  defer close(cp)
  f, err := os.Open(fileName)
  if err != nil { panic(err) }
  defer f.Close()

  gobDecoder := gob.NewDecoder(f)
  var length int
  gobDecoder.Decode(&length)

  log.Printf("Reading %d ranked pages", length)
  for i := 0; i < length-1; i++ {
    var page RankedPage
    err := gobDecoder.Decode(&page)
    if err != nil {
      panic(err)
    } else{
      cp <- &page
    }
  }
}

func WritePages(fileName string, numPages int, cp chan *Page, done chan bool) {
  outputFile, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0600)
  if err != nil { panic(err) }
  defer outputFile.Close()

  gobEncoder := gob.NewEncoder(outputFile)
  gobEncoder.Encode(numPages)

  for page := range cp {
    if(page != nil) {
      err = gobEncoder.Encode(page)
      if err != nil { panic(err) }
    } else {
      log.Printf("Null page while writing pages")
    }
  }
  done <- true
}

func WriteRankedPages(fileName string, numPages int, cp chan *RankedPage, done chan bool) {
  outputFile, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0600)
  if err != nil { panic(err) }
  defer outputFile.Close()

  gobEncoder := gob.NewEncoder(outputFile)
  gobEncoder.Encode(numPages)

  for rp := range cp {
    if rp != nil {
      gobEncoder.Encode(rp)
    } else {
      log.Printf("Null page while writing ranked pages")
    }
  }
  done <- true
}
