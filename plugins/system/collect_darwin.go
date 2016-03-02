package system

import (
	"bufio"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/customerio/monitor/plugins"
)

var splitter = regexp.MustCompile(" +")

func pullFloat64(str string, index int) float64 {
	f, _ := strconv.ParseFloat(splitter.Split(str, -1)[index], 64)
	return f
}

func (s *System) collect() {
	defer func() {
		if r := recover(); r != nil {
			plugins.Logger.Printf("panic: System: %v\n", r)
			s.clear()
		}
	}()

	// Collect the load average from the uptime command
	uptime, err := exec.Command("uptime").Output()
	if err != nil {
		panic(err)
	}

	load_avg, _ := strconv.ParseFloat(strings.Split(string(uptime), " ")[9], 64)
	s.updaters[loadAvgGauge].Update(load_avg)

	// Now some memory stats
	vmstat, err := exec.Command("vm_stat").Output()
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(strings.NewReader(string(vmstat)))

	var pages_free, pages_active, swap, pages_inactive, pages_speculative, pages_wired float64

	for scanner.Scan() {
		str := scanner.Text()
		if strings.HasPrefix(str, "Pages free") {
			pages_free = pullFloat64(str, 2)
		} else if strings.HasPrefix(str, "Pages active") {
			pages_active = pullFloat64(str, 2)
		} else if strings.HasPrefix(str, "Pages inactive") {
			pages_inactive = pullFloat64(str, 2)
		} else if strings.HasPrefix(str, "Pages speculative") {
			pages_speculative = pullFloat64(str, 2)
		} else if strings.HasPrefix(str, "Pages wired down") {
			pages_wired = pullFloat64(str, 2)
		} else if strings.HasPrefix(str, "Swapouts") {
			swap = pullFloat64(str, 1)
		}
	}

	if (pages_active + pages_free) != 0.0 {
		total := pages_free + pages_active + pages_inactive + pages_speculative + pages_wired
		s.updaters[memUsageGauge].Update(pages_active / total * 100)
	} else {
		s.updaters[memUsageGauge].Update(0)
	}
	s.updaters[swapUsageGauge].Update(swap)

}
