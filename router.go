package tra

import (
	"fmt"
	"net/http"
)

func NewRouter() *Router {
	r := &Router{}
	r.Group.router = r
	return r
}

type Router struct {
	Group
	root node
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.EscapedPath()
	node := router.root.findNode(path)

	if node == nil || node.routes == nil {
		http.NotFound(w, r)
		return
	}

	route, ok := node.routes[r.Method]

	if !ok {
		http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
		return
	}

	v := route.handler(&Context{
		Context: r.Context(),
		Request: r,
		Writer:  w,
	})

	switch v := v.(type) {
	case string:
		fmt.Fprint(w, v)
	}
}
