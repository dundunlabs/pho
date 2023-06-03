package tra

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type route struct {
	path    string
	method  string
	handler Handler
}

func parseParams(pattern string, path string) Params {
	before, after, found := strings.Cut(pattern, ":")

	if found {
		k, pattern, _ := strings.Cut(after, "/")
		v, path, _ := strings.Cut(strings.TrimPrefix(path, before), "/")
		params := parseParams(pattern, path)
		params[k] = v
		return params
	}

	return Params{}
}

func (route route) serveHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.EscapedPath()
	v := route.handler(&Context{
		Context: r.Context(),
		Request: r,
		Writer:  w,
		route:   route,
		Params:  parseParams(route.path, path),
	})

	switch v := v.(type) {
	case string:
		fmt.Fprintln(w, v)
	default:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(v)
	}
}
