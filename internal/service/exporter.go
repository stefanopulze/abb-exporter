package service

import (
	"abb-exporter/internal/exporter"
	"abb-exporter/internal/meter"
	"fmt"
	"log/slog"
)

const logMaxReadFails = 3

func NewExporter(exporters []exporter.Exporter, meters []meter.Meter) *Exporter {
	readFails := map[int]int{}
	for i := range meters {
		readFails[i] = 0
	}

	return &Exporter{
		exporters: exporters,
		meters:    meters,
		readFails: readFails,
	}
}

type Exporter struct {
	meters    []meter.Meter
	exporters []exporter.Exporter
	readFails map[int]int
}

func (e Exporter) QueryAndExportMetrics() {
	var data *meter.UsageStatus
	var err error
	for i := range e.meters {
		data, err = e.meters[i].QueryUsageStatus()
		if err != nil {
			if e.readFails[i] >= logMaxReadFails {
				continue
			}

			e.readFails[i]++
			slog.Error(fmt.Sprintf("cannot read registries on %s: %v", e.meters[i].Name(), err))
			continue
		}

		e.readFails[i] = 0
		go func() {
			err = e.Export(data, e.meters[i].Tags())
			if err != nil {
				slog.Error(fmt.Sprintf("cannot export data: %v", err))
			}
		}()
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
