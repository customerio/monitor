package system

import (
	"time"

	"github.com/rcrowley/go-metrics"
)

type System struct {
	loadAvg   metrics.GaugeFloat64
	memUsage  metrics.GaugeFloat64
	swapUsage metrics.GaugeFloat64
}

func gauge(registry metrics.Registry, name string) metrics.GaugeFloat64 {
	m := metrics.NewGaugeFloat64()
	registry.Register(name, m)
	return m
}

func New(registry metrics.Registry) *System {
	s := &System{}
	s.loadAvg = gauge(registry, "system.load")
	s.memUsage = gauge(registry, "system.mem_usage")
	s.swapUsage = gauge(registry, "system.swap_usage")
	return s
}

func (s *System) clear() {
	s.loadAvg.Update(0)
	s.memUsage.Update(0)
	s.swapUsage.Update(0)
}

func (s *System) Run(step time.Duration) {
	for _ = range time.Tick(step) {
		s.collect()
	}
}
