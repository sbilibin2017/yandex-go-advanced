package repositories

import (
	"sync"

	"github.com/sbilibin2017/yandex-go-advanced/internal/types"
)

var metrics map[types.MetricID]types.Metrics = make(map[types.MetricID]types.Metrics)
var mu sync.RWMutex
