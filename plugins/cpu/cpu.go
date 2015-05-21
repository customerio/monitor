package cpu

import (
	"time"

	"github.com/rcrowley/go-metrics"
)

type CPU struct {
	previous      map[string]int
	current       map[string]int
	currentTotal  int
	previousTotal int
	lastUpdate    time.Time

	user   metrics.GaugeFloat64
	system metrics.GaugeFloat64
	idle   metrics.GaugeFloat64
}

func gauge(registry metrics.Registry, name string) metrics.GaugeFloat64 {
	m := metrics.NewGaugeFloat64()
	registry.Register(name, m)
	return m
}

func New(registry metrics.Registry) *CPU {
	c := &CPU{}
	c.user = gauge(registry, "cpu.user")
	c.system = gauge(registry, "cpu.system")
	c.idle = gauge(registry, "cpu.idle")
	return c
}

func (c *CPU) Run(step time.Duration) {
	for _ = range time.Tick(step) {
		c.collect()

		c.user.Update(c.rate("user"))
		c.system.Update(c.rate("system"))
		c.idle.Update(c.rate("idle"))
	}
}

func (c *CPU) clear() {
	c.current = map[string]int{}
	c.previous = map[string]int{}
	c.currentTotal = 0
	c.previousTotal = 0
}

func (c *CPU) rate(name string) float64 {
	delta := c.current[name] - c.previous[name]
	total := c.currentTotal - c.previousTotal

	if total == 0.0 {
		return 0
	}
	return float64(delta) / float64(total) * 100
}
