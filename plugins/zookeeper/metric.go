package zookeeper

import (
	"time"
)

type metric struct {
	zk   *Zookeeper
	name string
}

func newMetric(e *Zookeeper, name string) *metric {
	e.paths = append(e.paths, name)
	return &metric{e, name}
}

func (m *metric) Run(step time.Duration) {
	m.zk.run(step)
}

func (m *metric) Report() float64 {
	return float64(m.zk.stats[m.name])
}
