package tra

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func BenchmarkRouter(b *testing.B) {
	router := setupRouter()
	tests := []struct {
		name   string
		path   string
		method string
	}{
		{"RootRoute", "/", http.MethodGet},
		{"SimpleRoute", "/a/b/c", http.MethodPost},
		{"DynamicRoute", "/api/users/1", http.MethodGet},
		{"DynamicRoute5Params", "/test/test/test/test/test", http.MethodPut},
		{"WildcardRouteShallow", "/api/test", http.MethodDelete},
		{"WildcardRouteDeep", "/public/assets/images/some.svg", http.MethodPatch},
		{"NotFoundRoute", "/not-found", http.MethodGet},
		{"MethodNotAllowedRoute", "/a/b/c", http.MethodGet},
	}

	for _, test := range tests {
		b.Run(test.name, func(b *testing.B) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(test.method, test.path, nil)

			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				router.ServeHTTP(w, r)
			}
		})
	}
}
