# xray-stats-exporter

Prometheus exporter for [XTLs/Xray-core](https://github.com/XTLS/Xray-core) stats API.

On each scrape, fetches the Xray expvars endpoint, updates internal counters, and exposes metrics in Prometheus text format.

## Metrics

| Metric | Type | Labels | Description |
|--------|------|--------|-------------|
| `xray_observatory_alive` | Gauge | `proxy` | Probe alive status (1/0) |
| `xray_observatory_delay_ms` | Gauge | `proxy` | Probe delay in milliseconds |
| `xray_observatory_last_seen_timestamp_seconds` | Gauge | `proxy` | Unix timestamp of last seen |
| `xray_inbound_uplink_bytes_total` | Counter | `inbound` | Total uplink bytes |
| `xray_inbound_downlink_bytes_total` | Counter | `inbound` | Total downlink bytes |
| `xray_outbound_uplink_bytes_total` | Counter | `outbound` | Total uplink bytes |
| `xray_outbound_downlink_bytes_total` | Counter | `outbound` | Total downlink bytes |

## Quick Start

```bash
# Build
go build -o xray-stats-exporter ./cmd/server

# Run
XRAY_STATS_URL=http://127.0.0.1:11111/debug/vars LISTEN_ADDR=:8080 ./xray-stats-exporter

# Scrape
curl http://localhost:8080/metrics
```

## Configuration

| Environment Variable | Default | Description |
|---------------------|---------|-------------|
| `XRAY_STATS_URL` | `http://127.0.0.1:11111/debug/vars` | Xray stats API endpoint |
| `LISTEN_ADDR` | `:8080` | Address to listen on |

## Prometheus Config

```yaml
scrape_configs:
  - job_name: xray
    static_configs:
      - targets: ['localhost:8080']
```

## Cross-compile for Raspberry Pi

```bash
GOOS=linux GOARCH=arm64 go build -o xray-stats-exporter-linux-arm64 ./cmd/server
```
