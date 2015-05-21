package elasticsearch

import (
	"time"

	"github.com/rcrowley/go-metrics"
)

const (
	RED    = 0
	YELLOW = 1
	GREEN  = 2
)

type Elasticsearch struct {
	server           string
	previousIndexes  int
	previousGets     int
	previousSearches int

	stats  map[string]int
	gauges map[string]metrics.GaugeFloat64
}

func gauge(registry metrics.Registry, name string) metrics.GaugeFloat64 {
	m := metrics.NewGaugeFloat64()
	registry.Register(name, m)
	return m
}

func New(registry metrics.Registry, srv string) *Elasticsearch {
	e := &Elasticsearch{
		server: srv,
		stats:  make(map[string]int),
		gauges: make(map[string]metrics.GaugeFloat64),
	}
	e.gauges["status"] = gauge(registry, "elastic.cluster")
	e.gauges["nodes"] = gauge(registry, "elastic.nodes")
	e.gauges["cpu"] = gauge(registry, "elastic.cpu")
	e.gauges["memory"] = gauge(registry, "elastic.memory")
	e.gauges["docs"] = gauge(registry, "elastic.docs")
	e.gauges["gets"] = gauge(registry, "elastic.indexes")
	e.gauges["indexes"] = gauge(registry, "elastic.gets")
	e.gauges["searches"] = gauge(registry, "elastic.searches")
	return e
}

func (e *Elasticsearch) clear() {
	e.stats = map[string]int{}
}

func (e *Elasticsearch) Run(step time.Duration) {
	for _ = range time.Tick(step) {
		e.collect()
		for k, g := range e.gauges {
			if v, ok := e.stats[k]; ok {
				g.Update(float64(v))
			}
		}
	}
}
