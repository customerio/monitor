package cpu

import (
	"time"
)

type metric struct {
	cpu  *CPU
	name string
}

func newMetric(c *CPU, name string) *metric {
	return &metric{c, name}
}

func (m *metric) Run(step time.Duration) {
	m.cpu.run(step)
}

func (m *metric) Report() float64 {
	return m.cpu.rate(m.name)
}
