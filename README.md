# Udadisi Engine

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
* localhost:8080 - returns simple home page

### Environment variables
The server uses the following environment variables:

* POSTGRES_DB - database host address
* DB_PASSWORD - db user password

### Sample Data
Currently the data is populated with 5 tweets that I randomly selected. Some of them have related terms.

### Setting up Postgres database
    createuser --createdb --login -P udadisi

Set password to udadisi

    createdb udadisi

### Building Docker version
    docker build -t udadis_postgresql .

### Running Docker
    docker run --rm -p 8080:8080 -e POSTGRES_DB='<database host address>' -e DB_PASSWORD='<password>' -P --name pg_test udadis_postgresql
