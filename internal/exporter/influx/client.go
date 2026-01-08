package influx

import (
	"abb-exporter/internal/exporter"
	"abb-exporter/internal/meter"
	"context"

	"github.com/InfluxCommunity/influxdb3-go/v2/influxdb3"
)

func NewClient(opts Options) (exporter.Exporter, error) {
	client, err := influxdb3.New(influxdb3.ClientConfig{
		Host:         opts.Host,
		Token:        opts.Token,
		Database:     opts.Database,
		Organization: opts.Organization,
	})
	if err != nil {
		return nil, err
	}

	return &Client{
		client: client,
	}, nil
}

type Client struct {
	client *influxdb3.Client
}

func (c *Client) Close() error {
	return c.client.Close()
}

func (c *Client) Name() string {
	return "influx"
}

func (c *Client) Metrics(data *meter.UsageStatus, tags *map[string]string) error {
	fields := map[string]any{
		"voltage":      data.Voltage,
		"current":      data.Current,
		"active_power": data.ActivePower,
		"frequency":    data.Frequency,
	}

	point := influxdb3.NewPoint("power_meter", *tags, fields, data.Timestamp)

	return c.client.WritePoints(context.Background(), []*influxdb3.Point{point})
}
