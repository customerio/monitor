package foundationdb

import (
	"sync"
	"time"
)

type FoundationDB struct {
	start      sync.Once
    port int

    // instance stats
    diskio int
    ram_used float64
    ram_total float64
    traffic float64
    cpu int

    // Cluster stats
    read_rate float64
    write_rate float64
    transaction_rate float64
    conflict_rate float64
}

func New(port int) *FoundationDB {
	return &FoundationDB{port: port}
}

func (f *FoundationDB) DiskIO() *metric {
    return newMetric(f, "diskio")
}

func (f *FoundationDB) Traffic() *metric {
    return newMetric(f, "traffic")
}

func (f *FoundationDB) CPU() *metric {
    return newMetric(f, "cpu")
}

func (f *FoundationDB) RAM() *metric {
    return newMetric(f, "ram")
}

func (f *FoundationDB) ReadRate() *metric {
    return newMetric(f, "read_rate")
}

func (f *FoundationDB) WriteRate() *metric {
    return newMetric(f, "write_rate")
}

func (f *FoundationDB) TransactionRate() *metric {
    return newMetric(f, "transaction_rate")
}

func (f *FoundationDB) ConflictRate() *metric {
    return newMetric(f, "conflict_rate")
}

func (f *FoundationDB) run(step time.Duration) {
	f.start.Do(func() {
		for _ = range time.NewTicker(step).C {
			f.collect()
		}
	})
}

func (f *FoundationDB) gather(name string) float64 {

    switch name {
        case "diskio": return float64(f.diskio)
        case "ram": return f.ram_used / f.ram_total * 100
        case "traffic": return f.traffic
        case "cpu": return float64(f.cpu)
        case "read_rate": return f.read_rate
        case "write_rate": return f.write_rate
        case "transaction_rate": return f.transaction_rate
        case "conflict_rate": return f.conflict_rate
    }
    return 0
}
