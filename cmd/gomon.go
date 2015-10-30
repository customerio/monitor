package main

import (
	"net/http"
	"os"

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

	"flag"
	"fmt"
	"strings"
	"time"

	"code.google.com/p/gcfg"
)

var config_file = flag.String("config", "", "Configuration file path")

type Config struct {
	Services struct {
		Librato string
		Slack   string
	}
	Metrics struct {
		Cpu           bool
		Riak          string
		Redis         string
		MySQL         string
		Zookeeper     []string
		Write         []string
		Elasticsearch string
		System        bool
		Disk          string
		Etcd          string
	}
	Options struct {
		Interval string
		Hostname string
		Logger   string
	}
}

func main() {
	flag.Parse()

	var cfg Config

	if *config_file != "" {
		err := gcfg.ReadFileInto(&cfg, *config_file)
		if err != nil {
			panic(err)
		}
	} else {
		panic(fmt.Errorf("Must specify a configuration file path"))
	}

	var duration time.Duration

	// Default value for interval is 1s
	if cfg.Options.Interval == "" {
		cfg.Options.Interval = "1s"
	}

	// We don't want long waits on http connections.
	http.DefaultClient.Timeout = time.Second * 5

	plugins.InitializeLogger(cfg.Options.Logger, "gomon")

	duration, err := time.ParseDuration(cfg.Options.Interval)
	if err != nil {
		panic(err)
	}

	if cfg.Metrics.Cpu {
		c := cpu.New()
		plugins.AddCollector(c)
	}

	if cfg.Metrics.Redis != "" {
		c := redis.New()
		plugins.AddCollector(c)
	}

	if cfg.Metrics.System {
		s := system.New()
		plugins.AddCollector(s)
	}

	if cfg.Metrics.Disk != "" {
		for i, diskname := range strings.Split(cfg.Metrics.Disk, ",") {
			d := disk.New(i, diskname)
			plugins.AddCollector(d)
		}
	}

	if cfg.Metrics.MySQL != "" {
		m := mysql.New(cfg.Metrics.MySQL)
		plugins.AddCollector(m)
	}

	if cfg.Metrics.Etcd != "" {
		m := etcd.New(cfg.Services.Slack, cfg.Metrics.Etcd)
		plugins.AddCollector(m)
	}

	if cfg.Metrics.Riak != "" {
		r := riak.New(cfg.Metrics.Riak)
		plugins.AddCollector(r)
	}

	if cfg.Metrics.Elasticsearch != "" {
		r := elasticsearch.New(cfg.Metrics.Elasticsearch)
		plugins.AddCollector(r)
	}

	if len(cfg.Metrics.Zookeeper) > 0 {
		z := zookeeper.New([]string{"localhost"})
		for _, m := range cfg.Metrics.Zookeeper {
			z.Add(m)
		}
		plugins.AddCollector(z)
	}

	if len(cfg.Metrics.Write) > 0 {
		z := write.New()
		for _, m := range cfg.Metrics.Write {
			z.Add(m)
		}
		plugins.AddCollector(z)
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

	plugins.Collect(host, email, token, duration)
}
