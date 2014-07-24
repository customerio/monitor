package riak

import (
	"sync"
	"time"
)

type Riak struct {
	start      sync.Once
	memory     int
	gets       int
	puts       int
	index_gets int
	server     string
}

func New(srv string) *Riak {
	return &Riak{server: srv}
}

func (r *Riak) MemUsage() *metric {
	return newMetric(r, "mem_usage")
}

func (r *Riak) Gets() *metric {
    return newMetric(r, "gets")
}

func (r *Riak) Puts() *metric {
    return newMetric(r, "puts")
}

func (r *Riak) IndexGets() *metric {
    return newMetric(r, "index_gets")
}

func (r *Riak) run(step time.Duration) {
	r.start.Do(func() {
		for _ = range time.NewTicker(step).C {
			r.collect()
		}
	})
}

func (r *Riak) gather(name string) float64 {

    switch name {
        case "mem_usage": return float64(r.memory)
        case "gets": return float64(r.gets)
        case "puts": return float64(r.puts)
        case "index_gets": return float64(r.index_gets)
        default: return 0
    }
}
