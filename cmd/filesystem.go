package cmd

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/NETWAYS/check_system_basics/internal/common/thresholds"
	"github.com/NETWAYS/check_system_basics/internal/filesystem"
	"github.com/NETWAYS/go-check"
	"github.com/NETWAYS/go-check/perfdata"
	"github.com/NETWAYS/go-check/result"
	"github.com/dustin/go-humanize"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/spf13/cobra"
)

var FsConfig filesystem.CheckConfig

var IncludeFsTypeDefaults = []string{
	"^ext2$",
	"^ext3$",
	"^ext4$",
	"^btrfs$",
	"^nfs$",
	"^ntfs$",
	"^reiserfs$",
	"^xfs$",
	"^zfs$",
}

var diskCmd = &cobra.Command{
	Use:   "filesystem",
	Short: "Submodule to check the usage of mounted file systems.",
	Example: `./check_system_basics filesystem --criticalPercentFreeSpace 0:30 --criticalAbsolutFreeSpace @0:20 --warningAbsolutFreeSpace @0:40 --criticalAbsolutUsedSpace 0:20 --warningAbsolutUsedSpace0:40 --warningPercentFreeSpace 0:60 --criticalPercentUsedSpace @0:20 --warningPercentUsedSpace @20:40 --warningAbsolutFreeInodes @0:200 --criticalAbsolutUsedInodes    @0:400
[CRITICAL] - states: critical=4
\_ [CRITICAL] / (44.25% used space, 82.69% free inodes)
		\_ [OK] Absolute free space: 31 GiB / 58 GiB
		\_ [CRITICAL] Absolute used space violates threshold: 24 GiB / 58 GiB
		\_ [OK] Absolute number of free inodes: 3231100 / 3907584
		\_ [OK] Absolute number of used inodes: 676484 / 3907584
		\_ [CRITICAL] Percentage of free space violates threshold: 55.75%
		\_ [OK] Percentage of used space: 44.25%
		\_ [OK] Percentage of used inodes: 17.31%
\_ [CRITICAL] /boot (62.98% used space, 99.72% free inodes)
		\_ [OK] Absolute free space: 165 MiB / 470 MiB
		\_ [CRITICAL] Absolute used space violates threshold: 281 MiB / 470 MiB
		\_ [OK] Absolute number of free inodes: 124576 / 124928
		\_ [CRITICAL] Absolute used inode number violates threshold: 352 / 124928
		\_ [CRITICAL] Percentage of free space violates threshold: 37.02%
		\_ [OK] Percentage of used space: 62.98%
		\_ [OK] Percentage of used inodes: 0.28%
\_ [CRITICAL] /var (40.42% used space, 92.83% free inodes)
		\_ [OK] Absolute free space: 132 GiB / 234 GiB
		\_ [CRITICAL] Absolute used space violates threshold: 90 GiB / 234 GiB
		\_ [OK] Absolute number of free inodes: 14510026 / 15630336
		\_ [OK] Absolute number of used inodes: 1120310 / 15630336
		\_ [CRITICAL] Percentage of free space violates threshold: 59.58%
		\_ [OK] Percentage of used space: 40.42%
		\_ [OK] Percentage of used inodes: 7.17%
\_ [CRITICAL] /home (48.05% used space, 94.77% free inodes)
		\_ [OK] Absolute free space: 231 GiB / 468 GiB
		\_ [CRITICAL] Absolute used space violates threshold: 214 GiB / 468 GiB
		\_ [OK] Absolute number of free inodes: 29617311 / 31252480
		\_ [OK] Absolute number of used inodes: 1635169 / 31252480
		\_ [CRITICAL] Percentage of free space violates threshold: 51.95%
		\_ [OK] Percentage of used space: 48.05%
		\_ [OK] Percentage of used inodes: 5.23%
|/_inodes_free_percentage=82.688% /_space_free=33146855424B;@40;@20;0;62669000704 /_space_used=26305536000B;40;20;0;62669000704 /_inodes_free=3231100;;@200;0;3907584 /_inodes_used=676484;;@400;0;3907584 /_space_free_percentage=55.754%;60;30 /_space_used_percentage=44.246%;@20:40;@20 /_inodes_used_percentage=17.312%;99;98 /boot_inodes_free_percentage=99.718% /boot_space_free=173091840B;@40;@20;0;493201408 /boot_space_used=294524928B;40;20;0;493201408 /boot_inodes_free=124576;;@200;0;124928 /boot_inodes_used=352;;@400;0;124928 /boot_space_free_percentage=37.016%;60;30 /boot_space_used_percentage=62.984%;@20:40;@20 /boot_inodes_used_percentage=0.282%;99;98 /var_inodes_free_percentage=92.832% /var_space_free=141826428928B;@40;@20;0;250843787264 /var_space_used=96200613888B;40;20;0;250843787264 /var_inodes_free=14510026;;@200;0;15630336 /var_inodes_used=1120310;;@400;0;15630336 /var_space_free_percentage=59.584%;60;30 /var_space_used_percentage=40.416%;@20:40;@20 /var_inodes_used_percentage=7.168%;99;98 /home_inodes_free_percentage=94.768% /home_space_free=247921197056B;@40;@20;0;502813065216 /home_space_used=229275156480B;40;20;0;502813065216 /home_inodes_free=29617311;;@200;0;31252480 /home_inodes_used=1635169;;@400;0;31252480 /home_space_free_percentage=51.954%;60;30 /home_space_used_percentage=48.046%;@20:40;@20 /home_inodes_used_percentage=5.232%;99;98`,
	Run: func(_ *cobra.Command, args []string) {

		overall := result.Overall{}

		err := validateOptions(&FsConfig)
		if err != nil {
			check.ExitError(err)
		}

		// Detect file systems
		filesystems, err := disk.Partitions(true)
		if err != nil {
			check.ExitError(err)
		}

		if debug {
			fmt.Printf("==== Detected filesystems: ====\n %v\n", filesystems)
		}

		filesystemList := make([]filesystem.FilesystemType, len(filesystems))

		for i := range filesystems {
			filesystemList[i].PartStats = filesystems[i]
		}

		if debug {
			fmt.Printf("==== Filesystem List: ====\n %v\n", filesystemList)
		}

		// Filter out unwanted
		filesystemList, err = filesystem.FilterFileSystem(filesystemList, &FsConfig.Filters)
		if err != nil {
			check.ExitError(err)
		}

		if debug {
			fmt.Printf("==== Filtered Filesystem List: ====\n %v\n", filesystemList)
		}

		if FsConfig.CriticalTotalCountOfFs.IsSet || FsConfig.WarningTotalCountOfFs.IsSet {
			countResult := result.PartialResult{}
			_ = countResult.SetDefaultState(check.OK)

			if len(filesystemList) == 1 {
				countResult.Output = "Found one matching filesystem"
			} else {
				countResult.Output = "Found " + strconv.Itoa(len(filesystemList)) + " matching filesystems"
			}

			if FsConfig.CriticalTotalCountOfFs.IsSet && FsConfig.CriticalTotalCountOfFs.Th.DoesViolate(float64(len(filesystemList))) {
				_ = countResult.SetState(check.Critical)
				countResult.Output += ". This violates the threshold of " + FsConfig.CriticalTotalCountOfFs.String()
			} else if FsConfig.WarningTotalCountOfFs.IsSet && FsConfig.WarningTotalCountOfFs.Th.DoesViolate(float64(len(filesystemList))) {
				_ = countResult.SetState(check.Warning)
				countResult.Output += ". This violates the threshold of " + FsConfig.WarningTotalCountOfFs.String()
			} else {
				countResult.Output += ". This number resides within the given thresholds"
			}

			overall.AddSubcheck(countResult)

		} else if len(filesystemList) == 0 {
			nullResult := result.PartialResult{}
			_ = nullResult.SetState(check.OK)
			nullResult.Output = "No filesystems remaining after applying filter expressions. Therefore all are OK"
			overall.AddSubcheck(nullResult)
			check.ExitRaw(overall.GetStatus(), overall.GetOutput())
		}

		// Retrieve stats
		internalTimeout := time.Duration(Timeout/2) * time.Second
		ctx := context.Background()

		err = filesystem.GetDiskUsage(ctx, internalTimeout, filesystemList, &FsConfig)
		if err != nil {
			check.ExitError(err)
		}

		// Compile the result
		for index := range filesystemList {
			sc := computeFsCheckResult(&filesystemList[index], &FsConfig)

			if filesystemList[index].Error == nil {
				sc.Output = fmt.Sprintf("%s (%.2f%% used space, %.2f%% free inodes)", sc.Output, filesystemList[index].UsageStats.UsedPercent, 100-filesystemList[index].UsageStats.InodesUsedPercent)
			}
			overall.AddSubcheck(sc)
		}

		// Output and Exit
		check.ExitRaw(overall.GetStatus(), overall.GetOutput())

	},
}

