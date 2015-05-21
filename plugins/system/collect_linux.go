package system

import (
	"bufio"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

var splitter = regexp.MustCompile(" +")

func pullFloat64(str string) float64 {
	f, _ := strconv.ParseFloat(splitter.Split(str, -1)[1], 64)
	return f
}

func (s *System) collect() {
	defer func() {
		if r := recover(); r != nil {
			s.clear()
		}
	}()

	// Collect load average
	data, err := ioutil.ReadFile("/proc/loadavg")
	if err != nil {
		panic(err)
	}
	load_avg, _ := strconv.ParseFloat(strings.Split(string(data), " ")[0], 64)
	s.loadAvg = load_avg

	// Now some memory stats
	meminfo, err := ioutil.ReadFile("/proc/meminfo")
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(strings.NewReader(string(meminfo)))

	var mem_total, mem_free, swap_free, swap_total float64

	for scanner.Scan() {
		str := scanner.Text()
		if strings.HasPrefix(str, "MemTotal") {
			mem_total = pullFloat64(str)
		} else if strings.HasPrefix(str, "MemFree") {
			mem_free = pullFloat64(str)
		} else if strings.HasPrefix(str, "SwapFree") {
			swap_free = pullFloat64(str)
		} else if strings.HasPrefix(str, "SwapTotal") {
			swap_total = pullFloat64(str)
		}
	}

	if mem_total != 0.0 {
		s.memUsage = mem_free / mem_total * 100
	} else {
		s.memUsage = 0
	}
	if swap_total != 0.0 {
		s.swapUsage = (swap_total - swap_free) / swap_total * 100
	} else {
		s.swapUsage = 0
	}
}
