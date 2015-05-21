package zookeeper

import (
	"log"
	"sync"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

type Zookeeper struct {
	start sync.Once
	conn  *zk.Conn
	paths []string
	stats map[string]int
}

func New(servers []string) *Zookeeper {

	conn, _, err := zk.Connect(servers, time.Second*5)
	if err != nil {
		log.Println("zookeeper error:", err)
		return nil
	}

	return &Zookeeper{
		conn:  conn,
		paths: []string{},
		stats: make(map[string]int),
	}
}

func (z *Zookeeper) PathCounter(path string) *metric {
	return newMetric(z, path)
}

func (z *Zookeeper) run(step time.Duration) {
	z.start.Do(func() {
		for _ = range time.Tick(step) {
			z.collect()
		}
	})
}
