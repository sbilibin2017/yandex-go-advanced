package facades

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/sbilibin2017/yandex-go-advanced/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestMetricUpdateFacade_AddsHTTPPrefixIfMissing(t *testing.T) {
	// Create a test HTTP server to check requests
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// The request URL's Host should contain "localhost" (testServer URL)
		assert.True(t, strings.HasPrefix(r.URL.String(), "/update/"))
		w.WriteHeader(http.StatusOK)
	}))
	defer testServer.Close()

	// Extract host without protocol, e.g. "localhost:12345"
	addrWithoutProtocol := testServer.Listener.Addr().String()

	// Create facade with serverAddress missing protocol prefix
	facade := NewMetricUpdateFacade(addrWithoutProtocol)

	value := 123.45
	metrics := []types.Metrics{
		{ID: "cpu", Type: types.Gauge, Value: &value},
	}

	err := facade.Update(context.Background(), metrics)
	assert.NoError(t, err)
}

func TestMetricUpdateFacade_Update_Success(t *testing.T) {
	// Create a test HTTP server to mock real server responses
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check method and path format
		assert.Equal(t, http.MethodPost, r.Method)

		// Example: /update/gauge/temperature/42.5
		assert.Contains(t, r.URL.Path, "/update/")

		w.WriteHeader(http.StatusOK)
	}))
	defer testServer.Close()

	facade := NewMetricUpdateFacade(testServer.URL)

	gaugeValue := 42.5
	counterValue := int64(7)

	metrics := []types.Metrics{
		{ID: "temperature", Type: types.Gauge, Value: &gaugeValue},
		{ID: "requests", Type: types.Counter, Delta: &counterValue},
	}

	err := facade.Update(context.Background(), metrics)
	assert.NoError(t, err)
}

func TestMetricUpdateFacade_Update_FailOnInvalidMetric(t *testing.T) {
	facade := NewMetricUpdateFacade("http://example.com")

	metrics := []types.Metrics{
		{ID: "invalid", Type: "unknown"},
	}

	err := facade.Update(context.Background(), metrics)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown metric type")
}

func TestMetricUpdateFacade_Update_FailOnNilValue(t *testing.T) {
	facade := NewMetricUpdateFacade("http://example.com")

	metrics := []types.Metrics{
		{ID: "missingvalue", Type: types.Gauge, Value: nil},
	}

	err := facade.Update(context.Background(), metrics)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "value is nil")

	metrics = []types.Metrics{
		{ID: "missingdelta", Type: types.Counter, Delta: nil},
	}

	err = facade.Update(context.Background(), metrics)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "delta value is nil")
}

func TestMetricUpdateFacade_Update_ServerReturnsError(t *testing.T) {
	// Setup a test server that returns 500 error
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}))
	defer testServer.Close()

	facade := NewMetricUpdateFacade(testServer.URL)

	value := 10.0
	metrics := []types.Metrics{
		{ID: "temp", Type: types.Gauge, Value: &value},
	}

	err := facade.Update(context.Background(), metrics)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to send metric")
}

func TestMetricUpdateFacade_Update_ServerNotReachable(t *testing.T) {
	facade := NewMetricUpdateFacade("http://localhost:9999") // assuming port not in use

	value := 10.0
	metrics := []types.Metrics{
		{ID: "temp", Type: types.Gauge, Value: &value},
	}

	err := facade.Update(context.Background(), metrics)
	assert.Error(t, err)
}
