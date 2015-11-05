package foundationdb

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"regexp"
	"strconv"

	"github.com/customerio/monitor/metrics"
)

// 127.0.0.1:4689         (  1% cpu;   6% machine; 0.000 Gbps;  0% disk IO; 2.4 GB / 11.1 GB RAM  )
// [ip]      [port]         [cpu]   [machine]      [gbps]     [diskio]   [ram_usage] [ram_total]
var detail_regex = regexp.MustCompile(`([0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}):([0-9]+)[ \(]+ ([0-9]+)% cpu; *([0-9]+)% machine; *([0-9\.]+) Gbps; *([0-9]+)% disk IO; *([0-9\.]+) GB +/ *([0-9\.]+)`)
var workload_regex = regexp.MustCompile(` +([a-zA-Z]+) rate +- +([0-9]+) Hz`)

func (f *FoundationDB) Collect(b *metrics.Batch) {
	f.collect()

	for _, u := range f.updaters {
		u.Fill(b)
	}
}
func (f *FoundationDB) collect() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("panic: FoundationDB: %v\n", r)
			f.clear()
		}
	}()

	usage, err := exec.Command("fdbcli", "--exec", "status details").Output()
	if err != nil {
		panic(err)
	}

	// Extract cluster information
	for _, matched := range workload_regex.FindAllSubmatch(usage, -1) {
		hz := bytesToFloat64(matched[2])
		switch string(matched[1]) {
		case "Read":
			f.updaters[readRateGauge].Update(hz)
		case "Write":
			f.updaters[writeRateGauge].Update(hz)
		case "Transaction":
			f.updaters[transactionRateGauge].Update(hz)
		case "Conflict":
			f.updaters[conflictRateGauge].Update(hz)
		}
	}

	var ram_used, ram_total float64
	// Extract machine-specific information
	ips := myIps()
	for _, matched := range detail_regex.FindAllSubmatch(usage, -1) {

		if ips[string(matched[1])] && bytesToInt(matched[2]) == f.port {

			f.updaters[cpuGauge].Update(float64(bytesToInt(matched[3])))
			f.updaters[trafficGauge].Update(bytesToFloat64(matched[5]))
			f.updaters[diskioGauge].Update(float64(bytesToInt(matched[6])))
			ram_used = bytesToFloat64(matched[7])
			ram_total = bytesToFloat64(matched[8])

			break

		}
	}

	if ram_total == 0.0 {
		f.updaters[ramGauge].Update(0)
	} else {
		f.updaters[ramGauge].Update(ram_used / ram_total * 100)
	}
}

func myIps() map[string]bool {

	// Defaults to just localhost
	output := map[string]bool{"127.0.0.1": true}

	name, err := os.Hostname()
	if err != nil {
		fmt.Printf("Oops: %v\n", err)
		return output
	}

	addrs, err := net.LookupHost(name)
	if err != nil {
		fmt.Printf("Oops: %v\n", err)
		return output
	}

	for _, addr := range addrs {
		output[addr] = true
	}

	return output
}

func bytesToInt(b []byte) int {
	s := string(b)
	i, _ := strconv.Atoi(s)
	return i
}

func bytesToFloat64(b []byte) float64 {
	s := string(b)
	f, _ := strconv.ParseFloat(s, 64)
	return f
}
