package cmd

import (
	"testing"

	"github.com/NETWAYS/check_system_basics/internal/common/thresholds"
	"github.com/NETWAYS/check_system_basics/internal/memory"
	"github.com/NETWAYS/go-check"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/stretchr/testify/assert"
)

var (
	testMemStats = memory.Mem{
		VirtMem: &mem.VirtualMemoryStat{
			Total:       512 * 1024, // 512 MiB
			Available:   256 * 1024,
			Used:        128 * 1024,
			UsedPercent: 25,
			Free:        128 * 1024,
		},
		MemAvailablePercentage: 50,
		SwapInfo:               &mem.SwapMemoryStat{},
	}
)

func TestComputeMemResultsWithoutThresholds(t *testing.T) {
	config := memory.MemConfig{
		MemAvailable: thresholds.Thresholds{
			Warn: thresholds.ThresholdWrapper{},
			Crit: thresholds.ThresholdWrapper{},
		},
	}

	memPartial := computeMemResults(&config, &testMemStats)

	assert.Equal(t, check.OK, memPartial.GetStatus())

	assert.Equal(t, 3, len(memPartial.PartialResults))
}

func TestComputeMemResultsWithThresholds(t *testing.T) {
	testConfig := memory.MemConfig{
		MemAvailablePercentage: thresholds.Thresholds{
			Warn: thresholds.ThresholdWrapper{
				IsSet: true,
				Th: check.Threshold{
					Inside: false,
					Lower:  0,
					Upper:  4,
				},
			},
			Crit: thresholds.ThresholdWrapper{},
		},
	}

	memPartial := computeMemResults(&testConfig, &testMemStats)

	assert.Equal(t, check.Warning, memPartial.GetStatus())

	assert.Equal(t, 3, len(memPartial.PartialResults))
}
