//go:build linux
// +build linux

package oslayer

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPercentage(t *testing.T) {
	tests := []struct {
		x        float64
		y        float64
		expected float64
	}{
		{x: 100, y: 10, expected: 90},
		{x: 100, y: 0, expected: 100},
		{x: 0, y: 10, expected: 0},
		{x: 100.10, y: 10, expected: 90.00999000999002},
		// we assume we never get negative number here
		// (what is negative inodes count? O_O)
		// this is why no additional hints here
		{x: 100, y: -10, expected: 110},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(fmt.Sprintf("%+v", tc), func(t *testing.T) {
			result := percentage(tc.x, tc.y)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestCheckIfRunsInDocker(t *testing.T) {
	t.Run("simple postitve test", func(t *testing.T) {
		// path that definately exist in linux
		tmp := dockerRootFSPrefix
		dockerRootFSPrefix = "/proc"
		defer func() { dockerRootFSPrefix = tmp }()

		expected := checkIfRunsInDocker()
		require.Equal(t, expected, "/proc")
	})

	t.Run("simple negative test", func(t *testing.T) {
		tmp := dockerRootFSPrefix
		dockerRootFSPrefix = "/path/to/nowhere"
		defer func() { dockerRootFSPrefix = tmp }()

		expected := checkIfRunsInDocker()
		require.Equal(t, expected, "")
	})
}

// no need to test internal func of package linux like getDevStats (except parseDevStats maybe)
// because this tests will be just Ctlr+C - Ctrl+V of tests for functions CalcDevStats
// unless we would mock sysfs & procfs filesystem

func TestParseDevStats(t *testing.T) {
	t.Run("simple postitve test", func(t *testing.T) {
		devs, err := ioutil.ReadDir(checkIfRunsInDocker() + BlockDevicesDir)
		require.NoError(t, err)
		for _, d := range devs {
			_, err := parseDevStats(d.Name(), BlockDevicesDir)
			require.NoError(t, err)
		}
	})

	t.Run("simple negative test", func(t *testing.T) {
		dir := "/path/to/nowhere"

		_, err := parseDevStats("not_exist_file", dir)
		require.Error(t, err)
	})
}

func TestCalcDevStats(t *testing.T) {
	t.Run("simple positive test", func(t *testing.T) {
		data, err := CalcDevStats()

		require.NoError(t, err)
		require.NotEmpty(t, data)
	})

	t.Run("negative test - broken dockerRootFSPrefix", func(t *testing.T) {
		tmp := dockerRootFSPrefix
		dockerRootFSPrefix = "/proc"
		defer func() { dockerRootFSPrefix = tmp }()

		_, err := CalcDevStats()

		require.Error(t, err)
	})

	t.Run("negative test - broken BlockDevicesDir", func(t *testing.T) {
		tmp := BlockDevicesDir
		BlockDevicesDir = "/path/to/nowhere"
		defer func() { BlockDevicesDir = tmp }()

		_, err := CalcDevStats()

		require.Error(t, err)
	})

	t.Run("negative test - broken DevStatsFilename", func(t *testing.T) {
		tmp := DevStatsFilename
		DevStatsFilename = "/path/to/nowhere"
		defer func() { DevStatsFilename = tmp }()

		_, err := CalcDevStats()

		require.Error(t, err)
	})
}

func TestCalcFsUtilisation(t *testing.T) {
	t.Run("simple positive test", func(t *testing.T) {
		data, err := CalcFsUtilisation()
		// FIXME cannot open /boot/efi fs: open /boot/efi: permission denied
		// when not in docker

		require.NoError(t, err)
		require.NotEmpty(t, data)
	})

	t.Run("negative test - broken dockerRootFSPrefix", func(t *testing.T) {
		tmp := dockerRootFSPrefix
		dockerRootFSPrefix = "/proc"
		defer func() { dockerRootFSPrefix = tmp }()

		_, err := CalcFsUtilisation()

		require.Error(t, err)
	})

	t.Run("negative test - broken MountinfoFile", func(t *testing.T) {
		tmp := MountinfoFile
		MountinfoFile = "/path/to/nowhere"
		defer func() { MountinfoFile = tmp }()

		_, err := CalcFsUtilisation()

		require.Error(t, err)
	})
}

func TestCalcCPUUsage(t *testing.T) {
	t.Run("simple positive test", func(t *testing.T) {
		data, err := CalcCPUUsage()

		require.NoError(t, err)
		require.NotEmpty(t, data)
	})

	t.Run("negative test - broken dockerRootFSPrefix", func(t *testing.T) {
		tmp := dockerRootFSPrefix
		dockerRootFSPrefix = "/proc"
		defer func() { dockerRootFSPrefix = tmp }()

		_, err := CalcCPUUsage()

		require.Error(t, err)
	})

	t.Run("negative test - broken LaFile", func(t *testing.T) {
		tmp := LaFile
		LaFile = "/path/to/nowhere"
		defer func() { LaFile = tmp }()

		_, err := CalcCPUUsage()

		require.Error(t, err)
	})

	t.Run("negative test - broken CPUStatsFile", func(t *testing.T) {
		tmp := CPUStatsFile
		CPUStatsFile = "/path/to/nowhere"
		defer func() { CPUStatsFile = tmp }()

		_, err := CalcCPUUsage()

		require.Error(t, err)
	})
}