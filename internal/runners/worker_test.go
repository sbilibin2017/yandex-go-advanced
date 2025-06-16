package runners

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRunWorker(t *testing.T) {
	type args struct {
		worker func(ctx context.Context) error
		ctx    context.Context
		cancel func()
	}

	tests := []struct {
		name            string
		args            args
		expectedErr     error
		cancelBeforeRun bool
	}{
		{
			name: "worker returns nil",
			args: args{
				worker: func(ctx context.Context) error {
					return nil
				},
				ctx: context.Background(),
			},
			expectedErr: nil,
		},
		{
			name: "worker returns error",
			args: args{
				worker: func(ctx context.Context) error {
					return errors.New("worker error")
				},
				ctx: context.Background(),
			},
			expectedErr: errors.New("worker error"),
		},
		{
			name: "context canceled before worker finishes",
			args: func() args {
				ctx, cancel := context.WithCancel(context.Background())
				return args{
					worker: func(ctx context.Context) error {
						time.Sleep(200 * time.Millisecond)
						return nil
					},
					ctx:    ctx,
					cancel: cancel,
				}
			}(),
			expectedErr:     context.Canceled,
			cancelBeforeRun: true,
		},
		{
			name: "context canceled after worker returns error",
			args: args{
				worker: func(ctx context.Context) error {
					time.Sleep(100 * time.Millisecond)
					return errors.New("worker error")
				},
				ctx: context.Background(),
			},
			expectedErr: errors.New("worker error"),
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.args.ctx
			if tt.cancelBeforeRun && tt.args.cancel != nil {
				tt.args.cancel()
			}

			err := RunWorker(ctx, tt.args.worker)

			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
