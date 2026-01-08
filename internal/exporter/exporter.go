package exporter

import (
	"abb-exporter/internal/meter"
)

type Exporter interface {
	Close() error
	Metrics(data *meter.UsageStatus, tags *map[string]string) error
}
