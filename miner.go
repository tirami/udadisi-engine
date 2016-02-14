package main

type Miner struct {
  Uid int `json:"id"`
  Name string `json:"name"`
  Location string `json:"location"`
  GeoCoord Point `json:"geo_coord"`
  Source string `json:"source"`
  Url string `json:"url"`
}

type Miners []Miner