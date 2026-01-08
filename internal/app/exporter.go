package app

import (
	"abb-exporter/internal/exporter"
	"abb-exporter/internal/exporter/clickhouse"
	"abb-exporter/internal/exporter/influx"
	"fmt"
)

func BuildExporter(cfg exporter.Config) (exporter.Exporter, error) {
	switch cfg.Type {
	case "influxdb":
		return influx.NewClient(influx.OptionsFromMap(cfg.Config))
	case "clickhouse":
		return clickhouse.NewClient(clickhouse.OptionsFromMap(cfg.Config))
	default:
		return nil, fmt.Errorf("unknown exporter type: %s", cfg.Type)
	}
}
