package metrics

import "github.com/customerio/librato"

type Updater interface {
	Update(v float64)
}

var gauges []*Gauge
var counters []*Counter

func FillBatch(batch *librato.Batch) {
	for _, g := range gauges {
		batch.Gauges = append(batch.Gauges, librato.Gauge{
			Name:  g.name,
			Value: g.value,
		})
	}
	for _, c := range counters {
		batch.Counters = append(batch.Counters, librato.Counter{
			Name:  c.name,
			Value: c.value,
		})
	}
}
