package load

import (
	"github.com/NETWAYS/check_system_basics/internal/common/thresholds"
)

// nolint
type LoadConfig struct {
	Load1Th  thresholds.Thresholds
	Load5Th  thresholds.Thresholds
	Load15Th thresholds.Thresholds
	PerCPU   bool
}
