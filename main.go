package main

import (
  "log"
  "net/http"
  "github.com/jmoiron/sqlx"
)

var db *sqlx.DB = ConnectToDatabase()

func main() {

  db.SetMaxOpenConns(1000) //tune this

  router := NewRouter()
  //BuildDatabase()
  defer db.Close()

  log.Fatal(http.ListenAndServe(":8080", router))
}