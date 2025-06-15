package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	gomock "github.com/golang/mock/gomock"
	"github.com/sbilibin2017/yandex-go-advanced/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestNewMetricGetPathHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := NewMockMetricPathGetter(ctrl)
	handler := NewMetricGetPathHandler(mockSvc)

	tests := []struct {
		name         string
		metricType   string
		metricName   string
		mockSetup    func()
		expectedCode int
		expectedBody string
	}{
		{
			name:         "Valid Gauge",
			metricType:   "gauge",
			metricName:   "temperature",
			expectedCode: http.StatusOK,
			expectedBody: "23.5",
			mockSetup: func() {
				val := 23.5
				mockSvc.EXPECT().
					Get(gomock.Any(), types.MetricID{ID: "temperature", MType: types.Gauge}).
					Return(&types.Metrics{
						ID:    "temperature",
						MType: types.Gauge,
						Value: &val,
					}, nil)
			},
		},
		{
			name:         "Valid Counter",
			metricType:   "counter",
			metricName:   "requests",
			expectedCode: http.StatusOK,
			expectedBody: "100",
			mockSetup: func() {
				val := int64(100)
				mockSvc.EXPECT().
					Get(gomock.Any(), types.MetricID{ID: "requests", MType: types.Counter}).
					Return(&types.Metrics{
						ID:    "requests",
						MType: types.Counter,
						Delta: &val,
					}, nil)
			},
		},
		{
			name:         "Metric Not Found",
			metricType:   "counter",
			metricName:   "missing",
			expectedCode: http.StatusNotFound,
			expectedBody: "metric not found\n",
			mockSetup: func() {
				mockSvc.EXPECT().
					Get(gomock.Any(), types.MetricID{ID: "missing", MType: types.Counter}).
					Return(nil, nil)
			},
		},
		{
			name:         "Service Get Error",
			metricType:   "counter",
			metricName:   "fail",
			expectedCode: http.StatusInternalServerError,
			expectedBody: "internal error\n",
			mockSetup: func() {
				mockSvc.EXPECT().
					Get(gomock.Any(), types.MetricID{ID: "fail", MType: types.Counter}).
					Return(nil, errors.New("db error"))
			},
		},
		{
			name:         "Unexpected Metric Type",
			metricType:   "gauge",
			metricName:   "weird",
			expectedCode: http.StatusBadRequest,
			expectedBody: "invalid metric type\n",
			mockSetup: func() {
				mockSvc.EXPECT().
					Get(gomock.Any(), types.MetricID{ID: "weird", MType: types.Gauge}).
					Return(&types.Metrics{
						ID:    "weird",
						MType: "unknown", // Unexpected type
					}, nil)
			},
		},
		{
			name:         "Invalid Type Param",
			metricType:   "banana",
			metricName:   "fruit",
			expectedCode: http.StatusBadRequest,
			expectedBody: "invalid metric type\n",
			mockSetup:    func() {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("type", tt.metricType)
			rctx.URLParams.Add("name", tt.metricName)

			req := httptest.NewRequest("GET", "/", nil)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedCode, rr.Code)
			assert.Equal(t, tt.expectedBody, rr.Body.String())
		})
	}
}

func Test_validateMetricGetPath(t *testing.T) {
	tests := []struct {
		name       string
		metricType string
		metricName string
		wantErr    error
	}{
		{"Valid Counter", "counter", "requests", nil},
		{"Valid Gauge", "gauge", "cpu", nil},
		{"Empty Name", "counter", "", errMetricGetNameMissing},
		{"Invalid Type", "invalid", "something", errMetricGetTypeInvalid},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := validateMetricGetPath(tt.metricType, tt.metricName)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func Test_formatGauge(t *testing.T) {
	val := 42.5
	assert.Equal(t, "42.5", metricGetPathFormatGauge(&val))
	assert.Equal(t, "0.0", metricGetPathFormatGauge(nil))
}

func Test_formatCounter(t *testing.T) {
	val := int64(99)
	assert.Equal(t, "99", metricGetPathFormatCounter(&val))
	assert.Equal(t, "0", metricGetPathFormatCounter(nil))
}

func Test_handleMetricGetPathError(t *testing.T) {
	tests := []struct {
		name         string
		err          error
		expectedCode int
		expectedBody string
	}{
		{"Metric Not Found", errMetricGetNotFound, http.StatusNotFound, "metric not found\n"},
		{"Invalid Type", errMetricGetTypeInvalid, http.StatusBadRequest, "invalid metric type\n"},
		{"Missing Name", errMetricGetNameMissing, http.StatusNotFound, "metric name is required\n"},
		{"Unknown Error", errors.New("oops"), http.StatusInternalServerError, "internal error\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			handleMetricGetPathError(rr, tt.err)
			assert.Equal(t, tt.expectedCode, rr.Code)
			assert.Equal(t, tt.expectedBody, rr.Body.String())
		})
	}
}
