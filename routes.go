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
    router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("static/"))))
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
        "TrendsRootIndex",
        "GET",
        "/v1/locations/{location}/trends",
        TrendsRootIndex,
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
        "WebStats",
        "GET",
        "web/stats",
        WebStats,
    },
    Route{
        "AdminLogin",
        "GET",
        "/admin/login",
        AdminLogin,
    },
    Route{
        "AdminLogin",
        "POST",
        "/admin/login",
        AdminLogin,
    },
    Route{
        "AdminLogout",
        "GET",
        "/admin/logout",
        AdminLogout,
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
        "AdminCreateIndexes",
        "GET",
        "/admin/createindexes",
        AdminCreateIndexes,
    },
    Route{
        "AdminClearData",
        "GET",
        "/admin/cleardata",
        AdminClearData,
    },
    Route{
        "AdminMiners",
        "GET",
        "/admin/miners",
        AdminMiners,
    },
    Route{
        "AdminResetMinersDatabase",
        "GET",
        "/admin/miners/resetdatabase",
        AdminMinersResetDatabase,
    },
    Route{
        "AdminMiners",
        "POST",
        "/admin/miners",
        AdminCreateMiner,
    },
    Route{
        "MinerPost",
        "POST",
        "/v1/minerpost",
        MinerPost,
    },
}