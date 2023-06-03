package tra

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupRouter() *Router {
	r := NewRouter()

	r.GET("/", func(ctx *Context) any {
		return nil
	})
	r.POST("/a/b/c", func(ctx *Context) any {
		return ctx.Route()
	})

	r.WithGroup("/api", func(g *Group) {
		g.WithGroup("/users", func(g *Group) {
			g.GET("", func(ctx *Context) any {
				return []string{"user 1", "user 2"}
			})
			g.GET("/:id", func(ctx *Context) any {
				id, _ := ctx.ParamInt("id")
				return fmt.Sprintf("user %d", id)
			})
		})

		g.DELETE("/*", func(c *Context) any {
			return "api not implemented"
		})
	})

	r.WithGroup("/admin", func(g *Group) {
		g.WithGroup("", func(g *Group) {
			g.GET("/users", func(ctx *Context) any {
				return []string{"user 1", "user 2"}
			})
		}, func(next Handler) Handler {
			return func(ctx *Context) any {
				if ctx.Request.Header.Get("Authorization") == "secret" {
					return ForbiddenError
				}
				return next(ctx)
			}
		})
	}, func(next Handler) Handler {
		return func(ctx *Context) any {
			if ctx.Request.Header.Get("Authorization") != "secret" {
				return UnauthorizedError
			}
			return next(ctx)
		}
	})

	r.WithGroup("", func(g *Group) {
		g.GET("/count", func(ctx *Context) any {
			return ctx.Get("count")
		})
	}, func(next Handler) Handler {
		return func(ctx *Context) any {
			ctx.Set("count", 0)
			return next(ctx)
		}
	}, func(next Handler) Handler {
		return func(ctx *Context) any {
			count := ctx.Get("count").(int)
			ctx.Set("count", count+1)
			return next(ctx)
		}
	})

	r.PUT("/:a/:b/:c/:d/:e", func(ctx *Context) any {
		return ctx.Params()
	})
	r.PATCH("/public/*asset", func(ctx *Context) any {
		return ctx.Param("asset")
	})
	r.POST("/error", func(ctx *Context) any {
		return errors.New("some error")
	})
	return r
}

type client struct {
	http.Handler
}

func (c client) Fetch(method string, target string, body io.Reader, opts ...func(req *http.Request)) *http.Response {
	req := httptest.NewRequest(method, target, body)
	for _, opt := range opts {
		opt(req)
	}
	w := httptest.NewRecorder()
	c.ServeHTTP(w, req)
	return w.Result()
}

var c = client{
	Handler: setupRouter(),
}

func TestRootRoute(t *testing.T) {
	res := c.Fetch(http.MethodGet, "/", nil)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestSimpleRoute(t *testing.T) {
	res := c.Fetch(http.MethodPost, "/a/b/c", nil)
	body, _ := io.ReadAll(res.Body)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "text/plain; charset=utf-8", res.Header.Get("Content-Type"))
	assert.Equal(t, "/a/b/c\n", string(body))
}

func TestMethodNotAllowedRoute(t *testing.T) {
	res := c.Fetch(http.MethodGet, "/a/b/c", nil)
	body, _ := io.ReadAll(res.Body)
	assert.Equal(t, http.StatusMethodNotAllowed, res.StatusCode)
	assert.Equal(t, "text/plain; charset=utf-8", res.Header.Get("Content-Type"))
	assert.Equal(t, "405 method not allowed\n", string(body))
}

func TestNotFoundRoute(t *testing.T) {
	res := c.Fetch(http.MethodGet, "/not-found", nil)
	body, _ := io.ReadAll(res.Body)
	assert.Equal(t, http.StatusNotFound, res.StatusCode)
	assert.Equal(t, "text/plain; charset=utf-8", res.Header.Get("Content-Type"))
	assert.Equal(t, "404 page not found\n", string(body))
}

func TestGroupedRoute(t *testing.T) {
	res := c.Fetch(http.MethodGet, "/api/users", nil)
	body, _ := io.ReadAll(res.Body)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", res.Header.Get("Content-Type"))
	assert.Equal(t, "[\"user 1\",\"user 2\"]\n", string(body))
}

