package cpu

import (
	"sync"
	"time"
)

type CPU struct {
	start         sync.Once
	previous      map[string]int
	current       map[string]int
	currentTotal  int
	previousTotal int
	lastUpdate    time.Time
}

func New() *CPU {
	return &CPU{}
}

func (c *CPU) User() *metric {
	return newMetric(c, "user")
}

func (c *CPU) System() *metric {
	return newMetric(c, "system")
}

func (c *CPU) Idle() *metric {
	return newMetric(c, "idle")
}

func (c *CPU) run(step time.Duration) {
	c.start.Do(func() {
		for _ = range time.NewTicker(step).C {
			c.collect()
		}
	})
}

func (c *CPU) rate(name string) float64 {
	delta := c.current[name] - c.previous[name]
	total := c.currentTotal - c.previousTotal

	return float64(delta) / float64(total) * 100
}
