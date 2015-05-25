package metrics

import (
	"bytes"
	"fmt"

	"github.com/customerio/librato"
)

type Updater interface {
	Update(v float64)
}

var gauges []*Gauge
var counters []*Counter

func Summarize() string {
	b := &bytes.Buffer{}
	fmt.Fprintf(b, "gauges: [")
	for i, g := range gauges {
		if i != 0 {
			fmt.Fprintf(b, ", ")
		}
		fmt.Fprintf(b, "%s", g.name)
	}
	fmt.Fprintf(b, "]\n")
	fmt.Fprintf(b, "counters: [")
	for i, c := range counters {
		if i != 0 {
			fmt.Fprintf(b, ", ")
		}
		fmt.Fprintf(b, "%s", c.name)
	}
	fmt.Fprintf(b, "]")
	return b.String()
}

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
