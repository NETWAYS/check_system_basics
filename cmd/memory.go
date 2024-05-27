package cmd

import (
	"fmt"

	"github.com/NETWAYS/check_system_basics/internal/common/thresholds"
	"github.com/NETWAYS/check_system_basics/internal/memory"
	"github.com/NETWAYS/go-check"
	"github.com/NETWAYS/go-check/perfdata"
	"github.com/NETWAYS/go-check/result"
	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"
)

var MemoryConfig memory.MemConfig

var memoryCmd = &cobra.Command{
	Use:   "memory",
	Short: "Submodule to check the current system memory and swap usage",
	Example: `check_system_basics memory --memory-available-warning 10:20 --memory-free-critical-percentage 50:80 --memory-used-warning @10 --percentage-in-perfdata --swap-free-warning-percentage @30:31
[CRITICAL] - states: critical=1 ok=1
\_ [CRITICAL] RAM
    \_ [WARNING] Available Memory (23 GiB/31 GiB, 74.36%)
    \_ [CRITICAL] Free Memory (16 GiB/31 GiB, 49.95%)
    \_ [OK] Used Memory (6.1 GiB/31 GiB, 19.57%)
\_ [OK] Swap Usage 0.00% (0 B / 36 GiB)
|available_memory_percentage=74.36%;15:100;5:100 available_memory=24856633344B;10:20;;0;33427595264 free_memory=16696102912B;;;0;33427595264 free_memory_percentage=49.947%;;50:80 used_memory=6542696448B;@10;;0;33427595264 free_memory_percentage=19.573% swap_usage_percent=0%;20;85 swap_used=0B;;;0;38654701568
`,
	Run: func(cmd *cobra.Command, args []string) {

		var overall result.Overall

		// ## RAM stuff
		memStats, err := memory.LoadMemStat()
		if err != nil {
			check.ExitError(err)
		}

		// Memory stuff
		partialMem := computeMemResults(&MemoryConfig, memStats)

		overall.AddSubcheck(partialMem)

		// Swap stuff
		if memStats.VirtMem.SwapTotal != 0 {
			partSwap := computeSwapResults(memStats)
			overall.AddSubcheck(*partSwap)
		}

		check.ExitRaw(overall.GetStatus(), overall.GetOutput())
	},
}

