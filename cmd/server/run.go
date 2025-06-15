package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/sbilibin2017/yandex-go-advanced/internal/handlers"
	"github.com/sbilibin2017/yandex-go-advanced/internal/repositories"
	"github.com/sbilibin2017/yandex-go-advanced/internal/services"
	"github.com/sbilibin2017/yandex-go-advanced/internal/storages"
	"github.com/sbilibin2017/yandex-go-advanced/internal/types"
)

var (
	metricRouter *chi.Mux
	srv          *http.Server
)

func run() error {
	setupMetricsRouter()
	setupServer()
	return runWithGracefulShutdown()
}

func setupServer() {
	srv = &http.Server{Addr: addressFlag, Handler: metricRouter}
}

func setupMetricsRouter() {
	memoryStorage := storages.NewMemoryStorage[types.MetricID, types.Metrics]()

	metricMemorySaveRepository := repositories.NewMetricMemorySaveRepository(memoryStorage)
	metricMemoryGetRepository := repositories.NewMetricMemoryGetRepository(memoryStorage)
	metricMemoryListRepository := repositories.NewMetricMemoryListRepository(memoryStorage)

	metricUpdateService := services.NewMetricUpdateService(metricMemorySaveRepository, metricMemoryGetRepository)
	metricGetService := services.NewMetricGetService(metricMemoryGetRepository)
	metricListService := services.NewMetricListService(metricMemoryListRepository)

	metricUpdatePathHandler := handlers.NewMetricUpdatePathHandler(metricUpdateService)
	metricGetPathHandler := handlers.NewMetricGetPathHandler(metricGetService)
	metricListHTMLHandler := handlers.NewMetricListHTMLHandler(metricListService)

	metricRouter = chi.NewRouter()

	metricRouter.Post("/update/{type}/{name}/{value}", metricUpdatePathHandler)
	metricRouter.Get("/value/{type}/{name}", metricGetPathHandler)
	metricRouter.Get("/list", metricListHTMLHandler)

}

func runWithGracefulShutdown() error {
	errChan := make(chan error, 1)
	defer close(errChan)

	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	defer stop()

	go func() {
		fmt.Println("Starting server...")
		err := srv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			errChan <- err
		}
	}()

	select {
	case <-ctx.Done():
		fmt.Println("Shutdown signal received...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("server shutdown failed: %w", err)
		}

		fmt.Println("Server shutdown gracefully")
		return nil

	case err := <-errChan:
		if err != nil {
			return fmt.Errorf("server error: %w", err)
		}
		return nil
	}
}
