package filesystem

import (
	"github.com/NETWAYS/check_system_basics/internal/common/thresholds"
)

type DualThresholdWrapper struct {
	Free thresholds.ThresholdWrapper
	Used thresholds.ThresholdWrapper
}

type Thresholds struct {
	Space  DualThresholdWrapper
	Inodes DualThresholdWrapper
}

type Filters struct {
	// Selection/File options
	//  Paths
	IncludeDevicePaths []string
	ExcludeDevicePaths []string

	//  Filesystem Types
	IncludeFsType []string
	ExcludeFsType []string

	// Mount Paths
	IncludeMountPaths []string
	ExcludeMountPaths []string
	// Filesystem Options
	//  Read-Only? Read-Write
	//  pass through general mount options?
	IncludeOptions []string
	ExcludeOptions []string
}

type CheckConfig struct {
	// Thresholds
	WarningAbsolutThreshold  Thresholds
	CriticalAbsolutThreshold Thresholds

	WarningPercentThreshold  Thresholds
	CriticalPercentThreshold Thresholds

	WarningTotalCountOfFs  thresholds.ThresholdWrapper
	CriticalTotalCountOfFs thresholds.ThresholdWrapper

	Filters Filters

	ReadonlyOption  bool
	ReadWriteOption bool

	// Output Verbosity
	Verbosity uint
}

func GetFilesystemsWithFixedNumberOfInodes() []string {
	return []string{
		"bfs",
		"ext2",
		"ext3",
		"ext4",
	}
}
