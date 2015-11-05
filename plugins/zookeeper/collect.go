package zookeeper

import (
	"github.com/customerio/monitor/metrics"
	"github.com/customerio/monitor/plugins"
)

func (z *Zookeeper) Collect(b *metrics.Batch) {
	z.collect()
	for _, u := range z.updaters {
		u.Fill(b)
	}
}

func (z *Zookeeper) collect() {
	for i, path := range z.paths {
		if children, _, err := z.conn.Children(path); err == nil {
			z.updaters[i].Update(float64(len(children)))
		} else {
			plugins.Logger.Printf("panic: Zookeeper: %v\n", err)
			z.updaters[i].Update(0)
		}
	}
}