func computeFsCheckResultInodes(fs *filesystem.FilesystemType, config *filesystem.CheckConfig) result.PartialResult {
	returnResult := result.PartialResult{
		Output: "Inodes",
	}
	_ = returnResult.SetDefaultState(check.OK)

	// One Perfdata point here with inodes free, warn, crit, total
	pdAbsoluteFreeInodes := perfdata.Perfdata{
		Min:   0,
		Max:   fs.UsageStats.InodesTotal,
		Uom:   "",
		Label: fs.PartStats.Mountpoint + "_inodes_free",
		Value: fs.UsageStats.InodesFree,
	}

	if config.WarningAbsolutThreshold.Inodes.Free.IsSet || config.CriticalAbsolutThreshold.Inodes.Free.IsSet {
		tmpPartialResult := result.PartialResult{}
		_ = tmpPartialResult.SetDefaultState(check.OK)

		if config.WarningAbsolutThreshold.Inodes.Free.IsSet {
			pdAbsoluteFreeInodes.Warn = &config.WarningAbsolutThreshold.Inodes.Free.Th

			if config.WarningAbsolutThreshold.Inodes.Free.Th.DoesViolate(float64(fs.UsageStats.InodesFree)) {
				_ = tmpPartialResult.SetState(check.Warning)
				tmpPartialResult.Output = fmt.Sprintf("Absolute free inode number violates threshold: %d / %d", fs.UsageStats.InodesFree, fs.UsageStats.InodesTotal)
			}
		}

		if config.CriticalAbsolutThreshold.Inodes.Free.IsSet {
			pdAbsoluteFreeInodes.Crit = &config.CriticalAbsolutThreshold.Inodes.Free.Th

			if config.CriticalAbsolutThreshold.Inodes.Free.Th.DoesViolate(float64(fs.UsageStats.InodesFree)) {
				_ = tmpPartialResult.SetState(check.Critical)
				tmpPartialResult.Output = fmt.Sprintf("Absolute free inode number violates threshold: %d / %d", fs.UsageStats.InodesFree, fs.UsageStats.InodesTotal)
			}
		}

		if tmpPartialResult.GetStatus() == check.OK {
			tmpPartialResult.Output = fmt.Sprintf("Absolute number of free inodes: %d / %d", fs.UsageStats.InodesFree, fs.UsageStats.InodesTotal)
		}

		tmpPartialResult.Perfdata.Add(&pdAbsoluteFreeInodes)
		returnResult.AddSubcheck(tmpPartialResult)
	} else {
		returnResult.Perfdata.Add(&pdAbsoluteFreeInodes)
	}

	// One Perfdata point here with inodes used, warn, crit, total
	pdAbsoluteUsedInodes := perfdata.Perfdata{
		Min:   0,
		Max:   fs.UsageStats.InodesTotal,
		Uom:   "",
		Label: fs.PartStats.Mountpoint + "_inodes_used",
		Value: fs.UsageStats.InodesUsed,
	}

	if config.WarningAbsolutThreshold.Inodes.Used.IsSet || config.CriticalAbsolutThreshold.Inodes.Used.IsSet {
		tmpPartialResult := result.PartialResult{}
		_ = tmpPartialResult.SetDefaultState(check.OK)

		if config.WarningAbsolutThreshold.Inodes.Used.IsSet {
			pdAbsoluteUsedInodes.Warn = &config.WarningAbsolutThreshold.Inodes.Used.Th

			if config.WarningAbsolutThreshold.Inodes.Used.Th.DoesViolate(float64(fs.UsageStats.InodesUsed)) {
				_ = tmpPartialResult.SetState(check.Warning)
				tmpPartialResult.Output = fmt.Sprintf("Absolute used inode number violates threshold: %d / %d", fs.UsageStats.InodesUsed, fs.UsageStats.InodesTotal)
			}
		}

		if config.CriticalAbsolutThreshold.Inodes.Used.IsSet {
			pdAbsoluteUsedInodes.Crit = &config.CriticalAbsolutThreshold.Inodes.Used.Th

			if config.CriticalAbsolutThreshold.Inodes.Used.Th.DoesViolate(float64(fs.UsageStats.InodesUsed)) {
				_ = tmpPartialResult.SetState(check.Critical)
				tmpPartialResult.Output = fmt.Sprintf("Absolute used inode number violates threshold: %d / %d", fs.UsageStats.InodesUsed, fs.UsageStats.InodesTotal)
			}
		}

		if tmpPartialResult.GetStatus() == check.OK {
			tmpPartialResult.Output = fmt.Sprintf("Absolute number of used inodes: %d / %d", fs.UsageStats.InodesUsed, fs.UsageStats.InodesTotal)
		}

		tmpPartialResult.Perfdata.Add(&pdAbsoluteUsedInodes)
		returnResult.AddSubcheck(tmpPartialResult)
	} else {
		returnResult.Perfdata.Add(&pdAbsoluteUsedInodes)
	}

	// One Perfdata point here with inodes free, warn, crit, total
	pdPercentageFreeInodes := perfdata.Perfdata{
		Uom:   "%",
		Label: fs.PartStats.Mountpoint + "_inodes_free_percentage",
		Value: 100 - fs.UsageStats.InodesUsedPercent,
	}

	if config.WarningPercentThreshold.Inodes.Free.IsSet || config.CriticalPercentThreshold.Inodes.Free.IsSet {
		tmpPartialResult := result.PartialResult{}
		_ = tmpPartialResult.SetDefaultState(check.OK)

		if config.WarningPercentThreshold.Inodes.Free.IsSet {
			pdPercentageFreeInodes.Warn = &config.WarningPercentThreshold.Inodes.Free.Th

			if config.WarningPercentThreshold.Inodes.Free.Th.DoesViolate(pdPercentageFreeInodes.Value.(float64)) {
				_ = tmpPartialResult.SetState(check.Warning)
				tmpPartialResult.Output = fmt.Sprintf("Percentage of free inodes violates threshold: %.2f%%", pdPercentageFreeInodes.Value)
			}
		}

		if config.CriticalPercentThreshold.Inodes.Free.IsSet {
			pdPercentageFreeInodes.Warn = &config.CriticalPercentThreshold.Inodes.Free.Th

			if config.CriticalPercentThreshold.Inodes.Free.Th.DoesViolate(pdPercentageFreeInodes.Value.(float64)) {
				_ = tmpPartialResult.SetState(check.Critical)
				tmpPartialResult.Output = fmt.Sprintf("Percentage of free inodes violates threshold: %.2f%%", pdPercentageFreeInodes.Value)
			}
		}

		if tmpPartialResult.GetStatus() == check.OK {
			tmpPartialResult.Output = fmt.Sprintf("Percentage of free inodes: %.2f%%", pdPercentageFreeInodes.Value)
		}

		tmpPartialResult.Perfdata.Add(&pdPercentageFreeInodes)
		returnResult.AddSubcheck(tmpPartialResult)
	} else {
		returnResult.Perfdata.Add(&pdPercentageFreeInodes)
	}

	// One Perfdata point here with inodes used, warn, crit, total
	pdPercentageUsedInodes := perfdata.Perfdata{
		Uom:   "%",
		Label: fs.PartStats.Mountpoint + "_inodes_used_percentage",
		Value: fs.UsageStats.InodesUsedPercent,
	}

	if config.WarningPercentThreshold.Inodes.Used.IsSet || config.CriticalPercentThreshold.Inodes.Used.IsSet {
		tmpPartialResult := result.PartialResult{}
		_ = tmpPartialResult.SetDefaultState(check.OK)

		if config.WarningPercentThreshold.Inodes.Used.IsSet {
			pdPercentageUsedInodes.Warn = &config.WarningPercentThreshold.Inodes.Used.Th

			if config.WarningPercentThreshold.Inodes.Used.Th.DoesViolate(fs.UsageStats.InodesUsedPercent) {
				_ = tmpPartialResult.SetState(check.Warning)
				tmpPartialResult.Output = fmt.Sprintf("Percentage of used inodes violates threshold: %.2f%%", fs.UsageStats.InodesUsedPercent)
			}
		}

		if config.CriticalPercentThreshold.Inodes.Used.IsSet {
			pdPercentageUsedInodes.Crit = &config.CriticalPercentThreshold.Inodes.Used.Th

			if config.CriticalPercentThreshold.Inodes.Used.Th.DoesViolate(fs.UsageStats.InodesUsedPercent) {
				_ = tmpPartialResult.SetState(check.Critical)
				tmpPartialResult.Output = fmt.Sprintf("Percentage of used inodes violates threshold: %.2f%%", fs.UsageStats.InodesUsedPercent)
			}
		}

		if tmpPartialResult.GetStatus() == check.OK {
			tmpPartialResult.Output = fmt.Sprintf("Percentage of used inodes: %.2f%%", fs.UsageStats.InodesUsedPercent)
		}

		tmpPartialResult.Perfdata.Add(&pdPercentageUsedInodes)
		returnResult.AddSubcheck(tmpPartialResult)
	} else {
		returnResult.Perfdata.Add(&pdPercentageUsedInodes)
	}

	return returnResult
}

