package main

import (
    "database/sql"
    "fmt"
    "strings"
    //"regexp"
    _ "github.com/lib/pq"
    "time"
    "bytes"
    "os"
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
    MinersTable: "CREATE TABLE IF NOT EXISTS miners(uid serial NOT NULL, name text, location text, url text)",
}

var DROP = map[int]string{
    Posts: "DROP TABLE IF EXISTS posts",
    Terms: "DROP TABLE IF EXISTS terms",
    SeedsTable: "DROP TABLE IF EXISTS seeds",
}



// A DatabaseError indicates an error with the database
type DatabaseError struct {
    Error error  // The raw error that precipitated this error, if any.
}

// String returns a human-readable error message.
func (e *DatabaseError) String() string {
    return fmt.Sprintf("%s", e.Error)
}


func BuildDatabase() {

    DropTable(DROP[Posts])
    DropTable(DROP[Terms])
    DropTable(DROP[SeedsTable])
    DropTable(DROP[MinersTable])
    CreateTable(CREATE[Posts])
    CreateTable(CREATE[Terms])
    CreateTable(CREATE[MinersTable])
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

func InsertMiner(name string, location string, url string) (lastInsertId int, err error) {
    defer func() {
        if r := recover(); r != nil {
            var ok bool
            err, ok = r.(error)
            if !ok {
                err = fmt.Errorf("Database: %v", r)
            }
        }
    }()

    err = db.QueryRow("INSERT INTO miners (name, location, url) VALUES($1,$2,$3) returning uid;", name, location, url).Scan(&lastInsertId)
    checkErr(err)

    return
}

func InsertSeed(miner string, location string, source string) {
    var lastInsertId int
    err := db.QueryRow("INSERT INTO seeds (minertype, location, source) VALUES($1,$2,$3) returning uid;", miner, location, source).Scan(&lastInsertId)
    checkErr(err)
}

func InsertTerm(location string, term string, wordcount int, postid int, posted time.Time) {
    var lastInsertId int
    err := db.QueryRow("INSERT INTO terms (postid, term, wordcount, posted, location) VALUES($1,$2,$3,$4,$5) returning uid;", postid, strings.ToLower(term), wordcount, posted.Format(time.RFC3339), location).Scan(&lastInsertId)
    checkErr(err)
}

func InsertPost(location string, sourceURI string, postedAt time.Time, minedAt time.Time) int {
    lastInsertId := 0

    // Check to see if we already have an entry for the sourceURI
    duplicate := false
    db.QueryRow("SELECT 1 FROM posts WHERE sourceURI=$1 AND location=$2", sourceURI, location).Scan(&duplicate)

    if duplicate == false {
       err := db.QueryRow("INSERT INTO posts (location, mined, posted, sourceURI) VALUES($1,$2,$3,$4) returning uid;", location, minedAt.Format(time.RFC3339), postedAt.Format(time.RFC3339), sourceURI).Scan(&lastInsertId)
        checkErr(err)
    }

    return lastInsertId
}

func QueryMiners() (rows *sql.Rows, err error) {

    defer func() {
        if r := recover(); r != nil {
            var ok bool
            err, ok = r.(error)
            if !ok {
                err = fmt.Errorf("Database: %v", r)
            }
        }
    }()

    rows, errDb := db.Query("SELECT * FROM miners")
    checkErr(errDb)
    return
}

func QueryMinerForId(minerId int) *sql.Rows {
    rows, err := db.Query("SELECT * FROM miners WHERE uid=$1", minerId)
    checkErr(err)

    return rows
}


func QueryTerms(location string, term string, fromDate string, toDate string) (rows *sql.Rows, err error) {
    defer func() {
        if r := recover(); r != nil {
            var ok bool
            err, ok = r.(error)
            if !ok {
                err = fmt.Errorf("Database: %v", r)
            }
        }
    }()

    fromTime, err := time.Parse("200601021504", fromDate)
    if err != nil {
        fmt.Errorf("invalid from date: %v", err)
    }

    toTime, err := time.Parse("200601021504", toDate)
    if err != nil {
        fmt.Errorf("invalid to date: %v", err)
    }

    fmt.Println("Searching for:", term, "in", location, "between", fromTime.Format(time.RFC3339), "and", toTime.Format(time.RFC3339))

    // rows, err = db.Query("SELECT * FROM terms WHERE LOWER(location) LIKE '%' || LOWER($4) || '%' AND posted between $1 AND $2 AND LOWER(term) LIKE '%' || LOWER($3) || '%' ORDER BY posted, term", fromTime.Format(time.RFC3339), toTime.Format(time.RFC3339), term, location)
    rows, err = db.Query("SELECT * FROM terms WHERE LOWER(location) LIKE '%' || LOWER($4) || '%' AND posted between $1 AND $2 AND LOWER(term) LIKE '%' || LOWER($3) || '%' ORDER BY posted, term", fromTime.Format(time.RFC3339), toTime.Format(time.RFC3339), term, location)
    checkErr(err)
    return
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
        panic(&DatabaseError{err})
    }
}