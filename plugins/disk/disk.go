package disk

import "time"

type Disk struct {
	filesystem string
	usage      float64
}

func New(fs string) *Disk {
	return &Disk{filesystem: fs}
}

func (d *Disk) Usage() float64 {
	return d.usage
}

func (d *Disk) clear() {
	d.usage = 0
}

func (d *Disk) Run(step time.Duration) {
	for _ = range time.Tick(step) {
		d.collect()
	}
}
