package main

import (
  "fmt"
  "log"
  "net/http"

  "github.com/gorilla/mux"
)

func main() {

  router := mux.NewRouter().StrictSlash(true)
  router.HandleFunc("/", Index)
  router.HandleFunc("/trends/{location}/{term}", Trends)
  log.Fatal(http.ListenAndServe(":8080", router))
}

func Index(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "Welcome to the TFW application server")
}

func Trends(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  location := vars["location"]
  term := vars["term"]
  fmt.Fprintln(w, "Trends for ", location, " term:", term)
}