package system

import (
	"sync"
	"time"
)

type System struct {
	start     sync.Once
	loadAvg   float64
	memUsage  float64
	swapUsage float64
}

func New() *System {
	return &System{}
}

func (s *System) LoadAverage() *metric {
	return newMetric(s, "load_avg")
}

func (s *System) MemUsage() *metric {
	return newMetric(s, "mem_usage")
}

func (s *System) SwapUsage() *metric {
	return newMetric(s, "swap_usage")
}

func (s *System) clear() {
	s.loadAvg = 0
	s.memUsage = 0
	s.swapUsage = 0
}

func (s *System) run(step time.Duration) {
	s.start.Do(func() {
		for _ = range time.Tick(step) {
			s.collect()
		}
	})
}

func (s *System) gather(name string) float64 {
	switch name {
	case "load_avg":
		return s.loadAvg
	case "mem_usage":
		return s.memUsage
	case "swap_usage":
		return s.swapUsage
	default:
		return 0
	}
}
