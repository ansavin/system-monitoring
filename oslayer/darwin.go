//go:build darwin && cgo
// +build darwin,cgo

package oslayer

import (
	"fmt"
	"unsafe"
)

// we need a separate import for C because otherwise the Go preprocessor
// won't notice it in bulk import => code won't compile

// #include <stdlib.h>
// #include <mach/mach_host.h>
// #include <mach/host_info.h>
import "C"

// we have to use cgo to get stats form macos here
func parseLA() (float64, error) {
	var laStats [1]C.double

	returnCode := C.getloadavg(&laStats[0], 1)
	if returnCode != 1 {
		return 0, fmt.Errorf("cannot get loadavg")
	}

	return float64(laStats[0]), nil
}

// we have to use cgo to get stats form macos here
// https://github.com/apple/darwin-xnu/blob/main/osfmk/kern/host.c#L484
func parseCPUStats() (cpuInfo, error) {
	var cpuLoadInfoData C.host_cpu_load_info_data_t
	var hostCpuLoadInfoCount C.mach_msg_type_number_t = C.HOST_CPU_LOAD_INFO_COUNT

	returnCode := C.host_statistics(
		C.host_t(C.mach_host_self()),
		C.HOST_CPU_LOAD_INFO,
		C.host_info_t(unsafe.Pointer(&cpuLoadInfoData)),
		&hostCpuLoadInfoCount,
	)
	if returnCode != C.KERN_SUCCESS {
		return cpuInfo{}, fmt.Errorf("error in host_statistics: %d", returnCode)
	}

	return cpuInfo{
		usr: float64(cpuLoadInfoData.cpu_ticks[C.CPU_STATE_USER] +
			cpuLoadInfoData.cpu_ticks[C.CPU_STATE_NICE]),
		sys:  float64(cpuLoadInfoData.cpu_ticks[C.CPU_STATE_SYSTEM]),
		idle: float64(cpuLoadInfoData.cpu_ticks[C.CPU_STATE_IDLE]),
	}, nil
}

// CalcCPUUsage returns CPU utilization data like
// load average & CPU time distribution (sys/usr/idle)
func CalcCPUUsage(_ string) (CPUstats, error) {
	la, err := parseLA()
	if err != nil {
		return CPUstats{}, fmt.Errorf("cannot calc cpu la: %s", err.Error())
	}

	stats, err := parseCPUStats()
	if err != nil {
		return CPUstats{}, fmt.Errorf("cannot calc cpu stats: %s", err.Error())
	}

	tot := (stats.usr + stats.sys + stats.idle) / 100

	return CPUstats{
		LA:              la,
		UsrUsagePercent: stats.usr / tot,
		SysUsagePercent: stats.sys / tot,
		IdlePercent:     stats.idle / tot,
	}, nil
}

// Not implemented yet
func CalcDevStats(_ string) ([]DevStats, error) {
	return nil, nil
}

// Not implemented yet
func CalcFsUtilization(_ string) ([]FsStats, error) {
	return nil, nil
}
