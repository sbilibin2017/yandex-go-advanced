package workers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/sbilibin2017/yandex-go-advanced/internal/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCollectRuntimeGaugeMetrics(t *testing.T) {
	metricsList := collectRuntimeGaugeMetrics()
	assert.Len(t, metricsList, 28, "Gauge metrics count should be 28")

}

func TestCollectRuntimeCounterMetrics(t *testing.T) {
	metricsList := collectRuntimeCounterMetrics()
	assert.Len(t, metricsList, 1, "Counter metrics count should be 1")
}

func TestPollMetrics(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second) // > 2 сек
	defer cancel()

	floatPtr := func(f float64) *float64 { return &f }
	intPtr := func(i int64) *int64 { return &i }

	mockCollector := func() []types.Metrics {
		return []types.Metrics{
			{ID: "test_gauge", MType: types.Gauge, Value: floatPtr(123.45)},
			{ID: "test_counter", MType: types.Counter, Delta: intPtr(7)},
		}
	}

	out := pollMetrics(ctx, 1, mockCollector) // pollInterval = 1 sec

	var collected []types.Metrics
	for m := range out {
		collected = append(collected, m)
	}

	require.GreaterOrEqual(t, len(collected), 2, "Expected at least one collector batch (2 metrics)")
	assert.True(t, len(collected)%2 == 0, "Each collector batch should return exactly 2 metrics")
}

func TestWaitForContextOrError(t *testing.T) {
	// Case 1: context done returns context error
	ctx1, cancel1 := context.WithCancel(context.Background())
	cancel1() // cancel immediately

	err := waitForContextOrError(ctx1, nil)
	assert.ErrorIs(t, err, context.Canceled)

	// Case 2: error received from error channel
	ctx2 := context.Background()
	errCh2 := make(chan error, 1)
	errCh2 <- errors.New("some error")
	close(errCh2)

	err = waitForContextOrError(ctx2, errCh2)
	assert.EqualError(t, err, "some error")

	// Case 3: error channel closed with no errors
	ctx3 := context.Background()
	errCh3 := make(chan error)
	close(errCh3)

	err = waitForContextOrError(ctx3, errCh3)
	assert.NoError(t, err)
}

func TestReportMetrics(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	floatPtr := func(f float64) *float64 { return &f }
	intPtr := func(i int64) *int64 { return &i }

	metricsCh := make(chan types.Metrics, 10)

	m1 := types.Metrics{ID: "metric1", MType: types.Gauge, Value: floatPtr(10.0)}
	m2 := types.Metrics{ID: "metric2", MType: types.Counter, Delta: intPtr(5)}

	var mu sync.Mutex
	calls := make(map[string]int)

	handler := func(ctx context.Context, serverAddress string, m types.Metrics) error {
		mu.Lock()
		defer mu.Unlock()
		calls[m.ID]++
		return nil
	}

	errCh := reportMetrics(ctx, "http://localhost", handler, 1, 2, metricsCh)

	metricsCh <- m1
	metricsCh <- m2
	close(metricsCh)

	var errs []error
	for err := range errCh {
		errs = append(errs, err)
	}

	require.Empty(t, errs, "expected no errors from handler")

	mu.Lock()
	defer mu.Unlock()
	assert.Equal(t, 1, calls[m1.ID], "handler should be called once for m1")
	assert.Equal(t, 1, calls[m2.ID], "handler should be called once for m2")
}

func TestReportMetrics_WorkerReportsErrors(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	floatPtr := func(f float64) *float64 { return &f }
	intPtr := func(i int64) *int64 { return &i }

	metricsCh := make(chan types.Metrics, 10)

	m1 := types.Metrics{ID: "metric1", MType: types.Gauge, Value: floatPtr(10.0)}
	m2 := types.Metrics{ID: "metric2", MType: types.Counter, Delta: intPtr(5)}

	handler := func(ctx context.Context, serverAddress string, m types.Metrics) error {
		if m.ID == "metric2" {
			return errors.New("update failed")
		}
		return nil
	}

	errCh := reportMetrics(ctx, "http://localhost", handler, 1, 1, metricsCh)

	metricsCh <- m1
	metricsCh <- m2
	close(metricsCh)

	var errs []error
	for err := range errCh {
		errs = append(errs, err)
	}

	require.Len(t, errs, 1)
	assert.EqualError(t, errs[0], "update failed")
}

