package memory

import (
	"github.com/NETWAYS/check_system_basics/internal/common/thresholds"
)

type MemConfig struct {
	MemAvailable           thresholds.Thresholds
	MemUsed                thresholds.Thresholds
	MemFree                thresholds.Thresholds
	MemAvailablePercentage thresholds.Thresholds
	MemUsedPercentage      thresholds.Thresholds
	MemFreePercentage      thresholds.Thresholds

	SwapUsed           thresholds.Thresholds
	SwapFree           thresholds.Thresholds
	SwapUsedPercentage thresholds.Thresholds
	SwapFreePercentage thresholds.Thresholds

	Verbose              bool
	PercentageInPerfdata bool
}
