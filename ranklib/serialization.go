package ranklib

import (
  "encoding/gob"
  "os"
  "io"
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

func ReadPreprocessedPages(fileName string, cp chan *PreprocessedPage) {
  defer close(cp)
  f, err := os.Open(fileName)
  if err != nil { panic(err) }
  defer f.Close()

  gobDecoder := gob.NewDecoder(f)
  for {
    var page PreprocessedPage
    err := gobDecoder.Decode(&page)
    if err == io.EOF {
      break
    } else if err != nil {
      panic(err)
    }
    cp <- &page
  }
}

func ReadPageRankedArticles(fileName string, cp chan *PageRankedArticle) {
  defer close(cp)
  f, err := os.Open(fileName)
  if err != nil { panic(err) }
  defer f.Close()

  gobDecoder := gob.NewDecoder(f)
  for {
    var page PageRankedArticle
    err := gobDecoder.Decode(&page)
    if err == io.EOF {
      break
    } else if err != nil {
      panic(err)
    }
    cp <- &page
  }
}

func WritePages(fileName string, numPages int, cp chan *Page, done chan bool) {
  outputFile, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
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
  outputFile, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
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

func WritePreprocessedPages(fileName string, cp chan *PreprocessedPage, done chan bool) {
  outputFile, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
  if err != nil { panic(err) }
  defer outputFile.Close()

  gobEncoder := gob.NewEncoder(outputFile)
  for rp := range cp {
    if rp != nil {
      gobEncoder.Encode(rp)
    } else {
      log.Printf("Null page while writing parsed pages")
    }
  }
  done <- true
}

func WritePageRankedArticles(fileName string, cp chan *PageRankedArticle, done chan bool) {
  outputFile, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
  if err != nil { panic(err) }
  defer outputFile.Close()

  gobEncoder := gob.NewEncoder(outputFile)
  for rp := range cp {
    if rp != nil {
      gobEncoder.Encode(rp)
    } else {
      log.Printf("Null page while writing parsed pages")
    }
  }
  done <- true
}
