package service

import (
	"abb-exporter/internal/exporter"
	"abb-exporter/internal/meter"
	"fmt"
	"log/slog"
)

func NewExporter(exporters []exporter.Exporter, meters []meter.Meter) *Exporter {
	return &Exporter{
		exporters: exporters,
		meters:    meters,
	}
}

type Exporter struct {
	meters    []meter.Meter
	exporters []exporter.Exporter
}

func (e Exporter) QueryAndExportMetrics() {
	var data *meter.UsageStatus
	var err error
	for _, m := range e.meters {
		data, err = m.QueryUsageStatus()
		if err != nil {
			slog.Error(fmt.Sprintf("cannot read registries on %s: %v", m.Name(), err))
			continue
		}

		go func(d *meter.UsageStatus) {
			_ = e.Export(d, m.Tags())
		}(data)
	}
}

func (e Exporter) Export(data *meter.UsageStatus, tags *map[string]string) error {
	var err error
	for _, exp := range e.exporters {
		err = exp.Metrics(data, tags)
		if err != nil {
			slog.Error(fmt.Sprintf("cannot export data: %v", err))
		}
	}

	return err
}
