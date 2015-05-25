package elasticsearch

import "github.com/customerio/monitor/metrics"

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

	updaters []metrics.Updater
}

func New(srv string) *Elasticsearch {
	e := &Elasticsearch{
		server: srv,
		updaters: []metrics.Updater{
			statusGauge:   metrics.NewGauge("elastic.cluster"),
			nodesGauge:    metrics.NewGauge("elastic.nodes"),
			cpuGauge:      metrics.NewGauge("elastic.cpu"),
			memoryGauge:   metrics.NewGauge("elastic.memory"),
			docsGauge:     metrics.NewGauge("elastic.docs"),
			indexesGauge:  metrics.NewGauge("elastic.indexes"),
			getsGauge:     metrics.NewGauge("elastic.gets"),
			searchesGauge: metrics.NewGauge("elastic.searches"),
		},
	}
	return e
}

func (e *Elasticsearch) clear() {
	for _, g := range e.updaters {
		g.Update(0)
	}
}
