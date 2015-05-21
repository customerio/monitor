package riak

import (
	"time"

	"github.com/customerio/monitor/plugins"
	"github.com/rcrowley/go-metrics"
)

type Riak struct {
	memory     metrics.GaugeFloat64
	gets       metrics.GaugeFloat64
	puts       metrics.GaugeFloat64
	index_gets metrics.GaugeFloat64
	server     string
}

func New(srv string) *Riak {
	r := &Riak{server: srv}
	r.memory = plugins.Gauge("riak.mem_usage")
	r.gets = plugins.Gauge("riak.gets")
	r.puts = plugins.Gauge("riak.puts")
	r.index_gets = plugins.Gauge("riak.index_gets")
	return r
}

func (r *Riak) clear() {
	r.memory.Update(0)
	r.gets.Update(0)
	r.puts.Update(0)
	r.index_gets.Update(0)
}

func (r *Riak) Run(step time.Duration) {
	for _ = range time.Tick(step) {
		r.collect()
	}
}
