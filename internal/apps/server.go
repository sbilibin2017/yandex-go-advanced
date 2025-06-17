// Package apps provides application-level components for starting and stopping services
// such as the HTTP server application. These components implement the Runnable interface
// to support lifecycle management via centralized runners.
package apps

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/yandex-go-advanced/internal/configs"
	"github.com/sbilibin2017/yandex-go-advanced/internal/handlers"
	"github.com/sbilibin2017/yandex-go-advanced/internal/middlewares"
	"github.com/sbilibin2017/yandex-go-advanced/internal/repositories"
	"github.com/sbilibin2017/yandex-go-advanced/internal/routers"
	"github.com/sbilibin2017/yandex-go-advanced/internal/services"
	"github.com/sbilibin2017/yandex-go-advanced/internal/validators"
)

// ServerApp represents the HTTP server application.
//
// It encapsulates the full setup and lifecycle of the server, including
// HTTP routing, middleware, service and repository initialization, and graceful shutdown.
//
// This struct is intended to be managed by a runner that supports the Runnable interface.
type ServerApp struct {
	server *http.Server
}

// NewServerApp creates and initializes a new instance of ServerApp using the provided server configuration.
//
// This function wires together repositories, services, validators, handlers, middleware, and the router.
//
// Parameters:
//   - config: Pointer to a ServerConfig that defines the server address and log level.
//
// Returns:
//   - A pointer to a ServerApp instance ready to be started.
//   - An error, if any setup fails.
func NewServerApp(config *configs.ServerConfig) (*ServerApp, error) {
	// Initialize repositories
	metricMemorySaveRepository := repositories.NewMetricMemorySaveRepository()
	metricMemoryGetRepository := repositories.NewMetricMemoryGetRepository()
	metricMemoryListRepository := repositories.NewMetricMemoryListRepository()

	// Initialize services
	metricUpdateService := services.NewMetricUpdateService(metricMemorySaveRepository, metricMemoryGetRepository)
	metricGetService := services.NewMetricGetService(metricMemoryGetRepository)
	metricListService := services.NewMetricListService(metricMemoryListRepository)

	// Initialize handlers with validation
	metricUpdatePathHandler := handlers.NewMetricUpdatePathHandler(
		validators.ValidateMetricAttributes,
		metricUpdateService,
	)
	metricUpdateBodyHandler := handlers.NewMetricUpdateBodyHandler(
		validators.ValidateMetric,
		metricUpdateService,
	)
	metricGetPathHandler := handlers.NewMetricGetPathHandler(
		validators.ValidateMetricIDAttributes,
		metricGetService,
	)
	metricGetBodyHandler := handlers.NewMetricGetBodyHandler(
		validators.ValidateMetricID,
		metricGetService,
	)
	metricListHTMLHandler := handlers.NewMetricListHTMLHandler(metricListService)

	// Register middleware
	middlewareList := []func(http.Handler) http.Handler{
		middlewares.LoggingMiddleware,
		middlewares.GzipMiddleware,
	}

	// Set up router
	metricRouter := routers.NewMetricRouter(
		metricUpdatePathHandler,
		metricUpdateBodyHandler,
		metricGetPathHandler,
		metricGetBodyHandler,
		metricListHTMLHandler,
		middlewareList...,
	)

	// Create HTTP server
	httpServer := &http.Server{
		Addr:    config.Address,
		Handler: metricRouter,
	}

	return &ServerApp{
		server: httpServer,
	}, nil
}

// Start runs the HTTP server and blocks until it shuts down or encounters an error.
//
// This method satisfies the Runnable interface.
//
// Parameters:
//   - ctx: Context for managing cancellation and timeout.
//
// Returns:
//   - An error if the server fails to start or crashes during runtime.
func (app *ServerApp) Start(ctx context.Context) error {
	return app.server.ListenAndServe()
}

// Stop gracefully shuts down the HTTP server using the provided context.
//
// This method satisfies the Runnable interface and allows the server to finish
// processing ongoing requests before terminating.
//
// Parameters:
//   - ctx: Context for controlling shutdown timeout and cancellation.
//
// Returns:
//   - An error if shutdown fails or times out.
func (app *ServerApp) Stop(ctx context.Context) error {
	return app.server.Shutdown(ctx)
}
