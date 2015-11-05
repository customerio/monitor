package metrics

import "github.com/customerio/librato"

type Counter struct {
	name  string
	value float64
}

func NewCounter(name string) *Counter {
	c := &Counter{name, 0}
	return c
}

func (g *Counter) Update(v float64) {
	g.value = v
}

func (g *Counter) Fill(b *Batch) {
	b.AddCounter(librato.Counter{
		Name:  g.name,
		Value: g.value,
	})
}
