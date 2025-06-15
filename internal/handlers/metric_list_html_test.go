package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sbilibin2017/yandex-go-advanced/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestNewMetricListHTMLHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	float64Ptr := func(v float64) *float64 { return &v }
	int64Ptr := func(v int64) *int64 { return &v }

	type testCase struct {
		name           string
		mockReturn     []types.Metrics
		mockError      error
		wantStatusCode int
		wantBodySubstr []string
	}

	tests := []testCase{
		{
			name: "success with gauge and counter",
			mockReturn: []types.Metrics{
				{ID: "gauge_metric", MType: types.Gauge, Value: float64Ptr(42.42)},
				{ID: "counter_metric", MType: types.Counter, Delta: int64Ptr(123)},
			},
			mockError:      nil,
			wantStatusCode: http.StatusOK,
			wantBodySubstr: []string{"gauge_metric", "42.42", "counter_metric", "123", "<html>", "</html>"},
		},
		{
			name:           "error from service",
			mockReturn:     nil,
			mockError:      assert.AnError,
			wantStatusCode: http.StatusInternalServerError,
			wantBodySubstr: []string{"internal error"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockLister := NewMockMetricHTMLLister(ctrl)
			mockLister.EXPECT().List(gomock.Any()).Return(tc.mockReturn, tc.mockError).Times(1)

			handler := NewMetricListHTMLHandler(mockLister)

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			assert.Equal(t, tc.wantStatusCode, res.StatusCode)
			for _, substr := range tc.wantBodySubstr {
				assert.Contains(t, rec.Body.String(), substr)
			}
		})
	}
}

func Test_newMetricsHTML(t *testing.T) {
	float64Ptr := func(v float64) *float64 { return &v }
	int64Ptr := func(v int64) *int64 { return &v }

	tests := []struct {
		name     string
		metrics  []types.Metrics
		expected []string // substrings expected in output
	}{
		{
			name: "gauge and counter metrics",
			metrics: []types.Metrics{
				{ID: "test_gauge", MType: types.Gauge, Value: float64Ptr(3.14)},
				{ID: "test_counter", MType: types.Counter, Delta: int64Ptr(42)},
				{ID: "unknown_metric", MType: "unknown"},
			},
			expected: []string{"test_gauge", "3.14", "test_counter", "42", "unknown_metric: "},
		},
		{
			name:     "empty metrics",
			metrics:  []types.Metrics{},
			expected: []string{"<html>", "<ul>", "</ul>", "</html>"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			html := newMetricsHTML(tc.metrics)
			for _, substr := range tc.expected {
				assert.Contains(t, html, substr)
			}
		})
	}
}

func Test_metricListHTMLFormatGauge(t *testing.T) {
	float64Ptr := func(v float64) *float64 { return &v }

	tests := []struct {
		name     string
		input    *float64
		expected string
	}{
		{"nil input", nil, "0.0"},
		{"valid input", float64Ptr(1.2345), "1.2345"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := metricListHTMLFormatGauge(tc.input)
			assert.Equal(t, tc.expected, got)
		})
	}
}

func Test_metricListHTMLFormatCounter(t *testing.T) {
	int64Ptr := func(v int64) *int64 { return &v }

	tests := []struct {
		name     string
		input    *int64
		expected string
	}{
		{"nil input", nil, "0"},
		{"valid input", int64Ptr(987), "987"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := metricListHTMLFormatCounter(tc.input)
			assert.Equal(t, tc.expected, got)
		})
	}
}
