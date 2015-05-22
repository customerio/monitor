package main

import (
	"github.com/customerio/monitor/plugins/cpu"
	"github.com/customerio/monitor/plugins/disk"
	"github.com/customerio/monitor/plugins/elasticsearch"
	"github.com/customerio/monitor/plugins/mysql"
	"github.com/customerio/monitor/plugins/riak"
	"github.com/customerio/monitor/plugins/system"
	"github.com/rcrowley/go-metrics"
	"github.com/rcrowley/go-metrics/librato"

	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"code.google.com/p/gcfg"
)

var config_file = flag.String("config", "", "Configuration file path")

type Config struct {
	Services struct {
		Librato string
	}
	Metrics struct {
		Cpu           bool
		System        bool
		Riak          string
		Elasticsearch string
		Disk          string
		MySQL         string
	}
	Options struct {
		Interval string
		Hostname string
	}
}

func setup_librato(interval time.Duration, owner string, token string) {

	host, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Binding to librato: %v (%v)", owner, host)

	go librato.Librato(metrics.DefaultRegistry,
		interval,      // interval
		owner,         // account owner email address
		token,         // Librato API token
		host,          // source
		[]float64{95}, // precentiles to send
		interval,      // time unit
	)
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

	duration, err := time.ParseDuration(cfg.Options.Interval)
	if err != nil {
		panic(err)
	}

	if cfg.Metrics.Cpu {
		c := cpu.New()
		go c.Run(time.Second)
	}

	if cfg.Metrics.System {
		s := system.New()
		go s.Run(time.Second)
	}

	if cfg.Metrics.Disk != "" {
		for i, diskname := range strings.Split(cfg.Metrics.Disk, ",") {
			d := disk.New(i, diskname)
			go d.Run(time.Second)
		}
	}

	if cfg.Metrics.MySQL != "" {
		m := mysql.New(cfg.Metrics.MySQL)
		go m.Run(time.Second)
	}

	if cfg.Metrics.Riak != "" {
		r := riak.New(cfg.Metrics.Riak)
		go r.Run(time.Second)
	}

	if cfg.Metrics.Elasticsearch != "" {
		r := elasticsearch.New(cfg.Metrics.Elasticsearch)
		go r.Run(time.Second)
	}
	if cfg.Services.Librato != "" {
		credentials := strings.Split(cfg.Services.Librato, ":")
		if len(credentials) != 2 {
			panic(fmt.Errorf("Bad librato credentials expected: EMAIL:TOKEN, got: %v", cfg.Services.Librato))
		}
		setup_librato(duration, credentials[0], credentials[1])
	}

	metrics.Log(metrics.DefaultRegistry, duration, log.New(os.Stdout, "metrics: ", log.Lmicroseconds))

}
