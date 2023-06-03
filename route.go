package tra

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type route struct {
	path    string
	method  string
	handler Handler
}

func (route route) serveHTTP(w http.ResponseWriter, r *http.Request) {
	v := route.handler(&Context{
		Context: r.Context(),
		Request: r,
		Writer:  w,
		route:   route,
	})

	switch v := v.(type) {
	case string:
		fmt.Fprintln(w, v)
	default:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(v)
	}
}
