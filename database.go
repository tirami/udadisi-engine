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
    Posts: "DROP TABLE IF EXISTS posts",
    Terms: "DROP TABLE IF EXISTS terms",
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
    CreateTable(CREATE[Posts])
    CreateTable(CREATE[Terms])

    // Seed with some sample data
    fmt.Println("# Populating with seed data...")
    anaconda.SetConsumerKey(TWITTER_CONSUMER_KEY)
    anaconda.SetConsumerSecret(TWITTER_CONSUMER_SECRET)

    PopulateWithTweets("symboticaandrew")
    PopulateWithTweets("digitalwestie")
    PopulateWithTweets("mashable")
    PopulateWithTweets("pluginadventure")
    PopulateWithTweets("kickstarter")

    fmt.Println("# Populated")

    /*
    searchResult, _ := api.GetSearch("solar", nil)
    for _ , tweet := range searchResult.Statuses {
        fmt.Println(tweet.Text)
    }
    */
}

func PopulateWithTweets(user string) {
    api := anaconda.NewTwitterApi(TWITTER_ACCESS_TOKEN, TWITTER_ACCESS_TOKEN_SECRET)

    v := url.Values{}
    v.Set("screen_name", user)
    v.Set("count", "200")
    timeline, _ := api.GetUserTimeline(v)
    for _ , tweet := range timeline {
        tweetUrl := fmt.Sprintf("https://twitter.com/%s/status/%s", user, tweet.IdStr)
        AddTweet(tweetUrl, tweet.Text)
    }
}

func AddTweet(address string, contents string) {
    lastInsertId := InsertPost(address)

    reg, err := regexp.Compile("[^A-Za-z @]+")
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