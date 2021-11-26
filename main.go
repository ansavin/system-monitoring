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
	for {
		select {
		case <-ctx.Done():
			return
		default:
			f()
		}
	}
}

func (s *server) getCPUStats(settings *protobuf.Settings, stats *protobuf.Stats) {
	avgTime := float64(settings.AveragingTime)

	// we should wait till enough data will be collected
	for int(settings.AveragingTime) > (len(s.cpuStatsSamples)) {
		time.Sleep(time.Duration(int(settings.AveragingTime)-len(s.cpuStatsSamples)) * time.Second)
	}

	// averaging our data for client according to his settings
	start := len(s.cpuStatsSamples) - (int(settings.AveragingTime))

	cpuStats := protobuf.CPUstats{}

	for _, sample := range s.cpuStatsSamples[start:] {
		cpuStats.La += sample.LA
		cpuStats.Sys += sample.SysUsagePercent
		cpuStats.Usr += sample.UsrUsagePercent
		cpuStats.Idle += sample.IdlePercent
	}
	cpuStats.La /= avgTime
	cpuStats.Sys /= avgTime
	cpuStats.Usr /= avgTime
	cpuStats.Idle /= avgTime

	stats.CPUstats = &cpuStats
}

func (s *server) getDevStats(settings *protobuf.Settings, stats *protobuf.Stats) error {
	avgTime := float64(settings.AveragingTime)

	// defensive programming: we assume slices are filled almost simultaneously
	// but despiting that fact, we will check if slice is ok before working
	for int(settings.AveragingTime) > (len(s.devStatsSamples)) {
		time.Sleep(time.Duration(int(settings.AveragingTime)-len(s.devStatsSamples)) * time.Second)
	}

	start := len(s.devStatsSamples) - (int(settings.AveragingTime))

	// here we definitely should have enough data
	if len(s.devStatsSamples[start:]) < int(settings.AveragingTime) {
		s.cancel()
		return status.Errorf(codes.DataLoss, "sending message error: device statistics corrupted")
	}

	// averaging our data for client according to his settings
	// limitation: we assume device number is constant. If not, there might be an error there
	for devIdx, dev := range s.devStatsSamples[0] {
		devStats := protobuf.DevStats{
			Name: dev.Name,
		}

		for _, sample := range s.devStatsSamples[start:] {

			spl := sample[devIdx]

			// we assume devices order is always the same, otherwise this won`t work
			if devStats.Name != spl.Name {
				s.cancel()
				return status.Errorf(codes.DataLoss, "sending message error: mixing different devices stats")
			}

			devStats.Tps += spl.TransPS
			devStats.Read += spl.ReadPS
			devStats.Write += spl.WritePS
		}

		devStats.Tps /= avgTime
		devStats.Read /= avgTime
		devStats.Write /= avgTime

		stats.DevStats = append(stats.DevStats, &devStats)
	}
	return nil
}

func (s *server) getFsStats(settings *protobuf.Settings, stats *protobuf.Stats) error {
	avgTime := float64(settings.AveragingTime)

	// defensive programming: we assume slices are filled almost simultaneously
	// but despiting that fact, we will check if slice is ok before working
	for int(settings.AveragingTime) > (len(s.fsStatsSamples)) {
		time.Sleep(time.Duration(int(settings.AveragingTime)-len(s.fsStatsSamples)) * time.Second)
	}

	start := len(s.fsStatsSamples) - (int(settings.AveragingTime))

	// here we definitely should have enough data
	if len(s.fsStatsSamples[start:]) < int(settings.AveragingTime) {
		s.cancel()
		return status.Errorf(codes.DataLoss, "sending message error: fs statistics corrupted")
	}
	// averaging our data for client according to his settings
	// limitation: we assume device number is constant. If not, there might be an error there
	for devIdx, dev := range s.fsStatsSamples[0] {
		fsStats := protobuf.FsStats{
			Name: dev.Name,
		}

		for _, sample := range s.fsStatsSamples[start:] {

			spl := sample[devIdx]

			// we assume fs order is always the same, otherwise this won`t work
			if fsStats.Name != spl.Name {
				s.cancel()
				return status.Errorf(codes.DataLoss, "sending message error: mixing different devices stats")
			}

			fsStats.Bytes += spl.UsedGBytes
			fsStats.BytesPercent += spl.UsedStoragePercent
			fsStats.Inode += spl.UsedInodes
			fsStats.InodePercent += spl.UsedInodesPercent
		}

		fsStats.Bytes /= avgTime
		fsStats.BytesPercent /= avgTime
		fsStats.Inode /= avgTime
		fsStats.InodePercent /= avgTime

		stats.FsStats = append(stats.FsStats, &fsStats)
	}
	return nil
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
	statsCalcFunctions := []func(){
		func() {
			cpu, err := oslayer.CalcCPUUsage()
			if err != nil {
				s.lastError = err
				s.cancel()
			}

			// here we want to preserve only maxSamples elements in slices
			// we have no generic in go, so we cant move this to a func
			// without overengineering
			if len(s.cpuStatsSamples) == maxSamples {
				s.cpuStatsSamples = s.cpuStatsSamples[1:]
			}
			s.cpuStatsSamples = append(s.cpuStatsSamples, cpu)

			time.Sleep(oslayer.SamplingTime)
		},
		func() {
			devs, err := oslayer.CalcDevStats()
			if err != nil {
				s.lastError = err
				s.cancel()
			}

			// here we want to preserve only maxSamples elements in slices
			// we have no generic in go, so we cant move this to a func
			// without overengineering
			if len(s.devStatsSamples) == maxSamples {
				s.devStatsSamples = s.devStatsSamples[1:]
			}
			s.devStatsSamples = append(s.devStatsSamples, devs)
			// we doesn`t need to sleep here - CalcDevStats use the same SamplingTime
			// for collecting stats
		},
		func() {
			fsystems, err := oslayer.CalcFsUtilization()
			if err != nil {
				s.lastError = err
				s.cancel()
			}
			// we have no generic in go, so we cant move this to a func
			// without overengineering
			if len(s.fsStatsSamples) == maxSamples {
				s.fsStatsSamples = s.fsStatsSamples[1:]
			}
			s.fsStatsSamples = append(s.fsStatsSamples, fsystems)
		},
	}

	// statistics is calculated in separate goroutines and stored in slices
	for idx, f := range statsCalcFunctions {
		s.once[idx].Do(func() {
			go startPeriodicalSampling(s.ctx, f)
		})
	}

	for {
		select {
		case <-s.ctx.Done():
			return status.Errorf(codes.Internal, "collecting statistics error: %s", s.lastError.Error())
		default:
			stats := protobuf.Stats{}

			s.getCPUStats(settings, &stats)

			err := s.getDevStats(settings, &stats)
			if err != nil {
				return err
			}

			err = s.getFsStats(settings, &stats)
			if err != nil {
				return err
			}

			if err := srv.Send(&stats); err != nil {
				s.cancel()
				return status.Errorf(codes.Internal, "sending message error: %s", err.Error())
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
		cancel()
		return
	}
}
