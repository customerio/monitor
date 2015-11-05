package system

import "github.com/customerio/monitor/metrics"

const (
	loadAvgGauge = iota
	memUsageGauge
	swapUsageGauge
)

type System struct {
	updaters []metrics.Updater
}

func New() *System {
	return &System{
		updaters: []metrics.Updater{
			loadAvgGauge:   metrics.NewGauge("system.load"),
			memUsageGauge:  metrics.NewGauge("system.mem_usage"),
			swapUsageGauge: metrics.NewGauge("system.swap_usage"),
		},
	}
}

func (s *System) clear() {
	for _, g := range s.updaters {
		g.Update(0)
	}
}

func (s *System) Collect(b *metrics.Batch) {
	s.collect()
	for _, u := range s.updaters {
		u.Fill(b)
	}
}
