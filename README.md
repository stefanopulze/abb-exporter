# ABB Exporter

Export ABB B21 power meters data into `ClickHouse` or `InfluxDB`.  
You can use it with `Grafana` to plot your power consumption.

It's simple, and now works with `ABB B21` power meters, but you can implement other meters.  
Meters data are exported every 2 seconds.



## Configuration
```yaml
logging:
  level: "debug"
  type: "console"

meters:
  - name: "home"
    slave: 1
    type: "b21"
  - name: "enel"
    slave: 2
    type: "b21"

exporters:
  - type: "influxdb"
    config:
      host: 
      token: 
      database:
      organization:
  - type: "clickhouse"
    config:
      host: "192.168.1.52:9000"
      database: 
      username: 
      password:
```

### ClickHouse table

```sql
CREATE TABLE power_meter
(
    timestamp DateTime,
    voltage   Float32,
    current   Float32,
    frequency Float32,
    power     Float32,
    name      String
) ENGINE = MergeTree()
PARTITION BY toYYYYMM(timestamp)
ORDER BY (timestamp);
```