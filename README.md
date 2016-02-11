# Udadisi Engine

### Server
    go get github.com/tirami/udadisi-engine

or if you already have an older version installed

    go get -u github.com/tirami/udadisi-engine

### Configuring and Populating Server

1. Go to localhost:8080/admin
2. Select Build Database (this will delete all existing data and setup the new database schema)
3. Select Miners to view the list of current Miners and to register a new Miner

### Posting from Miner to Engine

* POST JSON to localhost:8080/v1/minerpost

    {
        "posts": [{
            "terms": {
                "foo": 2,
                "bar": 1
            },
            "url": "http://www.twitter.com/post/123456",
            "datetime": 201508211014,
            "mined_at": 201508211530
        }],
        "miner_id": 1
    }

### Sample Data Viewer

1. Go to localhost:8080
2. Select various data queries for each of the locations


A trival server that responds to the following:

* localhost:8080 - returns HTML
* localhost:8080/v1/locations/{location}/trends?limit={limit} - returns top {limit} trends as JSON
* localhost:8080/v1/locations/{location}/trends/{term} - returns JSON
* localhost:8080/web/trends/{location} - returns HTML list of terms, source URI, word counts
* localhost:8080/web/trends/{location}/{term} - returns HTML list of for term, source URIs and word counts
* localhost:8080 - returns simple home page

API spec in Swagger:

* localhost:8080/v1/swagger.json

### Environment variables
The server uses the following environment variables:

* POSTGRES_DB - database host address (defaults to localhost if not set)
* DB_PASSWORD - db user password (defaults to udadisi if not set)
* ADMIN_USERNAME - username for logging into admin suite
* ADMIN_PASSWORD - password for logging into admin suite

### Setting up Postgres database
    createuser --createdb --login -P udadisi

Set password to udadisi

    createdb udadisi

### Building Docker version
    docker build -t udadis_postgresql .

### Running Docker
    docker run --rm -p 8080:8080 -e POSTGRES_DB='<database host address>' -e DB_PASSWORD='<password>' -P --name pg_test udadis_postgresql