func TestReportMetrics_FlushOnTicker(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1500*time.Millisecond)
	defer cancel()

	floatPtr := func(f float64) *float64 { return &f }
	intPtr := func(i int64) *int64 { return &i }

	metricsCh := make(chan types.Metrics, 10)

	m1 := types.Metrics{ID: "metric1", MType: types.Gauge, Value: floatPtr(42.0)}
	m2 := types.Metrics{ID: "metric2", MType: types.Counter, Delta: intPtr(8)}

	calls := make(map[string]int)
	var mu sync.Mutex

	handler := func(ctx context.Context, serverAddress string, m types.Metrics) error {
		mu.Lock()
		defer mu.Unlock()
		calls[m.ID]++
		return nil
	}

	errCh := reportMetrics(ctx, "http://localhost", handler, 1, 1, metricsCh)

	metricsCh <- m1
	metricsCh <- m2

	// Wait for ticker flush
	time.Sleep(1200 * time.Millisecond)

	cancel()

	var errs []error
	for err := range errCh {
		errs = append(errs, err)
	}

	require.Empty(t, errs)

	mu.Lock()
	defer mu.Unlock()
	assert.Equal(t, 1, calls[m1.ID])
	assert.Equal(t, 1, calls[m2.ID])
}

func TestSendMetrics(t *testing.T) {
	// Create a test server to mock metric receiving endpoint
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check that the request method is POST
		assert.Equal(t, http.MethodPost, r.Method)

		// Validate URL path format: /update/{type}/{id}/{value}
		pathParts := strings.Split(r.URL.Path, "/")
		assert.Len(t, pathParts, 5) // "", "update", MType, ID, Value

		metricType := pathParts[2]
		id := pathParts[3]
		value := pathParts[4]

		// Basic validation to make sure URL parts are not empty
		assert.NotEmpty(t, metricType)
		assert.NotEmpty(t, id)
		assert.NotEmpty(t, value)

		// Respond with 200 OK
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	ctx := context.Background()

	// Test Gauge metric
	gaugeValue := 123.45
	gaugeMetric := types.Metrics{
		MType: types.Gauge,
		ID:    "test_gauge",
		Value: &gaugeValue,
	}

	err := sendMetrics(ctx, ts.URL, gaugeMetric)
	assert.NoError(t, err, "sendMetrics should not return error for valid gauge metric")

	// Test Counter metric
	deltaValue := int64(10)
	counterMetric := types.Metrics{
		MType: types.Counter,
		ID:    "test_counter",
		Delta: &deltaValue,
	}

	err = sendMetrics(ctx, ts.URL, counterMetric)
	assert.NoError(t, err, "sendMetrics should not return error for valid counter metric")

	// Test error on unknown metric type
	unknownMetric := types.Metrics{
		MType: "unknown",
		ID:    "id",
	}

	err = sendMetrics(ctx, ts.URL, unknownMetric)
	assert.Error(t, err, "sendMetrics should return error for unknown metric type")

	// Test error when Gauge value is nil
	nilValueMetric := types.Metrics{
		MType: types.Gauge,
		ID:    "nil_value",
		Value: nil,
	}

	err = sendMetrics(ctx, ts.URL, nilValueMetric)
	assert.Error(t, err, "sendMetrics should return error when Gauge value is nil")

	// Test error when Counter delta is nil
	nilDeltaMetric := types.Metrics{
		MType: types.Counter,
		ID:    "nil_delta",
		Delta: nil,
	}

	err = sendMetrics(ctx, ts.URL, nilDeltaMetric)
	assert.Error(t, err, "sendMetrics should return error when Counter delta is nil")
}

func TestStartMetricAgentWorker_WithMockServer(t *testing.T) {
	// Create a test server that expects POST requests on /update/ paths
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST method, got %s", r.Method)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if !strings.HasPrefix(r.URL.Path, "/update/") {
			t.Errorf("Unexpected path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer testServer.Close()

	// Set up a context with a timeout to stop the worker automatically
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Start the metric agent worker with short intervals and 2 workers
	err := StartMetricAgentWorker(ctx, testServer.URL, 1, 1, 2)

	// Accept context cancellation or deadline exceeded as expected shutdown signals
	if err != nil && err != context.Canceled && err != context.DeadlineExceeded {
		t.Fatalf("Unexpected error from StartMetricAgentWorker: %v", err)
	}
}
