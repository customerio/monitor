package riak

import (
	"time"

	"github.com/rcrowley/go-metrics"
)

type Riak struct {
	memory     metrics.GaugeFloat64
	gets       metrics.GaugeFloat64
	puts       metrics.GaugeFloat64
	index_gets metrics.GaugeFloat64
	server     string
}

func gauge(registry metrics.Registry, name string) metrics.GaugeFloat64 {
	m := metrics.NewGaugeFloat64()
	registry.Register(name, m)
	return m
}

func New(registry metrics.Registry, srv string) *Riak {
	r := &Riak{server: srv}
	r.memory = gauge(registry, "riak.mem_usage")
	r.gets = gauge(registry, "riak.gets")
	r.puts = gauge(registry, "riak.puts")
	r.index_gets = gauge(registry, "riak.index_gets")
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
