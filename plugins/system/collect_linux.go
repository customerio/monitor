package system

import (
    "strconv"
    "io/ioutil"
    "strings"
)

func (s *System) collect() {

    // Collect load average
    data, err := ioutil.ReadFile("/proc/loadavg")
    if err != nil {
        panic(err)
    }
    load_avg, _ := strconv.ParseFloat(strings.split(data, " ")[0], 64)
    s.loadAvg = loadavg

    // Now some memory stats
    data, err := ioutil.ReadFile("/proc/meminfo")
    if err != nil {
        panic(err)
    }

    scanner := bufio.NewScanner(strings.NewReader(string(vmstat)))

    splitter := regexp.MustCompile(" +")

    var mem_total, mem_free float64

    for scanner.Scan() {
        str := scanner.Text()
        if strings.HasPrefix(str, "MemTotal"){
            mem_total, _ = strconv.ParseFloat(splitter.Split(str, -1)[1], 64)
        }else if strings.HasPrefix(str, "MemFree"){
            mem_free, _ = strconv.ParseFloat(splitter.Split(str, -1)[1], 64)
        }
    }

    s.memUsage = mem_free / mem_total * 100
}


