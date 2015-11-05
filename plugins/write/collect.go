package write

import (
	"os"

	"github.com/customerio/monitor/metrics"
	"github.com/customerio/monitor/plugins"
)

func (z *Write) Collect(b *metrics.Batch) {
	z.collect()
	for _, u := range z.updaters {
		u.Fill(b)
	}
}

func (z *Write) collect() {
	for i, path := range z.paths {
		status := 1
		file, err := os.Create(path)
		if err != nil {
			plugins.Logger.Printf("write: %s: %v\n", path, err)
			status = 0
		} else {
			file.Close()
			os.Remove(path)
		}
		z.updaters[i].Update(float64(status))
	}
}
