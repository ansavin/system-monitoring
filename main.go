package main

import (
	"context"
	"fmt"
	"net"
	"os/linux"

	"protobuf"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	protobuf.UnimplementedMonitorServer
}

// getStats implements protobuf.MonitorServer
func (s *server) GetStats(context.Context, *protobuf.StatType) (*protobuf.CPUstats, error) {
	cpu, err := linux.CalcCPUUsage()
	if err != nil {
		return nil, status.Errorf(codes.Aborted, "server error: %s", err.Error())
	}

	// devs, err := linux.CalcDevStats()
	// if err != nil {
	// 	return nil, status.Errorf(codes.Aborted, "server error: %s", err.Error())
	// }

	// fsystems, err := linux.CalcFsUtilisation()
	// if err != nil {
	// 	return nil, status.Errorf(codes.Aborted, "server error: %s", err.Error())
	// }

	fmt.Println("la:", cpu.LA)
	fmt.Printf("CPU usr: %.2f%%, sys: %.2f%%, ide: %.2f%%\n", cpu.UsrUsagePercent, cpu.SysUsagePercent, cpu.IdlePercent)
	// fmt.Println("Devices statistic:")
	// for _, dev := range devs {
	// 	fmt.Printf("Name: %s, Transactions per sec: %.3f, Read, Kbps: %.3f, Write, Kbps: %.3f\n",
	// 		dev.Name,
	// 		dev.TransPS,
	// 		dev.ReadPS,
	// 		dev.WritePS,
	// 	)
	// }

	// fmt.Println("Filesystems utilization:")
	// for _, fs := range fsystems {
	// 	fmt.Printf("Name: %s, Used storage, bytes: %d, Used storage persentage: %.2f%%, Used inodes: %d, Used inodes persentage: %.2f%%\n",
	// 		fs.Name,
	// 		fs.UsedBytes,
	// 		fs.UsedStoragePercent,
	// 		fs.UsedInodes,
	// 		fs.UsedInodesPercent,
	// 	)
	// }

	return &protobuf.CPUstats{
		La:   cpu.LA,
		Sys:  cpu.SysUsagePercent,
		Usr:  cpu.UsrUsagePercent,
		Idle: cpu.IdlePercent,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", "localhost:8088")
	if err != nil {
		fmt.Println("cannot listen TCP:", err.Error())
	}
	s := grpc.NewServer()
	protobuf.RegisterMonitorServer(s, &server{})
	fmt.Println("server listening at", lis.Addr())
	if err := s.Serve(lis); err != nil {
		fmt.Println("cannot handle reguest:", err.Error())
	}
}
