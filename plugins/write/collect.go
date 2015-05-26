package write

import (
	"os"

	"github.com/customerio/monitor/plugins"
)

func (z *Write) Collect() {
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
