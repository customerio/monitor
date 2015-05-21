package zookeeper

import (
	"log"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

type Zookeeper struct {
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

func (z *Zookeeper) PathCounter(path string) func() float64 {
	return func() float64 {
		return float64(z.stats[path])
	}
}

func (z *Zookeeper) Run(step time.Duration) {
	for _ = range time.Tick(step) {
		z.collect()
	}
}
