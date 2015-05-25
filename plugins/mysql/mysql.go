package mysql

import "github.com/customerio/monitor/metrics"

const (
	queriesGauge = iota
	slowGauge
)

type MySQL struct {
	cs       string
	updaters []metrics.Updater
}

func New(connection_string string) *MySQL {
	return &MySQL{cs: connection_string,
		updaters: []metrics.Updater{
			queriesGauge: metrics.NewCounter("mysql.queries"),
			slowGauge:    metrics.NewCounter("mysql.slow"),
		},
	}
}

func (m *MySQL) clear() {
	for i, _ := range m.updaters {
		m.updaters[i].Update(0)
	}
}
