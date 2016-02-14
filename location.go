package main

type Location struct {
  Name string `json:"name"`
  GeoCoord Point `json:"geo_coord"`
}

type Locations []Location