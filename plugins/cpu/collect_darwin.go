package cpu

import (
	"fmt"
	"time"
	"unsafe"
)

/*
#include <mach/mach_init.h>
#include <mach/mach_error.h>
#include <mach/mach_host.h>
#include <mach/mach_port.h>
#include <mach/host_info.h>
*/
import "C"

func (c *CPU) collect() {
	defer func() {
		if r := recover(); r != nil {
			c.clear()
		}
	}()

	// collect CPU stats for All cpus aggregated
	var cpuinfo C.host_cpu_load_info_data_t
	count := C.mach_msg_type_number_t(C.HOST_CPU_LOAD_INFO_COUNT)
	host := C.mach_host_self()

	ret := C.host_statistics(C.host_t(host), C.HOST_CPU_LOAD_INFO,
		C.host_info_t(unsafe.Pointer(&cpuinfo)), &count)

	if ret != C.KERN_SUCCESS {
		panic(fmt.Errorf("error: %d", ret))
	}

	c.lastUpdate = time.Now()
	c.previous = c.current

	c.current = []int{
		userGauge:   int(cpuinfo.cpu_ticks[C.CPU_STATE_USER]),
		niceGauge:   int(cpuinfo.cpu_ticks[C.CPU_STATE_NICE]),
		systemGauge: int(cpuinfo.cpu_ticks[C.CPU_STATE_SYSTEM]),
		idleGauge:   int(cpuinfo.cpu_ticks[C.CPU_STATE_IDLE]),
	}

	if c.previous == nil {
		c.previous = c.current
	}

	c.previousTotal = c.currentTotal
	c.currentTotal = c.current[userGauge] + c.current[niceGauge] + c.current[systemGauge] + c.current[idleGauge]
}
