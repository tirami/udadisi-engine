package main

import (
    "database/sql"
    "github.com/jmoiron/sqlx"
    "fmt"
    "strings"
    //"regexp"
    _ "github.com/lib/pq"
    "time"
    "bytes"
    "os"
    "hash/fnv"
)

const (
    DB_USER     = "udadisi"
    DB_PASSWORD = "udadisi"
    DB_NAME     = "udadisi"
)

const (
    Posts = iota
    Terms
    MinersTable
)

var tables = map[int]string{
    0: "Posts",
    1: "Terms",
    2: "MinersTable",
}

var datetime = time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)

// Post: uid serial mined:datetime posted:datetime sourceURI: string
var CREATE = map[int]string{
    Posts: "CREATE TABLE IF NOT EXISTS posts(uid serial NOT NULL, mined timestamp without time zone, posted timestamp without time zone, sourceURI text, location text, source text, locationhash bigint)",
    Terms: "CREATE TABLE IF NOT EXISTS terms(uid serial NOT NULL, postid integer, term text,  wordcount integer, posted timestamp without time zone, location text, locationhash bigint)",
    MinersTable: "CREATE TABLE IF NOT EXISTS miners(uid serial NOT NULL, name text, source text, location text, url text, geocoord point, locationhash bigint)",
}

var DROP = map[int]string{
    Posts: "DROP TABLE IF EXISTS posts",
    Terms: "DROP TABLE IF EXISTS terms",
    MinersTable: "DROP TABLE IF EXISTS miners",
}

// A DatabaseError indicates an error with the database
type DatabaseError struct {
    Error error  // The raw error that precipitated this error, if any.
}

// String returns a human-readable error message.
func (e *DatabaseError) String() string {
    return fmt.Sprintf("%s", e.Error)
}

func LocationHash(s string) uint32 {
    h := fnv.New32a()
    h.Write([]byte(s))
    return h.Sum32()
}

func BuildDatabase() {

    DropTable(DROP[Posts])
    DropTable(DROP[Terms])
    DropTable(DROP[MinersTable])
    CreateTable(CREATE[Posts])
    CreateTable(CREATE[Terms])
    CreateTable(CREATE[MinersTable])
    CreateIndexes()
    AddStopwords()
}

func AddStopwords() (err error){
    sql := "ALTER TABLE miners ADD stopwords varchar(255) DEFAULT '';"
    
    defer func() {
        if r := recover(); r != nil {
            var ok bool
            err, ok = r.(error)
            if !ok {
                err = fmt.Errorf("Database: %v", r)
            }
        }
    }()

    fmt.Println("# Adding stopwords " + sql)

    if _, err := db.Exec(sql); err != nil {
        checkErr(err)
    }

    return
}


func CreateIndexes() {
    CreateIndex("CREATE INDEX ON posts (uid);")
    CreateIndex("CREATE INDEX ON terms (uid);")
    CreateIndex("CREATE INDEX ON miners (uid);")
    CreateIndex("CREATE INDEX ON posts (locationhash);")
    CreateIndex("CREATE INDEX ON terms (locationhash);")
    CreateIndex("CREATE INDEX ON miners (locationhash);")
}

func ResetMinersDatabase() {
    DropTable(DROP[MinersTable])
    CreateTable(CREATE[MinersTable])
}

func ClearData() (err error) {

    if err := DropTable(DROP[Posts]); err != nil {
        checkErr(err)
    }
    if err := DropTable(DROP[Terms]); err != nil {
        checkErr(err)
    }
    if err := CreateTable(CREATE[Posts]); err != nil {
        checkErr(err)
    }
    if err := CreateTable(CREATE[Terms]); err != nil {
        checkErr(err)
    }

    return
}

func CountWords(s string) map[string]int {
  counts := make(map[string]int)
  fields := strings.Fields(s)
  for i := 0; i < len(fields); i++ {
    counts[fields[i]]++
  }
  return counts
}

func ConnectToDatabase() *sqlx.DB {
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
    db, err := sqlx.Open("postgres", dbinfo)
    checkErr(err)

    return db
}

