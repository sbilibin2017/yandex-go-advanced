package runners_test

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/sbilibin2017/yandex-go-advanced/internal/runners"
)

func TestRunServer_TableDriven_WithSetup(t *testing.T) {
	tests := []struct {
		name           string
		cancelCtxDelay time.Duration
		expectErr      bool
		expectErrMsg   string
		setup          func(ctrl *gomock.Controller) *runners.MockServer
	}{
		{
			name:           "graceful shutdown",
			cancelCtxDelay: 50 * time.Millisecond,
			expectErr:      false,
			setup: func(ctrl *gomock.Controller) *runners.MockServer {
				mock := runners.NewMockServer(ctrl)

				mock.EXPECT().
					ListenAndServe().
					DoAndReturn(func() error {
						time.Sleep(100 * time.Millisecond)
						return http.ErrServerClosed
					})

				mock.EXPECT().
					Shutdown(gomock.Any()).
					Return(nil)

				return mock
			},
		},
		{
			name:           "listenAndServe returns fatal error",
			cancelCtxDelay: 0,
			expectErr:      true,
			expectErrMsg:   "listen error",
			setup: func(ctrl *gomock.Controller) *runners.MockServer {
				mock := runners.NewMockServer(ctrl)

				mock.EXPECT().
					ListenAndServe().
					Return(errors.New("listen error"))

				return mock
			},
		},
		{
			name:           "shutdown returns error",
			cancelCtxDelay: 50 * time.Millisecond,
			expectErr:      true,
			expectErrMsg:   "shutdown failure",
			setup: func(ctrl *gomock.Controller) *runners.MockServer {
				mock := runners.NewMockServer(ctrl)

				mock.EXPECT().
					ListenAndServe().
					DoAndReturn(func() error {
						time.Sleep(100 * time.Millisecond)
						return http.ErrServerClosed
					})

				mock.EXPECT().
					Shutdown(gomock.Any()).
					Return(errors.New("shutdown failure"))

				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockServer := tt.setup(ctrl)

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel() // prevent context leak

			if tt.cancelCtxDelay > 0 {
				go func() {
					time.Sleep(tt.cancelCtxDelay)
					cancel()
				}()
			}

			err := runners.RunServer(ctx, mockServer)

			if tt.expectErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectErrMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
