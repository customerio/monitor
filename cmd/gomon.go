package main

import (
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/customerio/monitor/plugins"
	"github.com/customerio/monitor/plugins/cpu"
	"github.com/customerio/monitor/plugins/disk"
	"github.com/customerio/monitor/plugins/elasticsearch"
	"github.com/customerio/monitor/plugins/etcd"
	"github.com/customerio/monitor/plugins/mysql"
	"github.com/customerio/monitor/plugins/redis"
	"github.com/customerio/monitor/plugins/riak"
	"github.com/customerio/monitor/plugins/system"
	"github.com/customerio/monitor/plugins/write"
	"github.com/customerio/monitor/plugins/zookeeper"

	"fmt"
	"strings"
	"time"
)

func main() {

	stop := make(chan struct{})
	s := make(chan os.Signal, 1)
	signal.Notify(s, os.Interrupt, syscall.SIGTERM, os.Kill)

	cfg, err := parseConfig()
	if err != nil {
		panic(err)
	}

	// We don't want long waits on http connections.
	http.DefaultClient.Timeout = time.Second * 5

	plugins.InitializeLogger(cfg.Options.Logger, "gomon")

	if cfg.Metrics.Cpu {
		d := cfg.Intervals.Cpu
		if d.Duration == 0 {
			d = cfg.Options.Interval
		}
		c := cpu.New(&cpu.Config{
			Threshold:         cfg.Options.CpuThreshold,
			SampleRate:        cfg.Options.CpuSampleRate.Duration,
			ReportingInterval: d.Duration,
			SlackURL:          cfg.Services.Slack,
			SlackInterval:     cfg.Options.SlackInterval.Duration,
			Hostname:          cfg.Options.Hostname,
			SkipNotification:  cfg.SlackDisabled.Cpu,
		})
		plugins.AddCollector(c, cfg.Intervals.Cpu.Duration)
	}

	if cfg.Metrics.Redis != "" {
		c := redis.New()
		plugins.AddCollector(c, cfg.Intervals.Redis.Duration)
	}

	if cfg.Metrics.System {
		s := system.New()
		plugins.AddCollector(s, cfg.Intervals.System.Duration)
	}

	if cfg.Metrics.Disk != "" {
		for i, diskname := range strings.Split(cfg.Metrics.Disk, ",") {
			d := disk.New(i, diskname)
			plugins.AddCollector(d, cfg.Intervals.Disk.Duration)
		}
	}

	if cfg.Metrics.MySQL != "" {
		m := mysql.New(cfg.Metrics.MySQL)
		plugins.AddCollector(m, cfg.Intervals.MySQL.Duration)
	}

	if cfg.Metrics.Etcd != "" {
		m := etcd.New(cfg.Services.Slack, cfg.Metrics.Etcd, cfg.Options.Hostname, cfg.SlackDisabled.Etcd)
		plugins.AddCollector(m, cfg.Intervals.Etcd.Duration)
	}

	if cfg.Metrics.Riak != "" {
		r := riak.New(cfg.Metrics.Riak)
		plugins.AddCollector(r, cfg.Intervals.Riak.Duration)
	}

	if cfg.Metrics.Elasticsearch != "" {
		r := elasticsearch.New(cfg.Metrics.Elasticsearch)
		plugins.AddCollector(r, cfg.Intervals.Elasticsearch.Duration)
	}

	if len(cfg.Metrics.Zookeeper) > 0 {
		z := zookeeper.New([]string{"localhost"})
		for _, m := range cfg.Metrics.Zookeeper {
			z.Add(m)
		}
		plugins.AddCollector(z, cfg.Intervals.Zookeeper.Duration)
	}

	if len(cfg.Metrics.Write) > 0 {
		z := write.New()
		for _, m := range cfg.Metrics.Write {
			z.Add(m)
		}
		plugins.AddCollector(z, cfg.Intervals.Write.Duration)
	}

	var email, token string
	if cfg.Services.Librato != "" {
		credentials := strings.Split(cfg.Services.Librato, ":")
		if len(credentials) != 2 {
			panic(fmt.Errorf("Bad librato credentials expected: EMAIL:TOKEN, got: %v", cfg.Services.Librato))
		}
		email = credentials[0]
		token = credentials[1]
	}
	host, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	plugins.Collect(host, email, token, cfg.Options.Interval.Duration, &wg, stop)

	<-s
	close(stop)
	wg.Wait()
}
