package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	ObservatoryAlive = prometheus.NewDesc(
		"xray_observatory_alive",
		"Whether the observatory probe is alive (1 = alive, 0 = dead)",
		[]string{"proxy"},
		nil,
	)
	ObservatoryDelay = prometheus.NewDesc(
		"xray_observatory_delay_ms",
		"Observed delay to the observatory probe in milliseconds",
		[]string{"proxy"},
		nil,
	)
	ObservatoryLastSeen = prometheus.NewDesc(
		"xray_observatory_last_seen_timestamp_seconds",
		"Unix timestamp of the last time the observatory probe was seen alive",
		[]string{"proxy"},
		nil,
	)
	InboundUplink = prometheus.NewDesc(
		"xray_inbound_uplink_bytes_total",
		"Total uplink bytes for an inbound proxy",
		[]string{"inbound"},
		nil,
	)
	InboundDownlink = prometheus.NewDesc(
		"xray_inbound_downlink_bytes_total",
		"Total downlink bytes for an inbound proxy",
		[]string{"inbound"},
		nil,
	)
	OutboundUplink = prometheus.NewDesc(
		"xray_outbound_uplink_bytes_total",
		"Total uplink bytes for an outbound proxy",
		[]string{"outbound"},
		nil,
	)
	OutboundDownlink = prometheus.NewDesc(
		"xray_outbound_downlink_bytes_total",
		"Total downlink bytes for an outbound proxy",
		[]string{"outbound"},
		nil,
	)
)
