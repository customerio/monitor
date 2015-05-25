package plugins

import (
	"sync"
	"time"

	"github.com/customerio/librato"
	"github.com/customerio/monitor/metrics"
)

type Collector interface {
	Collect()
}

var collectors []Collector

func AddCollector(c Collector) {
	collectors = append(collectors, c)
}

func Collect(source, email, token string, duration time.Duration) {
	Logger.Printf("Reporting the following metrics to librato\n%s\n", metrics.Summarize())

	client := librato.NewClient(email, token)
	collect(client, source, collectors)
	for _ = range time.Tick(duration) {
		collect(client, source, collectors)
	}
}

var wg sync.WaitGroup

func collect(client *librato.Client, source string, collectors []Collector) {
	b := librato.NewBatch()
	b.Source = source

	start := time.Now()

	// Do each collection in parallel. As a future improvement we might
	// want to wait for up to a some time interval for collection to
	// complete before publishing events to account for slow plugins.
	wg.Add(len(collectors))
	for _, c := range collectors {
		go func(collector Collector) {
			collector.Collect()
			wg.Done()
		}(c)
	}
	wg.Wait()

	Logger.Printf("processed %d collectors in %v\n", len(collectors), time.Now().Sub(start))
	metrics.FillBatch(b)

	start = time.Now()
	err := client.Send(b)
	if err != nil {
		Logger.Printf("send to librato failed: %v\n", err)
	} else {
		Logger.Printf("sent %d metrics to librato in %v\n", len(b.Counters)+len(b.Gauges), time.Now().Sub(start))
	}
}
