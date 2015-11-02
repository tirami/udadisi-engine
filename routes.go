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
        "Locations",
        "GET",
        "/v1/locations",
        RenderLocationsJSON,
    },
    Route{
        "LocationStats",
        "GET",
        "/v1/locations/{location}/stats",
        RenderLocationStatsJSON,
    },
    Route{
        "TrendsIndex",
        "GET",
        "/v1/locations/{location}/trends/{term}",
        TrendsIndex,
    },
    Route{
        "TrendsRouteIndex",
        "GET",
        "/v1/locations/{location}/trends",
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
    Route{
        "AdminIndex",
        "GET",
        "/admin/",
        AdminIndex,
    },
    Route{
        "AdminBuildDatabase",
        "GET",
        "/admin/builddatabase",
        AdminBuildDatabase,
    },
    Route{
        "AdminMiners",
        "GET",
        "/admin/miners",
        AdminMiners,
    },
    Route{
        "AdminMiners",
        "POST",
        "/admin/miner",
        AdminCreateMiner,
    },
    Route{
        "MinerPost",
        "POST",
        "/v1/minerpost",
        MinerPost,
    },
}