package foundationdb

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"regexp"
	"strconv"
)

// 127.0.0.1:4689         (  1% cpu;   6% machine; 0.000 Gbps;  0% disk IO; 2.4 GB / 11.1 GB RAM  )
// [ip]      [port]         [cpu]   [machine]      [gbps]     [diskio]   [ram_usage] [ram_total]
var detail_regex = regexp.MustCompile(`([0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}):([0-9]+)[ \(]+ ([0-9]+)% cpu; *([0-9]+)% machine; *([0-9\.]+) Gbps; *([0-9]+)% disk IO; *([0-9\.]+) GB +/ *([0-9\.]+)`)
var workload_regex = regexp.MustCompile(` +([a-zA-Z]+) rate +- +([0-9]+) Hz`)

type machine struct {
	ip        string
	port      int
	cpu       int
	machine   int
	gbps      float64
	diskio    int
	ram_used  float64
	ram_total float64
}

func (f *FoundationDB) collect() {
	defer func() {
		if r := recover(); r != nil {
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
			f.read_rate = hz
		case "Write":
			f.write_rate = hz
		case "Transaction":
			f.transaction_rate = hz
		case "Conflict":
			f.conflict_rate = hz
		}
	}

	// Extract machine-specific information
	ips := myIps()
	for _, matched := range detail_regex.FindAllSubmatch(usage, -1) {

		if ips[string(matched[1])] && bytesToInt(matched[2]) == f.port {

			f.cpu = bytesToInt(matched[3])
			f.traffic = bytesToFloat64(matched[5])
			f.diskio = bytesToInt(matched[6])
			f.ram_used = bytesToFloat64(matched[7])
			f.ram_total = bytesToFloat64(matched[8])

			return

		}
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
