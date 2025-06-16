package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/sbilibin2017/yandex-go-advanced/internal/types"
)

func TestMetricListHTMLHandler(t *testing.T) {
	tests := []struct {
		name         string
		serviceReply []types.Metrics
		serviceError error
		expectedCode int
		expectedBody string
		expectedCT   string
	}{
		{
			name: "successful list returns HTML",
			serviceReply: []types.Metrics{
				*types.NewMetric("gauge", "temp", "123.45"),
				*types.NewMetric("counter", "requests", "100"),
			},
			serviceError: nil,
			expectedCode: http.StatusOK,
			expectedBody: types.NewMetricsHTML([]types.Metrics{
				*types.NewMetric("gauge", "temp", "123.45"),
				*types.NewMetric("counter", "requests", "100"),
			}),
			expectedCT: "text/html; charset=utf-8",
		},
		{
			name:         "service returns error",
			serviceReply: nil,
			serviceError: errors.New("db failure"),
			expectedCode: http.StatusInternalServerError,
			expectedBody: "internal server error\n", // note the newline
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockSvc := NewMockMetricHTMLLister(ctrl)
			mockSvc.EXPECT().
				List(gomock.Any()).
				Return(tt.serviceReply, tt.serviceError).
				Times(1)

			handler := NewMetricListHTMLHandler(mockSvc)

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedCode, rr.Code)
			assert.Equal(t, tt.expectedBody, rr.Body.String())
			if tt.expectedCode == http.StatusOK {
				assert.Equal(t, tt.expectedCT, rr.Header().Get("Content-Type"))
			}
		})
	}
}
