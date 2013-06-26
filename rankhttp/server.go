package rankhttp

import (
  "fmt"
  "log"
  "path"
  "regexp"
  "strconv"
  "net/http"
  "encoding/json"
  "html/template"
  "github.com/cosbynator/wikirank/ranklib"
)

const (
  publicDir = "/home/tdimson/go/src/github.com/cosbynator/wikirank/rankhttp/public"
  templateDir = "/home/tdimson/go/src/github.com/cosbynator/wikirank/rankhttp/templates"
)

type IndexData struct {
  Categories []IndexCategory
}

type IndexCategory struct {
  CategoryName string
  TopTen []IndexCategoryPage
}

type IndexCategoryPage struct {
  Page ranklib.RankedPage
  PercentDecrease float64
}

var helpers template.FuncMap = template.FuncMap {
  "indexToOrdinal": func(index int) int {return index+1 },
  "wikiLink": func(title string) template.HTML {
    return template.HTML(fmt.Sprintf("<a href='http://en.wikipedia.org/wiki/%s'>%s</a>", title, title))
  },
  "notFirst": func(index int) bool {
    return index > 0
  },
}

func index(pageResolver *ranklib.PageResolver, w http.ResponseWriter, r *http.Request) {
  templateName := fmt.Sprintf("%s/index.gotemplate.html", templateDir)
  t := template.New("index")
  t.Funcs(helpers)
  _, err := t.ParseFiles(templateName)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }

  log.Printf("%s", t)

  categories := make([]IndexCategory, 0, 1000)
  for _, categoryResolver := range pageResolver.GetCategories() {
    pages := make([]IndexCategoryPage, 0, 1000)
    for i, page := range categoryResolver.PageResolver.OrderedPageRange(0,1000) {
      percentDecrease := 0.0
      if i > 0 {
        percentDecrease = 100 * float64(page.Rank / pages[0].Page.Rank)
      }

      pages = append(pages, IndexCategoryPage{
        Page: page,
        PercentDecrease: percentDecrease,
      })
    }

    categories = append(categories, IndexCategory{
      CategoryName: categoryResolver.CategoryName,
      TopTen: pages,
    })
  }

  templateData := IndexData{categories}
  err = t.ExecuteTemplate(w, "index.gotemplate.html", templateData)
  if err != nil {
    log.Printf(err.Error())
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }

}

func things(pageResolver *ranklib.PageResolver, w http.ResponseWriter, r *http.Request) {
  things, ok := r.URL.Query()["things[]"]
  if !ok || len(things) != 2 {
    http.Error(w, "Bad search", http.StatusBadRequest)
    return
  }

  ret := make([]interface{}, 0, len(things))
  for _, name := range things {
    page, ok := pageResolver.PageByTitle(name)
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
      index(pageResolver, w, r)
    }
  })
  setupStatics(
    []string{"robots.txt"},
    []string {"css", "js", "img"},
  )
  log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
