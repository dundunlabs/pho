package tra

import (
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
	node := router.findNode(path)

	if node == nil || node.routes == nil {
		http.NotFound(w, r)
		return
	}

	route, ok := node.routes[r.Method]

	if !ok {
		http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
		return
	}

	route.serveHTTP(w, r)
}
