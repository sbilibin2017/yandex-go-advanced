package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	internalErrors "github.com/sbilibin2017/yandex-go-advanced/internal/errors"
	"github.com/sbilibin2017/yandex-go-advanced/internal/types"
)

func TestNewMetricUpdateBodyHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	float64Ptr := func(f float64) *float64 {
		return &f
	}

	type testCase struct {
		name             string
		body             interface{}
		valFunc          func(types.Metrics) error
		mockUpdateReturn func(m *MockMetricUpdaterBody, metrics []*types.Metrics)
		wantStatus       int
		wantBodyContains string
	}

	validMetric := types.Metrics{
		ID:    "testMetric",
		Type:  "gauge",
		Value: float64Ptr(123.45),
	}

	tests := []testCase{
		{
			name:             "invalid JSON body",
			body:             "{invalid json",
			valFunc:          func(m types.Metrics) error { return nil },
			mockUpdateReturn: nil,
			wantStatus:       http.StatusBadRequest,
			wantBodyContains: "invalid JSON format",
		},
		{
			name:             "validation error - ID invalid",
			body:             validMetric,
			valFunc:          func(m types.Metrics) error { return internalErrors.ErrMetricIDInvalid },
			wantStatus:       http.StatusNotFound,
			wantBodyContains: internalErrors.ErrMetricIDInvalid.Error(),
		},
		{
			name:             "validation error - type invalid",
			body:             validMetric,
			valFunc:          func(m types.Metrics) error { return internalErrors.ErrMetricTypeInvalid },
			wantStatus:       http.StatusBadRequest,
			wantBodyContains: internalErrors.ErrMetricTypeInvalid.Error(),
		},
		{
			name:    "service update returns error",
			body:    validMetric,
			valFunc: func(m types.Metrics) error { return nil },
			mockUpdateReturn: func(m *MockMetricUpdaterBody, metrics []*types.Metrics) {
				m.EXPECT().
					Update(gomock.Any(), gomock.Eq(metrics)).
					Return(nil, internalErrors.ErrInternalServerError).
					Times(1)
			},
			wantStatus:       http.StatusInternalServerError,
			wantBodyContains: internalErrors.ErrInternalServerError.Error(),
		},

		{
			name:    "success - valid metric",
			body:    validMetric,
			valFunc: func(m types.Metrics) error { return nil },
			mockUpdateReturn: func(m *MockMetricUpdaterBody, metrics []*types.Metrics) {
				m.EXPECT().
					Update(gomock.Any(), gomock.Eq(metrics)).
					Return([]*types.Metrics{metrics[0]}, nil).
					Times(1)
			},
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := NewMockMetricUpdaterBody(ctrl)

			var bodyBytes []byte
			var err error
			switch b := tt.body.(type) {
			case string:
				bodyBytes = []byte(b)
			default:
				bodyBytes, err = json.Marshal(b)
				require.NoError(t, err)
			}

			if tt.mockUpdateReturn != nil {
				metrics := []*types.Metrics{}
				if m, ok := tt.body.(types.Metrics); ok {
					metrics = []*types.Metrics{&m}
				}
				tt.mockUpdateReturn(mockSvc, metrics)
			}

			handler := NewMetricUpdateBodyHandler(tt.valFunc, mockSvc)

			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(bodyBytes))
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			require.Equal(t, tt.wantStatus, rr.Code)
			if tt.wantBodyContains != "" {
				require.Contains(t, rr.Body.String(), tt.wantBodyContains)
			}
		})
	}
}
