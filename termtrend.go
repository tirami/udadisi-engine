package main

type TermTrend struct {
  Term string `json:"term"`
  WordCounts []WordCount `json:"word_counts"`
  Sources []Source `json:"sources"`
}

type TermTrends []TermTrend