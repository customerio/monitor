package disk

import (
	"fmt"
	"time"

	"github.com/customerio/monitor/plugins"
	"github.com/rcrowley/go-metrics"
)

type Disk struct {
	filesystem string
	usage      metrics.GaugeFloat64
}

func New(i int, fs string) *Disk {
	return &Disk{
		filesystem: fs,
		usage:      plugins.Gauge(fmt.Sprintf("disk.%d.usage", i)),
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
