package main

import (
	"context"
	"flag"
	"fmt"
	"io"

	"protobuf"

	"google.golang.org/grpc"
)

func main() {
	avgTime := flag.Int("a", 3, "statistics averaging time in seconds")
	msgTime := flag.Int("m", 3, "time between messages in seconds")
	port := flag.Int("p", 8088, "port at which statistics server runs")

	flag.Parse()

	timeBetweenTicks := *msgTime
	if timeBetweenTicks <= 0 {
		fmt.Println("expected integer > 0 as 1-d argument, got", timeBetweenTicks)
		return
	}

	averagingTime := *avgTime
	if averagingTime <= 0 {
		fmt.Println("expected integer > 0 as 2-d argument, got", averagingTime)
		return
	}

	// changing port for startup is nit supported yet
	conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", *port), grpc.WithInsecure())
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
		fmt.Printf("la: %.2f\n", r.CPUstats.La)
		fmt.Printf("CPU usr: %.2f%%, sys: %.2f%%, idle: %.2f%%\n", r.CPUstats.Usr, r.CPUstats.Sys, r.CPUstats.Idle)

		fmt.Println("Devices statistics:")
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
			fmt.Printf("Name: %s, Used storage: %.3f Gb, Used storage percentage: %.2f%%, Used inodes: %.0f, Used inodes percentage: %.2f%%\n",
				fs.Name,
				fs.Bytes,
				fs.BytesPercent,
				fs.Inode,
				fs.InodePercent,
			)
		}
	}
}
