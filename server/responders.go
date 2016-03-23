package server

import (
  "encoding/json"
  "fmt"
  "howtotip/models"
  "html/template"
  "log"
  "net/http"
  "path"
  "reflect"
  "strings"
  "time"
)

func jsonResponder(w http.ResponseWriter, r *http.Request, data interface{}, err error) {
  w.Header().Set("Content-Type", "application/json")
  enc := json.NewEncoder(w)
  if err != nil {
    http.Error(w, "", http.StatusBadRequest)
    enc.Encode(map[string]string{"error": err.Error()})
  } else {
    kind := reflect.TypeOf(data).Kind()
    if (kind == reflect.Slice || kind == reflect.Map || kind == reflect.Ptr) && reflect.ValueOf(data).IsNil() {
      if kind == reflect.Slice {
        fmt.Fprintf(w, "[]")
      } else {
        fmt.Fprintf(w, "{}")
      }
    } else {
      enc.Encode(data)
    }
  }
}

func successResponder(w http.ResponseWriter, r *http.Request, err error) {
  jsonResponder(w, r, map[string]bool{"success": err == nil}, err)
}

func CountriesHandler(w http.ResponseWriter, r *http.Request) {
  var err error
  var data []models.Country

  defer func() { jsonResponder(w, r, data, err) }()

  data = models.GetCountries()
}

func CountryHandler(w http.ResponseWriter, r *http.Request) {
  var err error
  var data models.Country

  defer func() { jsonResponder(w, r, data, err) }()

  r.ParseForm()

  slug := r.FormValue("slug")
  if slug == "" {
    return
  }

  data = models.GetCountry(slug)
}

func GetCountriesHandler(w http.ResponseWriter, r *http.Request) {
  data := models.GetCountries()
  page := path.Join("templates", "countries.html")
  tmpl, _ := template.ParseFiles(page)
  tmpl.Execute(w, &data)
}

func GetCountryHandler(w http.ResponseWriter, r *http.Request) {
  r.ParseForm()
  slug := r.FormValue("slug")
  if slug == "" {
    return
  }

  data := models.GetCountry(slug)
  page := path.Join("templates", "form.html")
  tmpl, _ := template.ParseFiles(page)
  tmpl.Execute(w, &data)
}

func PostCountryHandler(w http.ResponseWriter, r *http.Request) {
  slug := r.PostFormValue("slug")
  name := r.PostFormValue("name")
  caption := r.PostFormValue("caption")
  log.Printf("%s %s %s", slug, name, caption)
  if slug == "" || name == "" || caption == "" {
    return
  }

  models.SaveCountry(slug, name, caption)

  data := models.GetCountry(slug)
  page := path.Join("templates", "form.html")
  tmpl, _ := template.ParseFiles(page)
  tmpl.Execute(w, &data)
}

type PageData struct {
  Countries []models.Country
  Data interface{}
  Year string
}

func PageHandler(w http.ResponseWriter, r *http.Request) {
  var name string
  var data interface{}

  if r.URL.Path == "/" {
    name = "home"
  } else {
    name = "country"
    slug := strings.Split(r.URL.Path, "/")[1]
    data = models.GetCountry(slug)
  }

  countries := models.GetCountries()

  t := time.Now()
  year := t.Format("2006")
  pageData := PageData{
    countries,
    data,
    year,
  }

  layout := path.Join("templates", "layout.html")
  toolbar := path.Join("templates", "_toolbar.html")
  page := path.Join("templates", fmt.Sprintf("%s.html", name))
  tmpl, _ := template.ParseFiles(layout, page, toolbar)
  tmpl.ExecuteTemplate(w, "layout", &pageData)
}

func RouteHandler(router RegexpRouter) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    start := time.Now()
    router.ServeHTTP(w, r)
    log.Printf("%s %s %s %v", r.RemoteAddr, r.Method, r.URL, time.Since(start))
  })
}

