package cmd

import (
	"fmt"

	"github.com/NETWAYS/check_system_basics/internal/load"
	"github.com/NETWAYS/go-check"
	"github.com/NETWAYS/go-check/perfdata"
	"github.com/NETWAYS/go-check/result"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/spf13/cobra"
)

const critThresMsg = " exceeds critical threshold"
const warnThresMsg = " exceeds warning threshold"

var LoadConfig load.LoadConfig

var loadCmd = &cobra.Command{
	Use:   "load",
	Short: "Submodule to check the current system load average",
	Example: `./check_system_basics load --load1-warning 0:2
[OK] - states: ok=3
\_ [OK] 1 minute average: 0.10
\_ [OK] 5 minute average: 0.21
\_ [OK] 15 minute average: 0.25
|load1=0.1;2;;0 load5=0.21;;;0 load15=0.25;;;0`,
	Run: func(_ *cobra.Command, _ []string) {
		loadStats, err := load.GetActualLoadValues()
		if err != nil {
			check.ExitError(err)
		}

		cpuCount, err := cpu.Counts(true)
		if err != nil {
			check.ExitError(fmt.Errorf("could not get CPU count: %w", err))
		}

		var originalLoad [3]float64

		if LoadConfig.PerCPU {
			originalLoad[0] = loadStats.LoadAvg.Load1
			loadStats.LoadAvg.Load1 /= float64(cpuCount)
			originalLoad[1] = loadStats.LoadAvg.Load5
			loadStats.LoadAvg.Load5 /= float64(cpuCount)
			originalLoad[2] = loadStats.LoadAvg.Load15
			loadStats.LoadAvg.Load15 /= float64(cpuCount)
		}

		var overall result.Overall

		// 1 Minute average
		var partialLoad1 result.PartialResult
		_ = partialLoad1.SetDefaultState(check.OK)

		// TODO Use strings.Builder
		tmpOutput := fmt.Sprintf("1 minute average: %.2f", loadStats.LoadAvg.Load1)
		tmpPerfdata := &perfdata.Perfdata{
			Label: "load1",
			Value: loadStats.LoadAvg.Load1,
			Min:   0,
			Max:   nil,
		}

		if LoadConfig.Load1Th.Crit.IsSet {
			tmpPerfdata.Crit = &LoadConfig.Load1Th.Crit.Th
			if LoadConfig.Load1Th.Crit.Th.DoesViolate(loadStats.LoadAvg.Load1) {
				_ = partialLoad1.SetState(check.Critical)
				tmpOutput += critThresMsg
			}
		} else if LoadConfig.Load1Th.Warn.IsSet {
			tmpPerfdata.Warn = &LoadConfig.Load1Th.Warn.Th
			if LoadConfig.Load1Th.Warn.Th.DoesViolate(loadStats.LoadAvg.Load1) {
				_ = partialLoad1.SetState(check.Warning)
				tmpOutput += warnThresMsg
			}
		} else {
			_ = partialLoad1.SetState(check.OK)
		}
		if LoadConfig.PerCPU {
			tmpOutput += fmt.Sprintf(", system total: %.2f", originalLoad[0])
		}
		partialLoad1.Output = tmpOutput
		partialLoad1.Perfdata.Add(tmpPerfdata)

		// 5 Minute average
		var partialLoad5 result.PartialResult
		_ = partialLoad5.SetDefaultState(check.OK)

		tmpOutput = fmt.Sprintf("5 minute average: %.2f", loadStats.LoadAvg.Load5)
		tmpPerfdata = &perfdata.Perfdata{
			Label: "load5",
			Value: loadStats.LoadAvg.Load5,
			Min:   0,
			Max:   nil,
		}

		if LoadConfig.Load5Th.Crit.IsSet {
			tmpPerfdata.Crit = &LoadConfig.Load5Th.Crit.Th
			if LoadConfig.Load5Th.Crit.Th.DoesViolate(loadStats.LoadAvg.Load5) {
				_ = partialLoad5.SetState(check.Critical)
				tmpOutput += critThresMsg
			}
		} else if LoadConfig.Load5Th.Warn.IsSet {
			tmpPerfdata.Warn = &LoadConfig.Load5Th.Warn.Th
			if LoadConfig.Load5Th.Warn.Th.DoesViolate(loadStats.LoadAvg.Load5) {
				_ = partialLoad5.SetState(check.Warning)
				tmpOutput += warnThresMsg
			}
		} else {
			_ = partialLoad5.SetState(check.OK)
		}
		if LoadConfig.PerCPU {
			tmpOutput += fmt.Sprintf(", system total: %.2f", originalLoad[1])
		}
		partialLoad5.Output = tmpOutput
		partialLoad5.Perfdata.Add(tmpPerfdata)

		// 15 Minute average
		var partialLoad15 result.PartialResult
		_ = partialLoad15.SetDefaultState(check.OK)

		tmpOutput = fmt.Sprintf("15 minute average: %.2f", loadStats.LoadAvg.Load15)
		tmpPerfdata = &perfdata.Perfdata{
			Label: "load15",
			Value: loadStats.LoadAvg.Load15,
			Min:   0,
			Max:   nil,
		}

		if LoadConfig.Load15Th.Crit.IsSet {
			tmpPerfdata.Crit = &LoadConfig.Load15Th.Crit.Th
			if LoadConfig.Load15Th.Crit.Th.DoesViolate(loadStats.LoadAvg.Load15) {
				_ = partialLoad15.SetState(check.Critical)
				tmpOutput += critThresMsg
			}
		} else if LoadConfig.Load15Th.Warn.IsSet {
			tmpPerfdata.Warn = &LoadConfig.Load15Th.Warn.Th
			if LoadConfig.Load15Th.Warn.Th.DoesViolate(loadStats.LoadAvg.Load15) {
				_ = partialLoad15.SetState(check.Warning)
				tmpOutput += warnThresMsg
			}
		} else {
			_ = partialLoad15.SetState(check.OK)
		}
		if LoadConfig.PerCPU {
			tmpOutput += fmt.Sprintf(", system total: %.2f", originalLoad[2])
		}
		partialLoad15.Output = tmpOutput
		partialLoad15.Perfdata.Add(tmpPerfdata)

		overall.AddSubcheck(partialLoad1)
		overall.AddSubcheck(partialLoad5)
		overall.AddSubcheck(partialLoad15)

		check.ExitRaw(overall.GetStatus(), overall.GetOutput())
	},
}

func init() {
	rootCmd.AddCommand(loadCmd)
	loadCmd.DisableFlagsInUseLine = true

	loadFs := loadCmd.Flags()
	loadFs.Var(&LoadConfig.Load1Th.Warn, "load1-warning", "Warning threshold for the load 1 minute average.")
	loadFs.Var(&LoadConfig.Load1Th.Crit, "load1-critical", "Critical threshold for the load 1 minute average.")
	loadFs.Var(&LoadConfig.Load5Th.Warn, "load5-warning", "Warning threshold for the load 5 minute average.")
	loadFs.Var(&LoadConfig.Load5Th.Crit, "load5-critical", "Critical threshold for the load 5 minute average.")
	loadFs.Var(&LoadConfig.Load15Th.Warn, "load15-warning", "Warning threshold for the load 15 minute average.")
	loadFs.Var(&LoadConfig.Load15Th.Crit, "load15-critical", "Critical threshold for the load 15 minute average.")

	loadFs.BoolVarP(&LoadConfig.PerCPU, "per-cpu", "p", false,
		"Divide the load averages by the number of CPUs")

	loadFs.SortFlags = false
}
