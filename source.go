package main

import (
    "time"
)

type Source struct {
  Source string `json:"source"`
  Location string `json:"location"`
  SourceURI string `json:"source_uri"`
  Posted time.Time `json:"posted"`
  Mined time.Time `json:"mined"`
}

type Sources []Source