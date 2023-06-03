package tra

import (
	"context"
	"net/http"
	"strconv"
)

type Context struct {
	context.Context

	Request *http.Request
	Writer  http.ResponseWriter

	route  route
	params Params
}

func (c *Context) Route() string {
	return c.route.path
}

func (c *Context) Params() Params {
	if c.params == nil {
		c.params = parseParams(c.route.path, c.Request.URL.EscapedPath())
	}
	return c.params
}

func (c *Context) Param(key string) string {
	return c.Params()[key]
}

func (c *Context) ParamInt(key string) (int, error) {
	return strconv.Atoi(c.Param(key))
}
