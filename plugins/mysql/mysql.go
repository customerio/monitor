package mysql

import (
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
	cs      string
	queries timediff
	slow    timediff
}

func New(connection_string string) *MySQL {
	return &MySQL{cs: connection_string}
}

func (m *MySQL) Queries() float64 {
	return m.queries.gather()
}
func (m *MySQL) SlowQueries() float64 {
	return m.slow.gather()
}

func (m *MySQL) clear() {
	m.queries.set(0)
	m.slow.set(0)
}

func (m *MySQL) Run(step time.Duration) {
	for _ = range time.Tick(step) {
		m.collect()
	}
}
