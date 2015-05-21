package elasticsearch

import "time"

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
	stats            map[string]int
}

func New(srv string) *Elasticsearch {
	return &Elasticsearch{
		server: srv,
		stats:  make(map[string]int),
	}
}

func (e *Elasticsearch) Status() float64 {
	return float64(e.stats["status"])
}

func (e *Elasticsearch) Nodes() float64 {
	return float64(e.stats["nodes"])
}

func (e *Elasticsearch) CPU() float64 {
	return float64(e.stats["cpu"])
}

func (e *Elasticsearch) Memory() float64 {
	return float64(e.stats["memory"])
}

func (e *Elasticsearch) Docs() float64 {
	return float64(e.stats["docs"])
}

func (e *Elasticsearch) Indexes() float64 {
	return float64(e.stats["indexes"])
}

func (e *Elasticsearch) Gets() float64 {
	return float64(e.stats["gets"])
}

func (e *Elasticsearch) Searches() float64 {
	return float64(e.stats["searches"])
}

func (e *Elasticsearch) clear() {
	e.stats = map[string]int{}
}

func (e *Elasticsearch) Run(step time.Duration) {
	for _ = range time.Tick(step) {
		e.collect()
	}
}
