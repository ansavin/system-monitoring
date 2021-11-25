package main

import (
	"context"
	"fmt"
	"net"
	"oslayer"
	"sync"
	"time"

	"protobuf"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// maximum number of stored data in sample slices
const maxSamples = 300

type server struct {
	protobuf.UnimplementedMonitorServer
	cpuStatsSamples []oslayer.CPUstats
	fsStatsSamples  [][]oslayer.FsStats
	devStatsSamples [][]oslayer.DevStats
	once            [3]sync.Once
	ctx             context.Context
	cancel          context.CancelFunc
	lastError       error
}

func startPeriodicalSampling(ctx context.Context, f func()) {
	fmt.Printf("goroutine satrted with func %#v\n", f)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			f()
		}
	}
}

// getStats implements protobuf.MonitorServer
func (s *server) GetStats(settings *protobuf.Settings, srv protobuf.Monitor_GetStatsServer) error {
	if settings.AveragingTime > maxSamples {
		return status.Errorf(codes.Aborted, "server cannot average statistics for more then %d sec", maxSamples)
	}
	if settings.TimeBetweenTicks > maxSamples {
		return status.Errorf(codes.Aborted, "server cannot send statistics less then 1 time in %d sec", maxSamples)
	}

	// this stats will be shared between all clients
	// the main idea is to have 1 sample storage for each statistics
	// (like CPUstats) and then use data from it for averaging for
	// each client, like this:
	// [1, 2, 3, 4, 5, 6] <- data samples with interval 1 sec
	// clientData = 1, 2, 3, ... <- for client who wants data for 1 sec
	// anotherClientData = (2 + 1) / 2, (3 + 2) / 2, ... <- for client
	// who wants data for 2 sec
	s.once[0].Do(
		func() {
			go startPeriodicalSampling(s.ctx, func() {
				cpu, err := oslayer.CalcCPUUsage()
				if err != nil {
					s.lastError = err
					s.cancel()
				}

				// we have no generic in go, so we cant move this to a func
				// wuthout overengeneering
				if len(s.cpuStatsSamples) == maxSamples {
					s.cpuStatsSamples = s.cpuStatsSamples[1:]
				}
				s.cpuStatsSamples = append(s.cpuStatsSamples, cpu)

				time.Sleep(oslayer.SamplingTime)
			})
		},
	)
	s.once[1].Do(
		func() {
			go startPeriodicalSampling(s.ctx, func() {
				devs, err := oslayer.CalcDevStats()
				if err != nil {
					s.lastError = err
					s.cancel()
				}

				// we have no generic in go, so we cant move this to a func
				// wuthout overengeneering
				if len(s.devStatsSamples) == maxSamples {
					s.devStatsSamples = s.devStatsSamples[1:]
				}
				s.devStatsSamples = append(s.devStatsSamples, devs)
				// we doesn`t need to sleep here - CalcDevStats use the same SamplingTime
				// for collecting stats
			})
		},
	)
	s.once[2].Do(
		func() {
			go startPeriodicalSampling(s.ctx, func() {
				fsystems, err := oslayer.CalcFsUtilisation()
				if err != nil {
					s.lastError = err
					s.cancel()
				}

				// we have no generic in go, so we cant move this to a func
				// wuthout overengeneering
				if len(s.fsStatsSamples) == maxSamples {
					s.fsStatsSamples = s.fsStatsSamples[1:]
				}
				s.fsStatsSamples = append(s.fsStatsSamples, fsystems)

				time.Sleep(oslayer.SamplingTime)
			})
		},
	)

	for {
		select {
		case <-s.ctx.Done():
			return status.Errorf(codes.Aborted, "collecting statistics error: %s", s.lastError.Error())
		default:
			for int(settings.AveragingTime) > (len(s.cpuStatsSamples)) {
				time.Sleep(time.Duration(int(settings.AveragingTime)-len(s.cpuStatsSamples)) * time.Second)
			}

			start := len(s.cpuStatsSamples) - (int(settings.AveragingTime))
			cpuStats := protobuf.CPUstats{}

			for _, sample := range s.cpuStatsSamples[start:] {
				cpuStats.La += sample.LA
				cpuStats.Sys += sample.SysUsagePercent
				cpuStats.Usr += sample.UsrUsagePercent
				cpuStats.Idle += sample.IdlePercent
			}
			cpuStats.La /= float64(settings.AveragingTime)
			cpuStats.Sys /= float64(settings.AveragingTime)
			cpuStats.Usr /= float64(settings.AveragingTime)
			cpuStats.Idle /= float64(settings.AveragingTime)

			stats := protobuf.Stats{
				CPUstats: &cpuStats,
			}

			// defencive programming: we assume slices are filled almost simultaneously
			// but despiting that fact, we will check if slice is ok before working
			for int(settings.AveragingTime) > (len(s.devStatsSamples)) {
				fmt.Println("len s.devStatsSamples", len(s.devStatsSamples))
				time.Sleep(time.Duration(int(settings.AveragingTime)-len(s.devStatsSamples)) * time.Second)
			}

			fmt.Println("cpu, dev", len(s.cpuStatsSamples), len(s.devStatsSamples))
			start = len(s.devStatsSamples) - (int(settings.AveragingTime))

			if len(s.devStatsSamples[start:]) < int(settings.AveragingTime) {
				return status.Errorf(codes.Aborted, "sending message error: device statistics corrupted")
			}

			// limitation: we assume device number is constant. If not, there might be an error there
			for devIdx, dev := range s.devStatsSamples[0] {
				fmt.Println("devIdx, devName", devIdx, dev.Name)
				devStats := protobuf.DevStats{
					Name: dev.Name,
				}

				for idx, sample := range s.devStatsSamples[start:] {

					spl := sample[devIdx]
					fmt.Println("idx, spl.Name", idx, spl.Name)
					if devStats.Name != spl.Name {
						return status.Errorf(codes.Aborted, "sending message error: mixing different devices stats")
					}

					devStats.Tps += spl.TransPS
					devStats.Read += spl.ReadPS
					devStats.Write += spl.WritePS
				}
				fmt.Println("devStats.Tps, devStats.Read, devStats.Write", devStats.Tps, devStats.Read, devStats.Write)

				devStats.Tps /= float64(settings.AveragingTime)
				devStats.Read /= float64(settings.AveragingTime)
				devStats.Write /= float64(settings.AveragingTime)
				fmt.Println("devStats.Tps, devStats.Read, devStats.Write", devStats.Tps, devStats.Read, devStats.Write)

				stats.DevStats = append(stats.DevStats, &devStats)
			}

			// defencive programming: we assume slices are filled almost simultaneously
			// but despiting that fact, we will check if slice is ok before working
			for int(settings.AveragingTime) > (len(s.fsStatsSamples)) {
				fmt.Println("len s.fsStatsSamples", len(s.fsStatsSamples))
				time.Sleep(time.Duration(int(settings.AveragingTime)-len(s.fsStatsSamples)) * time.Second)
			}

			start = len(s.fsStatsSamples) - (int(settings.AveragingTime))

			if len(s.fsStatsSamples[start:]) < int(settings.AveragingTime) {
				return status.Errorf(codes.Aborted, "sending message error: fs statistics corrupted")
			}

			// limitation: we assume device number is constant. If not, there might be an error there
			for devIdx, dev := range s.fsStatsSamples[0] {
				fmt.Println("devIdx, devName", devIdx, dev.Name)
				fsStats := protobuf.FsStats{
					Name: dev.Name,
				}

				for idx, sample := range s.fsStatsSamples[start:] {

					spl := sample[devIdx]
					fmt.Println("idx, spl.Name", idx, spl.Name)
					if fsStats.Name != spl.Name {
						return status.Errorf(codes.Aborted, "sending message error: mixing different devices stats")
					}

					fsStats.Bytes += spl.UsedGBytes
					fsStats.BytesPercent += spl.UsedStoragePercent
					fsStats.Inode += spl.UsedInodes
					fsStats.InodePercent += spl.UsedInodesPercent
				}
				fmt.Println("fsStats.Tps, fsStats.Read, fsStats.Write", fsStats.Bytes, fsStats.BytesPercent, fsStats.Inode, fsStats.InodePercent)

				fsStats.Bytes /= float64(settings.AveragingTime)
				fsStats.BytesPercent /= float64(settings.AveragingTime)
				fsStats.Inode /= float64(settings.AveragingTime)
				fsStats.InodePercent /= float64(settings.AveragingTime)

				fmt.Println("fsStats.Tps, fsStats.Read, fsStats.Write", fsStats.Bytes, fsStats.BytesPercent, fsStats.Inode, fsStats.InodePercent)

				stats.FsStats = append(stats.FsStats, &fsStats)
			}

			if err := srv.Send(&stats); err != nil {
				return status.Errorf(codes.Aborted, "sending message error: %s", err.Error())
			}
			time.Sleep(time.Duration(settings.TimeBetweenTicks) * time.Second)
		}
	}
}

func main() {
	lis, err := net.Listen("tcp", "localhost:8088")
	if err != nil {
		fmt.Println("cannot listen TCP:", err.Error())
		return
	}

	s := grpc.NewServer()

	ctx, cancel := context.WithCancel(context.Background())
	cpuStatsSamples := make([]oslayer.CPUstats, 0)
	fsStatsSamples := make([][]oslayer.FsStats, 0)
	devStatsSamples := make([][]oslayer.DevStats, 0)

	protobuf.RegisterMonitorServer(s, &server{
		cpuStatsSamples: cpuStatsSamples,
		fsStatsSamples:  fsStatsSamples,
		devStatsSamples: devStatsSamples,
		ctx:             ctx,
		cancel:          cancel,
	})
	fmt.Println("server listening at", lis.Addr())
	if err := s.Serve(lis); err != nil {
		fmt.Println("cannot handle request:", err.Error())
		return
	}
}
