package riak

import (
    "time"
    "sync"
)

type Riak struct {
    start  sync.Once
    memory int
    server string
}

func New(srv string) *Riak {
    return &Riak{server: srv}
}

func (r *Riak) MemUsage() *metric {
    return newMetric(r, "mem_usage")
}


func (r *Riak) run(step time.Duration) {
    r.start.Do(func() {
        for _ = range time.NewTicker(step).C {
            r.collect()
        }
    })
}

func (r *Riak) gather(name string) float64 {
    if name == "mem_usage" {
        return float64(r.memory)
    }else{
        return 0
    }
}
