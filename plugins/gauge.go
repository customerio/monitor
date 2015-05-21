package plugins

import "github.com/rcrowley/go-metrics"

func Gauge(name string) metrics.GaugeFloat64 {
	m := metrics.NewGaugeFloat64()
	metrics.Register(name, m)
	return m
}
