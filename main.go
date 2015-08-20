package main

import (
  "log"
  "net/http"
  "database/sql"
)

var db *sql.DB = ConnectToDatabase()

func main() {

  router := NewRouter()
  //BuildDatabase()
  defer db.Close()

  log.Fatal(http.ListenAndServe(":8080", router))
}