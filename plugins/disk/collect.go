package disk

// via: https://github.com/measure/os

import (
	"bufio"
	"math"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type mount struct {
	name      string
	used      float64
	available float64
}

func (m mount) size() float64 {
	return m.used + m.available
}

func (m mount) usage() float64 {
	out := m.used / m.size() * 100
	if math.IsNaN(out) {
		return 0
	}
	return out
}

type mounts []mount

func (m mounts) Len() int           { return len(m) }
func (m mounts) Swap(i, j int)      { m[i], m[j] = m[j], m[i] }
func (m mounts) Less(i, j int) bool { return m[i].size() < m[j].size() }

func (m mounts) get(name string) *mount {

	if len(m) > 0 {

		if name == "largest" {
			sort.Sort(sort.Reverse(m))
			return &m[0]
		}
		for _, mount := range m {
			if mount.name == name {
				return &mount
			}
		}

	}
	return nil

}

func (d *Disk) collect() {

	// Collect load average
	usage, err := exec.Command("df").Output()
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(strings.NewReader(string(usage)))

	splitter := regexp.MustCompile(" +")

	var disk_available, disk_used float64

	all_disks := mounts{}

	for scanner.Scan() {
		split := splitter.Split(scanner.Text(), -1)
		disk_used, _ = strconv.ParseFloat(split[2], 64)
		disk_available, _ = strconv.ParseFloat(split[3], 64)
		all_disks = append(all_disks, mount{split[0], disk_used, disk_available})
	}

	drive := all_disks.get(d.filesystem)
	if drive != nil {
		d.usage = drive.usage()
	}
}
