package rankhttp

import (
  "io"
  "os"
  "bytes"
  "fmt"
  "log"
  "path"
  "regexp"
  "strconv"
  "net/http"
  "encoding/json"
  "github.com/cosbynator/wikirank/ranklib"
  "github.com/realistschuckle/gohaml"
)

const (
  publicDir = "/home/tdimson/go/src/github.com/cosbynator/wikirank/rankhttp/public"
  templateDir = "/home/tdimson/go/src/github.com/cosbynator/wikirank/rankhttp/templates"
)

func loadTemplate(templateName string) (engine *gohaml.Engine, err error) {
  templateLocation := fmt.Sprintf("%s/%s", templateDir, templateName)
  f, err := os.Open(templateLocation)
  if err != nil {
    return
  }
  defer f.Close()

  var bb bytes.Buffer
  if _, err = io.Copy(&bb, f); err != nil {
    return
  }

  return gohaml.NewEngine(bb.String())
}

var namedEntityRegex = regexp.MustCompile("/named_entity_suggestions([?].*)?")

func index(w http.ResponseWriter, r *http.Request) {
  var scope = make(map[string]interface{})
  engine, err := loadTemplate("index.haml")
  if err != nil {
    // TODO: something
  }

  response := engine.Render(scope)
  w.Header().Set("Content-Length", strconv.Itoa(len(response)))
  w.Header().Set("Content-Type", "text/html; charset=utf-8")
  fmt.Fprintf(w, response)
}


func namedEntitySuggestions(autocompleteIndex *ranklib.Trie, w http.ResponseWriter, r *http.Request) {
  suggestions := map[string]interface{}{}
  search, ok := r.URL.Query()["q"]
  log.Printf("Search is %s from URL %s", search, r.URL)
  if !ok || len(search) != 1 {
    //TODO: something
    http.Error(w, "Bad search prefix", http.StatusBadRequest)
    return
  }

  suggestions["suggestions"], ok = autocompleteIndex.GetTopSuggestions(ranklib.NormalizeSuggestionPrefix(search[0]), 6)
  if !ok {
    //TODO: something
  }

  responseObject, err := json.Marshal(suggestions)
  if err != nil {
    // TODO: something
  }
  response := string(responseObject)
  w.Header().Set("Content-Length", strconv.Itoa(len(response)))
  w.Header().Set("Content-Type", "application/json; charset=utf-8")
  fmt.Fprintf(w, response)
}


func setupStatics(rootStatics []string, exposedStaticDirs []string) {
  for _, filename := range rootStatics {
    fileLocation := path.Join(publicDir, filename)
    http.HandleFunc(fmt.Sprintf("/%s", filename), func(w http.ResponseWriter, r *http.Request) {
      http.ServeFile(w, r, fileLocation)
    })
  }

  for _, dir := range exposedStaticDirs {
    log.Printf("Serving %s from %s", fmt.Sprintf("/%s/", dir), http.Dir(path.Join(publicDir, dir)))
    http.Handle(fmt.Sprintf("/%s/", dir), http.StripPrefix(fmt.Sprintf("/%s", dir), http.FileServer(http.Dir(path.Join(publicDir, dir)))))
  }
  log.Printf("Static routes are set")
}

func Serve(autocompleteIndex *ranklib.Trie, port int) {
  log.Printf("Server: running on port %d", port)
  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    if namedEntityRegex.MatchString(r.URL.Path) {
      namedEntitySuggestions(autocompleteIndex, w, r)
    } else {
      index(w, r)
    }
  })
  setupStatics(
    []string{"robots.txt"},
    []string {"css", "js", "img"},
  )
  log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
