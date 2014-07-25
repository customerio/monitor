package foundationdb

import (
	"time"
)

type metric struct {
	fdb *FoundationDB
	name string
}

func newMetric(r *FoundationDB, name string) *metric {
	return &metric{r, name}
}

func (m *metric) Run(step time.Duration) {
	m.fdb.run(step)
}

func (m *metric) Report() float64 {
	return m.fdb.gather(m.name)
}
