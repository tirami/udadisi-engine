package main

type TermPackage struct {
  Term string `json:"term"`
  Velocity float64 `json:"velocity"`
  Series []int `json:"series"`
  Related []Related `json:"related"`
  SourceTypes []SourceType `json:"source_types"`
  Sources []Source `json:"sources"`
}