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

// Group collectors by their intervals to minimize number of long running
// goroutines spawned.
var collectors = make(map[time.Duration][]Collector)

func AddCollector(c Collector, d time.Duration) {
	collectors[d] = append(collectors[d], c)
}

func Collect(source, email, token string, defaultDuration time.Duration, wg *sync.WaitGroup, stopSignal <-chan struct{}) {
	var client = librato.NewClient(email, token)

	wg.Add(len(collectors))
	// We spawn long running goroutines for all different intervals that will
	// collect stats from the associated collectors. Goroutines exit cleanly
	// when kill signal is sent.
	for duration, plugins := range collectors {
		if duration == 0 {
			duration = defaultDuration
		}
		go doCollection(duration, plugins, client, source, wg, stopSignal)
	}
}

func doCollection(d time.Duration, collectors []Collector, client *librato.Client, source string, wg *sync.WaitGroup, stopSignal <-chan struct{}) {
	var (
		ticker = time.NewTicker(d)
		stats  = metrics.NewBatch(source)
	)
loop:
	for {
		select {
		case <-stopSignal:
			break loop
		case <-ticker.C:
			gatherAndPublish(client, stats, collectors)
			stats.Reset(source)
		}
	}
	ticker.Stop()
	wg.Done()
}

func gatherAndPublish(client *librato.Client, b *metrics.Batch, collectors []Collector) {
	start := time.Now()

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
