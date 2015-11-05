package metrics

import "github.com/customerio/librato"

type Gauge struct {
	name  string
	value float64
}

func NewGauge(name string) *Gauge {
	g := &Gauge{name, 0}
	return g
}

func (g *Gauge) Update(v float64) {
	g.value = v
}

func (g *Gauge) Fill(b *Batch) {
	b.AddGauge(librato.Gauge{
		Name:  g.name,
		Value: g.value,
	})
}
