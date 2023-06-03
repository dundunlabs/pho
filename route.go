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
	params := Params{}
	keys := strings.Split(pattern, "/")
	values := strings.Split(path, "/")

	for i, k := range keys {
		if key, ok := strings.CutPrefix(k, ":"); ok {
			params[key] = values[i]
		}
		if key, ok := strings.CutPrefix(k, "*"); ok {
			params[key] = strings.Join(values[i:], "/")
		}
	}

	return params
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
