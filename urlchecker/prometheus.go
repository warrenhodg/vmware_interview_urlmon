package urlchecker

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	urlUps           *prometheus.GaugeVec
	urlResponseTimes *prometheus.GaugeVec
)

// setupMetrics sets up the relevant gauges
func setupMetrics() {
	urlUps = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "sample",
		Subsystem: "external_url",
		Name:      "up",
		Help:      "0 if url is down, 1 if it is up",
	}, []string{"url"})

	urlResponseTimes = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "sample",
		Subsystem: "external_url",
		Name:      "response_ms",
		Help:      "number of milliseconds to receive response",
	}, []string{"url"})
}

// boolToFloat converts false/true to 0/1
func boolToFloat(v bool) float64 {
	if v {
		return 1
	}
	return 0
}

// observe records metric results
func observe(url string, up bool, d time.Duration) {
	urlUp := urlUps.WithLabelValues(url)
	urlUp.Set(boolToFloat(up))

	urlResponseTime := urlResponseTimes.WithLabelValues(url)
	urlResponseTime.Set(float64(d.Milliseconds()))
}
