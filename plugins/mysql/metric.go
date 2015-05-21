package mysql

import (
	"time"
)

type metric struct {
	mysql *MySQL
	name  string
}

func newMetric(m *MySQL, name string) *metric {
	return &metric{m, name}
}

func (m *metric) Run(step time.Duration) {
	m.mysql.run(step)
}

func (m *metric) Report() float64 {
	return m.mysql.gather(m.name)
}
