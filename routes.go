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
        "TrendSourcesCSV",
        "GET",
        "/v1/locations/{location}/trends/{term}/csv",
        TrendSourcesCSV,
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
        "/web/stats",
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
        "AdminAddStopwords",
        "GET",
        "/admin/addstopwords",
        AdminAddStopwords,
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
        "AdminNewMiner",
        "GET",
        "/admin/miners/new",
        AdminNewMiner,
    },
    Route{
        "AdminCreateMiner",
        "POST",
        "/admin/miners",
        AdminCreateMiner,
    },
    Route{
        "AdminEditMiner",
        "GET",
        "/admin/miners/{uid}/edit",
        AdminEditMiner,
    },
    Route{
        "AdminUpdateMiner",
        "PATCH",
        "/admin/miners/{uid}",
        AdminUpdateMiner,
    },
    Route{
        "AdminUpdateMiner",
        "PUT",
        "/admin/miners/{uid}",
        AdminUpdateMiner,
    },
    Route{
        "AdminDeleteMiner",
        "DELETE",
        "/admin/miners/{uid}",
        AdminDeleteMiner,
    },
    Route{
        "MinerPost",
        "POST",
        "/v1/minerpost",
        MinerPost,
    },
}