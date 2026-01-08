package clickhouse

import (
	"abb-exporter/internal/meter"
	"context"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

const pigsInsertQuery = `INSERT INTO power_meter(
timestamp,
voltage,
current,
frequency,
power,
name
) values (?,?,?,?,?,?);`

func NewClient(opts Options) (*Client, error) {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{opts.Addr},
		Auth: clickhouse.Auth{
			Database: opts.Database,
			Username: opts.Username,
			Password: opts.Password,
		},
	})
	if err != nil {
		return nil, err
	}

	return &Client{
		conn: conn,
	}, nil
}

type Client struct {
	conn driver.Conn
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) Name() string {
	return "clickHouse"
}

func (c *Client) Metrics(data *meter.UsageStatus, tags *map[string]string) error {
	name := (*tags)["name"]
	return c.conn.Exec(context.Background(), pigsInsertQuery,
		data.Timestamp,
		data.Voltage,
		data.Current,
		data.Frequency,
		data.ActivePower,
		name)
}
