package zookeeper

import (
	"log"
	"strings"
	"time"

	"github.com/customerio/monitor/plugins"
	"github.com/rcrowley/go-metrics"
	"github.com/samuel/go-zookeeper/zk"
)

type Zookeeper struct {
	conn *zk.Conn

	// Parallel arrays.
	paths  []string
	gauges []metrics.GaugeFloat64
}

func New(servers []string) *Zookeeper {
	conn, _, err := zk.Connect(servers, time.Second*5)
	if err != nil {
		log.Println("zookeeper error:", err)
		return nil
	}

	return &Zookeeper{
		conn:   conn,
		paths:  nil,
		gauges: nil,
	}
}

func (z *Zookeeper) Add(path string) {
	z.paths = append(z.paths, path)
	z.gauges = append(z.gauges, plugins.Gauge("zk."+strings.Trim(strings.Replace(path, "/", ".", -1), ".")))
}

func (z *Zookeeper) Run(step time.Duration) {
	for _ = range time.Tick(step) {
		z.collect()
	}
}
