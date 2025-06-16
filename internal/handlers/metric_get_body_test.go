package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	internalErrors "github.com/sbilibin2017/yandex-go-advanced/internal/errors"

	"github.com/sbilibin2017/yandex-go-advanced/internal/types"
)

func TestNewMetricGetBodyHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	float64Ptr := func(f float64) *float64 {
		return &f
	}

	mockSvc := NewMockMetricBodyGetter(ctrl)

	validate := func(expectedErr error) func(metric types.MetricID) error {
		return func(metric types.MetricID) error {
			return expectedErr
		}
	}

	tests := []struct {
		name         string
		requestBody  interface{} // will be marshaled to JSON
		validateFunc func(types.MetricID) error
		setupMock    func(metricID types.MetricID)
		wantCode     int
		wantBody     string
	}{
		{
			name:         "invalid JSON body",
			requestBody:  `{"id": "name", "type":}`, // invalid JSON
			validateFunc: validate(nil),
			wantCode:     http.StatusBadRequest,
			wantBody:     "invalid JSON format\n",
		},
		{
			name:         "validation error - ID invalid",
			requestBody:  types.MetricID{ID: "name", Type: "gauge"},
			validateFunc: validate(internalErrors.ErrMetricIDInvalid),
			wantCode:     http.StatusNotFound,
			wantBody:     internalErrors.ErrMetricIDInvalid.Error() + "\n",
		},
		{
			name:         "validation error - type invalid",
			requestBody:  types.MetricID{ID: "name", Type: "invalid-type"},
			validateFunc: validate(internalErrors.ErrMetricTypeInvalid),
			wantCode:     http.StatusBadRequest,
			wantBody:     internalErrors.ErrMetricTypeInvalid.Error() + "\n",
		},
		{
			name:         "service returns error",
			requestBody:  types.MetricID{ID: "name", Type: "gauge"},
			validateFunc: validate(nil),
			setupMock: func(metricID types.MetricID) {
				mockSvc.EXPECT().
					Get(gomock.Any(), metricID).
					Return(nil, internalErrors.ErrInternalServerError)
			},
			wantCode: http.StatusInternalServerError,
			wantBody: internalErrors.ErrInternalServerError.Error() + "\n",
		},
		{
			name:         "service returns nil metric",
			requestBody:  types.MetricID{ID: "name", Type: "gauge"},
			validateFunc: validate(nil),
			setupMock: func(metricID types.MetricID) {
				mockSvc.EXPECT().
					Get(gomock.Any(), metricID).
					Return(nil, nil)
			},
			wantCode: http.StatusNotFound,
			wantBody: internalErrors.ErrMetricNotFound.Error() + "\n",
		},
		{
			name:         "success returns metric JSON",
			requestBody:  types.MetricID{ID: "name", Type: "gauge"},
			validateFunc: validate(nil),
			setupMock: func(metricID types.MetricID) {
				m := &types.Metrics{
					ID:    "name",
					Type:  "gauge",
					Value: float64Ptr(123.45),
					Delta: nil,
					Hash:  "", // empty string omitted
				}
				mockSvc.EXPECT().
					Get(gomock.Any(), metricID).
					Return(m, nil)
			},
			wantCode: http.StatusOK,
			wantBody: `{"id":"name","type":"gauge","value":123.45}` + "\n", // updated expectation
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var bodyBytes []byte
			var err error

			switch v := tt.requestBody.(type) {
			case string:
				bodyBytes = []byte(v)
			default:
				bodyBytes, err = json.Marshal(tt.requestBody)
				assert.NoError(t, err)
			}

			req := httptest.NewRequest(http.MethodPost, "/metric", bytes.NewReader(bodyBytes))
			w := httptest.NewRecorder()

			if tt.setupMock != nil {
				if metricID, ok := tt.requestBody.(types.MetricID); ok {
					tt.setupMock(metricID)
				}
			}

			handler := NewMetricGetBodyHandler(tt.validateFunc, mockSvc)
			handler.ServeHTTP(w, req)

			resp := w.Result()
			defer resp.Body.Close() // <-- добавь это

			body := w.Body.String()

			assert.Equal(t, tt.wantCode, resp.StatusCode)
			assert.Equal(t, tt.wantBody, body)
		})
	}
}
