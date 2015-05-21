package zookeeper

import (
	"log"
	"strings"
	"time"

	"github.com/rcrowley/go-metrics"
	"github.com/samuel/go-zookeeper/zk"
)

type Zookeeper struct {
	registry metrics.Registry
	conn     *zk.Conn
	paths    []string

	stats  map[string]int
	gauges map[string]metrics.GaugeFloat64
}

func gauge(registry metrics.Registry, name string) metrics.GaugeFloat64 {
	m := metrics.NewGaugeFloat64()
	registry.Register(name, m)
	return m
}

func New(registry metrics.Registry, servers []string) *Zookeeper {
	conn, _, err := zk.Connect(servers, time.Second*5)
	if err != nil {
		log.Println("zookeeper error:", err)
		return nil
	}

	return &Zookeeper{
		registry: registry,
		conn:     conn,
		paths:    []string{},
		stats:    make(map[string]int),
		gauges:   make(map[string]metrics.GaugeFloat64),
	}
}

func (z *Zookeeper) Add(path string) {
	z.gauges[path] = gauge(z.registry, "zk."+strings.Trim(strings.Replace(path, "/", ".", -1), "."))
}

func (z *Zookeeper) Run(step time.Duration) {
	for _ = range time.Tick(step) {
		z.collect()
		for k, g := range z.gauges {
			if v, ok := z.stats[k]; ok {
				g.Update(float64(v))
			}
		}
	}
}
