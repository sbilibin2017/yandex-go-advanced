package runners_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/sbilibin2017/yandex-go-advanced/internal/runners"
)

func TestRun_SuccessfulStartStop(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRunnable := runners.NewMockRunnable(ctrl)

	// Start returns nil, simulating successful start and graceful shutdown
	mockRunnable.EXPECT().Start(gomock.Any()).Return(nil).Times(1)

	// Stop returns nil on graceful shutdown
	mockRunnable.EXPECT().Stop(gomock.Any()).Return(nil).Times(1)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Cancel context shortly to trigger shutdown flow
	go func() {
		time.Sleep(10 * time.Millisecond)
		cancel()
	}()

	err := runners.Run(ctx, mockRunnable)
	assert.NoError(t, err)
}

func TestRun_StartReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRunnable := runners.NewMockRunnable(ctrl)

	wantErr := errors.New("start error")

	// Start returns an error immediately
	mockRunnable.EXPECT().Start(gomock.Any()).Return(wantErr).Times(1)

	ctx := context.Background()

	err := runners.Run(ctx, mockRunnable)
	assert.Equal(t, wantErr, err)
}

func TestRun_StopReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRunnable := runners.NewMockRunnable(ctrl)

	// Start blocks until context is done, then returns nil to simulate graceful exit
	mockRunnable.EXPECT().Start(gomock.Any()).DoAndReturn(func(ctx context.Context) error {
		<-ctx.Done()
		return nil
	}).Times(1)

	wantErr := errors.New("stop error")

	// Stop returns error on shutdown
	mockRunnable.EXPECT().Stop(gomock.Any()).Return(wantErr).Times(1)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(10 * time.Millisecond)
		cancel()
	}()

	err := runners.Run(ctx, mockRunnable)
	assert.Equal(t, wantErr, err)
}
