package riak

import (
	"time"

	"github.com/customerio/monitor/plugins"
	"github.com/rcrowley/go-metrics"
)

const (
	memoryGauge = iota
	getsGauge
	putsGauge
	indexGetsGauge
)

type Riak struct {
	server string
	gauges []metrics.GaugeFloat64
}

func New(srv string) *Riak {
	return &Riak{server: srv,
		gauges: []metrics.GaugeFloat64{
			memoryGauge:    plugins.Gauge("riak.mem_usage"),
			getsGauge:      plugins.Gauge("riak.gets"),
			putsGauge:      plugins.Gauge("riak.puts"),
			indexGetsGauge: plugins.Gauge("riak.index_gets"),
		},
	}
}

func (r *Riak) clear() {
	for _, g := range r.gauges {
		g.Update(0)
	}
}

func (r *Riak) Run(step time.Duration) {
	for _ = range time.Tick(step) {
		r.collect()
	}
}
