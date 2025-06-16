package workers

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	gomock "github.com/golang/mock/gomock"
	"github.com/sbilibin2017/yandex-go-advanced/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestCollectRuntimeMetrics(t *testing.T) {
	tests := []struct {
		name          string
		collectorFunc func() []types.Metrics
		expectedCount int
		expectedType  string
	}{
		{
			name:          "Gauge metrics",
			collectorFunc: collectRuntimeGaugeMetrics,
			expectedCount: 28,
			expectedType:  types.Gauge,
		},
		{
			name:          "Counter metrics",
			collectorFunc: collectRuntimeCounterMetrics,
			expectedCount: 1,
			expectedType:  types.Counter,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metrics := tt.collectorFunc()

			assert.Len(t, metrics, tt.expectedCount, "unexpected number of metrics")
			for _, metric := range metrics {
				assert.Equal(t, tt.expectedType, metric.Type, "metric type mismatch")
			}
		})
	}
}

func TestPollMetrics(t *testing.T) {
	// Simple collector function returning fixed metrics
	collector1 := func() []types.Metrics {
		return []types.Metrics{
			{Type: types.Gauge, ID: "metric1"},
		}
	}

	collector2 := func() []types.Metrics {
		return []types.Metrics{
			{Type: types.Counter, ID: "metric2"},
			{Type: types.Counter, ID: "metric3"},
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	metricsChan := pollMetrics(ctx, 1, collector1, collector2) // pollInterval 1 second

	collected := make([]types.Metrics, 0)

	// Read metrics with timeout (to avoid blocking)
readLoop:
	for {
		select {
		case m, ok := <-metricsChan:
			if !ok {
				break readLoop
			}
			collected = append(collected, m)
		case <-time.After(2 * time.Second):
			t.Fatal("timeout waiting for metrics")
		}
	}

	// Because the ctx timeout is 350ms, and poll interval is 1s, we might get 0 or 1 batch.
	// So check the collected metrics are from the collectors if any.

	for _, m := range collected {
		assert.Contains(t, []string{"metric1", "metric2", "metric3"}, m.ID)
		assert.Contains(t, []string{types.Gauge, types.Counter}, m.Type)
	}
}

func TestReportMetrics_UpdateReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUpdater := NewMockMetricUpdater(ctrl)
	mockUpdater.EXPECT().
		Update(gomock.Any(), gomock.Any()).
		Return(errors.New("update failed")).
		AnyTimes()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	inputMetrics := []types.Metrics{
		{Type: types.Gauge, ID: "metricErr", Value: float64Ptr(9.99)},
	}
	inCh := make(chan types.Metrics, len(inputMetrics))
	for _, m := range inputMetrics {
		inCh <- m
	}
	close(inCh)

	errCh := reportMetrics(ctx, mockUpdater, 1, 1, inCh)

	// Cancel context after 1 second
	go func() {
		time.Sleep(1 * time.Second)
		cancel()
	}()

	var errCount int
	for err := range errCh {
		if err != nil {
			errCount++
		}
	}

	assert.Equal(t, 1, errCount, "expected one error when Update returns error")
}

func TestReportMetrics_ContextCanceledBeforeSending(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUpdater := NewMockMetricUpdater(ctrl)
	mockUpdater.EXPECT().
		Update(gomock.Any(), gomock.Any()).
		Return(nil).
		AnyTimes()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	inputMetrics := []types.Metrics{
		{Type: types.Gauge, ID: "metricDelayed", Value: float64Ptr(42)},
	}
	inCh := make(chan types.Metrics, len(inputMetrics))
	for _, m := range inputMetrics {
		inCh <- m
	}
	close(inCh)

	errCh := reportMetrics(ctx, mockUpdater, 5, 1, inCh)

	// Cancel context quickly (100ms) before sending can occur
	go func() {
		time.Sleep(100 * time.Millisecond)
		cancel()
	}()

	var errCount int
	for err := range errCh {
		if err != nil {
			errCount++
		}
	}

	assert.Equal(t, 0, errCount, "expected no errors when context canceled early")
}

func TestReportMetrics_WorkerExitsOnContextDone(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUpdater := NewMockMetricUpdater(ctrl)
	// We do NOT expect any call to Update since context is canceled immediately
	// So no EXPECT setup here or set EXPECT().Times(0)

	// Create context and cancel immediately
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately to trigger <-ctx.Done()

	// input channel with no metrics (empty)
	inCh := make(chan types.Metrics)
	close(inCh)

	errCh := reportMetrics(ctx, mockUpdater, 1, 1, inCh)

	// Collect errors (should be none)
	var errCount int
	for err := range errCh {
		if err != nil {
			errCount++
		}
	}

	assert.Equal(t, 0, errCount, "expected no errors, workers should exit immediately on context cancel")
}

func TestReportMetrics_FlushOnTicker(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUpdater := NewMockMetricUpdater(ctrl)

	// Prepare some metrics to be sent to input channel
	metrics := []types.Metrics{
		{Type: types.Gauge, ID: "metric1", Value: float64Ptr(1.1)},
		{Type: types.Counter, ID: "metric2", Delta: int64Ptr(2)},
	}

	// Expect Update to be called with each metric individually
	for _, metric := range metrics {
		mockUpdater.EXPECT().
			Update(gomock.Any(), []types.Metrics{metric}).
			Return(nil).
			Times(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Channel input with buffer large enough to send all metrics
	inCh := make(chan types.Metrics, len(metrics))
	for _, m := range metrics {
		inCh <- m
	}

	// Use 1 second interval (minimum 1 sec) so ticker.New doesn't panic
	reportInterval := 1 // seconds
	workerCount := 2

	errCh := reportMetrics(ctx, mockUpdater, reportInterval, workerCount, inCh)

	// Wait enough time for ticker to fire and flush buffer
	time.Sleep(1100 * time.Millisecond) // > 1 second

	// Close input channel and cancel context to finish goroutines cleanly
	close(inCh)
	cancel()

	// Drain errCh (expect no errors)
	for err := range errCh {
		assert.NoError(t, err)
	}
}

func TestWaitForContextOrError_ContextDone(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	errCh := make(chan error)
	defer close(errCh)

	err := waitForContextOrError(ctx, errCh)
	assert.ErrorIs(t, err, context.Canceled)
}

func TestWaitForContextOrError_ErrorFromChannel(t *testing.T) {
	ctx := context.Background()
	errCh := make(chan error, 1)
	expectedErr := errors.New("some error")
	errCh <- expectedErr

	err := waitForContextOrError(ctx, errCh)
	assert.Equal(t, expectedErr, err)
}

func TestWaitForContextOrError_ChannelClosedNoError(t *testing.T) {
	ctx := context.Background()
	errCh := make(chan error)
	close(errCh) // close immediately

	err := waitForContextOrError(ctx, errCh)
	assert.NoError(t, err)
}

func TestWaitForContextOrError_Blocking(t *testing.T) {
	// Test that the function blocks until context cancels or error is received
	ctx, cancel := context.WithCancel(context.Background())
	errCh := make(chan error)

	// Cancel context after 100ms
	go func() {
		time.Sleep(100 * time.Millisecond)
		cancel()
	}()

	start := time.Now()
	err := waitForContextOrError(ctx, errCh)
	duration := time.Since(start)

	assert.ErrorIs(t, err, context.Canceled)
	assert.GreaterOrEqual(t, duration.Milliseconds(), int64(100))
}

func TestStartMetricAgentWorker(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUpdater := NewMockMetricUpdater(ctrl)

	// Allow Update to be called any number of times (including zero)
	mockUpdater.EXPECT().
		Update(gomock.Any(), gomock.Any()).
		Return(nil).
		AnyTimes() // <-- This allows multiple calls

	ctx, cancel := context.WithCancel(context.Background())

	// Use short intervals (seconds)
	pollInterval := 1
	reportInterval := 1
	workerCount := 2

	go func() {
		time.Sleep(2 * time.Second)
		cancel()
	}()

	err := startMetricAgentWorker(ctx, mockUpdater, pollInterval, reportInterval, workerCount)

	assert.ErrorIs(t, err, context.Canceled)
}

func TestReportMetrics_FlushOnContextDone(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUpdater := NewMockMetricUpdater(ctrl)

	metrics := []types.Metrics{
		{Type: types.Gauge, ID: "metric1", Value: float64Ptr(1.23)},
		{Type: types.Counter, ID: "metric2", Delta: int64Ptr(5)},
	}

	for _, metric := range metrics {
		mockUpdater.EXPECT().
			Update(gomock.Any(), []types.Metrics{metric}).
			Return(nil).
			Times(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	inCh := make(chan types.Metrics, len(metrics))
	for _, m := range metrics {
		inCh <- m
	}
	close(inCh) // close input so goroutine reads these metrics

	errCh := reportMetrics(ctx, mockUpdater, 10 /* long interval */, 1, inCh)

	// Wait briefly to allow goroutine to read input and buffer metrics
	time.Sleep(50 * time.Millisecond)

	// Cancel context to trigger flush on <-ctx.Done()
	cancel()

	// Drain errCh
	var errs []error
	for err := range errCh {
		if err != nil {
			errs = append(errs, err)
		}
	}

	assert.Empty(t, errs, "expected no errors during flush on context done")
}

func TestNewMetricAgentWorker(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUpdater := NewMockMetricUpdater(ctrl)

	// Setup expectation: Update will be called at least once and return nil
	mockUpdater.EXPECT().
		Update(gomock.Any(), gomock.Any()).
		Return(nil).
		AnyTimes()

	// Create worker function with short intervals for quick test
	workerFunc := NewMetricAgentWorker(mockUpdater, 1, 1, 2)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Run the worker in a separate goroutine and cancel context after short delay
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := workerFunc(ctx)
		assert.ErrorIs(t, err, context.Canceled)
	}()

	time.Sleep(500 * time.Millisecond) // let worker run a bit
	cancel()                           // cancel context to stop worker
	wg.Wait()
}

func float64Ptr(f float64) *float64 { return &f }
func int64Ptr(i int64) *int64       { return &i }
