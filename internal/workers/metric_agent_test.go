package workers

import (
	"context"
	"errors"
	"strconv"
	"testing"
	"time"

	gomock "github.com/golang/mock/gomock"
	"github.com/sbilibin2017/yandex-go-advanced/internal/types"
	"github.com/stretchr/testify/require"
)

// --- Tests for collectRuntimeMetrics ---

func TestCollectRuntimeMetrics_EmitsMetrics(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	ch := collectRuntimeMetrics(ctx, 1) // 1 second poll interval

	count := 0
loop:
	for {
		select {
		case m, ok := <-ch:
			if !ok {
				break loop
			}
			require.NotNil(t, m)
			require.NotEmpty(t, m.ID)
			count++
		case <-time.After(4 * time.Second):
			t.Fatal("timeout waiting for metrics")
		}
	}

	require.Greater(t, count, 0, "Expected some metrics emitted")
}

func float64Ptr(f float64) *float64 {
	return &f
}

func TestUpdateMetrics_BufferingAndPeriodicFlush(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUpdater := NewMockMetricUpdater(ctrl)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	metricsIn := make(chan *types.Metrics)

	// Expect Update to be called at least once with some metrics
	mockUpdater.EXPECT().
		Update(gomock.Any(), gomock.Any()).
		Return(nil).
		MinTimes(1)

	errCh := updateMetrics(ctx, 1, mockUpdater, metricsIn) // 1 second flush interval

	go func() {
		for i := 0; i < 3; i++ {
			metricsIn <- &types.Metrics{
				ID:    "metric" + strconv.Itoa(i),
				Type:  "gauge",
				Value: float64Ptr(float64(i)),
			}
		}
		close(metricsIn)
	}()

	for err := range errCh {
		require.NoError(t, err)
	}
}

func TestUpdateMetrics_FlushOnInputChannelClose(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUpdater := NewMockMetricUpdater(ctrl)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	metricsIn := make(chan *types.Metrics)

	// Expect Update exactly once when channel closes
	mockUpdater.EXPECT().
		Update(gomock.Any(), gomock.Any()).
		Return(nil).
		Times(1)

	errCh := updateMetrics(ctx, 10, mockUpdater, metricsIn) // long flush interval, so flush only on close

	go func() {
		metricsIn <- &types.Metrics{
			ID:    "testMetric",
			Type:  "gauge",
			Value: float64Ptr(123),
		}
		close(metricsIn)
	}()

	for err := range errCh {
		require.NoError(t, err)
	}
}

func TestUpdateMetrics_UpdaterReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUpdater := NewMockMetricUpdater(ctrl)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	errExample := errors.New("update error")

	mockUpdater.EXPECT().
		Update(gomock.Any(), gomock.Any()).
		Return(errExample).
		AnyTimes()

	metricsIn := make(chan *types.Metrics)

	errCh := updateMetrics(ctx, 1, mockUpdater, metricsIn)

	go func() {
		metricsIn <- &types.Metrics{
			ID:    "errorMetric",
			Type:  "gauge",
			Value: float64Ptr(42),
		}
		close(metricsIn)
	}()

	errReceived := false
	for err := range errCh {
		require.Error(t, err)
		require.Equal(t, errExample, err)
		errReceived = true
		break
	}
	require.True(t, errReceived, "expected at least one error")
}

func TestUpdateMetrics_ContextDoneFlushesBuffer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUpdater := NewMockMetricUpdater(ctrl)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	metricsIn := make(chan *types.Metrics)

	errCh := updateMetrics(ctx, 10, mockUpdater, metricsIn) // Long interval to avoid periodic flush

	mockUpdater.EXPECT().
		Update(gomock.Any(), gomock.AssignableToTypeOf([]*types.Metrics{})).
		DoAndReturn(func(ctx context.Context, buffer []*types.Metrics) error {
			require.Len(t, buffer, 1)
			require.Equal(t, "bufferedMetric", buffer[0].ID)
			return nil
		}).
		Times(1)

	metricsIn <- &types.Metrics{
		ID:    "bufferedMetric",
		Type:  "gauge",
		Value: float64Ptr(42),
	}

	// Wait 2 seconds to ensure the metric is buffered before canceling the context
	time.Sleep(2 * time.Second)

	cancel()

	close(metricsIn)

	for err := range errCh {
		require.NoError(t, err)
	}
}

func TestUpdateMetrics_PeriodicFlushOnTicker(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUpdater := NewMockMetricUpdater(ctrl)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	metricsIn := make(chan *types.Metrics)

	// Use a short interval so ticker fires quickly
	reportInterval := 1

	errCh := updateMetrics(ctx, reportInterval, mockUpdater, metricsIn)

	// Expect Update to be called at least once due to ticker firing
	mockUpdater.EXPECT().
		Update(gomock.Any(), gomock.AssignableToTypeOf([]*types.Metrics{})).
		DoAndReturn(func(ctx context.Context, buffer []*types.Metrics) error {
			require.Greater(t, len(buffer), 0)
			return nil
		}).
		MinTimes(1)

	// Send some metrics to buffer
	for i := 0; i < 3; i++ {
		metricsIn <- &types.Metrics{
			ID:    "metric" + strconv.Itoa(i),
			Type:  "gauge",
			Value: float64Ptr(float64(i)),
		}
	}

	// Wait a bit to let ticker fire and flush the buffer
	time.Sleep(1500 * time.Millisecond)

	// Cancel context and close input to end the goroutine
	cancel()
	close(metricsIn)

	// Check errors (expect none)
	for err := range errCh {
		require.NoError(t, err)
	}
}

