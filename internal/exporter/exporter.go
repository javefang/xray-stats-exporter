package exporter

import (
	"context"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/xfang/xray-stats-exporter/internal/client"
	"github.com/xfang/xray-stats-exporter/internal/metrics"
)

type fetcher interface {
	Fetch(ctx context.Context, url string) (*client.XrayStats, error)
}

type promCollector struct {
	fetcher fetcher
	url     string
}

func MustRegister(f fetcher, statsURL string) {
	prometheus.MustRegister(&promCollector{fetcher: f, url: statsURL})
}

func NewHandler(mux *http.ServeMux) {
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Xray Stats Exporter\n"))
	})
}

func (c *promCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- metrics.ObservatoryAlive
	ch <- metrics.ObservatoryDelay
	ch <- metrics.ObservatoryLastSeen
	ch <- metrics.InboundUplink
	ch <- metrics.InboundDownlink
	ch <- metrics.OutboundUplink
	ch <- metrics.OutboundDownlink
}

func (c *promCollector) Collect(ch chan<- prometheus.Metric) {
	stats, err := c.fetcher.Fetch(context.Background(), c.url)
	if err != nil {
		return
	}

	// Observatory metrics
	for name, probe := range stats.Observatory {
		alive := 0.0
		if probe.Alive {
			alive = 1.0
		}
		ch <- prometheus.MustNewConstMetric(
			metrics.ObservatoryAlive, prometheus.GaugeValue, alive, name,
		)
		ch <- prometheus.MustNewConstMetric(
			metrics.ObservatoryDelay, prometheus.GaugeValue, float64(probe.Delay), name,
		)
		ch <- prometheus.MustNewConstMetric(
			metrics.ObservatoryLastSeen, prometheus.GaugeValue, float64(probe.LastSeenTime), name,
		)
	}

	// Inbound stats
	for inbound, counters := range stats.Stats.Inbound {
		if v, ok := counters["uplink"]; ok {
			ch <- prometheus.MustNewConstMetric(
				metrics.InboundUplink, prometheus.CounterValue, float64(v), inbound,
			)
		}
		if v, ok := counters["downlink"]; ok {
			ch <- prometheus.MustNewConstMetric(
				metrics.InboundDownlink, prometheus.CounterValue, float64(v), inbound,
			)
		}
	}

	// Outbound stats
	for outbound, counters := range stats.Stats.Outbound {
		if v, ok := counters["uplink"]; ok {
			ch <- prometheus.MustNewConstMetric(
				metrics.OutboundUplink, prometheus.CounterValue, float64(v), outbound,
			)
		}
		if v, ok := counters["downlink"]; ok {
			ch <- prometheus.MustNewConstMetric(
				metrics.OutboundDownlink, prometheus.CounterValue, float64(v), outbound,
			)
		}
	}
}
