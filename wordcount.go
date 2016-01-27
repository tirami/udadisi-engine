package main

//import "time"
import (
  "sort"
)

type WordCount struct {
  Term string `json:"term"`
  Occurrences  int `json:"occurrences"`
  Velocity float64 `json:"velocity"`
  Series []int `json:"series"`
  Sequence int `json:"sequence"`
}

type WordCounts []WordCount

type sortedWordCountMap struct {
  m map[string]WordCount
  s []string
}

func (sm *sortedWordCountMap) Len() int {
  return len(sm.m)
}

func (sm *sortedWordCountMap) Less(i, j int) bool {
    a, b := sm.m[sm.s[i]].Velocity, sm.m[sm.s[j]].Velocity
    if a != b {
        // Order by decreasing value.
        return a > b
    } else {
        // Otherwise, alphabetical order.
        return sm.s[j] > sm.s[i]
    }
}

func (sm *sortedWordCountMap) Swap(i, j int) {
  sm.s[i], sm.s[j] = sm.s[j], sm.s[i]
}

func sortedWordCountKeys(m map[string]WordCount) []string {
  sm := new(sortedWordCountMap)
  sm.m = m
  sm.s = make([]string, len(m))
  i := 0
  for key, _ := range m {
    sm.s[i] = key
    i++
  }
  sort.Sort(sm)
  return sm.s
}