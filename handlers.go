package main

import (
  "fmt"
  "net/http"
  "time"
  "encoding/json"
  "log"

  "github.com/gorilla/mux"
)

func Index(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "Welcome to the TFW application server")
}

func TrendsIndex(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  location := vars["location"]
  term := vars["term"]

  log.Printf("%s:%s", location, term)

  const jsonForm = "2015-08-04T10:20:30Z"
  time1, _ := time.Parse(jsonForm, "2015-08-04T14:34:00Z")
  time2, _ := time.Parse(jsonForm, "2015-08-06T03:23:00Z")

  trends := Trends {
    Trend{Term: "GPS", Source: "http://www.example.com/gps_posts", Mined: time1 },
    Trend{Term: "Water", Source: "http://www.h2o.com", Mined: time2  },
  }

  json.NewEncoder(w).Encode(trends)
}