package main

import (
	"fmt"
	"os/linux"
)

func main() {
	cpu, err := linux.CalcCPUUsage()
	if err != nil {
		fmt.Println(err.Error())
	}

	devs, err := linux.CalcDevStats()
	if err != nil {
		fmt.Println(err.Error())
	}

	fsystems, err := linux.CalcFsUtilisation()
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println("la:", cpu.LA)
	fmt.Printf("CPU usr: %.2f%%, sys: %.2f%%, ide: %.2f%%\n", cpu.UsrUsagePercent, cpu.SysUsagePercent, cpu.IdlePercent)
	fmt.Println("Devices statistic:")
	for _, dev := range devs {
		fmt.Printf("Name: %s, Transactions per sec: %.3f, Read, Kbps: %.3f, Write, Kbps: %.3f\n",
			dev.Name,
			dev.TransPS,
			dev.ReadPS,
			dev.WritePS,
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
