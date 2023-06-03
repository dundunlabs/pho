package tra

import (
	"net/http"
	"strings"
)

type Group struct {
	router     *Router
	path       string
	middleware Middleware
}

func (g *Group) Group(path string, middlewares ...Middleware) *Group {
	middleware := func(next Handler) Handler {
		n := next
		for i := range middlewares {
			mdw := middlewares[len(middlewares)-i-1]
			n = mdw(n)
		}

		if g.middleware != nil {
			return g.middleware(n)
		}
		return n
	}

	return &Group{
		router:     g.router,
		path:       g.path + path,
		middleware: middleware,
	}
}

func (g *Group) WithGroup(path string, fn func(g *Group), middlewares ...Middleware) {
	fn(g.Group(path, middlewares...))
}

func (g *Group) findNode(path string) *node {
	return g.router.root.findNode(strings.TrimPrefix(path, "/"))
}

func (g *Group) handle(method string, path string, handler Handler) {
	hdr := handler
	if g.middleware != nil {
		hdr = g.middleware(hdr)
	}

	route := route{
		method:  method,
		path:    g.path + path,
		handler: hdr,
	}

	g.router.root.addRoute(route)
}

func (g *Group) GET(path string, handler Handler) {
	g.handle(http.MethodGet, path, handler)
}

func (g *Group) POST(path string, handler Handler) {
	g.handle(http.MethodPost, path, handler)
}

func (g *Group) PUT(path string, handler Handler) {
	g.handle(http.MethodPut, path, handler)
}

func (g *Group) PATCH(path string, handler Handler) {
	g.handle(http.MethodPatch, path, handler)
}

func (g *Group) DELETE(path string, handler Handler) {
	g.handle(http.MethodDelete, path, handler)
}
