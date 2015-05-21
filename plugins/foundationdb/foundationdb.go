package foundationdb

import (
	"time"

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

func gauge(registry metrics.Registry, name string) metrics.GaugeFloat64 {
	m := metrics.NewGaugeFloat64()
	registry.Register(name, m)
	return m
}

func New(registry metrics.Registry, port int) *FoundationDB {
	f := &FoundationDB{port: port}

	f.diskio = gauge(registry, "fdb.diskio")
	f.traffic = gauge(registry, "fdb.traffic")
	f.cpu = gauge(registry, "fdb.cpu")
	f.ram = gauge(registry, "fdb.ram")
	f.read_rate = gauge(registry, "fdb.rate.read")
	f.write_rate = gauge(registry, "fdb.rate.write")
	f.transaction_rate = gauge(registry, "fdb.rate.transaction")
	f.conflict_rate = gauge(registry, "fdb.rate.conflict")

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
