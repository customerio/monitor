package redis

import (
	"os/exec"
	"strconv"
	"strings"

	"github.com/customerio/monitor/metrics"
	"github.com/customerio/monitor/plugins"
)

var updaterMap = map[string]int{
	"connected_clients":         connectedClientsGauge,
	"used_memory":               usedMemoryGauge,
	"used_memory_peak":          usedMemoryPeakGauge,
	"used_cpu_sys":              usedCpuSysCounter,
	"used_cpu_user":             usedCpuUserCounter,
	"total_commands_processed":  totalCommandsProcessedCounter,
	"instantaneous_ops_per_sec": instantaneousOpsPerSecGauge,
	"aof_rewrite_in_progress":   aofRewriteInProgressGauge,
}

func (f *Redis) Collect(b *metrics.Batch) {
	f.collect()

	for _, u := range f.updaters {
		u.Fill(b)
	}
}
func (f *Redis) collect() {
	defer func() {
		if r := recover(); r != nil {
			plugins.Logger.Printf("panic: Redis: %v\n", r)
			f.clear()
		}
	}()
	usage, err := exec.Command("redis-cli", "info").Output()
	if err != nil {
		panic(err)
	}

	vals := map[string]string{}
	for _, v := range strings.Split(string(usage), "\r\n") {
		index := strings.Index(v, ":")
		if index == -1 {
			continue
		}
		vals[v[:index]] = v[index+1:]
	}

	for k, v := range updaterMap {
		f.updaters[v].Update(parse(vals, k))
	}
}

func parse(vals map[string]string, metric string) float64 {
	v, err := strconv.ParseFloat(vals[metric], 64)
	if err != nil {
		panic(err)
	}
	return v
}
