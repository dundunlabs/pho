package main

import (
	"net/http"

	"github.com/dundunlabs/httplog"
	"github.com/dundunlabs/tra"
)

var router = tra.NewRouter()

func init() {
	router.GET("/", homeHandler)

	router.WithGroup("/api", func(g *tra.Group) {
		g.GET("/users", userListHandler)
	}, apiMiddleware)
}

func main() {
	server := httplog.NewHandler(router)
	http.ListenAndServe(":8080", server)
}
