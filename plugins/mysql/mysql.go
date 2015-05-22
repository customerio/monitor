package mysql

import (
	"time"

	"github.com/customerio/monitor/plugins"
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

const (
	queriesGauge = iota
	slowGauge
)

type MySQL struct {
	cs     string
	values []timediff
	gauges []metrics.GaugeFloat64
}

func New(connection_string string) *MySQL {
	return &MySQL{cs: connection_string,
		values: make([]timediff, 2),
		gauges: []metrics.GaugeFloat64{
			queriesGauge: plugins.Gauge("mysql.queries"),
			slowGauge:    plugins.Gauge("mysql.slow"),
		},
	}
}

func (m *MySQL) clear() {
	for i, _ := range m.values {
		m.values[i].set(0)
	}
}

func (m *MySQL) Run(step time.Duration) {
	for _ = range time.Tick(step) {
		m.collect()
		for i, v := range m.values {
			m.gauges[i].Update(v.gather())
		}
	}
}
