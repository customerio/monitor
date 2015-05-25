package cpu

import (
	"time"

	"github.com/customerio/monitor/metrics"
)

const (
	userGauge = iota
	systemGauge
	idleGauge
	niceGauge
)

type CPU struct {
	previous      []int
	current       []int
	currentTotal  int
	previousTotal int
	lastUpdate    time.Time

	updaters []metrics.Updater
}

func New() *CPU {
	return &CPU{
		updaters: []metrics.Updater{
			userGauge:   metrics.NewGauge("cpu.user"),
			systemGauge: metrics.NewGauge("cpu.system"),
			idleGauge:   metrics.NewGauge("cpu.idle"),
		},
	}
}

func (c *CPU) Collect() {
	c.collect()
	for _, i := range []int{userGauge, systemGauge, idleGauge} {
		c.updaters[i].Update(c.rate(i))
	}
}

func (c *CPU) clear() {
	c.current = nil
	c.previous = nil
	c.currentTotal = 0
	c.previousTotal = 0
}

func (c *CPU) rate(name int) float64 {
	if name >= len(c.current) {
		return 0
	}
	delta := c.current[name] - c.previous[name]
	total := c.currentTotal - c.previousTotal

	if total == 0.0 {
		return 0
	}
	return float64(delta) / float64(total) * 100
}
