package main

import (
  "fmt"
  "net/http"
  "io/ioutil"
  "os"
)

func Swagger(w http.ResponseWriter, r *http.Request) {
  file, e := ioutil.ReadFile("./swagger.json")
    if e != nil {
        fmt.Printf("File error: %v\n", e)
        os.Exit(1)
    }

  s := string(file)

  w.Header().Add("Access-Control-Allow-Origin", "*")
  w.Header().Add("Access-Control-Allow-Methods", "GET")
  w.Header().Add("Access-Control-Allow-Headers", "Content-Type, api_key, Authorization")
  fmt.Fprintf(w, s)
}