func CreateIndex(sql string) (err error) {
    defer func() {
        if r := recover(); r != nil {
            var ok bool
            err, ok = r.(error)
            if !ok {
                err = fmt.Errorf("Database: %v", r)
            }
        }
    }()

    fmt.Println("# Creating index " + sql)

    if _, err := db.Exec(sql); err != nil {
        checkErr(err)
    }

    return
}

func CreateTable(sql string) (err error) {
    defer func() {
        if r := recover(); r != nil {
            var ok bool
            err, ok = r.(error)
            if !ok {
                err = fmt.Errorf("Database: %v", r)
            }
        }
    }()

    fmt.Println("# Creating table " + sql)

    if _, err := db.Exec(sql); err != nil {
        checkErr(err)
    }

    return
}

func DropTable(sql string) (err error) {
    defer func() {
        if r := recover(); r != nil {
            var ok bool
            err, ok = r.(error)
            if !ok {
                err = fmt.Errorf("Database: %v", r)
            }
        }
    }()

    fmt.Println("# Dropping table " + sql)
    db := ConnectToDatabase()
    if _, err := db.Exec(sql); err != nil {
        checkErr(err)
    }

    return
}

func InsertMiner(name string, location string, latitude string, longitude string, source string, url string, stopwords string) (lastInsertId int, err error) {
    defer func() {
        if r := recover(); r != nil {
            var ok bool
            err, ok = r.(error)
            if !ok {
                err = fmt.Errorf("Database: %v", r)
            }
        }
    }()

    if latitude == "" {
        latitude = "0"
    }
    if longitude == "" {
        longitude = "0"
    }

    fmt.Println(LocationHash(location))

    err = db.QueryRow("INSERT INTO miners (name, location, geocoord, source, url, locationhash, stopwords) VALUES($1,$2,POINT($3,$4),$5,$6,$7,$8) returning uid;", name, location, latitude, longitude, source, url, LocationHash(location), stopwords).Scan(&lastInsertId)
    checkErr(err)

    return
}

func InsertTerm(location string, term string, wordcount int, postid int, posted time.Time) {
    var lastInsertId int
    err := db.QueryRow("INSERT INTO terms (postid, term, wordcount, posted, location, locationhash) VALUES($1,$2,$3,$4,$5,$6) returning uid;", postid, strings.ToLower(term), wordcount, posted.Format(time.RFC3339), location,LocationHash(location)).Scan(&lastInsertId)
    checkErr(err)
}

func InsertPost(source string, location string, sourceURI string, postedAt time.Time, minedAt time.Time) int {
    lastInsertId := 0

    // Check to see if we already have an entry for the sourceURI
    duplicate := false
    db.QueryRow("SELECT 1 FROM posts WHERE sourceURI=$1 AND locationhash=$2", sourceURI, LocationHash(location)).Scan(&duplicate)

    if duplicate == false {
       err := db.QueryRow("INSERT INTO posts (source, location, mined, posted, sourceURI, locationhash) VALUES($1,$2,$3,$4,$5,$6) returning uid;", source, location, minedAt.Format(time.RFC3339), postedAt.Format(time.RFC3339), sourceURI,LocationHash(location)).Scan(&lastInsertId)
        checkErr(err)
    } else {
        fmt.Println("We already have", sourceURI, "for", location)
    }

    return lastInsertId
}

func DatabasePostsCount(location string) (count int, err error) {
    defer func() {
        if r := recover(); r != nil {
            var ok bool
            err, ok = r.(error)
            if !ok {
                err = fmt.Errorf("Database: %v", r)
            }
        }
    }()

    if location != "all" {
        locationhash := LocationHash(location)
        errDb := db.QueryRow("SELECT count(uid) as count FROM posts where locationhash=$1", locationhash).Scan(&count)
        checkErr(errDb)
    } else {
        errDb := db.QueryRow("SELECT count(uid) as count FROM posts").Scan(&count)
        checkErr(errDb)    }
    return
}

