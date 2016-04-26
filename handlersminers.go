package main

import (
  "fmt"
  "net/http"
  "encoding/json"
  "bytes"
  "strconv"
)

func GetMiner(id int) (miner Miner, err error) {
  
  rows := QueryMinerForId(id)
  for rows.Next() {
    var uid int
    var name string
    var source string
    var location string
    var url string
    var geoCoord Point
    var locationHash int
    var stopwords string
    err := rows.Scan(&uid, &name, &source, &location, &url, &geoCoord, &locationHash, &stopwords)
    checkErr(err)
    miner = Miner {
      Uid: uid,
      Name: name,
      Source: source,
      Location: location,
      GeoCoord: geoCoord,
      Url: url,
      Stopwords: stopwords,
    }
  }
  return
}

// Miners admin home page
func AdminMiners(w http.ResponseWriter, r *http.Request) {

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
}

func AdminNewMiner(w http.ResponseWriter, r *http.Request) {
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
    content["Title"] = "Miners Admin: Add New Miner"
    // fmt.Fprintf(w, "<p>NEW MINER</p>")
    renderTemplate(w, "admin/miners/new", content)
  }
}

func AdminEditMiner(w http.ResponseWriter, r *http.Request) {
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
    content["Title"] = "Miners Admin: Add New Miner"

    uidParam := r.URL.Query().Get("uid")
    uidConv, _ := strconv.ParseInt(uidParam, 10, 0)
    miner, err := GetMiner(int(uidConv))

    if err != nil {
      content["Error"] = "Could not retrieve miner"
    } else {
      content["Miner"] = miner
    }

    // fmt.Fprintf(w, "<p>EDIT MINER {{ .Miner.uid }} </p>")
    renderTemplate(w, "admin/miners/edit", content)
  }
}

func AdminUpdateMiner(w http.ResponseWriter, r *http.Request) {
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
    content["Title"] = "Miners Admin: Update Miner"

    uidParam := r.URL.Query().Get("uid")
    uidConv, _ := strconv.ParseInt(uidParam, 10, 0)
    miner, err := GetMiner(int(uidConv))

    if err != nil {
      content["Error"] = "Could not retrieve miner"
    } else {
      content["Miner"] = miner
    }

    fmt.Fprintf(w, "<p>UPDATE MINER AND REDIRECT</p>")
  }
}

func AdminDeleteMiner(w http.ResponseWriter, r *http.Request) {
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
    content["Title"] = "Miners Admin: Delete Miner"

    uidParam := r.URL.Query().Get("uid")
    uidConv, _ := strconv.ParseInt(uidParam, 10, 0)
    miner, err := GetMiner(int(uidConv))

    if err != nil {
      content["Error"] = "Could not retrieve miner"
    } else {
      content["Miner"] = miner
    }

    fmt.Fprintf(w, "<p>DELETE MINER AND REDIRECT</p>")
  }
}

func AdminMinersResetDatabase(w http.ResponseWriter, r *http.Request) {

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
    ResetMinersDatabase()

    AdminMiners(w, r)
  }
}

// Creates a new miner
func AdminCreateMiner(w http.ResponseWriter, r *http.Request) {
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

    err := r.ParseForm()
    if err != nil {
      fmt.Println(err)
    }

    content := make(map[string]interface{})

    name := r.PostFormValue("name")
    url := r.PostFormValue("url")
    location := r.PostFormValue("location")
    latitude := r.PostFormValue("latitude")
    longitude := r.PostFormValue("longitude")
    source := r.PostFormValue("source")
    stopwords := r.PostFormValue("stopwords")
    lastInsertId, err := InsertMiner(name, location, latitude, longitude, source, url, stopwords)
    if err != nil {
      content["MinerError"] = err
    }
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
}

// Handles receipt of post from a Miner
func MinerPost(w http.ResponseWriter, r *http.Request) {
  decoder := json.NewDecoder(r.Body)
  var posts MinerPostsJSON
  err := decoder.Decode(&posts)
  if err != nil {
    fmt.Println("Error:", err)
    http.Error(w, err.Error(), 500)
    return
  }

  http.Error(w, "OK", 200)



  //fmt.Println("Post from Miner Id: ", posts.MinerId)
  //fmt.Println(posts)

  minerConv, _ := strconv.ParseInt(posts.MinerId, 10, 0)
  var miner Miner
  rows := QueryMinerForId(int(minerConv))
  for rows.Next() {
    var uid int
    var name string
    var source string
    var location string
    var url string
    var geoCoord Point
    var locationHash int
    var stopwords string
    err := rows.Scan(&uid, &name, &source, &location, &url, &geoCoord, &locationHash, &stopwords)
    checkErr(err)
    miner = Miner {
      Uid: uid,
      Name: name,
      Source: source,
      Location: location,
      GeoCoord: geoCoord,
      Url: url,
      Stopwords: stopwords,
    }
  }

  for _, post := range posts.Posts {
    url := post.Url
    posted := post.Datetime
    mined := post.MinedAt
    /*
    fmt.Println("SourceURI ", url)
    fmt.Println("Posted ", posted)
    fmt.Println("Mined ", mined)
    fmt.Println("Terms and their counts")
    fmt.Println("Miner location ", miner.Location)
    fmt.Println("Miner source", miner.Source)
    */
    lastInsertId := InsertPost(miner.Source, miner.Location, url, posted.Time, mined.Time)
    if lastInsertId != 0 {
      for k, v := range post.Terms {
        //fmt.Println(k, v)
        InsertTerm(miner.Location, k, v, lastInsertId, posted.Time)
      }
    }
  }
}