func computeMemResults(config *memory.MemConfig, memStats *memory.Mem) result.PartialResult {
	partialMem := result.PartialResult{
		Output: "RAM",
	}
	_ = partialMem.SetDefaultState(check.OK)

	// # Available
	var partialMemAvailable result.PartialResult
	_ = partialMemAvailable.SetDefaultState(check.OK)

	partialMemAvailable.Output = fmt.Sprintf("Available Memory (%s/%s, %.2f%%)",
		humanize.IBytes(memStats.VirtMem.Available),
		humanize.IBytes(memStats.VirtMem.Total),
		memStats.MemAvailablePercentage)

	// perfdata
	pdMemAvailable := perfdata.Perfdata{
		Label: "available_memory",
		Value: memStats.VirtMem.Available,
		Uom:   "B",
		Min:   0,
		Max:   memStats.VirtMem.Total,
	}

	pdMemAvailablePrcnt := perfdata.Perfdata{
		Label: "available_memory_percentage",
		Value: float64(memStats.VirtMem.Available) / float64(memStats.VirtMem.Total/100),
		Uom:   "%",
	}

	if config.MemAvailable.Warn.IsSet {
		pdMemAvailable.Warn = &config.MemAvailable.Warn.Th

		if config.MemAvailable.Warn.Th.DoesViolate(float64(memStats.VirtMem.Available)) {
			_ = partialMemAvailable.SetState(check.Warning)
		}
	}

	if config.MemAvailablePercentage.Warn.IsSet {
		pdMemAvailablePrcnt.Warn = &config.MemAvailablePercentage.Warn.Th

		if config.MemAvailablePercentage.Warn.Th.DoesViolate(memStats.MemAvailablePercentage) {
			_ = partialMemAvailable.SetState(check.Warning)
		}
	}

	if config.MemAvailable.Crit.IsSet {
		pdMemAvailable.Crit = &config.MemAvailable.Crit.Th

		if config.MemAvailable.Crit.Th.DoesViolate(float64(memStats.VirtMem.Available)) {
			_ = partialMemAvailable.SetState(check.Critical)
		}
	}

	if config.MemAvailablePercentage.Crit.IsSet {
		pdMemAvailablePrcnt.Crit = &config.MemAvailablePercentage.Crit.Th

		if config.MemAvailablePercentage.Crit.Th.DoesViolate(memStats.MemAvailablePercentage) {
			_ = partialMemAvailable.SetState(check.Critical)
		}
	}

	if config.PercentageInPerfdata {
		partialMemAvailable.Perfdata.Add(&pdMemAvailablePrcnt)
	}

	partialMemAvailable.Perfdata.Add(&pdMemAvailable)
	partialMem.AddSubcheck(partialMemAvailable)

	if (partialMemAvailable.GetStatus() > partialMem.GetStatus()) &&
		partialMemAvailable.GetStatus() != check.Unknown {
		_ = partialMem.SetState(partialMemAvailable.GetStatus())
	}

	// # Free
	var partialMemFree result.PartialResult
	_ = partialMemFree.SetDefaultState(check.OK)

	pdMemFree := perfdata.Perfdata{
		Label: "free_memory",
		Uom:   "B",
		Value: memStats.VirtMem.Free,
		Min:   0,
		Max:   memStats.VirtMem.Total,
	}

	MemFreePercentage := float64(memStats.VirtMem.Free) / (float64(memStats.VirtMem.Total) / 100)

	pdMemFreePercentage := perfdata.Perfdata{
		Label: "free_memory_percentage",
		Value: MemFreePercentage,
		Uom:   "%",
	}

	partialMemFree.Output = fmt.Sprintf("Free Memory (%s/%s, %.2f%%)",
		humanize.IBytes(memStats.VirtMem.Free),
		humanize.IBytes(memStats.VirtMem.Total),
		MemFreePercentage)

	if config.MemFree.Warn.IsSet {
		pdMemFree.Warn = &config.MemFree.Warn.Th

		if config.MemFree.Warn.Th.DoesViolate(float64(memStats.VirtMem.Free)) {
			_ = partialMemFree.SetState(check.Warning)
		}
	}

	if config.MemFree.Crit.IsSet {
		pdMemFree.Crit = &config.MemFree.Crit.Th

		if config.MemFree.Crit.Th.DoesViolate(float64(memStats.VirtMem.Free)) {
			_ = partialMemFree.SetState(check.Critical)
		}
	}

	if config.MemFreePercentage.Warn.IsSet {
		pdMemFreePercentage.Warn = &config.MemFreePercentage.Warn.Th

		if config.MemFreePercentage.Warn.Th.DoesViolate(MemFreePercentage) {
			_ = partialMemFree.SetState(check.Warning)
		}
	}

	if config.MemFreePercentage.Crit.IsSet {
		pdMemFreePercentage.Crit = &config.MemFreePercentage.Crit.Th

		if config.MemFreePercentage.Crit.Th.DoesViolate(MemFreePercentage) {
			_ = partialMemFree.SetState(check.Critical)
		}
	}

	partialMemFree.Perfdata.Add(&pdMemFree)

	if config.PercentageInPerfdata {
		partialMemFree.Perfdata.Add(&pdMemFreePercentage)
	}

	partialMem.AddSubcheck(partialMemFree)

	if (partialMemFree.GetStatus() > partialMem.GetStatus()) &&
		partialMemFree.GetStatus() != check.Unknown {
		_ = partialMem.SetState(partialMemFree.GetStatus())
	}

	// Used Memory
	var partialMemUsed result.PartialResult
	_ = partialMemUsed.SetDefaultState(check.OK)

	partialMemUsed.Output = fmt.Sprintf("Used Memory (%s/%s, %.2f%%)",
		humanize.IBytes(memStats.VirtMem.Used),
		humanize.IBytes(memStats.VirtMem.Total),
		memStats.VirtMem.UsedPercent)

	pdMemUsed := perfdata.Perfdata{
		Label: "used_memory",
		Uom:   "B",
		Value: memStats.VirtMem.Used,
		Min:   0,
		Max:   memStats.VirtMem.Total,
	}

	MemUsedPercentage := float64(memStats.VirtMem.Used) / (float64(memStats.VirtMem.Total) / 100)
	pdMemUsedPercentage := perfdata.Perfdata{
		Label: "free_memory_percentage",
		Value: MemUsedPercentage,
		Uom:   "%",
	}

	if config.MemUsed.Warn.IsSet {
		pdMemUsed.Warn = &config.MemUsed.Warn.Th

		if config.MemUsed.Warn.Th.DoesViolate(float64(memStats.VirtMem.Used)) {
			_ = partialMemUsed.SetState(check.Warning)
		}
	}

	if config.MemUsedPercentage.Warn.IsSet {
		pdMemUsedPercentage.Warn = &config.MemUsedPercentage.Warn.Th

		if config.MemUsedPercentage.Warn.Th.DoesViolate(memStats.VirtMem.UsedPercent) {
			_ = partialMemUsed.SetState(check.Warning)
		}
	}

	if config.MemUsed.Crit.IsSet {
		pdMemUsed.Crit = &config.MemUsed.Crit.Th

		if config.MemUsed.Crit.Th.DoesViolate(float64(memStats.VirtMem.Used)) {
			_ = partialMemUsed.SetState(check.Critical)
		}
	}

	if config.MemUsedPercentage.Crit.IsSet {
		pdMemUsedPercentage.Crit = &config.MemUsedPercentage.Crit.Th

		if config.MemUsedPercentage.Crit.Th.DoesViolate(memStats.VirtMem.UsedPercent) {
			_ = partialMemUsed.SetState(check.Critical)
		}
	}

	partialMemUsed.Perfdata.Add(&pdMemUsed)

	if config.PercentageInPerfdata {
		partialMemUsed.Perfdata.Add(&pdMemUsedPercentage)
	}

	partialMem.AddSubcheck(partialMemUsed)

	if (partialMemUsed.GetStatus() > partialMem.GetStatus()) &&
		partialMemUsed.GetStatus() != check.Unknown {
		_ = partialMem.SetState(partialMemUsed.GetStatus())
	}

	return partialMem
}

