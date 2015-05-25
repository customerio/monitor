package riak

import "github.com/customerio/monitor/metrics"

const (
	memoryGauge = iota
	getsGauge
	putsGauge
	indexGetsGauge
)

type Riak struct {
	server   string
	updaters []metrics.Updater
}

func New(srv string) *Riak {
	return &Riak{server: srv,
		updaters: []metrics.Updater{
			memoryGauge:    metrics.NewGauge("riak.mem_usage"),
			getsGauge:      metrics.NewGauge("riak.gets"),
			putsGauge:      metrics.NewGauge("riak.puts"),
			indexGetsGauge: metrics.NewGauge("riak.index_gets"),
		},
	}
}

func (r *Riak) clear() {
	for _, g := range r.updaters {
		g.Update(0)
	}
}
