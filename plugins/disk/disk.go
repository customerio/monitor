package disk

import (
	"fmt"

	"github.com/customerio/monitor/metrics"
)

type Disk struct {
	filesystem string
	usage      metrics.Updater
}

func New(i int, fs string) *Disk {
	return &Disk{
		filesystem: fs,
		usage:      metrics.NewGauge(fmt.Sprintf("disk.%d.usage", i)),
	}
}

func (d *Disk) clear() {
	d.usage.Update(0)
}
