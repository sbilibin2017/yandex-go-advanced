package facades_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/sbilibin2017/yandex-go-advanced/internal/facades"
	"github.com/sbilibin2017/yandex-go-advanced/internal/types"
	"github.com/stretchr/testify/assert"
)

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
				var m types.Metrics
				if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
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
				var m types.Metrics
				if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
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
					serverAddr = strings.TrimPrefix(ts.URL, "http://")
				} else {
					serverAddr = ts.URL
				}
			}

			facade := facades.NewMetricUpdateFacade(serverAddr)

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
