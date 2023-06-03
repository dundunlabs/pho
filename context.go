package tra

import (
	"context"
	"net/http"
)

type Context struct {
	context.Context
	Params

	Request *http.Request
	Writer  http.ResponseWriter

	route route
}

func (c *Context) Route() string {
	return c.route.path
}
