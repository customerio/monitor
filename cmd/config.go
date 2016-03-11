package main

import (
	"errors"
	"flag"
	"os"
	"time"

	"code.google.com/p/gcfg"
)

type DurationString struct {
	time.Duration
}

func (d *DurationString) UnmarshalText(text []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(string(text))
	return err
}

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
		Interval      DurationString
		CpuThreshold  float64        `gcfg:"cpu-threshold"`
		CpuSampleRate DurationString `gcfg:"cpu-sample-rate"`
		SlackInterval DurationString `gcfg:"slack-interval"`
		Hostname      string
		Logger        string
	}
	Intervals struct {
		Cpu           DurationString
		Riak          DurationString
		Redis         DurationString
		MySQL         DurationString
		Zookeeper     DurationString
		Write         DurationString
		Elasticsearch DurationString
		System        DurationString
		Disk          DurationString
		Etcd          DurationString
	}
}

func parseConfig() (*Config, error) {
	flag.Parse()

	if *config_file == "" {
		return nil, errors.New("Must specify a configuration file path")
	}

	var cfg Config

	if err := gcfg.ReadFileInto(&cfg, *config_file); err != nil {
		return nil, err
	}

	if cfg.Options.Interval.Duration == 0 {
		cfg.Options.Interval.Duration = time.Second
	}

	if cfg.Options.CpuSampleRate.Duration == 0 {
		cfg.Options.CpuSampleRate.Duration = time.Second
	}

	if cfg.Options.SlackInterval.Duration == 0 {
		cfg.Options.SlackInterval.Duration = time.Second
	}

	if cfg.Options.CpuThreshold == 0 {
		cfg.Options.CpuThreshold = 75
	}

	if cfg.Options.Hostname == "" {
		hostname, err := os.Hostname()
		if err != nil {
			hostname = "unknown"
		}
		cfg.Options.Hostname = hostname
	}

	return &cfg, nil
}
