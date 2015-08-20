package main

import (
  "fmt"
  "net/http"
  "time"
  "encoding/json"
  "strconv"
  "io/ioutil"
  "os"

  "github.com/gorilla/mux"
)


func APIDocs(w http.ResponseWriter, r *http.Request) {
  file, e := ioutil.ReadFile("./swagger.json")
    if e != nil {
        fmt.Printf("File error: %v\n", e)
        os.Exit(1)
    }

  s := string(file)

  w.Header().Add("Access-Control-Allow-Headers", "Content-Type, api_key, Authorization")
  fmt.Fprintf(w, s)
}

func Index(w http.ResponseWriter, r *http.Request) {
  w.Header().Add("Access-Control-Allow-Headers", "Content-Type, api_key, Authorization")

  fmt.Fprintf(w, "<h1>Welcome to the Udadisi Engine</h1>")
  fmt.Fprintf(w, "View <a href=\"web/trends/samplelocation\">basic</a> sample data viewer")

  fmt.Fprintf(w, "<h2>API Docs</h2>")
  fmt.Fprintf(w, "<a href=\"apidocs\">{hostname}/apidocs</a>")

  fmt.Fprintf(w, "<h2>JSON Output</h2>")
  fmt.Fprintf(w, "<p>Looking for JSON output? Use the following:</p>")
  fmt.Fprintf(w, "<ul>")
  fmt.Fprintf(w, "<li><a href=\"trends/samplelocation\">{hostname}/trends/samplelocation</a></li>")
  fmt.Fprintf(w, "<li><a href=\"trends/samplelocation?limit=10\">{hostname}/trends/samplelocation?limit=10</a> for top 10 results</li>")

  fmt.Fprintf(w, "<li><a href=\"trends/samplelocation/code\">{hostname}/trends/samplelocation/code</a></li>")
  fmt.Fprintf(w, "</ul>")
}

func TrendsRouteIndex(w http.ResponseWriter, r *http.Request) {
  //vars := mux.Vars(r)
  //location := vars["location"]
  limitParam := r.URL.Query().Get("limit")
  limit, _ := strconv.ParseInt(limitParam, 10, 0)
  wordCounts := WordCountRootCollection()

  totalCounts := map[string]int {}

  for _, wordcount := range wordCounts {
    count := totalCounts[wordcount.Term]
    count = count + wordcount.Occurrences
    totalCounts[wordcount.Term] = count
  }

  sortedCounts := WordCounts {}

  for _, res := range sortedKeys(totalCounts) {
    if limit == 0 || len(sortedCounts) < int(limit) {
      sortedCounts = append(sortedCounts, WordCount {
                Term: res,
                Occurrences: totalCounts[res],
              })
    }
  }

  json.NewEncoder(w).Encode(sortedCounts)
}

func TrendsIndex(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  //location := vars["location"]
  term := vars["term"]

  termTrends := TermTrends {}
  trends := TrendsCollection(term)

  last_trend_term := ""
  totalCounts := map[string]int {}
  thisTerm := ""
  sources := []string {}
  for _, trend := range trends {
    if trend.Term != last_trend_term {

      if thisTerm != "" {
        termTrend := BuildTrendsJSON(thisTerm, totalCounts, sources)
        termTrends = append(termTrends, termTrend)
      }


      last_trend_term = trend.Term
      totalCounts = map[string]int {}
      sources = []string {}
      thisTerm = trend.Term
    }
    sources = append(sources, trend.SourceURI)

    for _, wordcount := range trend.WordCounts {
      count := totalCounts[wordcount.Term]
      count = count + wordcount.Occurrences
      totalCounts[wordcount.Term] = count
    }
  }

  termTrend := BuildTrendsJSON(thisTerm, totalCounts, sources)
  termTrends = append(termTrends, termTrend)

  json.NewEncoder(w).Encode(termTrends)
}


// Web test stuff
func WebTrendsRouteIndex(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  location := vars["location"]
  wordCounts := WordCountRootCollection()

  fmt.Fprintf(w, "<a href=\"/\">Home</a>")
  fmt.Fprintf(w, "<h1>Main index</h1>")

  totalCounts := map[string]int {}

  for _, wordcount := range wordCounts {
    count := totalCounts[wordcount.Term]
    count = count + wordcount.Occurrences
    totalCounts[wordcount.Term] = count
  }

  DisplayRootCount(w, location, totalCounts)
}

func WebTrendsIndex(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  location := vars["location"]
  term := vars["term"]

  trends := TrendsCollection(term)

  fmt.Fprintf(w, "<a href=\"/\">Home</a>")
  fmt.Fprintf(w, "<h1><a href=\"../%s\">Main index</a></h1>", location)

  last_trend_term := ""
  totalCounts := map[string]int {}
  thisTerm := ""
  sources := []string {}
  for _, trend := range trends {
    if trend.Term != last_trend_term {

      if thisTerm != "" {
        DisplayCount(w, thisTerm, totalCounts, sources)
      }

      last_trend_term = trend.Term
      totalCounts = map[string]int {}
      sources = []string {}
      thisTerm = trend.Term
    }
    sources = append(sources, trend.SourceURI)

    for _, wordcount := range trend.WordCounts {
      count := totalCounts[wordcount.Term]
      count = count + wordcount.Occurrences
      totalCounts[wordcount.Term] = count
    }
  }

  DisplayCount(w, thisTerm, totalCounts, sources)
}

func DisplayRootCount(w http.ResponseWriter, location string, totalCounts map[string]int) {
  for _, res := range sortedKeys(totalCounts) {
    fmt.Fprintf(w, "<li><a href=\"%s/%s\">%s</a> : %d</li>", location, res, res, totalCounts[res])
  }
}

func BuildTrendsJSON(term string, totalCounts map[string]int, sources []string) TermTrend {
  termTrend := TermTrend {}

  termTrend.Term = term

  termWordCounts := WordCounts {}
  for _, res := range sortedKeys(totalCounts) {
    termWordCount := WordCount {
      Term: res,
      Occurrences: totalCounts[res],
    }
    termWordCounts = append(termWordCounts, termWordCount)
  }
  termTrend.WordCounts = termWordCounts

  termSources := Sources {}
  for _, source := range sources {
    termSource := Source {
      Source: "Twitter",
      SourceURI: source,
    }
    termSources = append(termSources, termSource)
  }
  termTrend.Sources = termSources

  return termTrend
}

func DisplayCount(w http.ResponseWriter, term string, totalCounts map[string]int, sources []string) {
  fmt.Fprintf(w, "<h2>%s (%d)</h2>", term, totalCounts[term])

  for _, res := range sortedKeys(totalCounts) {
    fmt.Fprintf(w, "<li><a href=\"%s\">%s</a> : %d</li>", res, res, totalCounts[res])
  }

  fmt.Fprintf(w, "<h3>Sources</h3>")
  for _, source := range sources {
    fmt.Fprintf(w, "<li><a href=\"%s\" target=\"_new\">%s</a></li>", source, source)
  }
}


// Internal stuff
func WordCountRootCollection() WordCounts {
  wordCounts := WordCounts {}

  rows := QueryTerms("")
    for rows.Next() {
        var uid int
        var postid int
        var term string
        var wordcount int
        err := rows.Scan(&uid, &postid, &term, &wordcount)
        checkErr(err)
        wordCount := WordCount {
          Term: term,
          Occurrences: wordcount,
        }
        wordCounts = append(wordCounts, wordCount)
    }

    return wordCounts
}

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