package cpu
// via: https://github.com/measure/os

import (
    "time"
    "unsafe"
    "strconv"
)

func (c *CPU) collect() {
    c.lastUpdate = time.Now()
    c.previous = c.current

    file, err := os.Open("/proc/stat")
    if err != nil {
        return
    }

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        f := regexp.MustCompile("\\s+").Split(scanner.Text(), -1)

        is_cpu, err := regexp.MatchString("^cpu\\d*", f[0])
        if err == nil && is_cpu {

            user, _ := strconv.ParseUint(f[1], 10, 64)
            nice, _ := strconv.ParseUint(f[2], 10, 64)
            system, _ := strconv.ParseUint(f[3], 10, 64)
            idle, _ := strconv.ParseUint(f[4], 10, 64)

            c.current = map[string]int{
                "user":   int(user),
                "nice":   int(nice),
                "system": int(system),
                "idle":   int(idle),
            }
            break
        }
    }

    if c.previous == nil {
        c.previous = c.current
    }

    c.previousTotal = c.currentTotal
    c.currentTotal = c.current["user"] + c.current["nice"] + c.current["system"] + c.current["idle"]
}


