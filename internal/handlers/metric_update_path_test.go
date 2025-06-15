package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"

	"github.com/sbilibin2017/yandex-go-advanced/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestMetricUpdatePathHandler(t *testing.T) {
	int64Ptr := func(i int64) *int64 { return &i }
	float64Ptr := func(f float64) *float64 { return &f }

	type want struct {
		code int
	}

	tests := []struct {
		name        string
		metricType  string
		metricName  string
		metricValue string
		mockSetup   func(m *MockMetricUpdater)
		want        want
	}{
		{
			name:        "valid counter update",
			metricType:  "counter",
			metricName:  "requests",
			metricValue: "42",
			mockSetup: func(m *MockMetricUpdater) {
				m.EXPECT().
					Update(gomock.Any(), gomock.Eq([]types.Metrics{
						{
							ID:    "requests",
							MType: types.Counter,
							Delta: int64Ptr(42),
						},
					})).
					Return(nil)
			},
			want: want{code: http.StatusOK},
		},
		{
			name:        "valid gauge update",
			metricType:  "gauge",
			metricName:  "load",
			metricValue: "3.14",
			mockSetup: func(m *MockMetricUpdater) {
				m.EXPECT().
					Update(gomock.Any(), gomock.Eq([]types.Metrics{
						{
							ID:    "load",
							MType: types.Gauge,
							Value: float64Ptr(3.14),
						},
					})).
					Return(nil)
			},
			want: want{code: http.StatusOK},
		},
		{
			name:        "missing name",
			metricType:  "gauge",
			metricName:  "",
			metricValue: "1.0",
			mockSetup:   func(m *MockMetricUpdater) {},
			want:        want{code: http.StatusNotFound},
		},
		{
			name:        "invalid type",
			metricType:  "invalid",
			metricName:  "foo",
			metricValue: "100",
			mockSetup:   func(m *MockMetricUpdater) {},
			want:        want{code: http.StatusBadRequest},
		},
		{
			name:        "invalid counter value",
			metricType:  "counter",
			metricName:  "bar",
			metricValue: "not-an-int",
			mockSetup:   func(m *MockMetricUpdater) {},
			want:        want{code: http.StatusBadRequest},
		},
		{
			name:        "invalid gauge value",
			metricType:  "gauge",
			metricName:  "cpu",
			metricValue: "not-a-float",
			mockSetup:   func(m *MockMetricUpdater) {},
			want:        want{code: http.StatusBadRequest},
		},
		{
			name:        "internal update error",
			metricType:  "gauge",
			metricName:  "mem",
			metricValue: "9.99",
			mockSetup: func(m *MockMetricUpdater) {
				m.EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Return(errors.New("db error"))
			},
			want: want{code: http.StatusInternalServerError},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUpdater := NewMockMetricUpdater(ctrl)
			tt.mockSetup(mockUpdater)

			r := chi.NewRouter()
			r.MethodFunc("POST", "/update/{type}/{name}/{value}", NewMetricUpdatePathHandler(mockUpdater))

			req := httptest.NewRequest("POST", "/update/"+tt.metricType+"/"+tt.metricName+"/"+tt.metricValue, nil)
			rec := httptest.NewRecorder()

			r.ServeHTTP(rec, req)

			assert.Equal(t, tt.want.code, rec.Code)
		})
	}
}

func TestHandleMetricUpdatePathError_InternalError(t *testing.T) {
	// given
	rr := httptest.NewRecorder()
	unexpectedErr := errors.New("something unexpected")

	// when
	handleMetricUpdatePathError(rr, unexpectedErr)

	// then
	res := rr.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)

	buf := make([]byte, 512)
	n, _ := res.Body.Read(buf)
	body := string(buf[:n])

	assert.True(t, strings.Contains(body, "internal error"))
}
