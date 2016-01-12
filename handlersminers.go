package main

import (
  "fmt"
  "net/http"
  "encoding/json"
  "bytes"
  "io/ioutil"
)

// Miners admin home page
func AdminMiners(w http.ResponseWriter, r *http.Request) {
  miners, err :=  MinersCollection()
  content := make(map[string]interface{})

  content["Title"] = "Admin Miners"
  if err != nil {
    content["Error"] = "Miners database table not yet created"
  } else {
    content["Miners"] = miners
  }
  renderTemplate(w, "admin/miners/index", content)
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
  lastInsertId, err := InsertMiner(name, location, url)
  sendIdUrl := fmt.Sprintf("%s/id", url)
  idData := fmt.Sprintf("{\"id\":\"%d\"}", lastInsertId)

  var jsonStr = []byte(idData)
    req, err := http.NewRequest("POST", sendIdUrl, bytes.NewBuffer(jsonStr))
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    fmt.Println("response Status:", resp.Status)
    fmt.Println("response Headers:", resp.Header)
    body, _ := ioutil.ReadAll(resp.Body)
    fmt.Println("response Body:", string(body))


  AdminMiners(w, r)
}

// Handles receipt of post from a Miner
func MinerPost(w http.ResponseWriter, r *http.Request) {
  decoder := json.NewDecoder(r.Body)

  var posts MinerPostsJSON
  err := decoder.Decode(&posts)
  if err != nil {
    fmt.Println("Error:", err)
  }

  fmt.Println("Post from Miner Id: ", posts.MinerId)
  fmt.Println(posts)

  var miner Miner
  rows := QueryMinerForId(posts.MinerId)
  for rows.Next() {
    var uid int
    var name string
    var location string
    var url string
    err := rows.Scan(&uid, &name, &location, &url)
    checkErr(err)
    miner = Miner {
      Uid: uid,
      Name: name,
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

    lastInsertId := InsertPost(miner.Location, url, posted.Time, mined.Time)
    if lastInsertId != 0 {
      for k, v := range post.Terms {
        fmt.Println(k, v)
        InsertTerm(miner.Location, k, v, lastInsertId, posted.Time)
      }
    }
  }
}