package mysql

import (
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/rcrowley/go-metrics"
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
	cs           string
	queries      timediff
	slow         timediff
	queriesGauge metrics.GaugeFloat64
	slowGauge    metrics.GaugeFloat64
}

func gauge(registry metrics.Registry, name string) metrics.GaugeFloat64 {
	m := metrics.NewGaugeFloat64()
	registry.Register(name, m)
	return m
}

func New(registry metrics.Registry, connection_string string) *MySQL {
	m := &MySQL{cs: connection_string}
	m.queriesGauge = gauge(registry, "mysql.queries")
	m.slowGauge = gauge(registry, "mysql.slow")
	return m
}

func (m *MySQL) clear() {
	m.queries.set(0)
	m.slow.set(0)
}

func (m *MySQL) Run(step time.Duration) {
	for _ = range time.Tick(step) {
		m.collect()
		m.queriesGauge.Update(m.queries.gather())
		m.slowGauge.Update(m.slow.gather())
	}
}
