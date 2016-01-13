package main

import (
  "fmt"
  "net/http"
  "time"
  "encoding/json"
  "strconv"
  "html/template"

  "github.com/gorilla/mux"
)

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
  }
  renderTemplate(w, "index", content)
}

// Generates JSON list of locations
func RenderLocationsJSON(w http.ResponseWriter, r *http.Request) {
  w.Header().Add("Access-Control-Allow-Origin", "*")
  w.Header().Add("Access-Control-Allow-Methods", "GET")
  w.Header().Add("Access-Control-Allow-Headers", "Content-Type, api_key, Authorization")
  locations, _ := BuildLocationsList()
  json.NewEncoder(w).Encode(locations)
}

// Generates JSON stats for a location
func RenderLocationStatsJSON(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  location := vars["location"]
  intervalParam := r.URL.Query().Get("interval")
  interval, _ := strconv.ParseInt(intervalParam, 10, 0)
  fromParam := r.URL.Query().Get("from")
  toParam := r.URL.Query().Get("to")
  wordCounts := WordCountRootCollection(location, fromParam, toParam, int(interval), 0)

  totalCounts := map[string]int {}

  for _, wordcount := range wordCounts {
    count := totalCounts[wordcount.Term]
    count = count + wordcount.Occurrences
    totalCounts[wordcount.Term] = count
  }

  stats := map[string]string {}
  stats["trendscount"] = strconv.Itoa(len(totalCounts))

  w.Header().Add("Access-Control-Allow-Origin", "*")
  w.Header().Add("Access-Control-Allow-Methods", "GET")
  w.Header().Add("Access-Control-Allow-Headers", "Content-Type, api_key, Authorization")
  json.NewEncoder(w).Encode(stats)
}

// Generates JSON for root list of trends
func TrendsRootIndex(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  location := vars["location"]
  limitParam := r.URL.Query().Get("limit")
  limit, _ := strconv.ParseInt(limitParam, 10, 0)
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
  sortedCounts := WordCountRootCollection(location, fromParam, toParam, int(interval), int(limit))

  w.Header().Add("Access-Control-Allow-Origin", "*")
  w.Header().Add("Access-Control-Allow-Methods", "GET")
  w.Header().Add("Access-Control-Allow-Headers", "Content-Type, api_key, Authorization")
  json.NewEncoder(w).Encode(sortedCounts)
}

