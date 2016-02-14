package main

import (
  "fmt"
  "net/http"
  "time"
  "encoding/json"
  "strconv"
  "github.com/gorilla/mux"
)

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
  fmt.Println(vars)
  location := vars["location"]
  source := r.URL.Query().Get("source")
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
  wordCounts, _ := WordCountRootCollection(location, source, fromParam, toParam, int(interval), 0)

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
  sortedCounts, _ := WordCountRootCollection(location, source, fromParam, toParam, int(interval), int(limit))

  w.Header().Add("Access-Control-Allow-Origin", "*")
  w.Header().Add("Access-Control-Allow-Methods", "GET")
  w.Header().Add("Access-Control-Allow-Headers", "Content-Type, api_key, Authorization")
  json.NewEncoder(w).Encode(sortedCounts)
}

// Generates JSON list of trends for a term
func TrendsIndex(w http.ResponseWriter, r *http.Request) {
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
          GeoCoord: miner.GeoCoord,
        }
        locations = append(locations, locationsAdded[miner.Location])
      }
  }

  return
}



func WordCountRootCollection(location string, source string, fromParam string, toParam string, interval int, limit int) (sortedCounts WordCounts,  collectionErr error) {

  defer func() {
        if r := recover(); r != nil {
            var ok bool
            collectionErr, ok = r.(error)
            if !ok {
                collectionErr = fmt.Errorf("WordCountRootCollection: %v", r)
            }
        }
    }()

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
  //fmt.Println("duration ", duration.Minutes())
  duration  = duration / time.Duration(interval)

  for i := 0; i < interval; i++ {

    toTime = fromTime.Add(duration)
    toParam = toTime.Format("200601021504")

    rows, err := QueryTerms(source, location, "", fromParam, toParam)
    checkErr(err)

    for rows.Next() {
      var uid int
      var postid int
      var term string
      var wordcount int
      var posted time.Time
      var termLocation string
      var locationHash int
      var termSource string
      err := rows.Scan(&uid, &postid, &term, &wordcount, &posted, &termLocation, &locationHash, &termSource)
      checkErr(err)
      wordCount := WordCount {
        Term: term,
        Occurrences: wordcount,
        Series: []int{wordcount},
        Sequence: i,
      }
      wordCounts = append(wordCounts, wordCount)
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

  velocityCounts := map[string]WordCount {}


  // For ordering by velocity
  for key, _ := range totalCounts {
    if totalCounts[key] > 1 {

      // Calculate the velocity
      //seriesAverage := float64(totalCounts[key]) / float64(interval - 1)
      velocity := 0.0
      if serieses[key][interval - 2] == 0 {
        velocity = float64(serieses[key][interval - 1])
      } else {
        velocity = (float64(serieses[key][interval - 1]) - float64(serieses[key][interval - 2])) / float64(serieses[key][interval - 2])
      }

      velocityCounts[key] = WordCount {
                Term: key,
                Occurrences: totalCounts[key],
                Series: serieses[key],
                Velocity: velocity,
              }
      }
  }

  rankings := map[float64]int {}

  sortedCounts = WordCounts {}
  for _, res := range sortedWordCountKeys(velocityCounts) {
    if _, ok := rankings[velocityCounts[res].Velocity]; ok {
        rankings[velocityCounts[res].Velocity]++
      } else {
        rankings[velocityCounts[res].Velocity] = 1
      }

      if len(rankings) > limit {
        break
      }
    sortedCounts = append(sortedCounts, velocityCounts[res])
  }

  // For ordering by occurrances
  /*
  sortedCounts = WordCounts {}
  for _, res := range sortedKeys(totalCounts) {
    if limit == 0 || len(sortedCounts) < int(limit) {
      if totalCounts[res] > 1 {
        // Calculate the velocity
        seriesAverage := float64(totalCounts[res]) / float64(interval)
        //fmt.Println("Total:", totalCounts[res], "interval:", interval)
        //fmt.Println("Average:", seriesAverage)

        sortedCounts = append(sortedCounts, WordCount {
                  Term: res,
                  Occurrences: totalCounts[res],
                  Series: serieses[res],
                  Velocity: float64(serieses[res][interval - 1]) / seriesAverage,
                })
      }
    }
  }
  */

  return
}

func MinersCollection() (miners Miners, err error) {
  miners = Miners {}

  rows, err := QueryMiners()
  if err != nil {

  } else {
    for rows.Next() {
      var uid int
      var name string
      var source string
      var location string
      var url string
      var geoCoord Point
      var locationHash int
      err := rows.Scan(&uid, &name, &source, &location, &url, &geoCoord, &locationHash)
      checkErr(err)
      miner := Miner {
        Uid: uid,
        Name: name,
        Source: source,
        Location: location,
        GeoCoord: geoCoord,
        Url: url,
      }
      miners = append(miners, miner)
    }
  }

  return
}

func TrendsCollection(source string, location string, term string, fromParam string, toParam string, interval int, velocityInterval float64, minimumVelocity float64) TermPackage {

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
  //fmt.Println("duration ", duration.Minutes())
  duration  = duration / time.Duration(interval)

  termPackage := TermPackage {
    Term: term,
    Series: make([]int, interval),
    Sources: make([]Source, 0),
    SourceTypes: make([]SourceType, 0),
  }

  related := map[string]int {}

  totalOccurrences := 0

  sourceSerieses := map[string][]int {}

  for i := 0; i < interval; i++ {
    toTime = fromTime.Add(duration)
    toParam = toTime.Format("200601021504")
    rows, err := QueryTerms(source, location, term, fromParam, toParam)
    if err != nil {
    } else {
      for rows.Next() {
        var uid int
        var postid int
        var term string
        var wordcount int
        var posted time.Time
        var location string
        var locationHash int
        var source string
        err := rows.Scan(&uid, &postid, &term, &wordcount, &posted, &location, &locationHash, &source)
        //fmt.Println(uid, postid, term, wordcount, posted, location)
        checkErr(err)
        termPackage.Series[i] = termPackage.Series[i] + wordcount
        totalOccurrences = totalOccurrences + wordcount

        if _, ok := sourceSerieses[source]; ok {
        } else {
          sourceSerieses[source] = make([]int, int(interval))
        }
        sourceSerieses[source][i] = sourceSerieses[source][i] + wordcount

        postRows := QueryPosts(fmt.Sprintf(" WHERE uid=%d", postid))
        for postRows.Next() {
          var thisPostuid int
          var mined time.Time
          var postPosted time.Time
          var sourceURI string
          var postLocation string
          var postSource string
          var postLocationHash int
          err = postRows.Scan(&thisPostuid, &mined, &postPosted, &sourceURI, &postLocation, &postSource, &postLocationHash)
          checkErr(err)
          source := Source {
            Source: postSource,
            Location: postLocation,
            SourceURI: sourceURI,
            Posted: postPosted,
            Mined: mined,
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
              var wcLocationHash int
              var wcSource string
              err := termsRows.Scan(&wcuid, &wcpostid, &wcTerm, &wordcount, &wcPosted, &wcLocation, &wcLocationHash, &wcSource)
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

  for key, value := range sourceSerieses {
    termPackage.SourceTypes = append(termPackage.SourceTypes, SourceType {
      Name: key,
      Series: value,
      })
  }

  for _, res := range sortedKeys(related) {
    // Exclude the term we are finding related terms for
    if term != res {
      termPackage.Related = append(termPackage.Related, Related {
        Term: res,
        Occurrences: related[res],
        })
    }
  }

  // Calculate the velocity
  seriesAverage := float64(totalOccurrences) / float64(interval)
  //fmt.Println("Total:", totalOccurrences, "interval:", interval)
  //fmt.Println("Average:", seriesAverage)
  if seriesAverage != 0 {
    termPackage.Velocity = float64(termPackage.Series[interval - 1]) / seriesAverage
  }

  /*
  fmt.Println("Term:", termPackage.Term)
  fmt.Println("Series:", termPackage.Series)
  fmt.Println("SourceTypes:", termPackage.SourceTypes)
  fmt.Println("Related:", termPackage.Related)
  fmt.Println("Sources:", termPackage.Sources)
  fmt.Println(related)
  fmt.Println(termPackage)
  */

  return termPackage
}