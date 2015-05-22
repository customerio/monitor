package foundationdb

import (
	"time"

	"github.com/customerio/monitor/plugins"
	"github.com/rcrowley/go-metrics"
)

const (
	// instance stats
	diskioGauge = iota
	ramGauge
	trafficGauge
	cpuGauge

	// Cluster stats
	readRateGauge
	writeRateGauge
	transactionRateGauge
	conflictRateGauge
)

type FoundationDB struct {
	port int

	gauges []metrics.GaugeFloat64
}

func New(port int) *FoundationDB {
	return &FoundationDB{port: port,
		gauges: []metrics.GaugeFloat64{
			diskioGauge:          plugins.Gauge("fdb.diskio"),
			trafficGauge:         plugins.Gauge("fdb.traffic"),
			cpuGauge:             plugins.Gauge("fdb.cpu"),
			ramGauge:             plugins.Gauge("fdb.ram"),
			readRateGauge:        plugins.Gauge("fdb.rate.read"),
			writeRateGauge:       plugins.Gauge("fdb.rate.write"),
			transactionRateGauge: plugins.Gauge("fdb.rate.transaction"),
			conflictRateGauge:    plugins.Gauge("fdb.rate.conflict"),
		},
	}
}

func (f *FoundationDB) clear() {
	for _, v := range f.gauges {
		v.Update(0)
	}
}

func (f *FoundationDB) Run(step time.Duration) {
	for _ = range time.Tick(step) {
		f.collect()
	}
}
