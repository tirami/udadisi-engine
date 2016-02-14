package main

import (
  "fmt"
  "net/http"
  "time"
  "strconv"
  "html/template"
  "github.com/gorilla/mux"
)

// Internal stuff
func renderTemplate(w http.ResponseWriter, tmpl string, content map[string]interface{}) {
  t, err := template.ParseFiles("views/" + tmpl + ".html")

  if err != nil {
    fmt.Println("%s", err)
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  err = t.Execute(w, content)
  if err != nil {
    fmt.Println("%s", err)
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }
}

// Main home page
func Index(w http.ResponseWriter, r *http.Request) {
  w.Header().Add("Access-Control-Allow-Headers", "Content-Type, api_key, Authorization")

  locations, err := BuildLocationsList()

  content := make(map[string]interface{})
  content["Title"] = "Welcome to the Udadisi Engine"
  if err != nil {
    content["Error"] = "Miners database table not yet created"
  } else {
    content["Locations"] = locations
    t := time.Now()
    lastWeek := t.Add(-24 * time.Hour * 7)
    content["LastWeek"] = lastWeek.Format("200601021504")
    lastMonth := t.Add(-24 * time.Hour * 7 * 4)
    content["LastMonth"] = lastMonth.Format("200601021504")
  }
  renderTemplate(w, "index", content)
}

// Diagonistic web pages
func WebStats(w http.ResponseWriter, r *http.Request) {

}

func WebTrendsRouteIndex(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  location := vars["location"]
  source := r.URL.Query().Get("source")
  limitParam := r.URL.Query().Get("limit")
  limit, _ := strconv.ParseInt(limitParam, 10, 0)
  if limit < 1 {
    limit = 10
  }
  intervalParam := r.URL.Query().Get("interval")
  interval, _ := strconv.ParseInt(intervalParam, 10, 0)
  fromParam := r.URL.Query().Get("from")
  toParam := r.URL.Query().Get("to")
  t := time.Now()
  if fromParam == "" {
    from := t.Add(-24 * time.Hour)
    fromParam = from.Format("200601021504")
  }
  if toParam == "" {
    toParam = t.Format("200601021504")
  }
  if interval < 1 {
    interval = 2
  }
  sortedCounts, err := WordCountRootCollection(location, source, fromParam, toParam, int(interval), int(limit))

  content := make(map[string]interface{})
  if err != nil {
    content["Error"] = err
  }
  content["Title"] = "Main Index"
  content["Location"] = location
  content["FromParam"] = fromParam
  content["ToParam"] = toParam
  content["Interval"] = int(interval)
  content["SortedCounts"] = sortedCounts
  content["VelocityMidPoint"] = 0.0

  renderTemplate(w, "termsindex", content)
}

func WebTrendsIndex(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  location := vars["location"]
  source := r.URL.Query().Get("source")
  term := vars["term"]
  fromParam := r.URL.Query().Get("from")
  toParam := r.URL.Query().Get("to")
  intervalParam := r.URL.Query().Get("interval")
  intervalConv, _ := strconv.ParseInt(intervalParam, 10, 0)
  interval := int(intervalConv)
  t := time.Now()
  if fromParam == "" {
    from := t.Add(-24 * time.Hour)
    fromParam = from.Format("200601021504")
  }
  if toParam == "" {
    toParam = t.Format("200601021504")
  }
  if interval < 1 {
    interval = 2
  }

  termPackage := TrendsCollection(source, location, term, fromParam, toParam, interval, 1.0, 0.0)

  content := make(map[string]interface{})
  content["Location"] = location
  content["FromParam"] = fromParam
  content["ToParam"] = toParam
  content["Interval"] = int(interval)
  content["TermPackage"] = termPackage

  renderTemplate(w, "term", content)
}