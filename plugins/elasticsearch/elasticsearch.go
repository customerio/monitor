package elasticsearch

import (
	"time"

	"github.com/customerio/monitor/plugins"
	"github.com/rcrowley/go-metrics"
)

const (
	statusRed    = 0
	statusYellow = 1
	statusGreen  = 2
)

const (
	statusGauge = iota
	nodesGauge
	cpuGauge
	memoryGauge
	docsGauge
	getsGauge
	indexesGauge
	searchesGauge
)

type Elasticsearch struct {
	server           string
	previousIndexes  int
	previousGets     int
	previousSearches int

	gauges []metrics.GaugeFloat64
}

func New(srv string) *Elasticsearch {
	e := &Elasticsearch{
		server: srv,
		gauges: []metrics.GaugeFloat64{
			statusGauge:   plugins.Gauge("elastic.cluster"),
			nodesGauge:    plugins.Gauge("elastic.nodes"),
			cpuGauge:      plugins.Gauge("elastic.cpu"),
			memoryGauge:   plugins.Gauge("elastic.memory"),
			docsGauge:     plugins.Gauge("elastic.docs"),
			indexesGauge:  plugins.Gauge("elastic.indexes"),
			getsGauge:     plugins.Gauge("elastic.gets"),
			searchesGauge: plugins.Gauge("elastic.searches"),
		},
	}
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
