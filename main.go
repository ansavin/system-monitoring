package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"syscall"
	"time"
)

type DevInfo struct {
	Transactions uint64
	ReadBytes    uint64
	WriteBytes   uint64
}

type DevStats struct {
	Name    string
	TransPS float64
	ReadPS  float64
	WritePS float64
}

type FsStats struct {
	Name               string
	UsedBytes          uint64
	UsedStoragePercent float64
	UsedInodes         uint64
	UsedInodesPercent  float64
}

type CPUInfo struct {
	usr  float64
	sys  float64
	idle float64
}

type CPUstats struct {
	LA              float64
	usrUsagePercent float64
	sysUsagePercent float64
	idlePercent     float64
}

// UNIX sector size
const SECTOR_SIZE = 512

// num of bytes in KB
const BYTES_IN_KB = 1000

// prefix for virtual dev (for ex, sysfs) in Linux mountinfo
const VIRTUAL_DEVICE_TYPE = "0"

// procfs file in Linux that shows LA
const LA_FILE = "/proc/loadavg"

// procfs file in Linux that shows CPU verbose stats
const CPU_STATS_FILE = "/proc/stat"

// procfs file in Linux that shows mounted filesystems
const MOUNTINFO_FILE = "/proc/self/mountinfo"

// procfs dir in Linux that shows block devices
const BLOCK_DEVICES_DIR = "/sys/block"

// file in /sys/block/<dev_name> that shows this device stats
const DEV_STATS_FILENAME = "stat"

func persentage(x, y float64) float64 {
	if x == 0 {
		return 0
	}
	return 100 * (x - y) / x
}

func getDevStats() (map[string]DevInfo, error) {
	res := make(map[string]DevInfo)
	devs, err := ioutil.ReadDir(BLOCK_DEVICES_DIR)
	if err != nil {
		return nil, fmt.Errorf("cannot read %s: %s", BLOCK_DEVICES_DIR, err.Error())
	}
	for _, d := range devs {
		name := d.Name()
		data, err := parseDevStats(name, BLOCK_DEVICES_DIR)
		if err != nil {
			return nil, fmt.Errorf("cannot parse dev stats: %s", err.Error())
		}
		res[name] = data
	}
	return res, nil
}

func parseDevStats(name, basePath string) (DevInfo, error) {
	var tmp, dscdReq, readReq, readSect, writeReq, writeSect uint64

	path := fmt.Sprintf("%s/%s/%s", basePath, name, DEV_STATS_FILENAME)

	dev, err := ioutil.ReadFile(path)
	if err != nil {
		return DevInfo{}, fmt.Errorf("cannot read %s: %s", path, err.Error())
	}

	_, err = fmt.Sscanf(string(dev), "%d %d %d %d %d %d %d %d %d %d %d %d",
		&writeReq, &tmp, &writeSect, &tmp,
		&readReq, &tmp, &readSect, &tmp, &tmp, &tmp, &tmp,
		&dscdReq,
	) // in go we cannot skip entry in Sscanf :(
	if err != nil {
		return DevInfo{}, fmt.Errorf("cannot parse %s: %s", path, err.Error())
	}

	return DevInfo{
		Transactions: dscdReq + readReq + writeReq,
		ReadBytes:    readSect * SECTOR_SIZE,
		WriteBytes:   writeSect * SECTOR_SIZE,
	}, nil

}

func calcDevStats() ([]DevStats, error) {
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
			ReadPS:  float64(stats.ReadBytes) - float64(devFirstSnapshot[name].ReadBytes),
			WritePS: float64(stats.WriteBytes) - float64(devFirstSnapshot[name].WriteBytes),
		})
	}
	return res, nil
}

func parseMounts() ([]string, error) {
	var tmp, devType, mountPoint string
	res := make([]string, 0)

	mounts, err := ioutil.ReadFile(MOUNTINFO_FILE)
	if err != nil {
		return nil, fmt.Errorf("cannot open %s: %s", MOUNTINFO_FILE, err.Error())
	}

	rows := strings.Split(string(mounts), "\n")
	for _, row := range rows {
		// in go we cannot skip entry in Sscanf :(
		_, err := fmt.Sscanf(row, "%s %s %s %s %s", &tmp, &tmp, &devType, &tmp, &mountPoint)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("cannot parse %s: %s", MOUNTINFO_FILE, err.Error())
		}

		// we check if FS here is not virtual (for ex., sysfs)
		// we doesn't want to check storage for such a FS
		if !strings.HasPrefix(devType, VIRTUAL_DEVICE_TYPE) {
			res = append(res, mountPoint)
		}
	}
	return res, nil
}

func calcFsUtilisation() ([]FsStats, error) {
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

		res = append(res, FsStats{
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

	laStats, err := ioutil.ReadFile(LA_FILE)
	if err != nil {
		return 0, fmt.Errorf("cannot read %s: %s", LA_FILE, err.Error())
	}

	_, err = fmt.Sscanf(string(laStats), "%f", &la)
	if err != nil {
		return 0, fmt.Errorf("cannot parse /proc/loadavg: %s", err.Error())
	}

	return la, nil
}

func parseCPUStats() (CPUInfo, error) {
	var usr, sys, idle float64
	var tmp string

	cpu, err := ioutil.ReadFile(CPU_STATS_FILE)
	if err != nil {
		return CPUInfo{}, fmt.Errorf("cannot read %s: %s", CPU_STATS_FILE, err.Error())
	}

	// in go we cannot skip entry in Sscanf :(
	_, err = fmt.Sscanf(string(cpu), "%s %f %s %f %f", &tmp, &usr, &tmp, &sys, &idle)
	if err != nil {
		return CPUInfo{}, fmt.Errorf("cannot parse %s: %s", CPU_STATS_FILE, err.Error())
	}

	return CPUInfo{
		usr:  usr,
		sys:  sys,
		idle: idle,
	}, nil
}

func calcCPUUsage() (CPUstats, error) {
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
		usrUsagePercent: stats.usr / tot,
		sysUsagePercent: stats.sys / tot,
		idlePercent:     stats.idle / tot,
	}, nil
}

func main() {
	cpu, err := calcCPUUsage()
	if err != nil {
		fmt.Println(err.Error())
	}

	devs, err := calcDevStats()
	if err != nil {
		fmt.Println(err.Error())
	}

	fsystems, err := calcFsUtilisation()
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println("la:", cpu.LA)
	fmt.Printf("CPU usr: %.2f%%, sys: %.2f%%, ide: %.2f%%\n", cpu.usrUsagePercent, cpu.sysUsagePercent, cpu.idlePercent)
	fmt.Println("Devices statistic:")
	for _, dev := range devs {
		fmt.Printf("Name: %s, Transactions per sec: %.3f, Read, Kbps: %.3f, Write, Kbps: %.3f\n",
			dev.Name,
			dev.TransPS,
			dev.ReadPS/BYTES_IN_KB,
			dev.WritePS/BYTES_IN_KB,
		)
	}

	fmt.Println("Filesystems utilization:")
	for _, fs := range fsystems {
		fmt.Printf("Name: %s, Used storage, bytes: %d, Used storage persentage: %.2f%%, Used inodes: %d, Used inodes persentage: %.2f%%\n",
			fs.Name,
			fs.UsedBytes,
			fs.UsedStoragePercent,
			fs.UsedInodes,
			fs.UsedInodesPercent,
		)
	}
}
