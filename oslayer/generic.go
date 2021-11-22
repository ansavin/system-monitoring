package oslayer

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
	UsedGBytes         uint64
	UsedStoragePercent float64
	UsedInodes         uint64
	UsedInodesPercent  float64
}
