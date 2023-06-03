package tra

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupRouter() *Router {
	r := NewRouter()
	r.GET("/", func(ctx *Context) any {
		return "hello tra"
	})
	r.POST("/a/b/c", func(ctx *Context) any {
		return "abc"
	})
	return r
}

type client struct {
	http.Handler
}

func (c client) Fetch(method string, target string, body io.Reader) *http.Response {
	req := httptest.NewRequest(method, target, body)
	w := httptest.NewRecorder()
	c.ServeHTTP(w, req)
	return w.Result()
}

var c = client{
	Handler: setupRouter(),
}

func TestRootRoute(t *testing.T) {
	res := c.Fetch(http.MethodGet, "/", nil)
	body, _ := io.ReadAll(res.Body)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "hello tra", string(body))
}

func TestSimpleRoute(t *testing.T) {
	res := c.Fetch(http.MethodPost, "/a/b/c", nil)
	body, _ := io.ReadAll(res.Body)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "abc", string(body))
}

func TestMethodNotAllowedRoute(t *testing.T) {
	res := c.Fetch(http.MethodGet, "/a/b/c", nil)
	body, _ := io.ReadAll(res.Body)
	assert.Equal(t, http.StatusMethodNotAllowed, res.StatusCode)
	assert.Equal(t, "405 method not allowed\n", string(body))
}

func TestNotFoundRoute(t *testing.T) {
	res := c.Fetch(http.MethodGet, "/not-found", nil)
	body, _ := io.ReadAll(res.Body)
	assert.Equal(t, http.StatusNotFound, res.StatusCode)
	assert.Equal(t, "404 page not found\n", string(body))
}
