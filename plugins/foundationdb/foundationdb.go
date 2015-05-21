package foundationdb

import (
	"time"

	"github.com/customerio/monitor/plugins"
	"github.com/rcrowley/go-metrics"
)

type FoundationDB struct {
	port int

	// instance stats
	diskio  metrics.GaugeFloat64
	ram     metrics.GaugeFloat64
	traffic metrics.GaugeFloat64
	cpu     metrics.GaugeFloat64

	// Cluster stats
	read_rate        metrics.GaugeFloat64
	write_rate       metrics.GaugeFloat64
	transaction_rate metrics.GaugeFloat64
	conflict_rate    metrics.GaugeFloat64
}

func New(port int) *FoundationDB {
	f := &FoundationDB{port: port}

	f.diskio = plugins.Gauge("fdb.diskio")
	f.traffic = plugins.Gauge("fdb.traffic")
	f.cpu = plugins.Gauge("fdb.cpu")
	f.ram = plugins.Gauge("fdb.ram")
	f.read_rate = plugins.Gauge("fdb.rate.read")
	f.write_rate = plugins.Gauge("fdb.rate.write")
	f.transaction_rate = plugins.Gauge("fdb.rate.transaction")
	f.conflict_rate = plugins.Gauge("fdb.rate.conflict")

	return f
}

func (f *FoundationDB) clear() {
	f.diskio.Update(0)
	f.ram.Update(0)
	f.traffic.Update(0)
	f.cpu.Update(0)
	f.read_rate.Update(0)
	f.write_rate.Update(0)
	f.transaction_rate.Update(0)
	f.conflict_rate.Update(0)
}

func (f *FoundationDB) Run(step time.Duration) {
	for _ = range time.Tick(step) {
		f.collect()
	}
}
