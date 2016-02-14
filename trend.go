package main

import (
    "time"
)

type Trend struct {
  Term string `json:"term"`
  SourceURI string `json:"source_uri"`
  Mined time.Time `json:"mined"`
  Posted time.Time `json:"posted"`
  WordCounts []WordCount `json:"word_counts"`
}

type Trends []Trend