func TestDynamicRoute(t *testing.T) {
	t.Run("SingleParams", func(t *testing.T) {
		res := c.Fetch(http.MethodGet, "/api/users/1", nil)
		body, _ := io.ReadAll(res.Body)
		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, "text/plain; charset=utf-8", res.Header.Get("Content-Type"))
		assert.Equal(t, "user 1\n", string(body))
	})

	t.Run("MultipleParams", func(t *testing.T) {
		res := c.Fetch(http.MethodPut, "/1/2/3/4/5", nil)
		body, _ := io.ReadAll(res.Body)
		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, "application/json; charset=utf-8", res.Header.Get("Content-Type"))
		assert.Equal(t, "{\"a\":\"1\",\"b\":\"2\",\"c\":\"3\",\"d\":\"4\",\"e\":\"5\"}\n", string(body))
	})
}

func TestWildcardRoute(t *testing.T) {
	t.Run("WithoutParams", func(t *testing.T) {
		t.Run("Shallow", func(t *testing.T) {
			res := c.Fetch(http.MethodDelete, "/api/test", nil)
			body, _ := io.ReadAll(res.Body)
			assert.Equal(t, http.StatusOK, res.StatusCode)
			assert.Equal(t, "text/plain; charset=utf-8", res.Header.Get("Content-Type"))
			assert.Equal(t, "api not implemented\n", string(body))
		})

		t.Run("Deep", func(t *testing.T) {
			res := c.Fetch(http.MethodDelete, "/api/blogs/1", nil)
			body, _ := io.ReadAll(res.Body)
			assert.Equal(t, http.StatusOK, res.StatusCode)
			assert.Equal(t, "text/plain; charset=utf-8", res.Header.Get("Content-Type"))
			assert.Equal(t, "api not implemented\n", string(body))
		})
	})

	t.Run("WithParams", func(t *testing.T) {
		res := c.Fetch(http.MethodPatch, "/public/images/some.svg", nil)
		body, _ := io.ReadAll(res.Body)
		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, "text/plain; charset=utf-8", res.Header.Get("Content-Type"))
		assert.Equal(t, "images/some.svg\n", string(body))
	})
}

func TestDuplicateRoute(t *testing.T) {
	r := NewRouter()
	r.GET("/:userId", func(ctx *Context) any { return nil })
	assert.Panics(t, func() {
		r.GET("/:blogId", func(ctx *Context) any { return nil })
	})
}

func TestMiddleware(t *testing.T) {
	t.Run("Simple", func(t *testing.T) {
		res := c.Fetch(http.MethodGet, "/admin/users", nil)
		body, _ := io.ReadAll(res.Body)
		assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
		assert.Equal(t, "text/plain; charset=utf-8", res.Header.Get("Content-Type"))
		assert.Equal(t, "401 unauthorized\n", string(body))
	})

	t.Run("Multiple", func(t *testing.T) {
		res := c.Fetch(http.MethodGet, "/count", nil)
		body, _ := io.ReadAll(res.Body)
		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, "application/json; charset=utf-8", res.Header.Get("Content-Type"))
		assert.Equal(t, "1\n", string(body))
	})

	t.Run("Nested", func(t *testing.T) {
		res := c.Fetch(http.MethodGet, "/admin/users", nil, func(req *http.Request) {
			req.Header.Set("Authorization", "secret")
		})
		body, _ := io.ReadAll(res.Body)
		assert.Equal(t, http.StatusForbidden, res.StatusCode)
		assert.Equal(t, "text/plain; charset=utf-8", res.Header.Get("Content-Type"))
		assert.Equal(t, "403 forbidden\n", string(body))
	})
}

func TestDefaultError(t *testing.T) {
	res := c.Fetch(http.MethodPost, "/error", nil)
	body, _ := io.ReadAll(res.Body)
	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	assert.Equal(t, "text/plain; charset=utf-8", res.Header.Get("Content-Type"))
	assert.Equal(t, "some error\n", string(body))
}
