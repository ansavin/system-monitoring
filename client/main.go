package main

import (
	"context"
	"fmt"
	"time"

	"protobuf"

	"google.golang.org/grpc"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial("localhost:8088", grpc.WithInsecure())
	if err != nil {
		fmt.Println("cannot connect to gRPC server:", err.Error())
	}
	defer conn.Close()
	c := protobuf.NewMonitorClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.GetStats(ctx, &protobuf.StatType{})
	if err != nil {
		fmt.Println("cannot get stats from server:", err.Error())
	}

	fmt.Println("la:", r.La)
	fmt.Printf("CPU usr: %.2f%%, sys: %.2f%%, ide: %.2f%%\n", r.Usr, r.Sys, r.Idle)
}
