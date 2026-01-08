package main

import (
	"abb-exporter/internal/api"
	"abb-exporter/internal/api/handler"
	"abb-exporter/internal/app"
	"abb-exporter/internal/config"
	"abb-exporter/internal/exporter"
	"abb-exporter/internal/infrastructure/http"
	"abb-exporter/internal/meter"
	"abb-exporter/internal/meter/abb_b21"
	"abb-exporter/internal/scheduler"
	"abb-exporter/internal/service"
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/simonvetter/modbus"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("cannot load config: %v", err)
	}

	slog.Info("ABB Exporter")
	mc, err := newModbusClient()
	if err != nil {
		log.Fatalf("cannot create modbus client: %v", err)
	}

	exporters, err := buildExportersClient(cfg.Exporters)
	if err != nil {
		slog.Error("cannot build exporters: %v", err)
		os.Exit(1)
	}
	slog.Info(fmt.Sprintf("Loaded %d exporets", len(exporters)))

	meters, err := buildMeters(mc, cfg.Meters)
	if err != nil {
		slog.Error("cannot build meters: %v", err)
		os.Exit(1)
	}
	slog.Info(fmt.Sprintf("Loaded %d meters", len(meters)))

	mg := meter.NewGroup(meters...)

	es := service.NewExporter(exporters, meters)

	sc := scheduler.NewScheduler(2 * time.Second)
	sc.Tick(es.QueryAndExportMetrics)
	sc.Start()

	// Server
	server := http.NewServer()
	meterApiHandler := handler.NewMeter(mg)
	api.BindApi(server.Router(), meterApiHandler)
	server.Start()

	// Listen for the interrupt signal
	<-ctx.Done()
	slog.Info("Shutting down...")
	_ = sc.Stop()
	_ = mc.Close()
	_ = server.Stop()
	for _, e := range exporters {
		_ = e.Close()
	}
}

func buildMeters(modbus *modbus.ModbusClient, mc []meter.Config) ([]meter.Meter, error) {
	meters := make([]meter.Meter, len(mc))

	for i, m := range mc {
		meters[i] = abb_b21.NewClient(modbus, m.Name, m.Slave)
	}

	return meters, nil
}

func newModbusClient() (*modbus.ModbusClient, error) {
	client, err := modbus.NewClient(&modbus.ClientConfiguration{
		URL:     "tcp://192.168.20.62:502",
		Timeout: 1 * time.Second,
	})
	if err != nil {
		return nil, err
	}
	err = client.Open()
	if err != nil {
		return nil, err
	}

	return client, nil
}

func buildExportersClient(cfg []exporter.Config) ([]exporter.Exporter, error) {
	var err error
	clients := make([]exporter.Exporter, len(cfg))

	for i, c := range cfg {
		clients[i], err = app.BuildExporter(c)
		if err != nil {
			return nil, err
		}
	}

	return clients, nil
}