func DatabaseLastMined(location string) (mined time.Time, err error) {
    defer func() {
        if r := recover(); r != nil {
            var ok bool
            err, ok = r.(error)
            if !ok {
                err = fmt.Errorf("Database: %v", r)
            }
        }
    }()

    if location != "all" {
        locationhash := LocationHash(location)
        errDb := db.QueryRow("SELECT mined FROM posts where locationhash=$1 ORDER BY mined DESC LIMIT 1", locationhash).Scan(&mined)
        checkErr(errDb)
    } else {
        errDb := db.QueryRow("SELECT mined FROM posts ORDER BY mined DESC LIMIT 1").Scan(&mined)
        checkErr(errDb)
    }
    return
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

func QueryStopwordsFor(location string, source string) (rows *sql.Rows, err error) {
    
    locationCondition := ""
    if (location == "all") || (location == "") { 
        locationCondition =  "locationhash > 0"
    } else {
        lh := LocationHash(location)
        locationhash := fmt.Sprint(lh)
        locationCondition =  "locationhash = " + locationhash
    }

    sourceCondition := "source = '" + source + "'"
    if (source == "all") || (source == "") { 
        sourceCondition = "source LIKE '%'"
    }

    statement := "SELECT stopwords FROM miners WHERE " + locationCondition + " AND " + sourceCondition + ";"
    rows, err = db.Query(statement)
    
    checkErr(err);
    return
}

func QueryTerms(source string, location string, term string, fromDate string, toDate string) (rows *sql.Rows, err error) {
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

    if location != "" {
        locationhash := LocationHash(location)

        if term != "" {
            rows, err = db.Query("SELECT terms.*, posts.source FROM terms, posts WHERE terms.postid=posts.uid AND posts.locationhash = $4 AND terms.posted between $1 AND $2 AND LOWER(term) LIKE LOWER($3) AND (LOWER(source) = LOWER($5) OR $5 = '') ORDER BY terms.posted, term", fromTime.Format(time.RFC3339), toTime.Format(time.RFC3339), term, locationhash, source)
        } else {
            rows, err = db.Query("SELECT terms.*, posts.source FROM terms, posts WHERE terms.postid=posts.uid AND posts.locationhash = $3 AND terms.posted between $1 AND $2 AND (LOWER(posts.source) = LOWER($4) OR $4 = '') ORDER BY terms.posted, term", fromTime.Format(time.RFC3339), toTime.Format(time.RFC3339), locationhash, source)
        }
    } else {
        if term != "" {
            rows, err = db.Query("SELECT terms.*, posts.source FROM terms, posts WHERE terms.postid=posts.uid AND terms.posted between $1 AND $2 AND LOWER(term) LIKE LOWER($3) AND (LOWER(source) = LOWER($4) OR $4 = '') ORDER BY terms.posted, term", fromTime.Format(time.RFC3339), toTime.Format(time.RFC3339), term, source)
        } else {
            rows, err = db.Query("SELECT terms.*, posts.source FROM terms, posts WHERE terms.postid=posts.uid AND terms.posted between $1 AND $2 AND (LOWER(posts.source) = LOWER($3) OR $3 = '') ORDER BY terms.posted, term", fromTime.Format(time.RFC3339), toTime.Format(time.RFC3339), source)
        }
    }
    checkErr(err)
    return
}

func QueryTermsForPost(postid int) *sql.Rows {
    rows, err := db.Query("SELECT terms.*, posts.source FROM terms, posts WHERE terms.postid=posts.uid AND postid=$1", postid)
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

func DeleteMiner(uid int) (affected int64, err error) {
    stmt, err := db.Prepare("DELETE FROM miners where uid=$1")
    checkErr(err)
    
    res, err := stmt.Exec(uid)
    checkErr(err)
    
    affected, err = res.RowsAffected()
    checkErr(err)
    
    return
}

func DeletePost(db *sqlx.DB, uid int) {
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
        fmt.Println("Error:", &DatabaseError{err})
        panic(&DatabaseError{err})
    }
}