package riak

import "time"

type Riak struct {
	memory     int
	gets       int
	puts       int
	index_gets int
	server     string
}

func New(srv string) *Riak {
	return &Riak{server: srv}
}

func (r *Riak) MemUsage() float64 {
	return float64(r.memory)
}

func (r *Riak) Gets() float64 {
	return float64(r.gets)
}

func (r *Riak) Puts() float64 {
	return float64(r.puts)
}

func (r *Riak) IndexGets() float64 {
	return float64(r.index_gets)
}
func (r *Riak) clear() {
	r.memory = 0
	r.gets = 0
	r.puts = 0
	r.index_gets = 0
}

func (r *Riak) Run(step time.Duration) {
	for _ = range time.Tick(step) {
		r.collect()
	}
}
