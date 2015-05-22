package cpu

import (
	"time"

	"github.com/customerio/monitor/plugins"
	"github.com/rcrowley/go-metrics"
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

	gauges []metrics.GaugeFloat64
}

func New() *CPU {
	return &CPU{
		gauges: []metrics.GaugeFloat64{
			userGauge:   plugins.Gauge("cpu.user"),
			systemGauge: plugins.Gauge("cpu.system"),
			idleGauge:   plugins.Gauge("cpu.idle"),
		},
	}
}

func (c *CPU) Run(step time.Duration) {
	for _ = range time.Tick(step) {
		c.collect()

		for _, i := range []int{userGauge, systemGauge, idleGauge} {
			c.gauges[i].Update(c.rate(i))
		}
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
