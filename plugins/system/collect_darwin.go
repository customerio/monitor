package system

import (
	"bufio"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

var splitter = regexp.MustCompile(" +")

func pullFloat64(str string, index int) float64 {
	f, _ := strconv.ParseFloat(splitter.Split(str, -1)[index], 64)
	return f
}

func (s *System) collect() {

	// Collect the load average from the uptime command
	uptime, err := exec.Command("uptime").Output()
	if err != nil {
		panic(err)
	}

	load_avg, _ := strconv.ParseFloat(strings.Split(string(uptime), " ")[9], 64)
	s.loadAvg = load_avg

	// Now some memory stats
	vmstat, err := exec.Command("vm_stat").Output()
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(strings.NewReader(string(vmstat)))

	var pages_free, pages_active, swap float64

	for scanner.Scan() {
		str := scanner.Text()
		if strings.HasPrefix(str, "Pages free") {
			pages_free = pullFloat64(str, 2)
		} else if strings.HasPrefix(str, "Pages active") {
			pages_active = pullFloat64(str, 2)
		} else if strings.HasPrefix(str, "Swapouts") {
			swap = pullFloat64(str, 1)
		}
	}

	s.memUsage = pages_active / (pages_active + pages_free) * 100
	s.swapUsage = swap

}