func computeFsCheckResultSpace(fs *filesystem.FilesystemType, config *filesystem.CheckConfig) result.PartialResult {
	returnResult := result.PartialResult{
		Output: "Space usage",
	}
	_ = returnResult.SetDefaultState(check.OK)

	// Absolute numbers
	// One Perfdata point here with bytes free, warn, crit, total
	pdAbsoluteFreeSpace := perfdata.Perfdata{
		Min:   0,
		Max:   fs.UsageStats.Total,
		Uom:   "B",
		Label: fs.PartStats.Mountpoint + "_space_free",
		Value: fs.UsageStats.Free,
	}

	if config.WarningAbsolutThreshold.Space.Free.IsSet || config.CriticalAbsolutThreshold.Space.Free.IsSet {
		tmpPartialResult := result.PartialResult{}
		_ = tmpPartialResult.SetDefaultState(check.OK)

		if config.WarningAbsolutThreshold.Space.Free.IsSet {
			pdAbsoluteFreeSpace.Warn = &config.WarningAbsolutThreshold.Space.Free.Th

			if config.WarningAbsolutThreshold.Space.Free.Th.DoesViolate(float64(fs.UsageStats.Free)) {
				_ = tmpPartialResult.SetState(check.Warning)
				tmpPartialResult.Output = fmt.Sprintf("Absolute free space violates threshold: %s / %s", humanize.IBytes(fs.UsageStats.Free), humanize.IBytes(fs.UsageStats.Total))
			}
		}

		if config.CriticalAbsolutThreshold.Space.Free.IsSet {
			pdAbsoluteFreeSpace.Crit = &config.CriticalAbsolutThreshold.Space.Free.Th

			if config.CriticalAbsolutThreshold.Space.Free.Th.DoesViolate(float64(fs.UsageStats.Free)) {
				_ = tmpPartialResult.SetState(check.Critical)
				tmpPartialResult.Output = fmt.Sprintf("Absolute free space violates threshold: %s / %s", humanize.IBytes(fs.UsageStats.Free), humanize.IBytes(fs.UsageStats.Total))
			}
		}

		if tmpPartialResult.GetStatus() == check.OK {
			tmpPartialResult.Output = fmt.Sprintf("Absolute free space: %s / %s", humanize.IBytes(fs.UsageStats.Free), humanize.IBytes(fs.UsageStats.Total))
		}

		tmpPartialResult.Perfdata.Add(&pdAbsoluteFreeSpace)
		returnResult.AddSubcheck(tmpPartialResult)
	} else {
		returnResult.Perfdata.Add(&pdAbsoluteFreeSpace)
	}

	// One Perfdata point here with bytes used, warn, crit, total
	pdAbsoluteUsedSpace := perfdata.Perfdata{
		Min:   0,
		Max:   fs.UsageStats.Total,
		Uom:   "B",
		Label: fs.PartStats.Mountpoint + "_space_used",
		Value: fs.UsageStats.Used,
	}

	if config.WarningAbsolutThreshold.Space.Used.IsSet || config.CriticalAbsolutThreshold.Space.Used.IsSet {
		tmpPartialResult := result.PartialResult{}
		_ = tmpPartialResult.SetDefaultState(check.OK)

		if config.WarningAbsolutThreshold.Space.Used.IsSet {
			pdAbsoluteUsedSpace.Warn = &config.WarningAbsolutThreshold.Space.Used.Th

			if config.WarningAbsolutThreshold.Space.Used.Th.DoesViolate(float64(fs.UsageStats.Used)) {
				_ = tmpPartialResult.SetState(check.Warning)
				tmpPartialResult.Output = fmt.Sprintf("Absolute used space violates threshold: %s / %s", humanize.IBytes(fs.UsageStats.Used), humanize.IBytes(fs.UsageStats.Total))
			}
		}

		if config.CriticalAbsolutThreshold.Space.Used.IsSet {
			pdAbsoluteUsedSpace.Crit = &config.CriticalAbsolutThreshold.Space.Used.Th

			if config.CriticalAbsolutThreshold.Space.Used.Th.DoesViolate(float64(fs.UsageStats.Used)) {
				_ = tmpPartialResult.SetState(check.Critical)
				tmpPartialResult.Output = fmt.Sprintf("Absolute used space violates threshold: %s / %s", humanize.IBytes(fs.UsageStats.Used), humanize.IBytes(fs.UsageStats.Total))
			}
		}

		if tmpPartialResult.GetStatus() == check.OK {
			tmpPartialResult.Output = fmt.Sprintf("Absolute used space: %s / %s", humanize.IBytes(fs.UsageStats.Used), humanize.IBytes(fs.UsageStats.Total))
		}

		tmpPartialResult.Perfdata.Add(&pdAbsoluteUsedSpace)
		returnResult.AddSubcheck(tmpPartialResult)
	} else {
		returnResult.Perfdata.Add(&pdAbsoluteUsedSpace)
	}

	// Percentage numbers
	//  Space

	// One Perfdata point here with bytes free, warn, crit, total
	pdPercentageFreeSpace := perfdata.Perfdata{
		Uom:   "%",
		Label: fs.PartStats.Mountpoint + "_space_free_percentage",
		Value: 100 - fs.UsageStats.UsedPercent,
	}

	if config.WarningPercentThreshold.Space.Free.IsSet || config.CriticalPercentThreshold.Space.Free.IsSet {
		tmpPartialResult := result.PartialResult{}
		_ = tmpPartialResult.SetDefaultState(check.OK)

		if config.WarningPercentThreshold.Space.Free.IsSet {
			pdPercentageFreeSpace.Warn = &config.WarningPercentThreshold.Space.Free.Th

			if config.WarningPercentThreshold.Space.Free.Th.DoesViolate(pdPercentageFreeSpace.Value.(float64)) {
				_ = tmpPartialResult.SetState(check.Warning)
				tmpPartialResult.Output = fmt.Sprintf("Percentage of free space violates threshold: %.2f%%", pdPercentageFreeSpace.Value)
			}
		}

		if config.CriticalPercentThreshold.Space.Free.IsSet {
			pdPercentageFreeSpace.Crit = &config.CriticalPercentThreshold.Space.Free.Th

			if config.CriticalPercentThreshold.Space.Free.Th.DoesViolate(pdPercentageFreeSpace.Value.(float64)) {
				_ = tmpPartialResult.SetState(check.Critical)
				tmpPartialResult.Output = fmt.Sprintf("Percentage of free space violates threshold: %.2f%%", pdPercentageFreeSpace.Value)
			}
		}

		if tmpPartialResult.GetStatus() == check.OK {
			tmpPartialResult.Output = fmt.Sprintf("Percentage of free space: %.2f%%", pdPercentageFreeSpace.Value)
		}

		tmpPartialResult.Perfdata.Add(&pdPercentageFreeSpace)
		returnResult.AddSubcheck(tmpPartialResult)
	} else {
		returnResult.Perfdata.Add(&pdPercentageFreeSpace)
	}

	// One Perfdata point here with bytes used, warn, crit, total
	pdPercentageUsedSpace := perfdata.Perfdata{
		Uom:   "%",
		Label: fs.PartStats.Mountpoint + "_space_used_percentage",
		Value: fs.UsageStats.UsedPercent,
	}

	if config.WarningPercentThreshold.Space.Used.IsSet || config.CriticalPercentThreshold.Space.Used.IsSet {
		tmpPartialResult := result.PartialResult{}
		_ = tmpPartialResult.SetDefaultState(check.OK)

		if config.WarningPercentThreshold.Space.Used.IsSet {
			pdPercentageUsedSpace.Warn = &config.WarningPercentThreshold.Space.Used.Th

			if config.WarningPercentThreshold.Space.Used.Th.DoesViolate(fs.UsageStats.UsedPercent) {
				_ = tmpPartialResult.SetState(check.Warning)
				tmpPartialResult.Output = fmt.Sprintf("Percentage of used space violates threshold: %.2f%%", fs.UsageStats.UsedPercent)
			}
		}

		if config.CriticalPercentThreshold.Space.Used.IsSet {
			pdPercentageUsedSpace.Crit = &config.CriticalPercentThreshold.Space.Used.Th

			if config.CriticalPercentThreshold.Space.Used.Th.DoesViolate(fs.UsageStats.UsedPercent) {
				_ = tmpPartialResult.SetState(check.Critical)
				tmpPartialResult.Output = fmt.Sprintf("Percentage of used space violates threshold: %.2f%%", fs.UsageStats.UsedPercent)
			}
		}

		if tmpPartialResult.GetStatus() == check.OK {
			tmpPartialResult.Output = fmt.Sprintf("Percentage of used space: %.2f%%", fs.UsageStats.UsedPercent)
		}

		tmpPartialResult.Perfdata.Add(&pdPercentageUsedSpace)
		returnResult.AddSubcheck(tmpPartialResult)
	} else {
		returnResult.Perfdata.Add(&pdPercentageUsedSpace)
	}

	return returnResult
}

