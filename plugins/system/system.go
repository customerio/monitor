package system

import (
    "time"
    "sync"
)

type System struct {
    start  sync.Once
    loadAvg float64
    memUsage float64
    swap float64
}

func New() *System {
    return &System{}
}

func (s *System) LoadAverage() *metric {
    return newMetric(s, "load_avg")
}

func (s *System) MemUsage() *metric {
    return newMetric(s, "mem_usage")
}

func (s *System) SwapOuts() *metric {
    return newMetric(s, "swap")
}


func (s *System) run(step time.Duration) {
    s.start.Do(func() {
        for _ = range time.NewTicker(step).C {
            s.collect()
        }
    })
}

func (s *System) gather(name string) float64 {
    switch name {
        case "load_avg": return s.loadAvg
        case "mem_usage": return s.memUsage
        case "swap": return s.swap
        default: return 0
    }
}
