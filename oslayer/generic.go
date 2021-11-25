package oslayer

import "time"

type cpuInfo struct {
	usr  float64
	sys  float64
	idle float64
}

// CPUstats represend CPU stats
type CPUstats struct {
	LA              float64
	UsrUsagePercent float64
	SysUsagePercent float64
	IdlePercent     float64
}

// DevStats represent block devises stats
type DevStats struct {
	Name    string
	TransPS float64
	ReadPS  float64
	WritePS float64
}

// FsStats represents filesystem stats
type FsStats struct {
	Name               string
	UsedGBytes         float64
	UsedStoragePercent float64
	UsedInodes         float64
	UsedInodesPercent  float64
}

// SamplingTime is time interval used for calculating some
// statistics. We cant send client statistics with higher granularity
const SamplingTime = time.Second