// Generates JSON list of trends for a term
func TrendsIndex(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  location := vars["location"]
  term := vars["term"]

  //termTrends := TermTrends {}
  fromParam := r.URL.Query().Get("from")
  toParam := r.URL.Query().Get("to")
  velocityParam := r.URL.Query().Get("velocity")
  minimumVelocity, _ := strconv.ParseFloat(velocityParam, 64)
  intervalParam := r.URL.Query().Get("interval")
  intervalConv, _ := strconv.ParseInt(intervalParam, 10, 2)
  interval := int(intervalConv)
  velocityInterval := float64(interval)
    if velocityInterval == 0.0 {
      velocityInterval = 1.0
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

  termPackage := TrendsCollection(location, term, fromParam, toParam, interval, velocityInterval, minimumVelocity)
  /*
  last_trend_term := ""
  totalCounts := map[string]int {}
  thisTerm := ""
  sources := Sources {}
  for _, trend := range trends {
    if trend.Term != last_trend_term {
      if thisTerm != "" {
        termTrend := BuildTrendsJSON(thisTerm, totalCounts, sources, velocityInterval)
        termTrends = append(termTrends, termTrend)
      }

      last_trend_term = trend.Term
      totalCounts = map[string]int {}
      sources = Sources {}

      thisTerm = trend.Term
    }
    source := Source {
      Source: "Twitter",
      SourceURI: trend.SourceURI,
      Posted: trend.Posted,
    }
    sources = append(sources, source)

    for _, wordcount := range trend.WordCounts {
      count := totalCounts[wordcount.Term]
      count = count + wordcount.Occurrences
      totalCounts[wordcount.Term] = count
    }
  }

  termTrend := BuildTrendsJSON(thisTerm, totalCounts, sources, velocityInterval)
  termTrends = append(termTrends, termTrend)
  */

  w.Header().Add("Access-Control-Allow-Origin", "*")
  w.Header().Add("Access-Control-Allow-Methods", "GET")
  w.Header().Add("Access-Control-Allow-Headers", "Content-Type, api_key, Authorization")
  json.NewEncoder(w).Encode(termPackage)
}

func BuildLocationsList() (locations Locations, err error) {
  miners, err := MinersCollection()
  locations = Locations{}
  locationsAdded := map[string]Location {}

  // Add the default 'all' location
  allLocations := Location{
    Name: "all",
  }
  locations = append(locations, allLocations)

  locationsAdded["all"] = allLocations
  for _, miner := range miners {
    if _, exists := locationsAdded[miner.Location]; !exists {
        locationsAdded[miner.Location] = Location{
          Name: miner.Location,
        }
        locations = append(locations, locationsAdded[miner.Location])
      }
  }

  return
}

// Diagonistic web pages
func WebTrendsRouteIndex(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  location := vars["location"]
  limitParam := r.URL.Query().Get("limit")
  limit, _ := strconv.ParseInt(limitParam, 10, 0)
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
  sortedCounts := WordCountRootCollection(location, fromParam, toParam, int(interval), int(limit))

  content := make(map[string]interface{})
  content["Title"] = "Main Index"
  content["Location"] = location
  content["FromParam"] = fromParam
  content["ToParam"] = toParam
  content["Interval"] = int(interval)
  content["SortedCounts"] = sortedCounts

  renderTemplate(w, "termsindex", content)
}

func WebTrendsIndex(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  location := vars["location"]
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

  termPackage := TrendsCollection(location, term, fromParam, toParam, interval, 1.0, 0.0)

  content := make(map[string]interface{})
  content["Location"] = location
  content["FromParam"] = fromParam
  content["ToParam"] = toParam
  content["Interval"] = int(interval)
  content["TermPackage"] = termPackage

  renderTemplate(w, "term", content)

}

func DisplayRootCount(w http.ResponseWriter, location string, fromParam string, interval int, totalCounts WordCounts) {
  fmt.Fprintf(w, "<h2>%d terms</h2>", len(totalCounts))
  fmt.Fprintf(w, "<div class=\"row\"><div class=\"col-md-4\">")
  fmt.Fprintf(w, "<ul class=\"list-group\">")
  for _, res := range totalCounts {
    fmt.Fprintf(w, "<li class=\"list-group-item\"><a href=\"%s/%s?from=%s&interval=%d\">%s</a> <span class=\"badge\">%d</span></li>",
      location,
      res.Term,
      fromParam,
      interval,
      res.Term, res.Occurrences)
  }
  fmt.Fprintf(w, "</ul></div></div>")
}

func BuildTrendsJSON(term string, totalCounts map[string]int, sources Sources, velocityInterval float64) TermTrend {
  termTrend := TermTrend {}

  termTrend.Term = term

  termWordCounts := WordCounts {}
  for _, res := range sortedKeys(totalCounts) {
    termWordCount := WordCount {
      Term: res,
      Occurrences: totalCounts[res],
      Velocity: float64(totalCounts[res]) / velocityInterval,
    }
    termWordCounts = append(termWordCounts, termWordCount)
  }
  termTrend.WordCounts = termWordCounts

  termSources := Sources {}
  for _, source := range sources {
    termSource := source
    termSources = append(termSources, termSource)
  }

  termTrend.Sources = termSources

  return termTrend
}

func DisplayCount(w http.ResponseWriter, fromParam string, interval int, term string, totalCounts map[string]int, sources Sources) {
  fmt.Fprintf(w, "<h2>%s (%d)</h2>", term, totalCounts[term])

  for _, res := range sortedKeys(totalCounts) {
    fmt.Fprintf(w, "<li><a href=\"%s?from=%s&interval=%d\">%s</a> : %d</li>", res, fromParam, interval, res, totalCounts[res])
  }

  fmt.Fprintf(w, "<h3>Sources</h3>")
  for _, source := range sources {
    fmt.Fprintf(w, "<li>%s : mined %s : posted %s <a href=\"%s\" target=\"_new\">%s</a></li>",
      source.Source,
      source.Mined,
      source.Posted,
      source.SourceURI, source.SourceURI)
  }
}


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

func SeedsCollection() Seeds {
  seeds := Seeds {}

  rows := QuerySeeds()
  for rows.Next() {
    var uid int
    var miner string
    var location string
    var source string
    err := rows.Scan(&uid, &miner, &location, &source)
    checkErr(err)
    seed := Seed {
      Miner: miner,
      Location: location,
      Source: source,
    }
    seeds = append(seeds, seed)
  }

  return seeds
}

func WordCountRootCollection(location string, fromParam string, toParam string, interval int, limit int) WordCounts {
  wordCounts := WordCounts {}

  if location == "all" {
    location = ""
  }

  t := time.Now()
  if fromParam == "" {
    from := t.Add(-24 * time.Hour)
    fromParam = from.Format("200601021504")
  }

  if toParam == "" {
    toParam = t.Format("200601021504")
  }

  fromTime, err := time.Parse("200601021504", fromParam)
  if err != nil {
      fmt.Errorf("invalid from date: %v", err)
  }

  toTime, err := time.Parse("200601021504", toParam)
  if err != nil {
      fmt.Errorf("invalid to date: %v", err)
  }

  duration := toTime.Sub(fromTime)
  fmt.Println("duration ", duration.Minutes())
  duration  = duration / time.Duration(interval)

  for i := 0; i < interval; i++ {

    toTime = fromTime.Add(duration)
    toParam = toTime.Format("200601021504")

    rows, err := QueryTerms(location, "", fromParam, toParam)

    if err != nil {

    } else {
      for rows.Next() {
        var uid int
        var postid int
        var term string
        var wordcount int
        var posted time.Time
        var termLocation string
        err := rows.Scan(&uid, &postid, &term, &wordcount, &posted, &termLocation)
        checkErr(err)
        wordCount := WordCount {
          Term: term,
          Occurrences: wordcount,
          Series: []int{wordcount},
          Sequence: i,
        }
        wordCounts = append(wordCounts, wordCount)
      }
    }

    fromTime = fromTime.Add(duration)
    fromParam = fromTime.Format("200601021504")
 }

 totalCounts := map[string]int {}
  serieses := map[string][]int {}

  for _, wordcount := range wordCounts {
    if _, ok := serieses[wordcount.Term]; ok {
    } else {
      serieses[wordcount.Term] = make([]int, int(interval))
    }
  }

  for _, wordcount := range wordCounts {
    count := totalCounts[wordcount.Term]
    count = count + wordcount.Occurrences
    totalCounts[wordcount.Term] = count
    serieses[wordcount.Term][wordcount.Sequence] = serieses[wordcount.Term][wordcount.Sequence] + 1
  }

  sortedCounts := WordCounts {}

  for _, res := range sortedKeys(totalCounts) {
    if limit == 0 || len(sortedCounts) < int(limit) {
      sortedCounts = append(sortedCounts, WordCount {
                Term: res,
                Occurrences: totalCounts[res],
                Series: serieses[res],
              })
    }
  }

  return sortedCounts
}

func MinersCollection() (miners Miners, err error) {
  miners = Miners {}

  rows, err := QueryMiners()
  if err != nil {

  } else {
    for rows.Next() {
      var uid int
      var name string
      var location string
      var url string
      err := rows.Scan(&uid, &name, &location, &url)
      checkErr(err)
      miner := Miner {
        Uid: uid,
        Name: name,
        Location: location,
        Url: url,
      }
      miners = append(miners, miner)
    }
  }

  return
}

func TrendsCollection(location string, term string, fromParam string, toParam string, interval int, velocityInterval float64, minimumVelocity float64) TermPackage {
  //trends := Trends {}

  if location == "all" {
    location = ""
  }

  t := time.Now()
  if fromParam == "" {
    from := t.Add(-24 * time.Hour)
    fromParam = from.Format("200601021504")
  }

  if toParam == "" {
    toParam = t.Format("200601021504")
  }

  fromTime, err := time.Parse("200601021504", fromParam)
  if err != nil {
      fmt.Errorf("invalid from date: %v", err)
  }

  toTime, err := time.Parse("200601021504", toParam)
  if err != nil {
      fmt.Errorf("invalid to date: %v", err)
  }

  duration := toTime.Sub(fromTime)
  fmt.Println("duration ", duration.Minutes())
  duration  = duration / time.Duration(interval)


  termPackage := TermPackage {
    Term: term,
    Series: make([]int, interval),
    Sources: make([]Source, 0),
  }

  related := map[string]int {}

  totalOccurrences := 0

  for i := 0; i < interval; i++ {

    toTime = fromTime.Add(duration)
    toParam = toTime.Format("200601021504")
    rows, err := QueryTerms(location, term, fromParam, toParam)
    if err != nil {
    } else {
      for rows.Next() {
        var uid int
        var postid int
        var term string
        var wordcount int
        var posted time.Time
        var location string
        err := rows.Scan(&uid, &postid, &term, &wordcount, &posted, &location)
        fmt.Println(uid, postid, term, wordcount, posted, location)
        checkErr(err)
        termPackage.Series[i] = termPackage.Series[i] + wordcount
        totalOccurrences = totalOccurrences + wordcount

        postRows := QueryPosts(fmt.Sprintf(" WHERE uid=%d", postid))
        for postRows.Next() {
          var thisPostuid int
          var mined time.Time
          var postPosted time.Time
          var sourceURI string
          var postLocation string
          err = postRows.Scan(&thisPostuid, &mined, &postPosted, &sourceURI, &postLocation)
          checkErr(err)
          source := Source {
            Source: "Twitter",
            SourceURI: sourceURI,
            Posted: postPosted,
          }
          termPackage.Sources = append(termPackage.Sources, source)
          termsRows := QueryTermsForPost(thisPostuid)
            for termsRows.Next() {
              var wcuid int
              var wcpostid int
              var wcTerm string
              var wordcount int
              var wcPosted time.Time
              var wcLocation string
              err := termsRows.Scan(&wcuid, &wcpostid, &wcTerm, &wordcount, &wcPosted, &wcLocation)
              checkErr(err)
              if _, ok := related[wcTerm]; ok {
                related[wcTerm] += wordcount
              } else {
                related[wcTerm] = wordcount
              }
            }
        }
      }
    }

    // TODO: Need to sort the related terms by velocity
    fromTime = fromTime.Add(duration)
    fromParam = fromTime.Format("200601021504")
  }

  for key, _ := range related {
    termPackage.Related = append(termPackage.Related, key)
  }


  fmt.Println("Term:", termPackage.Term)
  fmt.Println("Series:", termPackage.Series)
  fmt.Println("Related:", termPackage.Related)
  fmt.Println("Sources:", termPackage.Sources)
  fmt.Println(related)
  fmt.Println(termPackage)


/*
  if err != nil {

  } else {

    for rows.Next() {
          var uid int
          var postid int
          var term string
          var wordcount int
          var posted time.Time
          var location string
          err := rows.Scan(&uid, &postid, &term, &wordcount, &posted, &location)
          checkErr(err)
          postRows := QueryPosts(fmt.Sprintf(" WHERE uid=%d", postid))
          for postRows.Next() {
            var thisPostuid int
            var mined time.Time
            var postPosted time.Time
            var sourceURI string
            var postLocation string
            err = postRows.Scan(&thisPostuid, &mined, &postPosted, &sourceURI, &postLocation)
            checkErr(err)
            termsRows := QueryTermsForPost(thisPostuid)
            wordCounts := WordCounts {}
            for termsRows.Next() {
              var wcuid int
              var wcpostid int
              var wcTerm string
              var wordcount int
              var wcPosted time.Time
              var wcLocation string
              err := termsRows.Scan(&wcuid, &wcpostid, &wcTerm, &wordcount, &wcPosted, &wcLocation)
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
              Posted: postPosted,
              WordCounts: wordCounts,
            }
            trends = append(trends, trend)
          }
      }
    }
    */

    return termPackage
}