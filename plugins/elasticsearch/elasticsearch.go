package elasticsearch

import (
	"sync"
	"time"
)

const (
	RED    = 0
	YELLOW = 1
	GREEN  = 2
)

type Elasticsearch struct {
	start  sync.Once
	server string
	stats  map[string]int
}

func New(srv string) *Elasticsearch {
	return &Elasticsearch{
		server: srv,
		stats:  make(map[string]int),
	}
}

func (e *Elasticsearch) Status() *metric {
	return newMetric(e, "status")
}

func (e *Elasticsearch) Nodes() *metric {
	return newMetric(e, "nodes")
}

func (e *Elasticsearch) CPU() *metric {
	return newMetric(e, "cpu")
}

func (e *Elasticsearch) Memory() *metric {
	return newMetric(e, "memory")
}

func (e *Elasticsearch) Docs() *metric {
	return newMetric(e, "docs")
}

func (e *Elasticsearch) Indexes() *metric {
	return newMetric(e, "indexes")
}

func (e *Elasticsearch) Gets() *metric {
	return newMetric(e, "gets")
}

func (e *Elasticsearch) run(step time.Duration) {
	e.start.Do(func() {
		for _ = range time.NewTicker(step).C {
			e.collect()
		}
	})
}
