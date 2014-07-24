package disk

import (
	"sync"
	"time"
)

type Disk struct {
	start         sync.Once
	filesystem 	string
	usage float64
}

func New(fs string) *Disk {
	return &Disk{filesystem: fs}
}

func (d *Disk) Usage() *metric {
	return newMetric(d, "usage")
}

func (d *Disk) run(step time.Duration) {
	d.start.Do(func() {
		for _ = range time.NewTicker(step).C {
			d.collect()
		}
	})
}

func (d *Disk) rate(name string) float64 {
	if name == "usage" {
		return d.usage
	}
	return 0
}
