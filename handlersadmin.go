package main

import (
  "fmt"
  "net/http"
)

func AdminIndex(w http.ResponseWriter, r *http.Request) {
  content := make(map[string]interface{})
  content["Title"] = "Admin Home Page"
  renderTemplate(w, "admin/index", content)
}

func AdminBuildDatabase(w http.ResponseWriter, r *http.Request) {

  BuildDatabase()

  fmt.Fprintf(w, "<a href=\"/\">Home</a>")
  fmt.Fprintf(w, "<p>Database built</p>")
  fmt.Fprintf(w, "<a href=\"/admin/\">Admin Home</a>")
}