package main

import "time"

type WordCount struct {
  Term string `json:"term"`
  Occurrences  int `json:"occurrences"`
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
  SourceURI string `json:"source_uri"`
  Posted time.Time `json:"posted"`
}

type Sources []Source

type TermTrend struct {
  Term string `json:"term"`
  WordCounts []WordCount `json:"word_counts"`
  Sources []Source `json:"sources"`
}

type TermTrends []TermTrend

type Seed struct {
  Miner string `json:"miner"`
  Location string `json:"location"`
  Source string `json:"source"`
}

type Seeds []Seed

type Location struct {
  Name string `json:"name"`
}

type Locations []Location

type Miner struct {
  Name string `json:"name"`
  Location string `json:"location"`
  Url string `json:"url"`
}

type Miners []Miner