func init() {
	rootCmd.AddCommand(memoryCmd)

	memoryCmd.DisableFlagsInUseLine = true

	memPerFs := memoryCmd.Flags()

	memoryThresholds := []thresholds.ThresholdOption{
		{
			Th:          &MemoryConfig.MemFree.Warn,
			FlagString:  "memory-free-warning",
			Description: "Warning threshold for free memory",
		},
		{
			Th:          &MemoryConfig.MemFree.Crit,
			FlagString:  "memory-free-critical",
			Description: "Critical threshold for free memory",
		},
		{
			Th:          &MemoryConfig.MemFreePercentage.Crit,
			FlagString:  "memory-free-critical-percentage",
			Description: "Critical threshold for free memory",
		},
		{
			Th:          &MemoryConfig.MemUsed.Warn,
			FlagString:  "memory-used-warning",
			Description: "Warning threshold for used memory",
		},
		{
			Th:          &MemoryConfig.MemUsed.Crit,
			FlagString:  "memory-used-critical",
			Description: "Critical threshold for used memory",
		},
		{
			Th:          &MemoryConfig.MemUsedPercentage.Warn,
			FlagString:  "memory-used-warning-percentage",
			Description: "Warning threshold for used memory",
		},
		{
			Th:          &MemoryConfig.MemUsedPercentage.Crit,
			FlagString:  "memory-used-critical-percentage",
			Description: "Critical threshold for used memory",
		},
		{
			Th:          &MemoryConfig.MemAvailable.Warn,
			FlagString:  "memory-available-warning",
			Description: "Warning threshold for available memory",
		},
		{
			Th:          &MemoryConfig.MemAvailable.Crit,
			FlagString:  "memory-available-critical",
			Description: "Critical threshold for available memory",
		},
		{
			Th:          &MemoryConfig.MemAvailablePercentage.Warn,
			FlagString:  "memory-available-warning-percentage",
			Description: "Warning threshold for available memory",
			Default: thresholds.ThresholdWrapper{
				IsSet: true,
				Th: check.Threshold{
					Lower: 15,
					Upper: check.PosInf,
				},
			},
		},
		{
			Th:          &MemoryConfig.MemAvailablePercentage.Crit,
			FlagString:  "memory-available-critical-percentage",
			Description: "Critical threshold for available memory",
			Default: thresholds.ThresholdWrapper{
				IsSet: true,
				Th: check.Threshold{
					Lower: 5,
					Upper: check.PosInf,
				},
			},
		},
		{
			Th:          &MemoryConfig.SwapFree.Warn,
			FlagString:  "swap-free-warning",
			Description: "Warning threshold for free memory",
		},
		{
			Th:          &MemoryConfig.SwapFree.Crit,
			FlagString:  "swap-free-critical",
			Description: "Critical threshold for free memory",
		},
		{
			Th:          &MemoryConfig.SwapFreePercentage.Warn,
			FlagString:  "swap-free-warning-percentage",
			Description: "Warning threshold for free memory",
		},
		{
			Th:          &MemoryConfig.SwapFreePercentage.Crit,
			FlagString:  "swap-free-critical-percentage",
			Description: "Critical threshold for free memory",
		},

		{
			Th:          &MemoryConfig.SwapUsed.Warn,
			FlagString:  "swap-used-warning",
			Description: "Warning threshold for used memory",
		},
		{
			Th:          &MemoryConfig.SwapUsed.Crit,
			FlagString:  "swap-used-critical",
			Description: "Critical threshold for used memory",
		},
		{
			Th:          &MemoryConfig.SwapUsedPercentage.Warn,
			FlagString:  "swap-used-warning-percentage",
			Description: "Warning threshold for used memory",
			Default: thresholds.ThresholdWrapper{
				IsSet: true,
				Th: check.Threshold{
					Lower: 0,
					Upper: 20,
				},
			},
		},
		{
			Th:          &MemoryConfig.SwapUsedPercentage.Crit,
			FlagString:  "swap-used-critical-percentage",
			Description: "Critical threshold for used memory",
			Default: thresholds.ThresholdWrapper{
				IsSet: true,
				Th: check.Threshold{
					Lower: 0,
					Upper: 85,
				},
			},
		},
	}

	// Thresholds
	thresholds.AddFlags(memPerFs, &memoryThresholds)

	memPerFs.BoolVarP(&MemoryConfig.PercentageInPerfdata, "percentage-in-perfdata", "", false, "Add computed percentage values to perfdata, although they are technically redundant")

	memPerFs.SortFlags = false
}

