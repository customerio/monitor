package system

import (
	"time"

	"github.com/customerio/monitor/plugins"
	"github.com/rcrowley/go-metrics"
)

type System struct {
	loadAvg   metrics.GaugeFloat64
	memUsage  metrics.GaugeFloat64
	swapUsage metrics.GaugeFloat64
}

func New() *System {
	s := &System{}
	s.loadAvg = plugins.Gauge("system.load")
	s.memUsage = plugins.Gauge("system.mem_usage")
	s.swapUsage = plugins.Gauge("system.swap_usage")
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
