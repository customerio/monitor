package monitor

import (
	"time"
)

type Collector interface {
	Run(time.Duration)
	Report() float64
}

func Report(c Collector, r Reporter, step time.Duration) {
	go c.Run(step)

	for _ = range time.NewTicker(step).C {
		r.Update(c.Report())
	}
}
