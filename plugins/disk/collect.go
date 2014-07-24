package disk
// via: https://github.com/measure/os

import (
    "os/exec"
    "strconv"
    "bufio"
    "regexp"
    "strings"
    "math"
)

func (d *Disk) collect() {

    // Collect load average
    usage, err := exec.Command("df").Output()
    if err != nil {
        panic(err)
    }

    scanner := bufio.NewScanner(strings.NewReader(string(usage)))

    splitter := regexp.MustCompile(" +")

    var disk_available, disk_used float64

    for scanner.Scan() {
        str := scanner.Text()
        if strings.HasPrefix(str, d.filesystem){
            split := splitter.Split(str, -1)
            disk_used, _ = strconv.ParseFloat(split[2], 64)
            disk_available, _ = strconv.ParseFloat(split[3], 64)
            break
        }
    }

    d.usage = disk_used / (disk_used + disk_available) * 100

    if math.IsNaN(d.usage){
        d.usage = 0
    }
}




