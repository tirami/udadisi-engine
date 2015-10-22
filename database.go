package main

import (
    "database/sql"
    "fmt"
    "strings"
    "regexp"
    _ "github.com/lib/pq"
    "time"
    "bytes"
    "os"
    "github.com/ChimeraCoder/anaconda"
    "net/url"
)

const (
    DB_USER     = "udadisi"
    DB_PASSWORD = "udadisi"
    DB_NAME     = "udadisi"
)

const (
    Posts = iota
    Terms
    SeedsTable
    MinersTable
)

var tables = map[int]string{
    0: "Posts",
    1: "Terms",
    2: "SeedsTable",
    3: "MinersTable",
}

var datetime = time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)

// Post: uid serial mined:datetime posted:datetime sourceURI: string
var CREATE = map[int]string{
    Posts: "CREATE TABLE IF NOT EXISTS posts(uid serial NOT NULL, mined timestamp without time zone, posted timestamp without time zone, sourceURI text, location text)",
    Terms: "CREATE TABLE IF NOT EXISTS terms(uid serial NOT NULL, postid integer, term text,  wordcount integer, posted timestamp without time zone, location text)",
    SeedsTable: "CREATE TABLE IF NOT EXISTS seeds(uid serial NOT NULL, minertype text, location text, source text)",
    MinersTable: "CREATE TABLE IF NOT EXISTS miners(uid serial NOT NULL, name text, location text, url text)",
}

var DROP = map[int]string{
    Posts: "DROP TABLE IF EXISTS posts",
    Terms: "DROP TABLE IF EXISTS terms",
    SeedsTable: "DROP TABLE IF EXISTS seeds",
}

const (
    TWITTER_CONSUMER_KEY = "nPHuNmCFnYGAJQgr0fwB571bc"
    TWITTER_CONSUMER_SECRET = "5Xd8A5WcmP8avrKwlGKsaaLQiHspe87Qiz68AxSWQ83iARieJl"
    TWITTER_ACCESS_TOKEN = "3417359543-fss1uvyGrPc9k3GpfewWBF7EZaXm8Tw8c8boN6C"
    TWITTER_ACCESS_TOKEN_SECRET = "zXngb9iZ4rsugWHpgkZHYRtIo2qAQs7dBrsfn0kLiKQJP"
)

func BuildDatabase() {

    DropTable(DROP[Posts])
    DropTable(DROP[Terms])
    DropTable(DROP[SeedsTable])
    DropTable(DROP[MinersTable])
    CreateTable(CREATE[Posts])
    CreateTable(CREATE[Terms])
    CreateTable(CREATE[SeedsTable])
    CreateTable(CREATE[MinersTable])
}

func BuildSeeds() {
    InsertSeed("twitter", "Dhaka", "a2ztechnews")
    InsertSeed("twitter", "Scotland", "symboticaandrew")
    InsertSeed("twitter", "Scotland", "digitalwestie")
    InsertSeed("twitter", "Scotland", "pluginadventure")
    InsertSeed("twitter", "UK", "pluginadventure")
    InsertSeed("twitter", "UK", "kickstarter")
}

func BuildWithTweets() {
    // Seed with some sample data
    fmt.Println("# Populating with seed data...")
    anaconda.SetConsumerKey(TWITTER_CONSUMER_KEY)
    anaconda.SetConsumerSecret(TWITTER_CONSUMER_SECRET)

    rows := QuerySeeds()
    for rows.Next() {
        var uid int
        var miner string
        var location string
        var source string
        err := rows.Scan(&uid, &miner, &location, &source)
        checkErr(err)

        PopulateWithTweets(location, source)
    }

    fmt.Println("# Populated")
}

func PopulateWithTweets(location string, user string) {
    api := anaconda.NewTwitterApi(TWITTER_ACCESS_TOKEN, TWITTER_ACCESS_TOKEN_SECRET)

    v := url.Values{}
    v.Set("screen_name", user)
    v.Set("count", "200")
    timeline, _ := api.GetUserTimeline(v)
    for _ , tweet := range timeline {
        tweetUrl := fmt.Sprintf("https://twitter.com/%s/status/%s", user, tweet.IdStr)
        AddTweet(location, tweetUrl, tweet.Text, tweet.CreatedAt)
    }
}

func AddTweet(location string, address string, contents string, createdAt string) {
    lastInsertId := InsertPost(location, address, createdAt)

    reg, err := regexp.Compile("[^A-Za-z @]+")
    if err != nil {
        fmt.Println("%s", err)
    }
    cleanedContent := reg.ReplaceAllString(strings.ToLower(string(contents)), "")
    wordCounts := CountWords(cleanedContent)
    for k, v := range wordCounts {
      InsertTerm(location, k, v, lastInsertId, createdAt)
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
    host := os.Getenv("POSTGRES_DB")
    if host == "" {
        host = "localhost"
    }

    password := os.Getenv("DB_PASSWORD")
    if password == "" {
        password = DB_PASSWORD
    }

    dbinfo := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
        host, DB_USER, password, DB_NAME)
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
    fmt.Println("# Dropping table " + sql)
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

func InsertMiner(name string, location string, url string) {
    var lastInsertId int
    err := db.QueryRow("INSERT INTO miners (name, location, url) VALUES($1,$2,$3) returning uid;", name, location, url).Scan(&lastInsertId)
    checkErr(err)
}

func InsertSeed(miner string, location string, source string) {
    var lastInsertId int
    err := db.QueryRow("INSERT INTO seeds (minertype, location, source) VALUES($1,$2,$3) returning uid;", miner, location, source).Scan(&lastInsertId)
    checkErr(err)
}

func InsertTerm(location string, term string, wordcount int, postid int, posted string) {
    var lastInsertId int
    err := db.QueryRow("INSERT INTO terms (postid, term, wordcount, posted, location) VALUES($1,$2,$3,$4,$5) returning uid;", postid, term, wordcount, posted, location).Scan(&lastInsertId)
    checkErr(err)
}

func InsertPost(location string, sourceURI string, createdAt string) int {
    var lastInsertId int
    err := db.QueryRow("INSERT INTO posts (location, mined, posted, sourceURI) VALUES($1,$2,$3,$4) returning uid;", location, datetime.Format(time.RFC3339), createdAt, sourceURI).Scan(&lastInsertId)
    checkErr(err)

    return lastInsertId
}

func QueryMiners() *sql.Rows {
    rows, err := db.Query("SELECT * FROM miners")
    checkErr(err)

    return rows
}

func QueryTerms(location string, term string, fromDate string, interval int) *sql.Rows {
    t, err := time.Parse("20060102", fromDate)
    if err != nil {
        fmt.Errorf("invalid date: %v", err)
    }

    if interval > 0 {
        toDate := t.Add(time.Duration(interval) * time.Hour * 24)

        rows, err := db.Query("SELECT * FROM terms WHERE LOWER(location) LIKE '%' || LOWER($4) || '%' AND posted between $1 AND $2 AND LOWER(term) LIKE '%' || LOWER($3) || '%' ORDER BY term", t.Format(time.RFC3339), toDate.Format(time.RFC3339), term, location)
        checkErr(err)
        return rows
    } else {
        rows, err := db.Query("SELECT * FROM terms WHERE LOWER(location) LIKE '%' || LOWER($1) || '%' AND posted > $2 AND LOWER(term) LIKE '%' || LOWER($3) || '%' ORDER BY term", location, t.Format(time.RFC3339), term)
        checkErr(err)
        return rows
    }


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

func QuerySeeds(args ...string) *sql.Rows {
    query := "SELECT * from seeds"
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