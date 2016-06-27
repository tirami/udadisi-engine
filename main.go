package main

import (
  "log"
  "net/http"
  "github.com/jmoiron/sqlx"
  "github.com/astaxie/beego/session"
)

var db *sqlx.DB = ConnectToDatabase()
var globalSessions *session.Manager

func init() {
  globalSessions, _ = session.NewManager("file",`{"cookieName":"gosessionid","gclifetime":3600,"ProviderConfig":"./tmp"}`)
  go globalSessions.GC()
}

func main() {

  db.SetMaxOpenConns(20) //tune this

  router := NewRouter()
  defer db.Close()

  log.Fatal(http.ListenAndServe(":8080", router))
}