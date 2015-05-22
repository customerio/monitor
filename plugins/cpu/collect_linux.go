package cpu

// via: https://github.com/measure/os

import (
	"bufio"
	"os"
	"regexp"
	"strconv"
	"time"
)

func (c *CPU) collect() {
	defer func() {
		if r := recover(); r != nil {
			c.clear()
		}
	}()

	file, err := os.Open("/proc/stat")
	if err != nil {
		panic(err)
	}

	c.lastUpdate = time.Now()
	c.previous = c.current

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		f := regexp.MustCompile("\\s+").Split(scanner.Text(), -1)

		is_cpu, err := regexp.MatchString("^cpu\\d*", f[0])
		if err == nil && is_cpu {

			user, _ := strconv.ParseUint(f[1], 10, 64)
			nice, _ := strconv.ParseUint(f[2], 10, 64)
			system, _ := strconv.ParseUint(f[3], 10, 64)
			idle, _ := strconv.ParseUint(f[4], 10, 64)

			c.current = []int{
				userGauge:   int(user),
				niceGauge:   int(nice),
				systemGauge: int(system),
				idleGauge:   int(idle),
			}
			break
		}
	}

	if c.previous == nil {
		c.previous = c.current
	}

	c.previousTotal = c.currentTotal
	c.currentTotal = c.current[userGauge] + c.current[niceGauge] + c.current[systemGauge] + c.current[idleGauge]
}
