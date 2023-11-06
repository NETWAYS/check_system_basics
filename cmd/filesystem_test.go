package cmd

import (
	"testing"

	"github.com/NETWAYS/check_system_basics/internal/common/thresholds"
	"github.com/NETWAYS/check_system_basics/internal/filesystem"
	"github.com/NETWAYS/go-check"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/stretchr/testify/assert"
)

var (
	testFsHalfFull = filesystem.FilesystemType{
		PartStats: disk.PartitionStat{
			Device:     "/dev/testfilesystem",
			Mountpoint: "/testMountpoint",
			Fstype:     "test",
			Opts:       []string{},
		},
		UsageStats: disk.UsageStat{
			Path:              "/testMountpoint",
			Fstype:            "test",
			Total:             4194304, // 1024 pages aka  4096 * 1024
			Free:              2097152, // 512 pages, -> 50%
			Used:              2097152,
			UsedPercent:       50,
			InodesTotal:       1024,
			InodesUsed:        512,
			InodesFree:        512,
			InodesUsedPercent: 50,
		},
	}
)

func TestFsCheckResult0(t *testing.T) {
	// No config given
	config := filesystem.CheckConfig{}

	result := computeFsCheckResult(&testFsHalfFull, &config)

	assert.Equal(t, check.OK, result.GetStatus())
	assert.Equal(t, "/testMountpoint", result.Output)
}

func TestFsCheckResult1(t *testing.T) {

	// this config should produce warning
	config := filesystem.CheckConfig{
		WarningAbsolutThreshold: filesystem.Thresholds{
			Space: filesystem.DualThresholdWrapper{
				Free: thresholds.ThresholdWrapper{
					Th: check.Threshold{
						Inside: true,
						Lower:  2097152,
						Upper:  4194304,
					},
					IsSet: true,
				},
			},
		},
	}

	result := computeFsCheckResult(&testFsHalfFull, &config)

	assert.Equal(t, check.Warning, result.GetStatus())
}
