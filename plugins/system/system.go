package system

import (
	"time"

	"github.com/customerio/monitor/plugins"
	"github.com/rcrowley/go-metrics"
)

const (
	loadAvgGauge = iota
	memUsageGauge
	swapUsageGauge
)

type System struct {
	gauges []metrics.GaugeFloat64
}

func New() *System {
	return &System{
		gauges: []metrics.GaugeFloat64{
			loadAvgGauge:   plugins.Gauge("system.load"),
			memUsageGauge:  plugins.Gauge("system.mem_usage"),
			swapUsageGauge: plugins.Gauge("system.swap_usage"),
		},
	}
}

func (s *System) clear() {
	for _, g := range s.gauges {
		g.Update(0)
	}
}

func (s *System) Run(step time.Duration) {
	for _ = range time.Tick(step) {
		s.collect()
	}
}
