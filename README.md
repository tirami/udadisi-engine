# Udadisi Engine

### Server
    go get github.com/tirami/udadisi-engine

or if you already have an older version installed

    go get -u github.com/tirami/udadisi-engine
    
### Configuring and Populating Server

1. Go to localhost:8080
2. Go to Admin Home
3. Select Build Database (this will delete all existing data and setup the new database schema)
4. Select Build Seeds to populate the Seed sources
5. Select Build Data to populate the data based on the Seed sources
6. Select View Seeds to see a list of the Seed sources, their Miner type and Location

### Sample Data Viewer

1. Go to localhost:8080
2. Select various data queries for each of the locations


A trival server that responds to the following:

* localhost:8080 - returns HTML
* localhost:8080/v1/trends/{location}?limit={limit} - returns top {limit} trends as JSON
* localhost:8080/v1/trends/{location}/{term} - returns JSON
* localhost:8080/web/trends/{location} - returns HTML list of terms, source URI, word counts
* localhost:8080/web/trends/{location}/{term} - returns HTML list of for term, source URIs and word counts
* localhost:8080 - returns simple home page

API spec in Swagger:

* localhost:8080/v1/swagger.json

### Environment variables
The server uses the following environment variables:

* POSTGRES_DB - database host address (defaults to localhost if not set)
* DB_PASSWORD - db user password (defaults to udadisi if not set)

### Sample Data
Currently the data is populated with 200 tweets from 1 offical source Twitter account and 5 additional Twitter accounts that I randomly selected. Some of them have related terms.

### Setting up Postgres database
    createuser --createdb --login -P udadisi

Set password to udadisi

    createdb udadisi

### Building Docker version
    docker build -t udadis_postgresql .

### Running Docker
    docker run --rm -p 8080:8080 -e POSTGRES_DB='<database host address>' -e DB_PASSWORD='<password>' -P --name pg_test udadis_postgresql
