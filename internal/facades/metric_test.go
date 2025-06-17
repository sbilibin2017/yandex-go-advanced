package facades

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"

	"github.com/sbilibin2017/yandex-go-advanced/internal/types"
	"github.com/stretchr/testify/assert"
)

// Helper function to decompress gzip request body if needed
func decompressRequestBody(r *http.Request) (io.ReadCloser, error) {
	if r.Header.Get("Content-Encoding") == "gzip" {
		gr, err := gzip.NewReader(r.Body)
		if err != nil {
			return nil, err
		}
		return gr, nil
	}
	return r.Body, nil
}

func TestMetricUpdateFacade_Update(t *testing.T) {
	tests := []struct {
		name           string
		serverAddress  string
		handlerFunc    http.HandlerFunc
		metrics        []*types.Metrics
		wantErr        bool
		expectedErrMsg string
	}{
		{
			name: "success",
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				bodyReader, err := decompressRequestBody(r)
				if err != nil {
					http.Error(w, "bad gzip encoding", http.StatusBadRequest)
					return
				}
				defer bodyReader.Close()

				var m types.Metrics
				if err := json.NewDecoder(bodyReader).Decode(&m); err != nil {
					http.Error(w, "bad request", http.StatusBadRequest)
					return
				}
				w.WriteHeader(http.StatusOK)
			},
			metrics: []*types.Metrics{
				{
					ID:   "metric1",
					Type: types.Gauge,
					Value: func() *float64 {
						v := 10.5
						return &v
					}(),
				},
			},
			wantErr: false,
		},
		{
			name: "bad request response",
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "bad request", http.StatusBadRequest)
			},
			metrics: []*types.Metrics{
				{
					ID:   "metric2",
					Type: types.Counter,
					Delta: func() *int64 {
						v := int64(5)
						return &v
					}(),
				},
			},
			wantErr:        true,
			expectedErrMsg: "metrics update request failed: 400 Bad Request",
		},
		{
			name: "invalid json request",
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "invalid json", http.StatusBadRequest)
			},
			metrics: []*types.Metrics{
				{
					ID:   "metric3",
					Type: types.Gauge,
					Value: func() *float64 {
						v := 2.5
						return &v
					}(),
				},
			},
			wantErr:        true,
			expectedErrMsg: "metrics update request failed: 400 Bad Request",
		},
		{
			name: "network error",
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			},
			metrics: []*types.Metrics{
				{
					ID:   "metric4",
					Type: types.Counter,
					Delta: func() *int64 {
						v := int64(1)
						return &v
					}(),
				},
			},
			wantErr:        true,
			expectedErrMsg: "failed to send metrics update request:",
		},
		{
			name: "missing scheme in serverAddress adds http",
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				bodyReader, err := decompressRequestBody(r)
				if err != nil {
					http.Error(w, "bad gzip encoding", http.StatusBadRequest)
					return
				}
				defer bodyReader.Close()

				var m types.Metrics
				if err := json.NewDecoder(bodyReader).Decode(&m); err != nil {
					http.Error(w, "bad request", http.StatusBadRequest)
					return
				}
				w.WriteHeader(http.StatusOK)
			},
			metrics: []*types.Metrics{
				{
					ID:   "metric5",
					Type: types.Counter,
					Delta: func() *int64 {
						v := int64(100)
						return &v
					}(),
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ts *httptest.Server

			if tt.name == "network error" {
				// Create and immediately close server to simulate network error
				ts = httptest.NewServer(tt.handlerFunc)
				ts.Close()
			} else {
				r := chi.NewRouter()
				r.Post("/update/", tt.handlerFunc)
				ts = httptest.NewServer(r)
			}

			serverAddr := tt.serverAddress
			if serverAddr == "" {
				if tt.name == "missing scheme in serverAddress adds http" {
					// Remove http:// prefix to test scheme addition in facade
					serverAddr = strings.TrimPrefix(ts.URL, "http://")
				} else {
					serverAddr = ts.URL
				}
			}

			facade := NewMetricUpdateFacade(serverAddr)

			err := facade.Update(context.Background(), tt.metrics)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErrMsg)
			} else {
				assert.NoError(t, err)
			}

			ts.Close()
		})
	}
}

func TestCompressMetrics(t *testing.T) {
	metric := &types.Metrics{
		ID:    "testMetric",
		Type:  types.Gauge,
		Value: func() *float64 { v := 123.456; return &v }(),
	}

	compressedData, err := compressMetrics(metric)
	assert.NoError(t, err)
	assert.NotEmpty(t, compressedData)

	// Decompress and verify that it matches the original metric JSON
	reader, err := gzip.NewReader(bytes.NewReader(compressedData))
	assert.NoError(t, err)

	decompressedData, err := io.ReadAll(reader)
	assert.NoError(t, err)

	err = reader.Close()
	assert.NoError(t, err)

	// Unmarshal decompressed JSON
	var result types.Metrics
	err = json.Unmarshal(decompressedData, &result)
	assert.NoError(t, err)

	// Compare result with original metric
	assert.Equal(t, metric.ID, result.ID)
	assert.Equal(t, metric.Type, result.Type)
	assert.NotNil(t, result.Value)
	assert.Equal(t, *metric.Value, *result.Value)
}

func TestCompressMetrics_EmptyMetric(t *testing.T) {
	emptyMetric := &types.Metrics{}

	compressedData, err := compressMetrics(emptyMetric)
	assert.NoError(t, err)
	assert.NotEmpty(t, compressedData)

	// Decompress to verify valid gzip data
	reader, err := gzip.NewReader(bytes.NewReader(compressedData))
	assert.NoError(t, err)
	defer reader.Close()

	decompressedData, err := io.ReadAll(reader)
	assert.NoError(t, err)

	var result types.Metrics
	err = json.Unmarshal(decompressedData, &result)
	assert.NoError(t, err)
	assert.Equal(t, emptyMetric, &result)
}
