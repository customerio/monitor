package zookeeper

import (
	"log"
	"strings"
	"time"

	"github.com/customerio/monitor/metrics"
	"github.com/samuel/go-zookeeper/zk"
)

type Zookeeper struct {
	conn *zk.Conn

	// Parallel arrays.
	paths    []string
	updaters []metrics.Updater
}

func New(servers []string) *Zookeeper {
	conn, _, err := zk.Connect(servers, time.Second*5)
	if err != nil {
		log.Println("zookeeper error:", err)
		return nil
	}

	return &Zookeeper{
		conn:     conn,
		paths:    nil,
		updaters: nil,
	}
}

func (z *Zookeeper) Add(path string) {
	z.paths = append(z.paths, path)
	z.updaters = append(z.updaters, metrics.NewGauge("zk."+strings.Trim(strings.Replace(path, "/", ".", -1), ".")))
}
