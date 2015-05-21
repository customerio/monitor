package system

import "time"

type System struct {
	loadAvg   float64
	memUsage  float64
	swapUsage float64
}

func New() *System {
	return &System{}
}

func (s *System) LoadAverage() float64 {
	return s.loadAvg
}

func (s *System) MemUsage() float64 {
	return s.memUsage
}

func (s *System) SwapUsage() float64 {
	return s.swapUsage
}

func (s *System) clear() {
	s.loadAvg = 0
	s.memUsage = 0
	s.swapUsage = 0
}

func (s *System) Run(step time.Duration) {
	for _ = range time.Tick(step) {
		s.collect()
	}
}
