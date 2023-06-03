package tra

import (
	"context"
	"net/http"
)

type Context struct {
	context.Context
	Request *http.Request
	Writer  http.ResponseWriter
}
