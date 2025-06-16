package routers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMetricRouter(t *testing.T) {
	// Dummy handlers that respond with their name
	updatePathHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("updatePath"))
	})
	updateBodyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("updateBody"))
	})
	getPathHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("getPath"))
	})
	getBodyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("getBody"))
	})
	listHTMLHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("listHTML"))
	})

	// Middleware that adds a test header
	testMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Test-Middleware", "true")
			next.ServeHTTP(w, r)
		})
	}

	router := NewMetricRouter(
		updatePathHandler,
		updateBodyHandler,
		getPathHandler,
		getBodyHandler,
		listHTMLHandler,
		testMiddleware,
	)

	tests := []struct {
		method       string
		route        string
		expectedBody string
	}{
		{"POST", "/update/gauge/temp/42", "updatePath"},
		{"POST", "/update/", "updateBody"},
		{"GET", "/value/counter/hits", "getPath"},
		{"POST", "/value/", "getBody"},
		{"GET", "/", "listHTML"},
	}

	for _, tt := range tests {
		req := httptest.NewRequest(tt.method, tt.route, nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "true", rec.Header().Get("X-Test-Middleware"))
		assert.Equal(t, tt.expectedBody, rec.Body.String())
	}
}
