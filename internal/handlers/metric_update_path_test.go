package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	valErrors "github.com/sbilibin2017/yandex-go-advanced/internal/errors"
)

func TestMetricUpdatePathHandler(t *testing.T) {
	dummyValidator := func(expectedErr error) func(string, string, string) error {
		return func(mt, mn, mv string) error {
			return expectedErr
		}
	}

	tests := []struct {
		name           string
		urlParams      map[string]string
		validatorError error
		updateError    error
		expectedCode   int
	}{
		{
			name:           "valid request",
			urlParams:      map[string]string{"type": "gauge", "name": "temp", "value": "123.4"},
			validatorError: nil,
			updateError:    nil,
			expectedCode:   http.StatusOK,
		},
		{
			name:           "missing metric name",
			urlParams:      map[string]string{"type": "gauge", "name": "", "value": "123.4"},
			validatorError: valErrors.ErrMetricNameMissing,
			expectedCode:   http.StatusNotFound,
		},
		{
			name:           "invalid metric type",
			urlParams:      map[string]string{"type": "invalid", "name": "temp", "value": "123.4"},
			validatorError: valErrors.ErrMetricTypeInvalid,
			expectedCode:   http.StatusBadRequest,
		},
		{
			name:           "invalid metric value",
			urlParams:      map[string]string{"type": "gauge", "name": "temp", "value": "abc"},
			validatorError: valErrors.ErrMetricValueInvalid,
			expectedCode:   http.StatusBadRequest,
		},
		{
			name:           "update service error",
			urlParams:      map[string]string{"type": "gauge", "name": "temp", "value": "123.4"},
			validatorError: nil,
			updateError:    errors.New("update failure"),
			expectedCode:   http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUpdater := NewMockMetricUpdater(ctrl)

			// Set expectation on Update only if validator passes
			if tt.validatorError == nil {
				// Expect Update to be called once with any context and any metrics slice
				mockUpdater.EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Return(tt.updateError).
					Times(1)
			}

			handler := NewMetricUpdatePathHandler(
				dummyValidator(tt.validatorError),
				mockUpdater,
			)

			req := httptest.NewRequest(http.MethodPost, "/update/"+tt.urlParams["type"]+"/"+tt.urlParams["name"]+"/"+tt.urlParams["value"], nil)
			// Set URL params for chi
			rctx := chi.NewRouteContext()
			for k, v := range tt.urlParams {
				rctx.URLParams.Add(k, v)
			}
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedCode, rr.Code)
		})
	}
}
