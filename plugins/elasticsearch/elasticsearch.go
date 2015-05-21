package elasticsearch

import (
	"time"

	"github.com/customerio/monitor/plugins"
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

	gauges map[string]metrics.GaugeFloat64
}

func New(srv string) *Elasticsearch {
	e := &Elasticsearch{
		server: srv,
		gauges: make(map[string]metrics.GaugeFloat64),
	}
	e.gauges["status"] = plugins.Gauge("elastic.cluster")
	e.gauges["nodes"] = plugins.Gauge("elastic.nodes")
	e.gauges["cpu"] = plugins.Gauge("elastic.cpu")
	e.gauges["memory"] = plugins.Gauge("elastic.memory")
	e.gauges["docs"] = plugins.Gauge("elastic.docs")
	e.gauges["gets"] = plugins.Gauge("elastic.indexes")
	e.gauges["indexes"] = plugins.Gauge("elastic.gets")
	e.gauges["searches"] = plugins.Gauge("elastic.searches")
	return e
}

func (e *Elasticsearch) clear() {
	for _, g := range e.gauges {
		g.Update(0)
	}
}

func (e *Elasticsearch) Run(step time.Duration) {
	for _ = range time.Tick(step) {
		e.collect()
	}
}
