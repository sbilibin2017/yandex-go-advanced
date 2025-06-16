package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	internalErrors "github.com/sbilibin2017/yandex-go-advanced/internal/errors"

	"github.com/sbilibin2017/yandex-go-advanced/internal/types"
	"github.com/stretchr/testify/require"
)

func TestNewMetricUpdatePathHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		metricType  string
		metricName  string
		metricValue string
	}

	tests := []struct {
		name             string
		args             args
		valFunc          func(string, string, string) error
		mockUpdateReturn func(m *MockMetricUpdaterPath)
		wantStatus       int
		wantBodyContains string
	}{
		{
			name: "validation error - name missing",
			args: args{"gauge", "", "123"},
			valFunc: func(mt, mn, mv string) error {
				return internalErrors.ErrMetricNameMissing
			},
			mockUpdateReturn: nil,
			wantStatus:       http.StatusNotFound,
			wantBodyContains: internalErrors.ErrMetricNameMissing.Error(),
		},
		{
			name: "validation error - invalid type",
			args: args{"invalid", "cpu", "123"},
			valFunc: func(mt, mn, mv string) error {
				return internalErrors.ErrMetricTypeInvalid
			},
			mockUpdateReturn: nil,
			wantStatus:       http.StatusBadRequest,
			wantBodyContains: internalErrors.ErrMetricTypeInvalid.Error(),
		},
		{
			name: "validation error - invalid value",
			args: args{"gauge", "cpu", "abc"},
			valFunc: func(mt, mn, mv string) error {
				return internalErrors.ErrMetricValueInvalid
			},
			mockUpdateReturn: nil,
			wantStatus:       http.StatusBadRequest,
			wantBodyContains: internalErrors.ErrMetricValueInvalid.Error(),
		},
		{
			name: "service update returns error",
			args: args{"gauge", "cpu", "123"},
			valFunc: func(mt, mn, mv string) error {
				return nil
			},
			mockUpdateReturn: func(m *MockMetricUpdaterPath) {
				m.EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Return(nil, internalErrors.ErrInternalServerError).
					Times(1)
			},
			wantStatus:       http.StatusInternalServerError,
			wantBodyContains: internalErrors.ErrInternalServerError.Error(),
		},
		{
			name: "successful update",
			args: args{"gauge", "cpu", "123"},
			valFunc: func(mt, mn, mv string) error {
				return nil
			},
			mockUpdateReturn: func(m *MockMetricUpdaterPath) {
				m.EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Return([]*types.Metrics{}, nil).
					Times(1)
			},
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := NewMockMetricUpdaterPath(ctrl)
			if tt.mockUpdateReturn != nil {
				tt.mockUpdateReturn(mockSvc)
			}

			handler := NewMetricUpdatePathHandler(tt.valFunc, mockSvc)

			// Prepare request with URL params set via chi.RouteContext
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("type", tt.args.metricType)
			rctx.URLParams.Add("name", tt.args.metricName)
			rctx.URLParams.Add("value", tt.args.metricValue)

			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(""))
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			require.Equal(t, tt.wantStatus, rr.Code)
			if tt.wantBodyContains != "" {
				require.Contains(t, rr.Body.String(), tt.wantBodyContains)
			}
		})
	}
}
