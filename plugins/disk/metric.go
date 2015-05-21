package disk

import (
	"time"
)

type metric struct {
	disk *Disk
	name string
}

func newMetric(d *Disk, name string) *metric {
	return &metric{d, name}
}

func (m *metric) Run(step time.Duration) {
	m.disk.run(step)
}

func (m *metric) Report() float64 {
	return m.disk.rate(m.name)
}
