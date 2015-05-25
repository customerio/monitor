package redis

import "github.com/customerio/monitor/metrics"

const (
	connectedClientsGauge = iota
	usedMemoryGauge
	usedMemoryPeakGauge
	usedCpuSysCounter
	usedCpuUserCounter
	totalCommandsProcessedCounter
	instantaneousOpsPerSecGauge
	aofRewriteInProgressGauge
)

type Redis struct {
	updaters []metrics.Updater
}

func New() *Redis {
	return &Redis{
		updaters: []metrics.Updater{
			connectedClientsGauge:         metrics.NewGauge("redis.connected_clients"),
			usedMemoryGauge:               metrics.NewGauge("redis.used_memory"),
			usedMemoryPeakGauge:           metrics.NewGauge("redis.used_memory_peak"),
			usedCpuSysCounter:             metrics.NewCounter("redis.used_cpu_sys"),
			usedCpuUserCounter:            metrics.NewCounter("redis.used_cpu_user"),
			totalCommandsProcessedCounter: metrics.NewCounter("redis.total_commands_processed"),
			instantaneousOpsPerSecGauge:   metrics.NewGauge("redis.instantaneous_ops_per_sec"),
			aofRewriteInProgressGauge:     metrics.NewGauge("redis.aof_rewrite_in_progress"),
		},
	}
}

func (f *Redis) clear() {
	for _, v := range f.updaters {
		v.Update(0)
	}
}
