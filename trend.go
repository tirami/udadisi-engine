package main

import "time"

type Trend struct {
  Term string `jons:"term"`
  Source string `json:"source"`
  Mined time.Time `json:"mined"`
}

type Trends []Trend