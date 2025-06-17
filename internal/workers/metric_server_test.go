package workers

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/sbilibin2017/yandex-go-advanced/internal/types"
)

func TestLoadFromFile_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSaver := NewMockMetricContextSaver(ctrl)
	mockLister := NewMockMetricFileLister(ctrl)

	var delta1 int64 = 10
	var delta2 int64 = 20

	metrics := []types.Metrics{
		{ID: "metric1", Type: "counter", Delta: &delta1},
		{ID: "metric2", Type: "counter", Delta: &delta2},
	}

	mockLister.EXPECT().List(gomock.Any()).Return(metrics, nil)
	mockSaver.EXPECT().Save(gomock.Any(), metrics[0]).Return(nil)
	mockSaver.EXPECT().Save(gomock.Any(), metrics[1]).Return(nil)

	err := loadFromFile(context.Background(), mockSaver, mockLister)
	require.NoError(t, err)
}

func TestLoadFromFile_ListError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSaver := NewMockMetricContextSaver(ctrl)
	mockLister := NewMockMetricFileLister(ctrl)

	mockLister.EXPECT().List(gomock.Any()).Return(nil, errors.New("list error"))

	err := loadFromFile(context.Background(), mockSaver, mockLister)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list error")
}

func TestSaveToFile_ListError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLister := NewMockMetricContextLister(ctrl)
	mockSaver := NewMockMetricFileSaver(ctrl)

	mockLister.EXPECT().List(gomock.Any()).Return(nil, errors.New("list error"))

	err := saveToFile(context.Background(), mockLister, mockSaver)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list error")
}

func TestNewMetricServerWorker_StoreIntervalZero(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSaverContext := NewMockMetricContextSaver(ctrl)
	mockListerContext := NewMockMetricContextLister(ctrl)
	mockSaverFile := NewMockMetricFileSaver(ctrl)
	mockListerFile := NewMockMetricFileLister(ctrl)

	// No restore on startup
	// So no call to listerFile.List expected
	mockListerFile.EXPECT().List(gomock.Any()).Times(0)

	var delta int64 = 10
	metrics := []types.Metrics{{ID: "metric1", Type: "counter", Delta: &delta}}

	// Expect SaveToFile on shutdown
	mockListerContext.EXPECT().List(gomock.Any()).Return(metrics, nil)
	mockSaverFile.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil)

	ctx, cancel := context.WithCancel(context.Background())

	worker := NewMetricServerWorker(
		mockSaverContext,
		mockListerContext,
		mockSaverFile,
		mockListerFile,
		0, // storeInterval=0 disables periodic saving; save only on shutdown
		false,
	)

	go worker(ctx)

	// Trigger shutdown immediately
	cancel()

	time.Sleep(10 * time.Millisecond) // Let the worker finish
}

func TestNewMetricServerWorker_RestoreSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSaverContext := NewMockMetricContextSaver(ctrl)
	mockListerContext := NewMockMetricContextLister(ctrl)
	mockSaverFile := NewMockMetricFileSaver(ctrl)
	mockListerFile := NewMockMetricFileLister(ctrl)

	var delta int64 = 10
	metrics := []types.Metrics{{ID: "metric1", Type: "counter", Delta: &delta}}

	// Restore from file: load metrics and save to context
	mockListerFile.EXPECT().List(gomock.Any()).Return(metrics, nil)
	mockSaverContext.EXPECT().Save(gomock.Any(), metrics[0]).Return(nil)

	// During runtime: periodically save metrics to file
	mockListerContext.EXPECT().List(gomock.Any()).Return(metrics, nil).AnyTimes()
	mockSaverFile.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	ctx, cancel := context.WithCancel(context.Background())

	worker := NewMetricServerWorker(
		mockSaverContext,
		mockListerContext,
		mockSaverFile,
		mockListerFile,
		1, // 1 second interval for periodic saving
		true,
	)

	go worker(ctx)

	time.Sleep(1100 * time.Millisecond) // Let the periodic save run once

	cancel() // Trigger shutdown

	time.Sleep(10 * time.Millisecond) // Let shutdown complete
}

