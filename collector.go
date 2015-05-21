package monitor

import (
	"time"
)

type Collector interface {
	Report() float64
}

type pair struct {
	c Collector
	r Reporter
}

var pairs []pair

type wrapCollector struct {
	report func() float64
}

func (w *wrapCollector) Report() float64 {
	return w.report()
}

func CollectorFunc(report func() float64) Collector {
	return &wrapCollector{report}
}

func ReportFunc(collector func() float64, r Reporter) {
	Report(CollectorFunc(collector), r)
}

func Report(c Collector, r Reporter) {
	pairs = append(pairs, pair{c, r})
}

func Run(step time.Duration) {
	for _ = range time.Tick(step) {
		for _, p := range pairs {
			p.r.Update(p.c.Report())
		}
	}
}
