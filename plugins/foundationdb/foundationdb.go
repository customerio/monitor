package foundationdb

import "time"

type FoundationDB struct {
	port int

	// instance stats
	diskio    int
	ram_used  float64
	ram_total float64
	traffic   float64
	cpu       int

	// Cluster stats
	read_rate        float64
	write_rate       float64
	transaction_rate float64
	conflict_rate    float64
}

func New(port int) *FoundationDB {
	return &FoundationDB{port: port}
}

func (f *FoundationDB) DiskIO() float64 {
	return float64(f.diskio)
}

func (f *FoundationDB) Traffic() float64 {

	return f.traffic
}

func (f *FoundationDB) CPU() float64 {
	return float64(f.cpu)
}

func (f *FoundationDB) RAM() float64 {
	if f.ram_total == 0.0 {
		return 0
	}
	return f.ram_used / f.ram_total * 100
}

func (f *FoundationDB) ReadRate() float64 {
	return f.read_rate
}

func (f *FoundationDB) WriteRate() float64 {
	return f.write_rate
}

func (f *FoundationDB) TransactionRate() float64 {
	return f.transaction_rate
}

func (f *FoundationDB) ConflictRate() float64 {
	return f.conflict_rate
}

func (f *FoundationDB) clear() {
	f.diskio = 0
	f.ram_used = 0
	f.ram_total = 0
	f.traffic = 0
	f.cpu = 0
	f.read_rate = 0
	f.write_rate = 0
	f.transaction_rate = 0
	f.conflict_rate = 0
}

func (f *FoundationDB) Run(step time.Duration) {
	for _ = range time.Tick(step) {
		f.collect()
	}
}

func (f *FoundationDB) gather(name string) float64 {

	switch name {
	case "diskio":
		return float64(f.diskio)
	case "ram":
	case "traffic":
	case "cpu":
	case "read_rate":
	case "write_rate":
	case "transaction_rate":
	case "conflict_rate":
	}
	return 0
}
