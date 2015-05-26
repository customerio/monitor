package write

import (
	"strings"

	"github.com/customerio/monitor/metrics"
	"github.com/customerio/monitor/plugins"
)

type Write struct {
	// Parallel arrays.
	paths    []string
	updaters []metrics.Updater
}

func New() *Write {
	return &Write{
		paths:    nil,
		updaters: nil,
	}
}

func (z *Write) Add(path string) {
	index := strings.Index(path, ":")
	if index == -1 {
		plugins.Logger.Printf("write: invalid: '%s'\n", path)
		return
	}
	z.updaters = append(z.updaters, metrics.NewGauge("write."+path[0:index]))
	z.paths = append(z.paths, path[index+1:])
}
