package linux

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"syscall"
	"time"
)

type devInfo struct {
	Transactions uint64
	ReadBytes    uint64
	WriteBytes   uint64
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
	UsedBytes          uint64
	UsedStoragePercent float64
	UsedInodes         uint64
	UsedInodesPercent  float64
}

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

// SectorSize is UNIX sector size
const SectorSize = 512

// BytesInKb is num of bytes in KB
const BytesInKb = 1000

// VirtualDeviceType is prefix for virtual dev (for ex, sysfs) in Linux mountinfo
const VirtualDeviceType = "0"

// LaFile is procfs file in Linux that shows LA
const LaFile = "/proc/loadavg"

// CPUStatsFile is procfs file in Linux that shows CPU verbose stats
const CPUStatsFile = "/proc/stat"

// MountinfoFile is procfs file in Linux that shows mounted filesystems
const MountinfoFile = "/proc/self/mountinfo"

// BlockDevicesDir is procfs dir in Linux that shows block devices
const BlockDevicesDir = "/sys/block"

// DevStatsFilename is file in /sys/block/<dev_name> that shows this device stats
const DevStatsFilename = "stat"

func persentage(x, y float64) float64 {
	if x == 0 {
		return 0
	}
	return 100 * (x - y) / x
}

func getDevStats() (map[string]devInfo, error) {
	res := make(map[string]devInfo)
	devs, err := ioutil.ReadDir(BlockDevicesDir)
	if err != nil {
		return nil, fmt.Errorf("cannot read %s: %s", BlockDevicesDir, err.Error())
	}
	for _, d := range devs {
		name := d.Name()
		data, err := parseDevStats(name, BlockDevicesDir)
		if err != nil {
			return nil, fmt.Errorf("cannot parse dev stats: %s", err.Error())
		}
		res[name] = data
	}
	return res, nil
}

func parseDevStats(name, basePath string) (devInfo, error) {
	var tmp, dscdReq, readReq, readSect, writeReq, writeSect uint64

	path := fmt.Sprintf("%s/%s/%s", basePath, name, DevStatsFilename)

	dev, err := ioutil.ReadFile(path)
	if err != nil {
		return devInfo{}, fmt.Errorf("cannot read %s: %s", path, err.Error())
	}

	_, err = fmt.Sscanf(string(dev), "%d %d %d %d %d %d %d %d %d %d %d %d",
		&writeReq, &tmp, &writeSect, &tmp,
		&readReq, &tmp, &readSect, &tmp, &tmp, &tmp, &tmp,
		&dscdReq,
	) // in go we cannot skip entry in Sscanf :(
	if err != nil {
		return devInfo{}, fmt.Errorf("cannot parse %s: %s", path, err.Error())
	}

	return devInfo{
		Transactions: dscdReq + readReq + writeReq,
		ReadBytes:    readSect * SectorSize,
		WriteBytes:   writeSect * SectorSize,
	}, nil

}

// CalcDevStats returns us Read/Write per sec & transactions per sec for
// all blk devises
func CalcDevStats() ([]DevStats, error) {
	res := make([]DevStats, 0)

	devFirstSnapshot, err := getDevStats()
	if err != nil {
		return nil, fmt.Errorf("cannot get dev snapshot #1: %s", err.Error())
	}

	time.Sleep(1 * time.Second) //FIXME

	devSecondSnapshot, err := getDevStats()
	if err != nil {
		return nil, fmt.Errorf("cannot get dev snapshot #2: %s", err.Error())
	}

	if len(devFirstSnapshot) != len(devSecondSnapshot) {
		return nil, fmt.Errorf("dev snapshots lengths mismatch")
	}
	for name, stats := range devSecondSnapshot {
		res = append(res, DevStats{
			Name:    name,
			TransPS: float64(stats.Transactions) - float64(devFirstSnapshot[name].Transactions),
			ReadPS:  (float64(stats.ReadBytes) - float64(devFirstSnapshot[name].ReadBytes)) / BytesInKb,
			WritePS: (float64(stats.WriteBytes) - float64(devFirstSnapshot[name].WriteBytes)) / BytesInKb,
		})
	}
	return res, nil
}

