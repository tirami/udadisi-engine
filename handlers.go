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
  fmt.Fprintf(w, "<h1>Welcome to the TFW application server</h1>")
  fmt.Fprintf(w, "View <a href=\"web/trends/samplelocation\">basic</a> sample data viewer")
}

func TrendsRouteIndex(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  location := vars["location"]
  term := ""
  log.Printf("%s:%s", location, term)

  trends := TrendsCollection(term)

  json.NewEncoder(w).Encode(trends)
}

func TrendsIndex(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  location := vars["location"]
  term := vars["term"]

  log.Printf("%s:%s", location, term)

  trends := TrendsCollection(term)

  json.NewEncoder(w).Encode(trends)
}


// Web test stuff
func WebTrendsRouteIndex(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  location := vars["location"]
  term := ""
  log.Printf("%s:%s", location, term)

  trends := TrendsCollection(term)

  fmt.Fprintf(w, "<a href=\"/\">Home</a>")

  fmt.Fprintf(w, "<h1>Main index</h1>")
  for _, trend := range trends {
    fmt.Fprintf(w, "<h2>%s</h2>", trend.Term)
    fmt.Fprintf(w, "<h3><a href=\"%s\" target=\"_new\">%s</a></h3>", trend.SourceURI, trend.SourceURI)
    for _, wordCount := range trend.WordCounts {
      fmt.Fprintf(w, "(<a href=\"%s/%s\">%s</a> : %d), ", location, wordCount.Term, wordCount.Term, wordCount.Occurrences)
    }
  }
}

func WebTrendsIndex(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  location := vars["location"]
  term := vars["term"]
  log.Printf("%s:%s", location, term)

  trends := TrendsCollection(term)

  fmt.Fprintf(w, "<a href=\"/\">Home</a>")

  fmt.Fprintf(w, "<h1><a href=\"../%s\">Main index</a></h1>", location)
  for _, trend := range trends {
    fmt.Fprintf(w, "<h2>%s</h2>", trend.Term)
    fmt.Fprintf(w, "<h3><a href=\"%s\" target=\"_new\">%s</a></h3>", trend.SourceURI, trend.SourceURI)
    for _, wordCount := range trend.WordCounts {
      fmt.Fprintf(w, "(<a href=\"%s\">%s</a> : %d), ", wordCount.Term, wordCount.Term, wordCount.Occurrences)
    }
  }
}


// Internal stuff
func TrendsCollection(term string) Trends {
  trends := Trends {}

  rows := QueryTerms(term)
    for rows.Next() {
        var uid int
        var postid int
        var term string
        var wordcount int
        err := rows.Scan(&uid, &postid, &term, &wordcount)
        checkErr(err)
        postRows := QueryPosts(fmt.Sprintf(" WHERE uid=%d", postid))
        for postRows.Next() {
          var thisPostuid int
          var mined time.Time
          var posted time.Time
          var sourceURI string
          err = postRows.Scan(&thisPostuid, &mined, &posted, &sourceURI)
          checkErr(err)
          termsRows := QueryTermsForPost(thisPostuid)
          wordCounts := WordCounts {}
          for termsRows.Next() {
            var wcuid int
            var wcpostid int
            var wcTerm string
            var wordcount int
            err := termsRows.Scan(&wcuid, &wcpostid, &wcTerm, &wordcount)
            checkErr(err)
            wordCount := WordCount {
              Term: wcTerm,
              Source: "twitter",
              Occurrences: wordcount,
            }
            wordCounts = append(wordCounts, wordCount)
          }
          trend := Trend{
            Term: term,
            SourceURI: sourceURI,
            Mined: mined,
            WordCounts: wordCounts,
          }
          trends = append(trends, trend)
        }
    }

    return trends
}