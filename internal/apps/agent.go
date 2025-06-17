// Package apps provides application-level components for launching
// background agents and servers. These components implement the Runnable
// interface to integrate with centralized runners for start/stop control.
package apps

import (
	"context"

	"github.com/sbilibin2017/yandex-go-advanced/internal/configs"
	"github.com/sbilibin2017/yandex-go-advanced/internal/facades"
	"github.com/sbilibin2017/yandex-go-advanced/internal/workers"
)

// AgentApp represents the background metric collection agent.
//
// It runs a worker routine responsible for polling and reporting metrics
// to a remote server. Implements the Runnable interface for lifecycle management.
type AgentApp struct {
	worker func(ctx context.Context)
}

// NewAgentApp initializes and returns a new AgentApp.
//
// It creates a MetricUpdateFacade based on the provided configuration and
// constructs a worker function that handles polling and reporting intervals.
//
// Parameters:
//   - config: AgentConfig containing server address, polling, and reporting intervals.
//
// Returns:
//   - Pointer to an AgentApp instance ready to be started.
//   - An error if initialization fails (currently always nil).
func NewAgentApp(
	config *configs.AgentConfig,
) (*AgentApp, error) {
	metricUpdateFacade := facades.NewMetricUpdateFacade(config.ServerAddress)

	worker := workers.NewMetricAgentWorker(
		metricUpdateFacade,
		config.PollInterval,
		config.ReportInterval,
	)

	return &AgentApp{worker: worker}, nil
}

// Start launches the background metric agent worker.
//
// This method blocks until the provided context is canceled.
// It satisfies the Runnable interface.
//
// Parameters:
//   - ctx: Context to control cancellation and timeout of the worker.
//
// Returns:
//   - An error if the worker exits unexpectedly (currently always nil).
func (app *AgentApp) Start(ctx context.Context) error {
	app.worker(ctx)
	return nil
}

// Stop performs cleanup or shutdown of the agent.
//
// For AgentApp, Stop is a no-op because the worker respects context cancellation.
// It satisfies the Runnable interface.
//
// Parameters:
//   - ctx: Context to control timeout or cancellation.
//
// Returns:
//   - An error if shutdown fails (currently always nil).
func (app *AgentApp) Stop(ctx context.Context) error {
	return nil
}