func computeFsCheckResult(fs *filesystem.FilesystemType, config *filesystem.CheckConfig) result.PartialResult {
	returnResult := result.PartialResult{}
	returnResult.Output = fs.PartStats.Mountpoint
	_ = returnResult.SetDefaultState(check.OK)

	if fs.Error != nil {
		_ = returnResult.SetState(check.Unknown)
		returnResult.Output = fmt.Sprintf("Could not determine status of the filesystem  mounted at %s (%s) stats due to: %s", fs.PartStats.Mountpoint, fs.PartStats.Device, fs.Error)

		return returnResult
	}

	returnResult.AddSubcheck(computeFsCheckResultSpace(fs, config))

	filesystemsWithFixedNumberOfInodes := filesystem.GetFilesystemsWithFixedNumberOfInodes()

	for i := range filesystemsWithFixedNumberOfInodes {
		if fs.PartStats.Fstype == filesystemsWithFixedNumberOfInodes[i] {
			returnResult.AddSubcheck(computeFsCheckResultInodes(fs, config))
			break
		}
	}

	return returnResult
}

// nolint: funlen
func init() {
	rootCmd.AddCommand(diskCmd)

	fsThresholds := []thresholds.ThresholdOption{
		{
			Th:          &FsConfig.CriticalAbsolutThreshold.Space.Free,
			FlagString:  "criticalAbsolutFreeSpace",
			Description: "Absolute critical threshold for free filesystem space.",
		},
		{
			Th:          &FsConfig.WarningAbsolutThreshold.Space.Free,
			FlagString:  "warningAbsolutFreeSpace",
			Description: "Absolute warning threshold for free filesystem space.",
		},
		{
			Th:          &FsConfig.CriticalAbsolutThreshold.Space.Used,
			FlagString:  "criticalAbsolutUsedSpace",
			Description: "Absolute critical threshold for used filesystem space.",
		},
		{
			Th:          &FsConfig.WarningAbsolutThreshold.Space.Used,
			FlagString:  "warningAbsolutUsedSpace",
			Description: "Absolute warning threshold for used filesystem space.",
		},
		{
			Th:          &FsConfig.CriticalPercentThreshold.Space.Free,
			FlagString:  "criticalPercentFreeSpace",
			Description: "Percentage critical threshold for free filesystem space.",
			Default: thresholds.ThresholdWrapper{
				IsSet: true,
				Th: check.Threshold{
					Lower: 2,
					Upper: 100,
				},
			},
		},
		{
			Th:          &FsConfig.WarningPercentThreshold.Space.Free,
			FlagString:  "warningPercentFreeSpace",
			Description: "Percentage warning threshold for free filesystem space.",
			Default: thresholds.ThresholdWrapper{
				IsSet: true,
				Th: check.Threshold{
					Lower: 5,
					Upper: 100,
				},
			},
		},
		{
			Th:          &FsConfig.CriticalPercentThreshold.Space.Used,
			FlagString:  "criticalPercentUsedSpace",
			Description: "Percentage critical threshold for used filesystem space.",
		},
		{
			Th:          &FsConfig.WarningPercentThreshold.Space.Used,
			FlagString:  "warningPercentUsedSpace",
			Description: "Percentage warning threshold for used filesystem space.",
		},
		{
			Th:          &FsConfig.CriticalAbsolutThreshold.Inodes.Free,
			FlagString:  "warningAbsolutFreeInodes",
			Description: "Absolute warning threshold for number of free inodes",
		},
		{
			Th:          &FsConfig.CriticalAbsolutThreshold.Inodes.Used,
			FlagString:  "criticalAbsolutUsedInodes",
			Description: "Absolute critical threshold for number of used inodes",
		},
		{
			Th:          &FsConfig.CriticalPercentThreshold.Inodes.Free,
			FlagString:  "criticalPercentFreeInodes",
			Description: "Percentage critical threshold for percentage of free inodes",
		},
		{
			Th:          &FsConfig.WarningPercentThreshold.Inodes.Free,
			FlagString:  "warningPercentFreeInodes",
			Description: "Percentage warning threshold for percentage of free inodes",
		},
		{
			Th:          &FsConfig.CriticalPercentThreshold.Inodes.Used,
			FlagString:  "criticalPercentUsedInodes",
			Description: "Percentage critical threshold for percentage of used inodes",
			Default: thresholds.ThresholdWrapper{
				IsSet: true,
				Th: check.Threshold{
					Lower: 0,
					Upper: 98,
				},
			},
		},
		{
			Th:          &FsConfig.WarningPercentThreshold.Inodes.Used,
			FlagString:  "warningPercentUsedInodes",
			Description: "Percentage warning threshold for percentage of used inodes",
			Default: thresholds.ThresholdWrapper{
				IsSet: true,
				Th: check.Threshold{
					Lower: 0,
					Upper: 99,
				},
			},
		},
		{
			Th:          &FsConfig.WarningTotalCountOfFs,
			FlagString:  "warningTotalCountOfMatches",
			Description: "A warning threshold for the number of filesystems matching the filters",
		},
		{
			Th:          &FsConfig.CriticalTotalCountOfFs,
			FlagString:  "criticalTotalCountOfMatches",
			Description: "A critical threshold for the number of filesystems matching the filters",
		},
	}

	fs := diskCmd.Flags()

	// Thresholds
	thresholds.AddFlags(fs, &fsThresholds)

	// TODO change this directly to the regex expression type
	fs.StringSliceVar(&FsConfig.Filters.ExcludeFsType, "exclude-fs-type", nil,
		"Ignore all filesystems of indicated type (may be repeated). E.g. 'zfs', 'apfs'\nNote: The same thresholds will be applied to all filesystems")
	fs.StringSliceVar(&FsConfig.Filters.IncludeFsType, "include-fs-type", IncludeFsTypeDefaults,
		"Explicitly include only filesystems of indicated type (may be repeated). E.g. 'zfs', 'apfs'\nNote: The same thresholds will be applied to all filesystems")

	fs.StringSliceVar(&FsConfig.Filters.ExcludeDevicePaths, "exclude-device-path", nil,
		"Ignore the given device path regex (may be repeated). E.g. '/dev/sd.*'")
	fs.StringSliceVar(&FsConfig.Filters.IncludeDevicePaths, "include-device-path", nil,
		"Explicitly include only filesystems of indicated type (may be repeated). E.g. '/dev/sda'")

	fs.StringSliceVar(&FsConfig.Filters.ExcludeMountPaths, "exclude-mount-path", nil,
		"Ignore the given mount path regex (may be repeated). E.g. '^/srv/mount.*'")
	fs.StringSliceVar(&FsConfig.Filters.IncludeMountPaths, "include-mount-path", nil,
		"Explicitly include only filesystems of indicated type (may be repeated). E.g. '/dev/sda'")

	fs.StringSliceVar(&FsConfig.Filters.ExcludeOptions, "exclude-mount-options", nil,
		"Ignore the filesystems with this mount option (in form of a go regexp regex) (may be repeated). E.g. 'async' or '^sync$'")
	fs.StringSliceVar(&FsConfig.Filters.IncludeOptions, "include-mount-options", nil,
		"Explicitly include only filesystems mounted with the given option (in form of a go regexp regex) (may be repeated). E.g. 'async' or '^sync$'")

	fs.BoolVar(&FsConfig.ReadonlyOption, "readonly-filesystems", false,
		"Only list filesystem mounted as readonly. This is just a convenient shorthand for \"--include-mount-options '^ro$'\"")
	fs.BoolVar(&FsConfig.ReadWriteOption, "readwrite-filesystems", false,
		"Only list filesystem mounted as readwrite. This is just a convenient shorthand for \"--include-mount-options '^rw$'\"")

	fs.SortFlags = false
}

func validateOptions(config *filesystem.CheckConfig) error {
	if config.ReadonlyOption && config.ReadWriteOption {
		return errors.New("readonly and readwrite options are mutually exclusive. Please remove one of them")
	}

	if config.ReadonlyOption {
		config.Filters.IncludeOptions = append(config.Filters.IncludeOptions, "^ro$")
	}

	if config.ReadWriteOption {
		config.Filters.IncludeOptions = append(config.Filters.IncludeOptions, "^rw$")
	}

	return nil
}
