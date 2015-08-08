package main

import (
    "database/sql"
    "fmt"
    "strings"
    "regexp"
    _ "github.com/lib/pq"
    "time"
    "bytes"
)

const (
    DB_USER     = "tfw"
    DB_PASSWORD = "tfw"
    DB_NAME     = "tfw"
)

const (
    Posts = iota
    Terms
)

var tables = map[int]string{
    0: "Posts",
    1: "Terms",
}

var datetime = time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)

// Post: uid serial mined:datetime posted:datetime sourceURI: string
var CREATE = map[int]string{
    Posts: "CREATE TABLE IF NOT EXISTS posts(uid serial NOT NULL, mined timestamp without time zone, posted timestamp without time zone, sourceURI text)",
    Terms: "CREATE TABLE IF NOT EXISTS terms(uid serial NOT NULL, postid integer, term text,  wordcount integer)",
}

var DROP = map[int]string{
    Posts: "DROP TABLE posts",
    Terms: "DROP TABLE terms",
}

func BuildDatabase() {

    DropTable(DROP[Posts])
    DropTable(DROP[Terms])
    CreateTable(CREATE[Posts])
    CreateTable(CREATE[Terms])

    // Seed with some sample data
    AddTweet("https://twitter.com/assaadrazzouk/status/629956498349756416",
        "AssaadRazzouk: New York Ski Resorts Turn to #Solar Power. Because It's Cheaper and Cleaner. http://t.co/eYEXmdD6ql #climate #nyc http://t.co/myVa1yi15s")

    AddTweet("https://twitter.com/oldpicsarchive/status/630020696589139968",
        "oldpicsarchive: 1910s: etiquette warnings shown before silent movies http://t.co/rQsLenicva")

    AddTweet("https://twitter.com/news24hbgd/status/630041104453210112",
        "news24hbgd: Tech talk on solar PV electricity: The Department of Electrical and Electronic Engineering (EEE) of Daffodil I... http://t.co/FLU5PdThuq")

    AddTweet("https://twitter.com/psfk/status/630033528495960064",
        "PSFK: Solar-powered coffee, toast, waffles--or whatever else you can plug into an outlet http://t.co/LjediOVYno http://t.co/qBppQn2FbP")

    AddTweet("https://twitter.com/tilapya_/status/630043313245065216",
        "tilapya_: Confused about the kind of plug China uses. Sometimes it’s the U.S. plug (but always at 220V), sometimes it’s the Australian or UK type. :|")
}

func AddTweet(address string, contents string) {
    lastInsertId := InsertPost(address)

    reg, err := regexp.Compile("[^A-Za-z ]+")
    if err != nil {
        fmt.Println("%s", err)
    }
    cleanedContent := reg.ReplaceAllString(strings.ToLower(string(contents)), "")
    wordCounts := CountWords(cleanedContent)
    for k, v := range wordCounts {
      InsertTerm(k, v, lastInsertId)
    }
}


func CountWords(s string) map[string]int {
  counts := make(map[string]int)
  fields := strings.Fields(s)
  for i := 0; i < len(fields); i++ {
    counts[fields[i]]++
  }
  return counts
}

func ConnectToDatabase() *sql.DB {
    dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
        DB_USER, DB_PASSWORD, DB_NAME)
    db, err := sql.Open("postgres", dbinfo)
    checkErr(err)

    return db
}

func CreateTable(sql string) {
    fmt.Println("# Creating table " + sql)

    if _, err := db.Exec(sql); err != nil {
        checkErr(err)
    }
}

func DropTable(sql string) {
    fmt.Println("# Dropping table" + sql)
    db := ConnectToDatabase()
    if _, err := db.Exec(sql); err != nil {
        checkErr(err)
    }
}


func UpdatePost() {
    //db := ConnectToDatabase()
    /*
    stmt, err := db.Prepare("update userinfo set username=$1 where uid=$2")
    checkErr(err)

    res, err := stmt.Exec("astaxieupdate", lastInsertId)
    checkErr(err)

    affect, err := res.RowsAffected()
    checkErr(err)

    fmt.Println(affect, "rows changed")
    */
}

func InsertTerm(term string, wordcount int, postid int) {
    var lastInsertId int
    err := db.QueryRow("INSERT INTO terms (postid, term, wordcount) VALUES($1,$2,$3) returning uid;", postid, term, wordcount).Scan(&lastInsertId)
    checkErr(err)
}

func InsertPost(sourceURI string) int {
    var lastInsertId int
    err := db.QueryRow("INSERT INTO posts (mined, posted, sourceURI) VALUES($1,$2,$3) returning uid;", datetime.Format(time.RFC3339), datetime.Format(time.RFC3339), sourceURI).Scan(&lastInsertId)
    checkErr(err)

    return lastInsertId
}

func QueryTerms(term string) *sql.Rows {
    rows, err := db.Query("SELECT * FROM terms WHERE LOWER(term) LIKE '%' || LOWER($1) || '%' ORDER BY term", term)
    checkErr(err)

    return rows
}

func QueryTermsForPost(postid int) *sql.Rows {
    rows, err := db.Query("SELECT * FROM terms WHERE postid=$1", postid)
    checkErr(err)

    return rows
}

func QueryPosts(args ...string) *sql.Rows {
    query := "SELECT * FROM posts"
    var buffer = bytes.NewBufferString(query)
    for _, v := range args {
        buffer.WriteString(fmt.Sprint(v, " "))
    }
    query = buffer.String()
    rows, err := db.Query(query)
    checkErr(err)

    return rows
}

func QueryAll() {
    rows, err := db.Query("select Posts.*, Terms.term from posts inner join terms on terms.postid=posts.uid")
    checkErr(err)

    fmt.Println("uid | minded | posted | sourceURI | term")
    for rows.Next() {
        var uid int
        var minded time.Time
        var posted time.Time
        var sourceURI string
        var term string
        err = rows.Scan(&uid, &minded, &posted, &sourceURI, &term)
        checkErr(err)
        fmt.Printf("%3v | %6v | %6v | %6v | %6v\n", uid, minded, posted, sourceURI, term)
    }
}

func DeletePost(db *sql.DB, uid int) {
    stmt, err := db.Prepare("delete from posts where uid=$1")
    checkErr(err)

    res, err := stmt.Exec(uid)
    checkErr(err)

    affect, err := res.RowsAffected()
    checkErr(err)

    fmt.Println(affect, "rows changed")
}

func checkErr(err error) {
    if err != nil {
        panic(err)
    }
}