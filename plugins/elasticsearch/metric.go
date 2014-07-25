package elasticsearch

import (
	"time"
)

type metric struct {
	elastic *Elasticsearch
	name    string
}

func newMetric(e *Elasticsearch, name string) *metric {
	return &metric{e, name}
}

func (m *metric) Run(step time.Duration) {
	m.elastic.run(step)
}

func (m *metric) Report() float64 {
	return float64(m.elastic.stats[m.name])
}