func TestStartMetricServerWorker_FailureInSave(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSaverContext := NewMockMetricContextSaver(ctrl)
	mockListerContext := NewMockMetricContextLister(ctrl)
	mockSaverFile := NewMockMetricFileSaver(ctrl)
	mockListerFile := NewMockMetricFileLister(ctrl)

	// Restore fails, error is logged, but no panic
	mockListerFile.EXPECT().List(gomock.Any()).Return(nil, errors.New("restore error"))

	var delta int64 = 10
	metrics := []types.Metrics{{ID: "metric1", Type: "counter", Delta: &delta}}

	// SaveToFile fails during periodic save and shutdown — logged errors only
	mockListerContext.EXPECT().List(gomock.Any()).Return(metrics, nil).AnyTimes()
	mockSaverFile.EXPECT().Save(gomock.Any(), gomock.Any()).Return(errors.New("save error")).AnyTimes()

	ctx, cancel := context.WithCancel(context.Background())

	worker := NewMetricServerWorker(
		mockSaverContext,
		mockListerContext,
		mockSaverFile,
		mockListerFile,
		1,
		true,
	)

	go worker(ctx)

	time.Sleep(1100 * time.Millisecond)

	cancel()

	time.Sleep(10 * time.Millisecond)
}

func TestSaveToFile_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLister := NewMockMetricContextLister(ctrl)
	mockSaver := NewMockMetricFileSaver(ctrl)

	var value1 = 1.23
	var value2 = 4.56

	metrics := []types.Metrics{
		{ID: "metric1", Type: "gauge", Value: &value1},
		{ID: "metric2", Type: "gauge", Value: &value2},
	}

	mockLister.EXPECT().List(gomock.Any()).Return(metrics, nil)
	mockSaver.EXPECT().Save(gomock.Any(), metrics[0]).Return(nil)
	mockSaver.EXPECT().Save(gomock.Any(), metrics[1]).Return(nil)

	err := saveToFile(context.Background(), mockLister, mockSaver)
	require.NoError(t, err)
}

func TestLoadFromFile_SaveError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSaver := NewMockMetricContextSaver(ctrl)
	mockLister := NewMockMetricFileLister(ctrl)

	var delta int64 = 42

	metrics := []types.Metrics{
		{ID: "metric1", Type: "counter", Delta: &delta},
	}

	mockLister.EXPECT().List(gomock.Any()).Return(metrics, nil)
	mockSaver.EXPECT().Save(gomock.Any(), metrics[0]).Return(errors.New("save error"))

	// Call the function under test
	err := loadFromFile(context.Background(), mockSaver, mockLister)

	// The loadFromFile should NOT return error despite Save failure (just logs it)
	require.NoError(t, err)
}

func TestStartMetricServerWorker_RestoreError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSaverContext := NewMockMetricContextSaver(ctrl)
	mockListerContext := NewMockMetricContextLister(ctrl)
	mockSaverFile := NewMockMetricFileSaver(ctrl)
	mockListerFile := NewMockMetricFileLister(ctrl)

	// Restore fails with error
	mockListerFile.EXPECT().List(gomock.Any()).Return(nil, errors.New("restore error"))

	// For shutdown save - it should still call saveToFile (which calls listerContext.List)
	mockListerContext.EXPECT().List(gomock.Any()).Return([]types.Metrics{}, nil)
	mockSaverFile.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil).Times(0) // no metrics to save

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go startMetricServerWorker(
		ctx,
		mockSaverContext,
		mockListerContext,
		mockSaverFile,
		mockListerFile,
		0,    // storeInterval=0 disables ticker, only save on shutdown
		true, // restore enabled
	)

	// Cancel to trigger shutdown path
	cancel()

	// Give goroutine a moment to complete
	time.Sleep(10 * time.Millisecond)
}

func TestStartMetricServerWorker_ContextDoneWithSaveError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSaverContext := NewMockMetricContextSaver(ctrl)
	mockListerContext := NewMockMetricContextLister(ctrl)
	mockSaverFile := NewMockMetricFileSaver(ctrl)
	mockListerFile := NewMockMetricFileLister(ctrl)

	// No restore
	mockListerFile.EXPECT().List(gomock.Any()).Return(nil, nil).Times(0)

	// On shutdown, saveToFile is called and returns error (simulated by Save returning error)
	mockListerContext.EXPECT().List(gomock.Any()).Return([]types.Metrics{
		{ID: "m1", Type: "counter"},
	}, nil)
	mockSaverFile.EXPECT().Save(gomock.Any(), gomock.Any()).Return(errors.New("save error"))

	ctx, cancel := context.WithCancel(context.Background())

	go startMetricServerWorker(
		ctx,
		mockSaverContext,
		mockListerContext,
		mockSaverFile,
		mockListerFile,
		0,     // zero store interval, so only saves on shutdown
		false, // no restore
	)

	// Cancel to trigger shutdown path
	cancel()

	// Give goroutine a moment to complete
	time.Sleep(10 * time.Millisecond)
}
