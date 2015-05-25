package metrics

type Gauge struct {
	name  string
	value float64
}

func NewGauge(name string) *Gauge {
	g := &Gauge{name, 0}
	gauges = append(gauges, g)
	return g
}

func (g *Gauge) Update(v float64) {
	g.value = v
}
