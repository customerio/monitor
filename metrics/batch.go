package metrics

import (
	"sync"

	"github.com/customerio/librato"
)

type Updater interface {
	Update(v float64)
	Fill(b *Batch)
}

type Batch struct {
	mutex   sync.Mutex
	Librato *librato.Batch
}

func NewBatch(source string) *Batch {
	b := &Batch{
		Librato: librato.NewBatch(),
	}
	b.Librato.Source = source
	return b
}

func (b *Batch) AddGauge(g librato.Gauge) {
	b.mutex.Lock()
	b.Librato.Gauges = append(b.Librato.Gauges, g)
	b.mutex.Unlock()
}
func (b *Batch) AddCounter(g librato.Counter) {
	b.mutex.Lock()
	b.Librato.Counters = append(b.Librato.Counters, g)
	b.mutex.Unlock()
}
