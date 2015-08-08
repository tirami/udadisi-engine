# tfw-application-server

Tirami TFW Application Server

### Server
    go get github.com/tirami/tfw-application-server

or if you already have an older version installed
    go get -u github.com/tirami/tfw-application-server

A trival server that responds to the following:

* localhost:8080 - returns HTML
* localhost:8080/trends/{location} - returns JSON
* localhost:8080/trends/{location}/{term} - returns JSON
* localhost:8080/web/trends/{location} - returns HTML list of terms, source URI, word counts
* localhost:8080/web/trends/{location}/{term} - returns HTML list of for term, source URIs and word counts

### Sample Data
Currently the data is populated with 5 tweets that I randomly selected. Some of them have related terms.

### Setting up Postgres database
    createuser --createdb --login -P tfw
    createdb tfw