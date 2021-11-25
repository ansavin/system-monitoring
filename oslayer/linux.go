//go:build linux
// +build linux

package oslayer

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"syscall"
	"time"
)

type devInfo struct {
	Transactions uint64
	ReadBytes    uint64
	WriteBytes   uint64
}

// SectorSize is UNIX sector size
const SectorSize = 512

// BytesInKb is num of bytes in KB, in Linux its still 2^10
const BytesInKb = 1024

// BytesInGb is num of bytes in GB
const BytesInGb = 1073741824

// LaFile is procfs file in Linux that shows LA
var LaFile = "/proc/loadavg"

// CPUStatsFile is procfs file in Linux that shows CPU verbose stats
var CPUStatsFile = "/proc/stat"

// MountinfoFile is procfs file in Linux that shows mounted filesystems
var MountinfoFile = "/proc/1/mountinfo"

// BlockDevicesDir is procfs dir in Linux that shows block devices
var BlockDevicesDir = "/sys/block"

// DevStatsFilename is file in /sys/block/<dev_name> that shows this device stats
var DevStatsFilename = "stat"

// in docker we mount host root fs to /host
// we make it var, not const in order to have possibility to write
// unit-tests for checkIfRunsInDocker
var dockerRootFSPrefix = "/host"

// ValidDeviceTypes is prefixes for non-virtual block dev in Linux mountinfo
// See https://www.kernel.org/doc/Documentation/admin-guide/devices.txt
var ValidDeviceTypes = []string{"3", "8", "9", "22", "33", "34"}

// check if code runs in docker or in host OS
func checkIfRunsInDocker() string {
	_, err := os.Open(dockerRootFSPrefix)
	// if /host does not exist, then we run code in host OS
	// we does not check other errors here because in every place we call
	// this func we will check file open errors
	if errors.Is(err, os.ErrNotExist) {
		return ""
	}
	// if /host exist, then we run in docker
	// (we should bind-mount host / to /host in order to get info about FS)
	return dockerRootFSPrefix
}

func percentage(x, y float64) float64 {
	if x == 0 {
		return 0
	}
	return 100 * (x - y) / x
}

func getDevStats() (map[string]devInfo, error) {
	res := make(map[string]devInfo)

	devs, err := ioutil.ReadDir(checkIfRunsInDocker() + BlockDevicesDir)
	if err != nil {
		return nil, fmt.Errorf(
			"cannot read %s: %s",
			checkIfRunsInDocker()+BlockDevicesDir,
			err.Error())
	}
	for _, d := range devs {
		name := d.Name()
		if strings.HasPrefix(name, "loop") {
			continue
		}
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

	prefix := checkIfRunsInDocker()

	dev, err := ioutil.ReadFile(prefix + path)
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

	time.Sleep(SamplingTime)

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

	sort.Slice(res[:], func(i, j int) bool {
		return strings.Compare(res[i].Name, res[j].Name) > 0
	})

	return res, nil
}

func parseMounts() ([]string, error) {
	var tmp, devType, rootDentry, mountPoint string
	res := make([]string, 0)

	prefix := checkIfRunsInDocker()

	mounts, err := ioutil.ReadFile(prefix + MountinfoFile)
	if err != nil {
		return nil, fmt.Errorf("cannot open %s: %s", MountinfoFile, err.Error())
	}

	rows := strings.Split(string(mounts), "\n")
	for _, row := range rows {
		// in go we cannot skip entry in Sscanf :(
		_, err := fmt.Sscanf(row, "%s %s %s %s %s", &tmp, &tmp, &devType, &rootDentry, &mountPoint)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("cannot parse %s: %s", MountinfoFile, err.Error())
		}

		// we want to get info only for root mounts
		if rootDentry != "/" {
			continue
		}

		// we check if FS here is not virtual (for ex., sysfs)
		// we doesn't want to check storage for such a FS
		for _, prefix := range ValidDeviceTypes {
			if strings.HasPrefix(devType, prefix+":") {
				res = append(res, mountPoint)
				break
			}
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
		prefix := checkIfRunsInDocker()

		fd, err := os.Open(prefix + fs)

		if err != nil {
			fd.Close()
			return nil, fmt.Errorf("cannot open %s fs: %s", fs, err.Error())
		}
		err = syscall.Fstatfs(int(fd.Fd()), &stats)
		if err != nil {
			fd.Close()
			return nil, fmt.Errorf("syscall statfs returns error: %s", err.Error())
		}

		res = append(res, FsStats{
			Name:               fs,
			UsedGBytes:         float64(stats.Blocks-stats.Bfree) * float64(stats.Bsize) / BytesInGb,
			UsedStoragePercent: percentage(float64(stats.Blocks), float64(stats.Bfree)),
			UsedInodes:         float64(stats.Files - stats.Ffree),
			UsedInodesPercent:  percentage(float64(stats.Files), float64(stats.Ffree)),
		})

	}

	sort.Slice(res[:], func(i, j int) bool {
		return strings.Compare(res[i].Name, res[j].Name) > 0
	})

	return res, nil
}

func parseLA() (float64, error) {
	var la float64

	prefix := checkIfRunsInDocker()

	laStats, err := ioutil.ReadFile(prefix + LaFile)
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

	prefix := checkIfRunsInDocker()

	cpu, err := ioutil.ReadFile(prefix + CPUStatsFile)
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
