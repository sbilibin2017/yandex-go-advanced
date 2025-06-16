package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	gomock "github.com/golang/mock/gomock"
	internalErrors "github.com/sbilibin2017/yandex-go-advanced/internal/errors"
	"github.com/sbilibin2017/yandex-go-advanced/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestMetricGetPathHandler(t *testing.T) {
	dummyValidator := func(expectedErr error) func(string, string) error {
		return func(mt, mn string) error {
			return expectedErr
		}
	}

	tests := []struct {
		name           string
		urlParams      map[string]string
		validatorError error
		serviceReturn  *types.Metrics
		serviceError   error
		expectedCode   int
		expectedBody   string
	}{
		{
			name:           "valid request returns metric",
			urlParams:      map[string]string{"type": "gauge", "name": "temperature"},
			validatorError: nil,
			serviceReturn:  types.NewMetric("gauge", "temperature", "123.45"),
			serviceError:   nil,
			expectedCode:   http.StatusOK,
			expectedBody:   "123.45",
		},
		{
			name:           "validator error - metric name missing",
			urlParams:      map[string]string{"type": "gauge", "name": ""},
			validatorError: internalErrors.ErrMetricNameMissing,
			serviceReturn:  nil,
			serviceError:   nil,
			expectedCode:   http.StatusNotFound,
			expectedBody:   internalErrors.ErrMetricNameMissing.Error() + "\n",
		},
		{
			name:           "validator error - metric type invalid",
			urlParams:      map[string]string{"type": "invalid", "name": "temperature"},
			validatorError: internalErrors.ErrMetricTypeInvalid,
			serviceReturn:  nil,
			serviceError:   nil,
			expectedCode:   http.StatusBadRequest,
			expectedBody:   internalErrors.ErrMetricTypeInvalid.Error() + "\n",
		},
		{
			name:           "service returns error",
			urlParams:      map[string]string{"type": "gauge", "name": "temperature"},
			validatorError: nil,
			serviceReturn:  nil,
			serviceError:   errors.New("db failure"),
			expectedCode:   http.StatusInternalServerError,
			expectedBody:   internalErrors.ErrInternalServerError.Error() + "\n",
		},
		{
			name:           "service returns nil metric (not found)",
			urlParams:      map[string]string{"type": "gauge", "name": "temperature"},
			validatorError: nil,
			serviceReturn:  nil,
			serviceError:   nil,
			expectedCode:   http.StatusNotFound, // <-- fixed to NotFound here
			expectedBody:   internalErrors.ErrMetricNotFound.Error() + "\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockSvc := NewMockMetricPathGetter(ctrl)

			if tt.validatorError == nil {
				id := *types.NewMetricID(tt.urlParams["type"], tt.urlParams["name"])
				mockSvc.EXPECT().
					Get(gomock.Any(), id).
					Return(tt.serviceReturn, tt.serviceError).
					Times(1)
			}

			handler := NewMetricGetPathHandler(dummyValidator(tt.validatorError), mockSvc)

			req := httptest.NewRequest(http.MethodGet, "/value/"+tt.urlParams["type"]+"/"+tt.urlParams["name"], nil)

			rctx := chi.NewRouteContext()
			for k, v := range tt.urlParams {
				rctx.URLParams.Add(k, v)
			}
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedCode, rr.Code)
			assert.Equal(t, tt.expectedBody, rr.Body.String())
		})
	}
}
