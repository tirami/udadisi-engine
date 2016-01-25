package main

import (
  "bytes"
  "encoding/json"
  "fmt"
  "log"
  "strings"
  "strconv"
)

// Represents a Physical Point in geographic notation [lat, lng].
type Point struct {
  latitude float64
  longitude float64
}

// Returns a new Point populated by the passed in latitude (lat) and longitude (lng) values.
func NewPoint(latitude float64, longitude float64) *Point {
  return &Point{latitude: latitude, longitude: longitude}
}

// Returns Point p's latitude.
func (p *Point) Latitude() float64 {
  return p.latitude
}

// Returns Point p's longitude.
func (p *Point) Longitude() float64 {
  return p.longitude
}


// Renders the current Point to valid JSON.
// Implements the json.Marshaller Interface.
func (p *Point) MarshalJSON() ([]byte, error) {
  res := fmt.Sprintf(`{"latitude":%v, "longitude":%v}`, p.latitude, p.longitude)
  return []byte(res), nil
}

// Decodes the current Point from a JSON body.
// Throws an error if the body of the point cannot be interpreted by the JSON body
func (p *Point) UnmarshalJSON(data []byte) error {
  // TODO throw an error if there is an issue parsing the body.
  dec := json.NewDecoder(bytes.NewReader(data))
  var values map[string]float64
  err := dec.Decode(&values)

  if err != nil {
    log.Print(err)
    return err
  }

  *p = *NewPoint(values["latitude"], values["longitude"])

  return nil
}


func (p *Point) Scan(src interface{}) error {
  asBytes, ok := src.([]byte)
  if !ok {
    return fmt.Errorf("Scan source was not []bytes")
  }
  values := strings.Split(string(asBytes)[1:len(string(asBytes))- 2], ",")

  v1, _ := strconv.ParseFloat(values[0], 64)
  v2, _ := strconv.ParseFloat(values[1], 64)

  *p = *NewPoint(v1, v2)

  return nil
}

