package main

import "time"

type WordCount struct {
  Term string `json:"term"`
  Source string `json:"source"`
  Occurrences  int `json:"occurrences"`
}

type WordCounts []WordCount

type Trend struct {
  Term string `json:"term"`
  SourceURI string `json:"source_uri"`
  Mined time.Time `json:"mined"`
  WordCounts []WordCount `json:"word_counts"`
}

type Trends []Trend