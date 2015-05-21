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
	conn  *zk.Conn
	paths []string

	gauges map[string]metrics.GaugeFloat64
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
		gauges: map[string]metrics.GaugeFloat64{},
	}
}

func (z *Zookeeper) Add(path string) {
	z.gauges[path] = plugins.Gauge("zk." + strings.Trim(strings.Replace(path, "/", ".", -1), "."))
}

func (z *Zookeeper) Run(step time.Duration) {
	for _ = range time.Tick(step) {
		z.collect()
	}
}
