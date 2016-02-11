package main

import (
  "fmt"
  "net/http"
  "os"
)

func AdminIndex(w http.ResponseWriter, r *http.Request) {
  sess, err := globalSessions.SessionStart(w, r)
  if err != nil {
      //need logging here instead of print
      fmt.Printf("Error, could not start session %v\n", err)
      return
  }
  defer sess.SessionRelease(w)
  username := sess.Get("username")
  if username == nil {
    AdminLogin(w, r)
  } else {
    content := make(map[string]interface{})
    content["Title"] = "Admin Home Page"
    renderTemplate(w, "admin/index", content)
  }
}

func AdminBuildDatabase(w http.ResponseWriter, r *http.Request) {
  sess, err := globalSessions.SessionStart(w, r)
  if err != nil {
      //need logging here instead of print
      fmt.Printf("Error, could not start session %v\n", err)
      return
  }
  defer sess.SessionRelease(w)
  username := sess.Get("username")
  if username == nil {
    AdminLogin(w, r)
  } else {
    BuildDatabase()

    fmt.Fprintf(w, "<a href=\"/\">Home</a>")
    fmt.Fprintf(w, "<p>Database built</p>")
    fmt.Fprintf(w, "<a href=\"/admin/\">Admin Home</a>")
  }
}

func AdminClearData(w http.ResponseWriter, r *http.Request) {
  sess, err := globalSessions.SessionStart(w, r)
  if err != nil {
      //need logging here instead of print
      fmt.Printf("Error, could not start session %v\n", err)
      return
  }
  defer sess.SessionRelease(w)
  username := sess.Get("username")
  if username == nil {
    AdminLogin(w, r)
  } else {
    content := make(map[string]interface{})
    content["Title"] = "Admin Home Page"

    err := ClearData()
    if err != nil {
      content["Error"] = err
    }

    renderTemplate(w, "admin/index", content)
  }
}

func AdminLogin(w http.ResponseWriter, r *http.Request) {
    sess, err := globalSessions.SessionStart(w, r)
    if err != nil {
        //need logging here instead of print
        fmt.Printf("Error, could not start session %v\n", err)
        return
    }
    defer sess.SessionRelease(w)
    if r.Method == "GET" {
      content := make(map[string]interface{})
      renderTemplate(w, "admin/login", content)
    } else {
      username := r.PostFormValue("username")
      password := r.PostFormValue("password")
      admin_username := os.Getenv("ADMIN_USERNAME")
      admin_password := os.Getenv("ADMIN_PASSWORD")

      if username == admin_username && password == admin_password {
        sess.Set("username", username)
        content := make(map[string]interface{})
        content["Title"] = "Admin Home Page"
        renderTemplate(w, "admin/index", content)
      } else {
        content := make(map[string]interface{})
        content["Error"] = "Incorrect login details"
        renderTemplate(w, "admin/login", content)
      }
    }
}

func AdminLogout(w http.ResponseWriter, r *http.Request) {
    globalSessions.SessionDestroy(w, r)
    content := make(map[string]interface{})
    renderTemplate(w, "admin/login", content)
}