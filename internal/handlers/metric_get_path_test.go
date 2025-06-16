package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	internalErrors "github.com/sbilibin2017/yandex-go-advanced/internal/errors"

	"github.com/sbilibin2017/yandex-go-advanced/internal/types"
)

func TestNewMetricGetPathHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	int64Ptr := func(i int64) *int64 { return &i }
	float64Ptr := func(f float64) *float64 { return &f }

	mockSvc := NewMockMetricPathGetter(ctrl)

	validate := func(expectedErr error) func(string, string) error {
		return func(metricType, metricName string) error {
			return expectedErr
		}
	}

	tests := []struct {
		name         string
		metricType   string
		metricName   string
		validateFunc func(string, string) error
		setupMock    func()
		wantCode     int
		wantBody     string
	}{
		{
			name:         "validation error",
			metricType:   "gauge",
			metricName:   "name",
			validateFunc: validate(internalErrors.ErrMetricTypeInvalid),
			wantCode:     http.StatusBadRequest,
			wantBody:     internalErrors.ErrMetricTypeInvalid.Error() + "\n",
		},
		{
			name:         "service returns error",
			metricType:   "gauge",
			metricName:   "name",
			validateFunc: validate(nil),
			setupMock: func() {
				mockSvc.EXPECT().
					Get(gomock.Any(), types.MetricID{
						ID:   "name",
						Type: "gauge",
					}).
					Return(nil, internalErrors.ErrInternalServerError)
			},
			wantCode: http.StatusInternalServerError,
			wantBody: internalErrors.ErrInternalServerError.Error() + "\n",
		},
		{
			name:         "service returns nil metric",
			metricType:   "gauge",
			metricName:   "name",
			validateFunc: validate(nil),
			setupMock: func() {
				mockSvc.EXPECT().
					Get(gomock.Any(), types.MetricID{
						ID:   "name",
						Type: "gauge",
					}).
					Return(nil, nil)
			},
			wantCode: http.StatusNotFound,
			wantBody: internalErrors.ErrMetricNotFound.Error() + "\n",
		},
		{
			name:         "success returns metric string value (Value set)",
			metricType:   "gauge",
			metricName:   "name",
			validateFunc: validate(nil),
			setupMock: func() {
				m := &types.Metrics{
					ID:    "name",
					Type:  "gauge",
					Value: float64Ptr(42.42),
				}
				mockSvc.EXPECT().
					Get(gomock.Any(), types.MetricID{
						ID:   "name",
						Type: "gauge",
					}).
					Return(m, nil)
			},
			wantCode: http.StatusOK,
			wantBody: "42.42",
		},
		{
			name:         "success returns metric string value (Delta set)",
			metricType:   "counter",
			metricName:   "name",
			validateFunc: validate(nil),
			setupMock: func() {
				m := &types.Metrics{
					ID:    "name",
					Type:  "counter",
					Delta: int64Ptr(100),
				}
				mockSvc.EXPECT().
					Get(gomock.Any(), types.MetricID{
						ID:   "name",
						Type: "counter",
					}).
					Return(m, nil)
			},
			wantCode: http.StatusOK,
			wantBody: "100",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupMock != nil {
				tt.setupMock()
			}

			req := httptest.NewRequest(http.MethodGet,
				"/metric/"+tt.metricType+"/"+tt.metricName,
				nil)

			// Setup chi route context with URL params
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("type", tt.metricType)
			rctx.URLParams.Add("name", tt.metricName)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			w := httptest.NewRecorder()

			handler := NewMetricGetPathHandler(tt.validateFunc, mockSvc)
			handler.ServeHTTP(w, req)

			resp := w.Result()
			defer resp.Body.Close()
			body := w.Body.String()

			assert.Equal(t, tt.wantCode, resp.StatusCode)
			assert.Equal(t, tt.wantBody, body)
		})
	}
}
