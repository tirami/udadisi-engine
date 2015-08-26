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
