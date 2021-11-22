//go:build darwin && cgo
// +build darwin,cgo

package oslayer

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCalcCPUUsage(t *testing.T) {
	t.Run("simple positive test", func(t *testing.T) {
		data, err := CalcCPUUsage()

		require.NoError(t, err)
		require.NotEmpty(t, data)
	})
}
