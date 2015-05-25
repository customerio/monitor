package metrics

type Counter struct {
	name  string
	value float64
}

func NewCounter(name string) *Counter {
	c := &Counter{name, 0}
	counters = append(counters, c)
	return c
}

func (g *Counter) Update(v float64) {
	g.value = v
}
