package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sbilibin2017/yandex-go-advanced/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestNewMetricListHTMLHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := NewMockMetricHTMLLister(ctrl)

	handler := NewMetricListHTMLHandler(mockSvc)

	tests := []struct {
		name          string
		setupMock     func()
		wantCode      int
		wantBodyParts []string // parts expected in body string
	}{
		{
			name: "success returns HTML",
			setupMock: func() {
				mockSvc.EXPECT().
					List(gomock.Any()).
					Return([]types.Metrics{
						{ID: "m1", Type: "gauge"},
						{ID: "m2", Type: "counter"},
					}, nil)
			},
			wantCode:      http.StatusOK,
			wantBodyParts: []string{"m1", "m2", "n/a"}, // check IDs and the "n/a" values your HTML produces
		},
		{
			name: "service returns error",
			setupMock: func() {
				mockSvc.EXPECT().
					List(gomock.Any()).
					Return(nil, errors.New("fail"))
			},
			wantCode:      http.StatusInternalServerError,
			wantBodyParts: []string{"internal server error"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			w := httptest.NewRecorder()

			handler(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			body := w.Body.String()

			assert.Equal(t, tt.wantCode, resp.StatusCode)
			for _, part := range tt.wantBodyParts {
				assert.Contains(t, strings.ToLower(body), strings.ToLower(part))
			}
		})
	}
}
