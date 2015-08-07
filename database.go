package main

import (
    "database/sql"
    "fmt"
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
    Terms: "CREATE TABLE IF NOT EXISTS terms(uid serial NOT NULL, postid integer, term text)",
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
    lastInsertId := InsertPost("http://someblog.com")
    InsertTerm("GIS", lastInsertId)
    InsertTerm("earthquake", lastInsertId)
    lastInsertId = InsertPost("http://example.com")
    InsertTerm("GIS", lastInsertId)
    InsertTerm("water", lastInsertId)
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
    fmt.Println("# Updating")
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

func InsertTerm(term string, postid int) {
    fmt.Println("# Inserting term")

    db := ConnectToDatabase()
    var lastInsertId int
    err := db.QueryRow("INSERT INTO terms (postid, term) VALUES($1,$2) returning uid;", postid, term).Scan(&lastInsertId)
    checkErr(err)
}

func InsertPost(sourceURI string) int {
    fmt.Println("# Inserting values")
    db := ConnectToDatabase()
    var lastInsertId int
    err := db.QueryRow("INSERT INTO posts (mined, posted, sourceURI) VALUES($1,$2,$3) returning uid;", datetime.Format(time.RFC3339), datetime.Format(time.RFC3339), sourceURI).Scan(&lastInsertId)
    checkErr(err)

    return lastInsertId
}

func QueryTerms(term string) *sql.Rows {
    fmt.Println("# Querying")
    db := ConnectToDatabase()
    rows, err := db.Query("SELECT * FROM terms WHERE LOWER(term) LIKE '%' || LOWER($1) || '%'", term)
    checkErr(err)

    return rows
}

func QueryPosts(args ...string) *sql.Rows {
    fmt.Println("# Querying ", args)
    db := ConnectToDatabase()
    query := "SELECT * FROM posts"
    var buffer = bytes.NewBufferString(query)
    for _, v := range args {
        buffer.WriteString(fmt.Sprint(v, " "))
    }
    query = buffer.String()
    fmt.Println(query)
    rows, err := db.Query(query)
    checkErr(err)

    return rows
}

func QueryAll() {
    fmt.Println("# Querying")
    db := ConnectToDatabase()
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
    fmt.Println("# Deleting")
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