func computeSwapResults(stats *memory.Mem) *result.PartialResult {
	var partialSwap result.PartialResult
	_ = partialSwap.SetDefaultState(check.OK)

	_ = partialSwap.SetDefaultState(check.OK)

	if stats.VirtMem.SwapTotal == 0 {
		_ = partialSwap.SetState(check.Critical)
		partialSwap.Output = "Swap size is 0."

		return &partialSwap
	}

	partialSwap.Output = fmt.Sprintf("Swap Usage %.2f%% (%s / %s)", stats.SwapInfo.UsedPercent, humanize.IBytes(stats.SwapInfo.Used), humanize.IBytes(stats.SwapInfo.Total))

	pdSwapUsed := perfdata.Perfdata{
		Label: "swap_used",
		Value: stats.SwapInfo.Used,
		Uom:   "B",
		Min:   0,
		Max:   stats.SwapInfo.Total,
	}

	pdSwapPrcnt := perfdata.Perfdata{
		Label: "swap_usage_percent",
		Value: stats.SwapInfo.UsedPercent,
		Uom:   "%",
	}

	// Warning
	if MemoryConfig.SwapFree.Warn.IsSet {
		if MemoryConfig.SwapFree.Warn.Th.DoesViolate(float64(stats.SwapInfo.Free)) {
			_ = partialSwap.SetState(check.Warning)
		}
	}

	if MemoryConfig.SwapFreePercentage.Warn.IsSet {
		if MemoryConfig.SwapFreePercentage.Warn.Th.DoesViolate(1 - stats.SwapInfo.UsedPercent) {
			_ = partialSwap.SetState(check.Warning)
		}
	}

	if MemoryConfig.SwapUsed.Warn.IsSet {
		pdSwapUsed.Warn = &MemoryConfig.SwapUsed.Warn.Th

		if MemoryConfig.SwapUsed.Warn.Th.DoesViolate(float64(stats.SwapInfo.Used)) {
			_ = partialSwap.SetState(check.Warning)
		}
	}

	if MemoryConfig.SwapUsedPercentage.Warn.IsSet {
		pdSwapPrcnt.Warn = &MemoryConfig.SwapUsedPercentage.Warn.Th

		if MemoryConfig.SwapUsedPercentage.Warn.Th.DoesViolate(stats.SwapInfo.UsedPercent) {
			_ = partialSwap.SetState(check.Warning)
		}
	}

	// Critical
	if MemoryConfig.SwapFree.Crit.IsSet {
		if MemoryConfig.SwapFree.Crit.Th.DoesViolate(float64(stats.SwapInfo.Free)) {
			_ = partialSwap.SetState(check.Critical)
		}
	}

	if MemoryConfig.SwapFreePercentage.Crit.IsSet {
		if MemoryConfig.SwapFreePercentage.Crit.Th.DoesViolate(1 - stats.SwapInfo.UsedPercent) {
			_ = partialSwap.SetState(check.Critical)
		}
	}

	if MemoryConfig.SwapUsed.Crit.IsSet {
		pdSwapUsed.Crit = &MemoryConfig.SwapUsed.Crit.Th

		if MemoryConfig.SwapUsed.Warn.Th.DoesViolate(float64(stats.SwapInfo.Used)) {
			_ = partialSwap.SetState(check.Warning)
		}
	}

	if MemoryConfig.SwapUsedPercentage.Crit.IsSet {
		pdSwapPrcnt.Crit = &MemoryConfig.SwapUsedPercentage.Crit.Th

		if MemoryConfig.SwapUsedPercentage.Crit.Th.DoesViolate(stats.SwapInfo.UsedPercent) {
			_ = partialSwap.SetState(check.Critical)
		}
	}

	// Percentage to Perfdata
	if MemoryConfig.PercentageInPerfdata {
		partialSwap.Perfdata.Add(&pdSwapPrcnt)
	}

	partialSwap.Perfdata.Add(&pdSwapUsed)

	return &partialSwap
}
