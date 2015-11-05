package plugins

import (
	"sync"
	"time"

	"github.com/customerio/librato"
	"github.com/customerio/monitor/metrics"
)

type Collector interface {
	Collect(b *metrics.Batch)
}

var collectors []Collector

func AddCollector(c Collector) {
	collectors = append(collectors, c)
}

func Collect(source, email, token string, duration time.Duration) {
	client := librato.NewClient(email, token)
	collect(client, source, collectors)
	for _ = range time.Tick(duration) {
		collect(client, source, collectors)
	}
}

func collect(client *librato.Client, source string, collectors []Collector) {
	start := time.Now()

	b := metrics.NewBatch(source)

	// Do each collection in parallel. As a future improvement we might
	// want to wait for up to a some time interval for collection to
	// complete before publishing events to account for slow plugins.
	var wg sync.WaitGroup
	wg.Add(len(collectors))
	for _, c := range collectors {
		go func(collector Collector) {
			collector.Collect(b)
			wg.Done()
		}(c)
	}
	wg.Wait()

	Logger.Printf("processed %d collectors in %v\n", len(collectors), time.Now().Sub(start))

	start = time.Now()
	err := client.Send(b.Librato)
	if err != nil {
		Logger.Printf("send to librato failed: %v\n", err)
		return
	}

	Logger.Printf("sent %d metrics to librato in %v\n", len(b.Librato.Counters)+len(b.Librato.Gauges), time.Now().Sub(start))
}
