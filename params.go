package tra

import (
	"strings"
)

type Params map[string]string

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
