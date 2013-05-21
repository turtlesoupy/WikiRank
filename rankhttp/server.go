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


func index(w http.ResponseWriter, r *http.Request) {
  var scope = make(map[string]interface{})
  engine, err := loadTemplate("index.haml")
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }

  response := engine.Render(scope)
  w.Header().Set("Content-Length", strconv.Itoa(len(response)))
  w.Header().Set("Content-Type", "text/html; charset=utf-8")
  fmt.Fprintf(w, response)
}

func compareThings(autocompleteIndex *ranklib.Trie, w http.ResponseWriter, r *http.Request) {
  things, ok := r.URL.Query()["things[]"]
  if !ok || len(things) != 2 {
    http.Error(w, "Bad search", http.StatusBadRequest)
    return
  }
  log.Printf("Comparing %q", things)


  page1, ok := autocompleteIndex.GetEntry(ranklib.NormalizeSuggestionPrefix(things[0]))
  if !ok {
    http.Error(w, fmt.Sprintf("%s doesn't exist", things[0]), http.StatusBadRequest)
    return
  }

  page2, ok := autocompleteIndex.GetEntry(ranklib.NormalizeSuggestionPrefix(things[1]))
  if !ok {
    http.Error(w, fmt.Sprintf("%s doesn't exist", things[1]), http.StatusBadRequest)
    return
  }

  responseObject, err := json.Marshal(map[string]interface{}{
    "pages": []interface{}{page1, page2},
  })

  if err != nil {
    http.Error(w, "Unable to jsonify comparison", http.StatusInternalServerError)
    return
  }

  response := string(responseObject)
  w.Header().Set("Content-Length", strconv.Itoa(len(response)))
  w.Header().Set("Content-Type", "application/json; charset=utf-8")
  fmt.Fprintf(w, response)
}

func namedEntitySuggestions(autocompleteIndex *ranklib.Trie, w http.ResponseWriter, r *http.Request) {
  suggestions := map[string]interface{}{}
  search, ok := r.URL.Query()["q"]
  log.Printf("Search is %s from URL %s", search, r.URL)
  if !ok || len(search) != 1 {
    http.Error(w, "Bad search prefix", http.StatusBadRequest)
    return
  }

  suggestions["suggestions"], ok = autocompleteIndex.GetTopSuggestions(ranklib.NormalizeSuggestionPrefix(search[0]), 6)
  if !ok {
    http.Error(w, "Bad trie traversal", http.StatusInternalServerError)
    return
  }

  responseObject, err := json.Marshal(suggestions)
  if err != nil {
    http.Error(w, "Unable to jsonify suggestions", http.StatusInternalServerError)
    return
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

var namedEntityRegex = regexp.MustCompile("/named_entity_suggestions([?].*)?")
var compareRegex = regexp.MustCompile("/compare([?].*)?")
func Serve(autocompleteIndex *ranklib.Trie, port int) {
  log.Printf("Server: running on port %d", port)
  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    if namedEntityRegex.MatchString(r.URL.Path) {
      namedEntitySuggestions(autocompleteIndex, w, r)
    } else if compareRegex.MatchString(r.URL.Path) {
      compareThings(autocompleteIndex, w, r)
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
