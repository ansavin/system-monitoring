package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"strconv"

	"protobuf"

	"google.golang.org/grpc"
)

func main() {
	args := os.Args

	if len(args) < 3 {
		fmt.Println("usage: ./grpc-server <time_between_messages_in_sec> <time_for_stats_averaging_is_sec>")
		return
	}

	timeBetweenTicks, err := strconv.Atoi(args[1])
	if err != nil || timeBetweenTicks <= 0 {
		fmt.Println("expected integer > 0 as 1-d argument, got", args[1])
		return
	}
	averagingTime, err := strconv.Atoi(args[2])
	if err != nil || averagingTime <= 0 {
		fmt.Println("expected integer > 0 as 2-d argument, got", args[2])
		return
	}

	// changing port for startup is nit supported yet
	conn, err := grpc.Dial("localhost:8088", grpc.WithInsecure())
	if err != nil {
		fmt.Println("cannot connect to gRPC server:", err.Error())
		return
	}
	defer conn.Close()
	c := protobuf.NewMonitorClient(conn)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stream, err := c.GetStats(ctx, &protobuf.Settings{
		TimeBetweenTicks: uint32(timeBetweenTicks),
		AveragingTime:    uint32(averagingTime),
	})
	if err != nil {
		fmt.Println("cannot open stream:", err.Error())
		return
	}

	for {
		r, err := stream.Recv()
		if err == io.EOF {
			fmt.Println("stream is closed, finishing session")
			return
		}
		if err != nil {
			fmt.Println("error during streaming:", err.Error())
			return
		}

		fmt.Println("CPU statistics:")
		fmt.Println("la:", r.CPUstats.La)
		fmt.Printf("CPU usr: %.2f%%, sys: %.2f%%, ide: %.2f%%\n", r.CPUstats.Usr, r.CPUstats.Sys, r.CPUstats.Idle)

		fmt.Println("Devices statistic:")
		for _, dev := range r.DevStats {
			fmt.Printf("Name: %s, Transactions per sec: %.3f, Read: %.3f Kbps, Write: %.3f Kbps\n",
				dev.Name,
				dev.Tps,
				dev.Read,
				dev.Write,
			)
		}

		fmt.Println("Filesystems utilization:")
		for _, fs := range r.FsStats {
			fmt.Printf("Name: %s, Used storage: %.3f Gb, Used storage persentage: %.2f%%, Used inodes: %.0f, Used inodes persentage: %.2f%%\n",
				fs.Name,
				fs.Bytes,
				fs.BytesPercent,
				fs.Inode,
				fs.InodePercent,
			)
		}
	}
}
