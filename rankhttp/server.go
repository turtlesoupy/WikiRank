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
  "html/template"
  "github.com/cosbynator/wikirank/ranklib"
  "github.com/realistschuckle/gohaml"
)

const (
  publicDir = "/home/tdimson/go/src/github.com/cosbynator/wikirank/rankhttp/public"
  templateDir = "/home/tdimson/go/src/github.com/cosbynator/wikirank/rankhttp/templates"
)

type RankedPageWithInfluencers struct {
  Page ranklib.RankedPage
  Influencers []RankedPageInfluencer
}

type RankedPageInfluencer struct {
  Page ranklib.RankedPage
  Influence float32
}

func fetchPageWithInfluencers(title string, pageResolver *ranklib.PageResolver) (*RankedPageWithInfluencers, bool) {
  page, ok := pageResolver.PageByTitle(title)
  if !ok { return nil, false }

  influencers := make([]RankedPageInfluencer, 0, len(page.Influencers))
  for _, influencer := range page.Influencers {
    if iPage, iOk := pageResolver.PageById(influencer.Id); iOk {
      influencers = append(influencers, RankedPageInfluencer{Page: *iPage, Influence: influencer.Influence})
    }
  }

  return &RankedPageWithInfluencers{
    Page: *page,
    Influencers: influencers,
  }, true
}

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
  t, err := template.ParseFiles(fmt.Sprintf("%s/index.gotemplate.html", templateDir))
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }

  var i interface{}
  t.Execute(w, i)
}

func things(pageResolver *ranklib.PageResolver, w http.ResponseWriter, r *http.Request) {
  things, ok := r.URL.Query()["things[]"]
  if !ok || len(things) != 2 {
    http.Error(w, "Bad search", http.StatusBadRequest)
    return
  }

  ret := make([]interface{}, 0, len(things))
  for _, name := range things {
    page, ok := fetchPageWithInfluencers(name, pageResolver)
    if !ok {
      http.Error(w, fmt.Sprintf("%s doesn't exist", name), http.StatusBadRequest)
      return
    }
    ret = append(ret, page)
  }

  responseObject, err := json.Marshal(ret)

  if err != nil {
    http.Error(w, "Unable to jsonify all things", http.StatusInternalServerError)
    return
  }

  response := string(responseObject)
  w.Header().Set("Content-Length", strconv.Itoa(len(response)))
  w.Header().Set("Content-Type", "application/json; charset=utf-8")
  fmt.Fprintf(w, response)
}

func namedEntitySuggestions(pageResolver *ranklib.PageResolver, w http.ResponseWriter, r *http.Request) {
  suggestions := map[string]interface{}{}
  search, ok := r.URL.Query()["q"]
  log.Printf("Search is %s from URL %s", search, r.URL)
  if !ok || len(search) != 1 {
    http.Error(w, "Bad search prefix", http.StatusBadRequest)
    return
  }

  suggestions["suggestions"] = pageResolver.PrefixSuggestions(search[0], 10)

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
var compareRegex = regexp.MustCompile("/things/?([?].*)?")
func Serve(pageResolver *ranklib.PageResolver, port int) {
  log.Printf("Server: running on port %d", port)
  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    log.Printf("Request: %s", r.URL)
    if namedEntityRegex.MatchString(r.URL.Path) {
      namedEntitySuggestions(pageResolver, w, r)
    } else if compareRegex.MatchString(r.URL.Path) {
      things(pageResolver, w, r)
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
