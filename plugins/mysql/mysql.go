package mysql

import (
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type timediff struct {
	current int
	prev    int
}

func (t *timediff) set(v int) {
	if t.prev == 0 {
		t.prev = v
	} else {
		t.prev = t.current
	}
	t.current = v
}

func (t *timediff) gather() float64 {
	return float64(t.current - t.prev)
}

type MySQL struct {
	start   sync.Once
	cs      string
	queries timediff
	slow    timediff
}

func New(connection_string string) *MySQL {
	return &MySQL{cs: connection_string}
}

func (m *MySQL) Queries() *metric {
	return newMetric(m, "queries")
}
func (m *MySQL) SlowQueries() *metric {
	return newMetric(m, "slow")
}

func (m *MySQL) clear() {
	m.queries.set(0)
	m.slow.set(0)
}

func (m *MySQL) run(step time.Duration) {
	m.start.Do(func() {
		for _ = range time.Tick(step) {
			m.collect()
		}
	})
}

func (m *MySQL) gather(name string) float64 {
	switch name {
	case "queries":
		return m.queries.gather()
	case "slow":
		return m.slow.gather()
	}
	return 0
}
