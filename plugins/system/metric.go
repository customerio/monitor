package system

import (
    "time"
)

type metric struct {
    system  *System
    name string
}

func newMetric(r *System, name string) *metric {
    return &metric{r, name}
}

func (m *metric) Run(step time.Duration) {
    m.system.run(step)
}

func (m *metric) Report() float64 {
    return m.system.gather(m.name)
}
