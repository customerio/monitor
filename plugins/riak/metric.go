package riak

import (
	"time"
)

type metric struct {
	riak *Riak
	name string
}

func newMetric(r *Riak, name string) *metric {
	return &metric{r, name}
}

func (m *metric) Run(step time.Duration) {
	m.riak.run(step)
}

func (m *metric) Report() float64 {
	return m.riak.gather(m.name)
}
