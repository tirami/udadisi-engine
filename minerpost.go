package main

import (
    "strings"
    "time"
)

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
  tt, err := time.Parse("200601021504", strings.Trim(string(buf), `"`))
  if err != nil {
    return err
  }
  t.Time = tt
  return nil
}