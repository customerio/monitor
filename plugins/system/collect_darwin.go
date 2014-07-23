package system

import (
    "os/exec"
    "strconv"
    "strings"
    "bufio"
    "regexp"
)

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

    splitter := regexp.MustCompile(" +")

    var pages_free, pages_active float64

    for scanner.Scan() {
        str := scanner.Text()
        if strings.HasPrefix(str, "Pages free"){
            pages_free, _ = strconv.ParseFloat(splitter.Split(str, -1)[2], 64)
        }else if strings.HasPrefix(str, "Pages active"){
            pages_active, _ = strconv.ParseFloat(splitter.Split(str, -1)[2], 64)
        }
    }

    s.memUsage = pages_active / (pages_active + pages_free) * 100

}


