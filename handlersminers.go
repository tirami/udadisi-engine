package main

import (
  "fmt"
  "net/http"
  "encoding/json"
  "bytes"
  "strconv"
)

// Miners admin home page
func AdminMiners(w http.ResponseWriter, r *http.Request) {
  miners, err :=  MinersCollection()
  content := make(map[string]interface{})

  content["Title"] = "Miners Admin"
  if err != nil {
    content["Error"] = "Miners database table not yet created"
  } else {
    content["Miners"] = miners
  }
  renderTemplate(w, "admin/miners/index", content)
}

func AdminMinersResetDatabase(w http.ResponseWriter, r *http.Request) {

  ResetMinersDatabase()

  AdminMiners(w, r)
}

// Creates a new miner
func AdminCreateMiner(w http.ResponseWriter, r *http.Request) {
  err := r.ParseForm()
  if err != nil {
    fmt.Println(err)
  }

  name := r.PostFormValue("name")
  url := r.PostFormValue("url")
  location := r.PostFormValue("location")
  source := r.PostFormValue("source")
  lastInsertId, err := InsertMiner(name, location, source, url)
  sendIdUrl := fmt.Sprintf("%s/categories", url)
  idData := fmt.Sprintf("{\"id\":\"%d\"}", lastInsertId)

  var jsonStr = []byte(idData)
  req, err := http.NewRequest("POST", sendIdUrl, bytes.NewBuffer(jsonStr))
  req.Header.Set("Content-Type", "application/json")

  client := &http.Client{}
  resp, err := client.Do(req)
  if err == nil {
    defer resp.Body.Close()
  }

  content := make(map[string]interface{})

  content["Title"] = "Miners Admin"
  if err != nil {
    content["MinerError"] = err
  }

  miners, err :=  MinersCollection()
  if err != nil {
    content["Error"] = "Miners database table not yet created"
  } else {
    content["Miners"] = miners
  }

  renderTemplate(w, "admin/miners/index", content)
}

// Handles receipt of post from a Miner
func MinerPost(w http.ResponseWriter, r *http.Request) {
  fmt.Println(r.Body)
  decoder := json.NewDecoder(r.Body)
  var posts MinerPostsJSON
  err := decoder.Decode(&posts)
  if err != nil {
    fmt.Println("Error:", err)
  }

  fmt.Println("Post from Miner Id: ", posts.MinerId)
  fmt.Println(posts)

  minerConv, _ := strconv.ParseInt(posts.MinerId, 10, 0)
  var miner Miner
  rows := QueryMinerForId(int(minerConv))
  for rows.Next() {
    var uid int
    var name string
    var source string
    var location string
    var url string
    err := rows.Scan(&uid, &name, &source, &location, &url)
    checkErr(err)
    miner = Miner {
      Uid: uid,
      Name: name,
      Source: source,
      Location: location,
      Url: url,
    }
  }

  for _, post := range posts.Posts {
    url := post.Url
    posted := post.Datetime
    mined := post.MinedAt
    fmt.Println("SourceURI ", url)
    fmt.Println("Posted ", posted)
    fmt.Println("Mined ", mined)
    fmt.Println("Terms and their counts")
    fmt.Println("Miner location ", miner.Location)
    fmt.Println("Miner source", miner.Source)

    lastInsertId := InsertPost(miner.Source, miner.Location, url, posted.Time, mined.Time)
    if lastInsertId != 0 {
      for k, v := range post.Terms {
        fmt.Println(k, v)
        InsertTerm(miner.Location, k, v, lastInsertId, posted.Time)
      }
    }
  }
}