package main

import (
  "net/http"

  "github.com/gorilla/mux"
)

type Route struct {
    Name        string
    Method      string
    Pattern     string
    HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {

    router := mux.NewRouter().StrictSlash(true)
    for _, route := range routes {
        router.
            Methods(route.Method).
            Path(route.Pattern).
            Name(route.Name).
            Handler(route.HandlerFunc)
    }

    return router
}

var routes = Routes{
    Route{
        "Swagger",
        "GET",
        "/v1/swagger.json",
        Swagger,
    },
    Route{
        "Index",
        "GET",
        "/",
        Index,
    },
    Route{
        "TrendsIndex",
        "GET",
        "/v1/trends/{location}/{term}",
        TrendsIndex,
    },
    Route{
        "TrendsRouteIndex",
        "GET",
        "/v1/trends/{location}",
        TrendsRouteIndex,
    },
    Route{
        "WebTrendsIndex",
        "GET",
        "/web/trends/{location}",
        WebTrendsRouteIndex,
    },
    Route{
        "WebTrendsIndex",
        "GET",
        "/web/trends/{location}/{term}",
        WebTrendsIndex,
    },
}