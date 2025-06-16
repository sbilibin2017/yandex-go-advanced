package apps

import (
	"net/http"

	"github.com/sbilibin2017/yandex-go-advanced/internal/configs"
	"github.com/sbilibin2017/yandex-go-advanced/internal/handlers"
	"github.com/sbilibin2017/yandex-go-advanced/internal/middlewares"
	"github.com/sbilibin2017/yandex-go-advanced/internal/repositories"
	"github.com/sbilibin2017/yandex-go-advanced/internal/routers"
	"github.com/sbilibin2017/yandex-go-advanced/internal/services"
	"github.com/sbilibin2017/yandex-go-advanced/internal/validators"
)

func NewServerApp(config *configs.ServerConfig) (*http.Server, error) {
	metricMemorySaveRepository := repositories.NewMetricMemorySaveRepository()
	metricMemoryGetRepository := repositories.NewMetricMemoryGetRepository()
	metricMemoryListRepository := repositories.NewMetricMemoryListRepository()

	metricUpdateService := services.NewMetricUpdateService(metricMemorySaveRepository, metricMemoryGetRepository)
	metricGetService := services.NewMetricGetService(metricMemoryGetRepository)
	metricListService := services.NewMetricListService(metricMemoryListRepository)

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

	middlewareList := []func(next http.Handler) http.Handler{
		middlewares.LoggingMiddleware,
	}
	metricRouter := routers.NewMetricRouter(
		metricUpdatePathHandler,
		metricUpdateBodyHandler,
		metricGetPathHandler,
		metricGetBodyHandler,
		metricListHTMLHandler,
		middlewareList...,
	)

	return &http.Server{
		Addr:    config.Address,
		Handler: metricRouter,
	}, nil
}
