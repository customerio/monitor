package foundationdb

import "github.com/customerio/monitor/metrics"

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

	updaters []metrics.Updater
}

func New(port int) *FoundationDB {
	return &FoundationDB{port: port,
		updaters: []metrics.Updater{
			diskioGauge:          metrics.NewGauge("fdb.diskio"),
			trafficGauge:         metrics.NewGauge("fdb.traffic"),
			cpuGauge:             metrics.NewGauge("fdb.cpu"),
			ramGauge:             metrics.NewGauge("fdb.ram"),
			readRateGauge:        metrics.NewGauge("fdb.rate.read"),
			writeRateGauge:       metrics.NewGauge("fdb.rate.write"),
			transactionRateGauge: metrics.NewGauge("fdb.rate.transaction"),
			conflictRateGauge:    metrics.NewGauge("fdb.rate.conflict"),
		},
	}
}

func (f *FoundationDB) clear() {
	for _, v := range f.updaters {
		v.Update(0)
	}
}
