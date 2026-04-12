# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

A single-binary Go Prometheus exporter. On each scrape request it calls
XTLS/Xray-core's stats api (expvars format) to fetch application status,
updates internal counters, and returns metrics in Prometheus text format.

## Commands

```bash
go build -o xray-stats-exporter ./cmd/server
go test ./...
go test ./... -run TestName   # run a single test
go vet ./...
go fmt ./...
```

## Architecture

- `cmd/server/main.go` — entrypoint, wires up HTTP server and exporter
- `internal/exporter/` — Prometheus metric definitions and scraping logic
- `internal/client/` — HTTP client for the Xray expvars metrics API
- `internal/metrics/` — metric descriptors (counters, gauges, etc.)

The exporter exposes `/metrics`. The scrape flow: HTTP request → call xray stats api → update internal
counters → format Prometheus output. Call to external api happens at scrape time

## Configuration

Environment variables (no config file):
- `XRAY_STATS_URL` — full URL to the Xray status API (default: `http://127.0.0.1:11111/debug/vars`)
- `LISTEN_ADDR` — address to listen on (default `:8080`)

## Data format

File `xray-expvars-sample-output.json` contains a sample output from Xray's
metrics endpoint. Use that as an example to build schema of the source data.
You should understand the output to determine prometheus metric name, labels
and metric types.