func parseMounts() ([]string, error) {
	var tmp, devType, mountPoint string
	res := make([]string, 0)

	mounts, err := ioutil.ReadFile(MountinfoFile)
	if err != nil {
		return nil, fmt.Errorf("cannot open %s: %s", MountinfoFile, err.Error())
	}

	rows := strings.Split(string(mounts), "\n")
	for _, row := range rows {
		// in go we cannot skip entry in Sscanf :(
		_, err := fmt.Sscanf(row, "%s %s %s %s %s", &tmp, &tmp, &devType, &tmp, &mountPoint)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("cannot parse %s: %s", MountinfoFile, err.Error())
		}

		// we check if FS here is not virtual (for ex., sysfs)
		// we doesn't want to check storage for such a FS
		if !strings.HasPrefix(devType, VirtualDeviceType) {
			res = append(res, mountPoint)
		}
	}
	return res, nil
}

// CalcFsUtilisation return storage & inodes utilisation for all
// mounted non-virtual (like procfs) filesystems
func CalcFsUtilisation() ([]FsStats, error) {
	var stats syscall.Statfs_t

	res := make([]FsStats, 0)

	filesystems, err := parseMounts()
	if err != nil {
		return nil, fmt.Errorf("cannot get filesystems: %s", err.Error())
	}

	for _, fs := range filesystems {
		fd, err := os.Open(fs)

		if err != nil {
			fd.Close()
			return nil, fmt.Errorf("cannot open %s fs: %s", fs, err.Error())
		}
		err = syscall.Fstatfs(int(fd.Fd()), &stats)
		if err != nil {
			fd.Close()
			return nil, fmt.Errorf("syscall statfs returns error: %s", err.Error())
		}

		res = append(res, FsStats{ //FIXME
			Name:               fs,
			UsedBytes:          stats.Files - stats.Bfree,
			UsedStoragePercent: persentage(float64(stats.Blocks), float64(stats.Bfree)),
			UsedInodes:         stats.Files - stats.Ffree,
			UsedInodesPercent:  persentage(float64(stats.Files), float64(stats.Ffree)),
		})
	}
	return res, nil
}

func parseLA() (float64, error) {
	var la float64

	laStats, err := ioutil.ReadFile(LaFile)
	if err != nil {
		return 0, fmt.Errorf("cannot read %s: %s", LaFile, err.Error())
	}

	_, err = fmt.Sscanf(string(laStats), "%f", &la)
	if err != nil {
		return 0, fmt.Errorf("cannot parse /proc/loadavg: %s", err.Error())
	}

	return la, nil
}

func parseCPUStats() (cpuInfo, error) {
	var usr, sys, idle float64
	var tmp string

	cpu, err := ioutil.ReadFile(CPUStatsFile)
	if err != nil {
		return cpuInfo{}, fmt.Errorf("cannot read %s: %s", CPUStatsFile, err.Error())
	}

	// in go we cannot skip entry in Sscanf :(
	_, err = fmt.Sscanf(string(cpu), "%s %f %s %f %f", &tmp, &usr, &tmp, &sys, &idle)
	if err != nil {
		return cpuInfo{}, fmt.Errorf("cannot parse %s: %s", CPUStatsFile, err.Error())
	}

	return cpuInfo{
		usr:  usr,
		sys:  sys,
		idle: idle,
	}, nil
}

// CalcCPUUsage returns CPU utilization data like
// load average & CPU time distribution (sys/usr/idle)
func CalcCPUUsage() (CPUstats, error) {
	la, err := parseLA()
	if err != nil {
		return CPUstats{}, fmt.Errorf("cannot Calc cpu la: %s", err.Error())
	}

	stats, err := parseCPUStats()
	if err != nil {
		return CPUstats{}, fmt.Errorf("cannot Calc cpu stats: %s", err.Error())
	}

	tot := (stats.usr + stats.sys + stats.idle) / 100

	return CPUstats{
		LA:              la,
		UsrUsagePercent: stats.usr / tot,
		SysUsagePercent: stats.sys / tot,
		IdlePercent:     stats.idle / tot,
	}, nil
}
