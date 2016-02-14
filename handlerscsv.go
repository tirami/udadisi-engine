package main

import (
  //"fmt"
  "net/http"
  "time"
  "strconv"
  "github.com/gorilla/mux"
  "encoding/csv"
  "bytes"
)

// Generates CSV file of the sources for a trend
func TrendSourcesCSV(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  location := vars["location"]
  source := r.URL.Query().Get("source")
  term := vars["term"]

  //termTrends := TermTrends {}
  fromParam := r.URL.Query().Get("from")
  toParam := r.URL.Query().Get("to")
  velocityParam := r.URL.Query().Get("velocity")
  minimumVelocity, _ := strconv.ParseFloat(velocityParam, 64)
  intervalParam := r.URL.Query().Get("interval")
  intervalConv, _ := strconv.ParseInt(intervalParam, 10, 0)
  interval := int(intervalConv)
  velocityInterval := float64(interval)
    if velocityInterval == 0.0 {
      velocityInterval = 1.0
    }

  if(minimumVelocity < 0.0) {
    minimumVelocity = 0.0
  }

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

  termPackage := TrendsCollection(source,location, term, fromParam, toParam, interval, velocityInterval, minimumVelocity)

  b := &bytes.Buffer{} // creates IO Writer
  wr := csv.NewWriter(b) // creates a csv writer that uses the io buffer.

  wr.Write([]string{ "Source", "Location", "Posted", "Source URI"})
  for _, source := range termPackage.Sources {
    record := []string{ source.Source, source.Location, source.Posted.Format("2006-01-02 15:04:05"), source.SourceURI }
    wr.Write(record) // converts array of string to comma seperated values for 1 row.
  }
  wr.Flush() // writes the csv writer data to  the buffered data io writer(b(bytes.buffer))

  w.Header().Set("Content-Type", "text/csv")
  w.Header().Set("Content-Disposition", "attachment;filename=" + term +  ".csv")
  w.Write(b.Bytes())

  w.Header().Add("Access-Control-Allow-Origin", "*")
  w.Header().Add("Access-Control-Allow-Methods", "GET")
  w.Header().Add("Access-Control-Allow-Headers", "Content-Type, api_key, Authorization")
}