package types

const (
	Counter = "counter"
	Gauge   = "gauge"
)

type MetricID struct {
	ID    string `json:"id"`
	MType string `json:"type"`
}

type Metrics struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
	Hash  string   `json:"hash,omitempty"`
}