func TestUpdateMetrics_SendsErrorOnUpdateFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUpdater := NewMockMetricUpdater(ctrl)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	metricsIn := make(chan *types.Metrics)

	errExample := errors.New("update failed")

	// Expect Update to be called at least once and return an error
	mockUpdater.EXPECT().
		Update(gomock.Any(), gomock.AssignableToTypeOf([]*types.Metrics{})).
		Return(errExample).
		MinTimes(1)

	errCh := updateMetrics(ctx, 1, mockUpdater, metricsIn)

	// Send some metric to trigger buffering and update call
	go func() {
		metricsIn <- &types.Metrics{
			ID:    "metric_error",
			Type:  "gauge",
			Value: float64Ptr(42),
		}
		// Wait a bit to allow ticker or update call
		time.Sleep(500 * time.Millisecond)
		cancel()
		close(metricsIn)
	}()

	// Expect at least one error to be sent on errCh
	errReceived := false
	for err := range errCh {
		require.Error(t, err)
		require.Equal(t, errExample, err)
		errReceived = true
		break
	}
	require.True(t, errReceived, "expected an error sent on errCh")
}

func TestUpdateMetrics_TickerFlushSendsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUpdater := NewMockMetricUpdater(ctrl)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	metricsIn := make(chan *types.Metrics)

	testErr := errors.New("update failed")

	// We expect Update to be called when ticker flushes the buffer
	mockUpdater.EXPECT().
		Update(gomock.Any(), gomock.AssignableToTypeOf([]*types.Metrics{})).
		DoAndReturn(func(ctx context.Context, buffer []*types.Metrics) error {
			require.Greater(t, len(buffer), 0)
			return testErr
		}).
		Times(1)

	// Use a very short interval for the ticker to trigger flush quickly
	errCh := updateMetrics(ctx, 1, mockUpdater, metricsIn) // 1 second interval

	go func() {
		metricsIn <- &types.Metrics{
			ID:    "metric1",
			Type:  "gauge",
			Value: float64Ptr(1.0),
		}
		// Keep the channel open; ticker should trigger flush
	}()

	// Wait for the error from the ticker-triggered flush
	errReceived := false
	timeout := time.After(3 * time.Second)

	for !errReceived {
		select {
		case err, ok := <-errCh:
			if !ok {
				t.Fatal("errCh closed before error received")
			}
			require.Error(t, err)
			require.Equal(t, testErr, err)
			errReceived = true
		case <-timeout:
			t.Fatal("timeout waiting for error from ticker flush")
		}
	}
}

func TestStartMetricAgentWorker_RunAndStops(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUpdater := NewMockMetricUpdater(ctrl)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // отменяем сразу

	// Обновление не должно вызываться
	mockUpdater.EXPECT().
		Update(gomock.Any(), gomock.Any()).
		Times(0)

	done := make(chan struct{})

	go func() {
		startMetricAgentWorker(ctx, mockUpdater, 1, 1)
		close(done)
	}()

	select {
	case <-done:
		// воркер завершился успешно
	case <-time.After(100 * time.Millisecond):
		t.Fatal("воркер не завершился вовремя")
	}
}

func TestNewMetricAgentWorker_RunAndStops(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUpdater := NewMockMetricUpdater(ctrl)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// We allow Update to be called zero or more times, returning no error
	mockUpdater.EXPECT().
		Update(gomock.Any(), gomock.Any()).
		AnyTimes().
		Return(nil)

	worker := NewMetricAgentWorker(mockUpdater, 1, 1)

	done := make(chan struct{})

	go func() {
		worker(ctx)
		close(done)
	}()

	// Cancel context immediately to stop the worker
	cancel()

	select {
	case <-done:
		// worker stopped successfully
	case <-time.After(200 * time.Millisecond):
		t.Fatal("worker did not stop in time")
	}
}

func TestLogErrors_ExitsImmediatelyOnContextDone(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	errCh := make(chan error)

	logErrors(ctx, errCh)

	// Cancel context immediately to hit <-ctx.Done() case
	cancel()

	// Close channel to ensure goroutine can exit if it's still waiting on errCh
	close(errCh)

	// Wait briefly to allow goroutine to exit
	time.Sleep(50 * time.Millisecond)
}

func TestLogErrors_HandlesErrorsAndChannelClose(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errCh := make(chan error)

	logErrors(ctx, errCh)

	// Send some errors and nil to errCh
	errCh <- errors.New("test error")
	errCh <- nil

	// Close channel and expect goroutine to finish
	close(errCh)

	// Give some time for goroutine to process
	time.Sleep(50 * time.Millisecond)
}

func TestLogErrors_ContextDoneStopsGoroutine(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	errCh := make(chan error)

	logErrors(ctx, errCh)

	// Cancel context immediately
	cancel()

	// Give some time for goroutine to exit
	time.Sleep(50 * time.Millisecond)

	// Close errCh to unblock the goroutine if it waits for errors
	close(errCh)

	// Now do NOT send anything on errCh after closing it
	// Just wait a little to be sure no panic occurs
	time.Sleep(50 * time.Millisecond)
}
