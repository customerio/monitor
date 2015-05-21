package cpu

import "time"

type CPU struct {
	previous      map[string]int
	current       map[string]int
	currentTotal  int
	previousTotal int
	lastUpdate    time.Time
}

func New() *CPU {
	return &CPU{}
}

func (c *CPU) User() float64 {
	return c.rate("user")
}

func (c *CPU) System() float64 {
	return c.rate("system")
}

func (c *CPU) Idle() float64 {
	return c.rate("idle")
}

func (c *CPU) Run(step time.Duration) {
	for _ = range time.Tick(step) {
		c.collect()
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
