package main

//import "time"
import (
    "fmt"
    "strings"
    "time"
)


type WordCount struct {
  Term string `json:"term"`
  Occurrences  int `json:"occurrences"`
  Velocity float64 `json:"velocity"`
  Series []int `json:"series"`
  Sequence int `json:"sequence"`
}

type WordCounts []WordCount

type Trend struct {
  Term string `json:"term"`
  SourceURI string `json:"source_uri"`
  Mined time.Time `json:"mined"`
  Posted time.Time `json:"posted"`
  WordCounts []WordCount `json:"word_counts"`
}

type Trends []Trend

type Source struct {
  Source string `json:"source"`
  Location string `json:"location"`
  SourceURI string `json:"source_uri"`
  Posted time.Time `json:"posted"`
  Mined time.Time `json:"mined"`
}

type Sources []Source

type TermTrend struct {
  Term string `json:"term"`
  WordCounts []WordCount `json:"word_counts"`
  Sources []Source `json:"sources"`
}

type TermTrends []TermTrend

type TermPackage struct {
  Term string `json:"term"`
  Velocity float64 `json:"velocity"`
  Series []int `json:"series"`
  Related []string `json:"related"`
  SourceTypes []SourceType `json:"source_types"`
  Sources []Source `json:"sources"`
}

type SourceType struct {
  Name string `json:"name"`
  Series []int `json:"series"`
}

type Location struct {
  Name string `json:"name"`
  GeoCoord Point `json:"geo_coord"`
}

type Locations []Location

type Miner struct {
  Uid int `json:"id"`
  Name string `json:"name"`
  Location string `json:"location"`
  GeoCoord Point `json:"geo_coord"`
  Source string `json:"source"`
  Url string `json:"url"`
}

type Miners []Miner

type MinerPostsJSON struct {
  Posts []MinerPostJSON `json:"posts"`
  MinerId string `json:"miner_id"`
}

type MinerPostJSON struct {
  Terms map[string]int `json:"terms"`
  Url string `json:"url"`
  Datetime myTime `json:"datetime"`
  MinedAt myTime `json:"mined_at"`
}

type myTime struct {
  time.Time
}

func (t *myTime) UnmarshalJSON(buf []byte) error {
  fmt.Println(string(buf))

    tt, err := time.Parse("200601021504", strings.Trim(string(buf), `"`))
    if err != nil {
        return err
    }
    t.Time = tt
    return nil
}

