### WIP

Example:

    package main
    
    import (
    	"github.com/customerio/monitor"
    	"github.com/customerio/monitor/plugins/cpu"
    	"github.com/rcrowley/go-metrics"
    
    	"log"
    	"os"
    	"time"
    )
    
    func main() {
    	c := cpu.New()
    
    	go monitor.Report(c.User(), gauge("cpu.user"), time.Second)
    	go monitor.Report(c.System(), gauge("cpu.system"), time.Second)
    	go monitor.Report(c.Idle(), gauge("cpu.idle"), time.Second)
    
    	metrics.Log(metrics.DefaultRegistry, time.Second, log.New(os.Stdout, "metrics: ", log.Lmicroseconds))
    }
    
    func gauge(name string) metrics.GaugeFloat64 {
    	m := metrics.NewGaugeFloat64()
    	metrics.Register(name, m)
    	return m
    }

