package disk

import (
	"time"

	"github.com/rcrowley/go-metrics"
)

type Disk struct {
	filesystem string
	usage      metrics.GaugeFloat64
}

func New(fs string, usage metrics.GaugeFloat64) *Disk {
	return &Disk{
		filesystem: fs,
		usage:      usage,
	}
}

func (d *Disk) clear() {
	d.usage.Update(0)
}

func (d *Disk) Run(step time.Duration) {
	for _ = range time.Tick(step) {
		d.collect()
	}
}
