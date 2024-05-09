package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/NETWAYS/check_system_basics/internal/common/thresholds"
	"github.com/NETWAYS/check_system_basics/internal/psi"
	"github.com/NETWAYS/go-check"

	"github.com/NETWAYS/go-check/result"
	"github.com/spf13/cobra"
)

type psiConfig struct {
	IncludeCPU    bool
	IncludeMemory bool
	IncludeIO     bool

	WarningCPUSomeAvg10   thresholds.ThresholdWrapper
	WarningCPUSomeAvg60   thresholds.ThresholdWrapper
	WarningCPUSomeAvg300  thresholds.ThresholdWrapper
	CriticalCPUSomeAvg10  thresholds.ThresholdWrapper
	CriticalCPUSomeAvg60  thresholds.ThresholdWrapper
	CriticalCPUSomeAvg300 thresholds.ThresholdWrapper

	WarningIoSomeAvg10   thresholds.ThresholdWrapper
	WarningIoSomeAvg60   thresholds.ThresholdWrapper
	WarningIoSomeAvg300  thresholds.ThresholdWrapper
	CriticalIoSomeAvg10  thresholds.ThresholdWrapper
	CriticalIoSomeAvg60  thresholds.ThresholdWrapper
	CriticalIoSomeAvg300 thresholds.ThresholdWrapper

	WarningMemorySomeAvg10   thresholds.ThresholdWrapper
	WarningMemorySomeAvg60   thresholds.ThresholdWrapper
	WarningMemorySomeAvg300  thresholds.ThresholdWrapper
	CriticalMemorySomeAvg10  thresholds.ThresholdWrapper
	CriticalMemorySomeAvg60  thresholds.ThresholdWrapper
	CriticalMemorySomeAvg300 thresholds.ThresholdWrapper

	WarningCPUFullAvg10   thresholds.ThresholdWrapper
	WarningCPUFullAvg60   thresholds.ThresholdWrapper
	WarningCPUFullAvg300  thresholds.ThresholdWrapper
	CriticalCPUFullAvg10  thresholds.ThresholdWrapper
	CriticalCPUFullAvg60  thresholds.ThresholdWrapper
	CriticalCPUFullAvg300 thresholds.ThresholdWrapper

	WarningIoFullAvg10   thresholds.ThresholdWrapper
	WarningIoFullAvg60   thresholds.ThresholdWrapper
	WarningIoFullAvg300  thresholds.ThresholdWrapper
	CriticalIoFullAvg10  thresholds.ThresholdWrapper
	CriticalIoFullAvg60  thresholds.ThresholdWrapper
	CriticalIoFullAvg300 thresholds.ThresholdWrapper

	WarningMemoryFullAvg10   thresholds.ThresholdWrapper
	WarningMemoryFullAvg60   thresholds.ThresholdWrapper
	WarningMemoryFullAvg300  thresholds.ThresholdWrapper
	CriticalMemoryFullAvg10  thresholds.ThresholdWrapper
	CriticalMemoryFullAvg60  thresholds.ThresholdWrapper
	CriticalMemoryFullAvg300 thresholds.ThresholdWrapper

	WarningCPUAvg  thresholds.ThresholdWrapper
	CriticalCPUAvg thresholds.ThresholdWrapper

	WarningMemoryAvg  thresholds.ThresholdWrapper
	CriticalMemoryAvg thresholds.ThresholdWrapper

	WarningIoAvg  thresholds.ThresholdWrapper
	CriticalIoAvg thresholds.ThresholdWrapper
}

var (
	config = psiConfig{
		IncludeCPU:    false,
		IncludeIO:     false,
		IncludeMemory: false,

		WarningCPUAvg:     thresholds.ThresholdWrapper{Th: check.Threshold{Inside: true, Lower: 30, Upper: 100}, IsSet: false},
		CriticalCPUAvg:    thresholds.ThresholdWrapper{Th: check.Threshold{Inside: true, Lower: 95, Upper: 100}, IsSet: false},
		WarningMemoryAvg:  thresholds.ThresholdWrapper{Th: check.Threshold{Inside: true, Lower: 30, Upper: 100}, IsSet: false},
		CriticalMemoryAvg: thresholds.ThresholdWrapper{Th: check.Threshold{Inside: true, Lower: 95, Upper: 100}, IsSet: false},
		WarningIoAvg:      thresholds.ThresholdWrapper{Th: check.Threshold{Inside: true, Lower: 30, Upper: 100}, IsSet: false},
		CriticalIoAvg:     thresholds.ThresholdWrapper{Th: check.Threshold{Inside: true, Lower: 95, Upper: 100}, IsSet: false},

		WarningCPUSomeAvg10:   thresholds.ThresholdWrapper{Th: check.Threshold{}, IsSet: false},
		WarningCPUSomeAvg60:   thresholds.ThresholdWrapper{Th: check.Threshold{}, IsSet: false},
		WarningCPUSomeAvg300:  thresholds.ThresholdWrapper{Th: check.Threshold{}, IsSet: false},
		CriticalCPUSomeAvg10:  thresholds.ThresholdWrapper{Th: check.Threshold{}, IsSet: false},
		CriticalCPUSomeAvg60:  thresholds.ThresholdWrapper{Th: check.Threshold{}, IsSet: false},
		CriticalCPUSomeAvg300: thresholds.ThresholdWrapper{Th: check.Threshold{}, IsSet: false},
		WarningCPUFullAvg10:   thresholds.ThresholdWrapper{Th: check.Threshold{}, IsSet: false},
		WarningCPUFullAvg60:   thresholds.ThresholdWrapper{Th: check.Threshold{}, IsSet: false},
		WarningCPUFullAvg300:  thresholds.ThresholdWrapper{Th: check.Threshold{}, IsSet: false},
		CriticalCPUFullAvg10:  thresholds.ThresholdWrapper{Th: check.Threshold{}, IsSet: false},
		CriticalCPUFullAvg60:  thresholds.ThresholdWrapper{Th: check.Threshold{}, IsSet: false},
		CriticalCPUFullAvg300: thresholds.ThresholdWrapper{Th: check.Threshold{}, IsSet: false},

		WarningMemorySomeAvg10:   thresholds.ThresholdWrapper{Th: check.Threshold{}, IsSet: false},
		WarningMemorySomeAvg60:   thresholds.ThresholdWrapper{Th: check.Threshold{}, IsSet: false},
		WarningMemorySomeAvg300:  thresholds.ThresholdWrapper{Th: check.Threshold{}, IsSet: false},
		CriticalMemorySomeAvg10:  thresholds.ThresholdWrapper{Th: check.Threshold{}, IsSet: false},
		CriticalMemorySomeAvg60:  thresholds.ThresholdWrapper{Th: check.Threshold{}, IsSet: false},
		CriticalMemorySomeAvg300: thresholds.ThresholdWrapper{Th: check.Threshold{}, IsSet: false},
		WarningMemoryFullAvg10:   thresholds.ThresholdWrapper{Th: check.Threshold{}, IsSet: false},
		WarningMemoryFullAvg60:   thresholds.ThresholdWrapper{Th: check.Threshold{}, IsSet: false},
		WarningMemoryFullAvg300:  thresholds.ThresholdWrapper{Th: check.Threshold{}, IsSet: false},
		CriticalMemoryFullAvg10:  thresholds.ThresholdWrapper{Th: check.Threshold{}, IsSet: false},
		CriticalMemoryFullAvg60:  thresholds.ThresholdWrapper{Th: check.Threshold{}, IsSet: false},
		CriticalMemoryFullAvg300: thresholds.ThresholdWrapper{Th: check.Threshold{}, IsSet: false},

		WarningIoSomeAvg10:   thresholds.ThresholdWrapper{Th: check.Threshold{}, IsSet: false},
		WarningIoSomeAvg60:   thresholds.ThresholdWrapper{Th: check.Threshold{}, IsSet: false},
		WarningIoSomeAvg300:  thresholds.ThresholdWrapper{Th: check.Threshold{}, IsSet: false},
		CriticalIoSomeAvg10:  thresholds.ThresholdWrapper{Th: check.Threshold{}, IsSet: false},
		CriticalIoSomeAvg60:  thresholds.ThresholdWrapper{Th: check.Threshold{}, IsSet: false},
		CriticalIoSomeAvg300: thresholds.ThresholdWrapper{Th: check.Threshold{}, IsSet: false},
		WarningIoFullAvg10:   thresholds.ThresholdWrapper{Th: check.Threshold{}, IsSet: false},
		WarningIoFullAvg60:   thresholds.ThresholdWrapper{Th: check.Threshold{}, IsSet: false},
		WarningIoFullAvg300:  thresholds.ThresholdWrapper{Th: check.Threshold{}, IsSet: false},
		CriticalIoFullAvg10:  thresholds.ThresholdWrapper{Th: check.Threshold{}, IsSet: false},
		CriticalIoFullAvg60:  thresholds.ThresholdWrapper{Th: check.Threshold{}, IsSet: false},
		CriticalIoFullAvg300: thresholds.ThresholdWrapper{Th: check.Threshold{}, IsSet: false},
	}
)

var psiCmd = &cobra.Command{
	Use:   "psi",
	Short: "Submodule to check the pressure/stall information the linux kernel offers",
	Long: `This submodule tries to read the pressure stall information from /proc/pressure and compares ` +
		"it against the given thresholds,\nwhich should allow an user to identify overload situation and take " +
		"action accordingly.\n" +
		"This will not work on systems where this interface is not activated in the kernel. For example certain Red Hat (similar) systems.\n" +
		"In that case adding \"psi=1\" to the kernel cmdline might help and activate the PSI interface",
	Example: `./check_system_basics psi  --warning-cpu-avg 30:80 --critical-cpu-full-avg60 @31:81 --warning-io-avg 11:99 --critical-io-some-avg300 @49:51 --warning-memory-avg @00:23
[WARNING] - states: warning=3
\_ [WARNING] CPU Full Pressure - Avg10: 0.00, Avg60: 0.00, Avg300: 0.00
\_ [WARNING] IO Pressure - Avg10: 0.00, Avg60: 0.00, Avg300: 0.00
\_ [WARNING] Memory Pressure - Avg10: 0.00, Avg60: 0.00, Avg300: 0.00
|cpu-some-avg10=0%;30:80;@95:100;0;100 cpu-some-avg60=0.02%;30:80;@95:100;0;100 cpu-some-avg300=0.02%;30:80;@95:100;0;100 cpu-some-total=33046682c;30:80;@95:100;0 cpu-full-avg10=0%;30:80;@31:81;0;100 cpu-full-avg60=0%;30:80;@95:100;0;100 cpu-full-avg300=0%;;;0;100 cpu-full-total=0c;;;0 io-some-avg10=0%;11:99;@95:100;0;100 io-some-avg60=0%;11:99;@95:100;0;100 io-some-avg300=0%;11:99;@49:51;0;100 io-some-total=27037011c;11:99;@95:100;0 io-full-avg10=0%;11:99;@95:100;0;100 io-full-avg60=0%;11:99;@95:100;0;100 io-full-avg300=0%;;;0;100 io-full-total=26265104c;;;0 memory-some-avg10=0%;@23;@95:100;0;100 memory-some-avg60=0%;@23;@95:100;0;100 memory-some-avg300=0%;@23;@95:100;0;100 memory-some-total=354c;@23;@95:100;0 memory-full-avg10=0%;@23;@95:100;0;100 memory-full-avg60=0%;@23;@95:100;0;100 memory-full-avg300=0%;;;0;100 memory-full-total=193c;;;0`,
	Run: func(cmd *cobra.Command, args []string) {
		var overall result.Overall

		// If no mode is selected, select all
		if !config.IncludeCPU && !config.IncludeIO && !config.IncludeMemory {
			config.IncludeCPU = true
			config.IncludeIO = true
			config.IncludeMemory = true
		}

		// CPU Pressure
		if config.IncludeCPU {
			overall.AddSubcheck(checkPsiCPUPressure(&config))
		}

		// IO Pressure
		if config.IncludeIO {
			overall.AddSubcheck(checkPsiIoPressure(&config))
		}

		// Memory Pressure
		if config.IncludeMemory {
			overall.AddSubcheck(checkPsiMemoryPressure(&config))
		}

		check.ExitRaw(overall.GetStatus(), overall.GetOutput())
	},
}

func init() {
	rootCmd.AddCommand(psiCmd)
	psiFs := psiCmd.Flags()
	psiFs.SortFlags = false

	psiFs.Var(&config.WarningCPUAvg, "warning-cpu-avg", "Warning threshold for all the pressure/cpu values. Will be overwritten by more specific parameters.")
	psiFs.Var(&config.CriticalCPUAvg, "critical-cpu-avg", "Critical threshold for all the pressure/cpu values. Will be overwritten by more specific parameters.")

	psiFs.Var(&config.WarningMemoryAvg, "warning-memory-avg", "Warning threshold for all the pressure/memory values. Will be overwritten by more specific parameters.")
	psiFs.Var(&config.CriticalMemoryAvg, "critical-memory-avg", "Critical threshold for all the pressure/memory values. Will be overwritten by more specific parameters.")

	psiFs.Var(&config.WarningIoAvg, "warning-io-avg", "Warning threshold for all the pressure/io values. Will be overwritten by more specific parameters.")
	psiFs.Var(&config.CriticalIoAvg, "critical-io-avg", "Critical threshold for all the pressure/io values. Will be overwritten by more specific parameters.")

	psiFs.Var(&config.WarningCPUSomeAvg10, "warning-cpu-some-avg10", "Warning threshold for the pressure/cpu Some Avg10 value")
	psiFs.Var(&config.WarningCPUSomeAvg60, "warning-cpu-some-avg60", "Warning threshold for the pressure/cpu Some Avg60 value")
	psiFs.Var(&config.WarningCPUSomeAvg300, "warning-cpu-some-avg300", "Warning threshold for the pressure/cpu Some Avg300 value")

	psiFs.Var(&config.WarningCPUFullAvg10, "warning-cpu-full-avg10", "Warning threshold for the pressure/cpu Full Avg10 value")
	psiFs.Var(&config.WarningCPUFullAvg60, "warning-cpu-full-avg60", "Warning threshold for the pressure/cpu Full Avg60 value")
	psiFs.Var(&config.WarningCPUFullAvg300, "warning-cpu-full-avg300", "Warning threshold for the pressure/cpu Full Avg300 value")

	psiFs.Var(&config.WarningIoSomeAvg10, "warning-io-some-avg10", "Warning threshold for the pressure/io Some Avg10 value")
	psiFs.Var(&config.WarningIoSomeAvg60, "warning-io-some-avg60", "Warning threshold for the pressure/io Some Avg60 value")
	psiFs.Var(&config.WarningIoSomeAvg300, "warning-io-some-avg300", "Warning threshold for the pressure/io Some Avg300 value")

	psiFs.Var(&config.WarningIoFullAvg10, "warning-io-full-avg10", "Warning threshold for the pressure/io Full Avg10 value")
	psiFs.Var(&config.WarningIoFullAvg60, "warning-io-full-avg60", "Warning threshold for the pressure/io Full Avg60 value")
	psiFs.Var(&config.WarningIoFullAvg300, "warning-io-full-avg300", "Warning threshold for the pressure/io Full Avg300 value")

	psiFs.Var(&config.WarningMemorySomeAvg10, "warning-memory-some-avg10", "Warning threshold for the pressure/memory Some Avg10 value")
	psiFs.Var(&config.WarningMemorySomeAvg60, "warning-memory-some-avg60", "Warning threshold for the pressure/memory Some Avg60 value")
	psiFs.Var(&config.WarningMemorySomeAvg300, "warning-memory-some-avg300", "Warning threshold for the pressure/memory Some Avg300 value")

	psiFs.Var(&config.WarningMemoryFullAvg10, "warning-memory-full-avg10", "Warning threshold for the pressure/memory Full Avg10 value")
	psiFs.Var(&config.WarningMemoryFullAvg60, "warning-memory-full-avg60", "Warning threshold for the pressure/memory Full Avg60 value")
	psiFs.Var(&config.WarningMemoryFullAvg300, "warning-memory-full-avg300", "Warning threshold for the pressure/memory Full Avg300 value")

	psiFs.Var(&config.CriticalCPUSomeAvg10, "critical-cpu-some-avg10", "Critical threshold for the pressure/cpu Some Avg10 value")
	psiFs.Var(&config.CriticalCPUSomeAvg60, "critical-cpu-some-avg60", "Critical threshold for the pressure/cpu Some Avg60 value")
	psiFs.Var(&config.CriticalCPUSomeAvg300, "critical-cpu-some-avg300", "Critical threshold for the pressure/cpu Some Avg300 value")

	psiFs.Var(&config.CriticalCPUFullAvg10, "critical-cpu-full-avg10", "Critical threshold for the pressure/cpu Full Avg10 value")
	psiFs.Var(&config.CriticalCPUFullAvg60, "critical-cpu-full-avg60", "Critical threshold for the pressure/cpu Full Avg60 value")
	psiFs.Var(&config.CriticalCPUFullAvg300, "critical-cpu-full-avg300", "Critical threshold for the pressure/cpu Full Avg300 value")

	psiFs.Var(&config.CriticalIoSomeAvg10, "critical-io-some-avg10", "Critical threshold for the pressure/io Some Avg10 value")
	psiFs.Var(&config.CriticalIoSomeAvg60, "critical-io-some-avg60", "Critical threshold for the pressure/io Some Avg60 value")
	psiFs.Var(&config.CriticalIoSomeAvg300, "critical-io-some-avg300", "Critical threshold for the pressure/io Some Avg300 value")

	psiFs.Var(&config.CriticalIoFullAvg10, "critical-io-full-avg10", "Critical threshold for the pressure/io Full Avg10 value")
	psiFs.Var(&config.CriticalIoFullAvg60, "critical-io-full-avg60", "Critical threshold for the pressure/io Full Avg60 value")
	psiFs.Var(&config.CriticalIoFullAvg300, "critical-io-full-avg300", "Critical threshold for the pressure/io Full Avg300 value")

	psiFs.Var(&config.CriticalMemorySomeAvg10, "critical-memory-some-avg10", "Critical threshold for the pressure/memory Some Avg10 value")
	psiFs.Var(&config.CriticalMemorySomeAvg60, "critical-memory-some-avg60", "Critical threshold for the pressure/memory Some Avg60 value")
	psiFs.Var(&config.CriticalMemorySomeAvg300, "critical-memory-some-avg300", "Critical threshold for the pressure/memory Some Avg300 value")

	psiFs.Var(&config.CriticalMemoryFullAvg10, "critical-memory-full-avg10", "Critical threshold for the pressure/memory Full Avg10 value")
	psiFs.Var(&config.CriticalMemoryFullAvg60, "critical-memory-full-avg60", "Critical threshold for the pressure/memory Full Avg60 value")
	psiFs.Var(&config.CriticalMemoryFullAvg300, "critical-memory-full-avg300", "Critical threshold for the pressure/memory Full Avg300 value")

	psiFs.BoolVar(&config.IncludeCPU, "include-cpu", false, "Include CPU values explicitly (by default all are included)")
	psiFs.BoolVar(&config.IncludeMemory, "include-memory", false, "Include Memory values explicitly (by default all are included)")
	psiFs.BoolVar(&config.IncludeIO, "include-io", false, "Include IO values explicitly (by default all are included)")
}

func checkPsiCPUPressure(config *psiConfig) result.PartialResult {
	var cpuCheck result.PartialResult
	_ = cpuCheck.SetDefaultState(check.OK)
	cpuCheck.Output = "CPU"

	psiCPU, err := psi.ReadCPUPressure()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			_ = cpuCheck.SetState(check.Unknown)
			cpuCheck.Output = "CPU pressure file not found. Perhaps the PSI interface is not active on this system? It might be necessary to change the kernel config"

			return cpuCheck
		}

		check.ExitError(err)
	}

	cpuCheck.Perfdata = *psiCPU.Perfdata()

	//nolint:nestif
	if psiCPU.FullPresent {
		// Warn thresholds
		if config.WarningCPUFullAvg10.IsSet {
			cpuCheck.Perfdata[psi.CPUFullAvg10].Warn = &config.WarningCPUFullAvg10.Th
		} else {
			cpuCheck.Perfdata[psi.CPUFullAvg10].Warn = &config.WarningCPUAvg.Th
		}

		if config.WarningCPUFullAvg60.IsSet {
			cpuCheck.Perfdata[psi.CPUFullAvg60].Warn = &config.WarningCPUFullAvg60.Th
		} else {
			cpuCheck.Perfdata[psi.CPUFullAvg60].Warn = &config.WarningCPUAvg.Th
		}

		if config.WarningCPUFullAvg300.IsSet {
			cpuCheck.Perfdata[psi.CPUFullAvg300].Warn = &config.WarningCPUFullAvg300.Th
		} else {
			cpuCheck.Perfdata[psi.CPUFullAvg300].Warn = &config.WarningCPUAvg.Th
		}

		// Critical thresholds
		if config.CriticalCPUFullAvg10.IsSet {
			cpuCheck.Perfdata[psi.CPUFullAvg10].Crit = &config.CriticalCPUFullAvg10.Th
		} else {
			cpuCheck.Perfdata[psi.CPUFullAvg10].Crit = &config.CriticalCPUAvg.Th
		}

		if config.CriticalCPUFullAvg60.IsSet {
			cpuCheck.Perfdata[psi.CPUFullAvg60].Crit = &config.CriticalCPUFullAvg60.Th
		} else {
			cpuCheck.Perfdata[psi.CPUFullAvg60].Crit = &config.CriticalCPUAvg.Th
		}

		if config.CriticalCPUFullAvg300.IsSet {
			cpuCheck.Perfdata[psi.CPUFullAvg300].Crit = &config.CriticalCPUFullAvg300.Th
		} else {
			cpuCheck.Perfdata[psi.CPUFullAvg300].Crit = &config.CriticalCPUAvg.Th
		}

		cpuFullSc := result.PartialResult{}
		_ = cpuFullSc.SetDefaultState(check.OK)

		if cpuCheck.Perfdata[psi.CPUFullAvg10].Warn.DoesViolate(psiCPU.Full.Avg10) ||
			cpuCheck.Perfdata[psi.CPUFullAvg60].Warn.DoesViolate(psiCPU.Full.Avg60) ||
			cpuCheck.Perfdata[psi.CPUFullAvg300].Warn.DoesViolate(psiCPU.Full.Avg300) {
			_ = cpuFullSc.SetState(check.Warning)
		}

		if cpuCheck.Perfdata[psi.CPUFullAvg10].Crit.DoesViolate(psiCPU.Full.Avg10) ||
			cpuCheck.Perfdata[psi.CPUFullAvg60].Crit.DoesViolate(psiCPU.Full.Avg60) ||
			cpuCheck.Perfdata[psi.CPUFullAvg300].Crit.DoesViolate(psiCPU.Full.Avg300) {
			_ = cpuFullSc.SetState(check.Critical)
		}

		cpuFullSc.Output = fmt.Sprintf("Full - Avg10: %.2f, Avg60: %.2f, Avg300: %.2f", psiCPU.Full.Avg10, psiCPU.Full.Avg60, psiCPU.Full.Avg300)
		cpuCheck.AddSubcheck(cpuFullSc)
	}

	if config.WarningCPUSomeAvg10.IsSet {
		cpuCheck.Perfdata[psi.CPUSomeAvg10].Warn = &config.WarningCPUSomeAvg10.Th
	} else {
		cpuCheck.Perfdata[psi.CPUSomeAvg10].Warn = &config.WarningCPUAvg.Th
	}

	if config.WarningCPUSomeAvg60.IsSet {
		cpuCheck.Perfdata[psi.CPUSomeAvg60].Warn = &config.WarningCPUSomeAvg60.Th
	} else {
		cpuCheck.Perfdata[psi.CPUSomeAvg60].Warn = &config.WarningCPUAvg.Th
	}

	if config.WarningCPUSomeAvg300.IsSet {
		cpuCheck.Perfdata[psi.CPUSomeAvg300].Warn = &config.WarningCPUSomeAvg300.Th
	} else {
		cpuCheck.Perfdata[psi.CPUSomeAvg300].Warn = &config.WarningCPUAvg.Th
	}

	// Critical thresholds
	if config.CriticalCPUSomeAvg10.IsSet {
		cpuCheck.Perfdata[psi.CPUSomeAvg10].Crit = &config.CriticalCPUSomeAvg10.Th
	} else {
		cpuCheck.Perfdata[psi.CPUSomeAvg10].Crit = &config.CriticalCPUAvg.Th
	}

	if config.CriticalCPUSomeAvg60.IsSet {
		cpuCheck.Perfdata[psi.CPUSomeAvg60].Crit = &config.CriticalCPUSomeAvg60.Th
	} else {
		cpuCheck.Perfdata[psi.CPUSomeAvg60].Crit = &config.CriticalCPUAvg.Th
	}

	if config.CriticalCPUSomeAvg300.IsSet {
		cpuCheck.Perfdata[psi.CPUSomeAvg300].Crit = &config.CriticalCPUSomeAvg300.Th
	} else {
		cpuCheck.Perfdata[psi.CPUSomeAvg300].Crit = &config.CriticalCPUAvg.Th
	}

	cpuSomeSc := result.PartialResult{}
	_ = cpuSomeSc.SetDefaultState(check.OK)

	if (cpuCheck.GetStatus() != check.Critical) && (cpuCheck.GetStatus() != check.Warning) {
		if cpuCheck.Perfdata[psi.CPUSomeAvg10].Warn.DoesViolate(psiCPU.Some.Avg10) ||
			cpuCheck.Perfdata[psi.CPUSomeAvg60].Warn.DoesViolate(psiCPU.Some.Avg60) ||
			cpuCheck.Perfdata[psi.CPUSomeAvg300].Warn.DoesViolate(psiCPU.Some.Avg300) {
			_ = cpuSomeSc.SetState(check.Warning)
		}
	}

	if cpuCheck.GetStatus() != check.Critical {
		if cpuCheck.Perfdata[psi.CPUSomeAvg10].Crit.DoesViolate(psiCPU.Some.Avg10) ||
			cpuCheck.Perfdata[psi.CPUSomeAvg60].Crit.DoesViolate(psiCPU.Some.Avg60) ||
			cpuCheck.Perfdata[psi.CPUSomeAvg300].Crit.DoesViolate(psiCPU.Some.Avg300) {
			_ = cpuSomeSc.SetState(check.Critical)
		}
	}

	cpuSomeSc.Output = fmt.Sprintf("Some - Avg10: %.2f, Avg60: %.2f, Avg300: %.2f", psiCPU.Some.Avg10, psiCPU.Some.Avg60, psiCPU.Some.Avg300)
	cpuCheck.AddSubcheck(cpuSomeSc)

	return cpuCheck
}

func checkPsiIoPressure(config *psiConfig) result.PartialResult {
	var ioCheck result.PartialResult
	_ = ioCheck.SetDefaultState(check.OK)
	ioCheck.Output = "IO"

	psiIo, err := psi.ReadIoPressure()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			_ = ioCheck.SetState(check.Unknown)
			ioCheck.Output = "IO pressure file not found. Perhaps the PSI interface is not active on this system? It might be necessary to change the kernel config"

			return ioCheck
		}

		check.ExitError(err)
	}

	ioCheck.Perfdata = *psiIo.Perfdata()

	//nolint:nestif
	if psiIo.FullPresent {
		// Warn thresholds
		if config.WarningIoFullAvg10.IsSet {
			ioCheck.Perfdata[psi.IoFullAvg10].Warn = &config.WarningIoFullAvg10.Th
		} else {
			ioCheck.Perfdata[psi.IoFullAvg10].Warn = &config.WarningIoAvg.Th
		}

		if config.WarningIoFullAvg60.IsSet {
			ioCheck.Perfdata[psi.IoFullAvg60].Warn = &config.WarningIoFullAvg60.Th
		} else {
			ioCheck.Perfdata[psi.IoFullAvg60].Warn = &config.WarningIoAvg.Th
		}

		if config.WarningIoFullAvg300.IsSet {
			ioCheck.Perfdata[psi.IoFullAvg300].Warn = &config.WarningIoFullAvg300.Th
		} else {
			ioCheck.Perfdata[psi.IoFullAvg300].Warn = &config.WarningIoAvg.Th
		}

		// Critical thresholds
		if config.CriticalIoFullAvg10.IsSet {
			ioCheck.Perfdata[psi.IoFullAvg10].Crit = &config.CriticalIoFullAvg10.Th
		} else {
			ioCheck.Perfdata[psi.IoFullAvg10].Crit = &config.CriticalIoAvg.Th
		}

		if config.CriticalIoFullAvg60.IsSet {
			ioCheck.Perfdata[psi.IoFullAvg60].Crit = &config.CriticalIoFullAvg60.Th
		} else {
			ioCheck.Perfdata[psi.IoFullAvg60].Crit = &config.CriticalIoAvg.Th
		}

		if config.CriticalIoFullAvg300.IsSet {
			ioCheck.Perfdata[psi.IoFullAvg300].Crit = &config.CriticalIoFullAvg300.Th
		} else {
			ioCheck.Perfdata[psi.IoFullAvg300].Crit = &config.CriticalIoAvg.Th
		}

		ioFullSc := result.PartialResult{}
		_ = ioFullSc.SetDefaultState(check.OK)

		if ioCheck.Perfdata[psi.IoFullAvg10].Warn.DoesViolate(psiIo.Full.Avg10) ||
			ioCheck.Perfdata[psi.IoFullAvg60].Warn.DoesViolate(psiIo.Full.Avg60) ||
			ioCheck.Perfdata[psi.IoFullAvg300].Warn.DoesViolate(psiIo.Full.Avg300) {
			_ = ioFullSc.SetState(check.Warning)
		}

		if ioCheck.Perfdata[psi.IoFullAvg10].Crit.DoesViolate(psiIo.Full.Avg10) ||
			ioCheck.Perfdata[psi.IoFullAvg60].Crit.DoesViolate(psiIo.Full.Avg60) ||
			ioCheck.Perfdata[psi.IoFullAvg300].Crit.DoesViolate(psiIo.Full.Avg300) {
			_ = ioFullSc.SetState(check.Critical)
		}

		ioFullSc.Output = fmt.Sprintf("Full - Avg10: %.2f, Avg60: %.2f, Avg300: %.2f", psiIo.Full.Avg10, psiIo.Full.Avg60, psiIo.Full.Avg300)
		ioCheck.AddSubcheck(ioFullSc)
	}

	if config.WarningIoSomeAvg10.IsSet {
		ioCheck.Perfdata[psi.IoSomeAvg10].Warn = &config.WarningIoSomeAvg10.Th
	} else {
		ioCheck.Perfdata[psi.IoSomeAvg10].Warn = &config.WarningIoAvg.Th
	}

	if config.WarningIoSomeAvg60.IsSet {
		ioCheck.Perfdata[psi.IoSomeAvg60].Warn = &config.WarningIoSomeAvg60.Th
	} else {
		ioCheck.Perfdata[psi.IoSomeAvg60].Warn = &config.WarningIoAvg.Th
	}

	if config.WarningIoSomeAvg300.IsSet {
		ioCheck.Perfdata[psi.IoSomeAvg300].Warn = &config.WarningIoSomeAvg300.Th
	} else {
		ioCheck.Perfdata[psi.IoSomeAvg300].Warn = &config.WarningIoAvg.Th
	}

	if config.CriticalIoSomeAvg10.IsSet {
		ioCheck.Perfdata[psi.IoSomeAvg10].Crit = &config.CriticalIoSomeAvg10.Th
	} else {
		ioCheck.Perfdata[psi.IoSomeAvg10].Crit = &config.CriticalIoAvg.Th
	}

	if config.CriticalIoSomeAvg60.IsSet {
		ioCheck.Perfdata[psi.IoSomeAvg60].Crit = &config.CriticalIoSomeAvg60.Th
	} else {
		ioCheck.Perfdata[psi.IoSomeAvg60].Crit = &config.CriticalIoAvg.Th
	}

	if config.CriticalIoSomeAvg300.IsSet {
		ioCheck.Perfdata[psi.IoSomeAvg300].Crit = &config.CriticalIoSomeAvg300.Th
	} else {
		ioCheck.Perfdata[psi.IoSomeAvg300].Crit = &config.CriticalIoAvg.Th
	}

	ioSomeSc := result.PartialResult{}
	_ = ioSomeSc.SetDefaultState(check.OK)

	if (ioCheck.GetStatus() != check.Critical) && (ioCheck.GetStatus() != check.Warning) {
		if ioCheck.Perfdata[psi.IoSomeAvg10].Warn.DoesViolate(psiIo.Some.Avg10) ||
			ioCheck.Perfdata[psi.IoSomeAvg60].Warn.DoesViolate(psiIo.Some.Avg60) ||
			ioCheck.Perfdata[psi.IoSomeAvg300].Warn.DoesViolate(psiIo.Some.Avg300) {
			_ = ioSomeSc.SetState(check.Warning)
		}
	}

	if ioCheck.GetStatus() != check.Critical {
		if ioCheck.Perfdata[psi.IoSomeAvg10].Crit.DoesViolate(psiIo.Some.Avg10) ||
			ioCheck.Perfdata[psi.IoSomeAvg60].Crit.DoesViolate(psiIo.Some.Avg60) ||
			ioCheck.Perfdata[psi.IoSomeAvg300].Crit.DoesViolate(psiIo.Some.Avg300) {
			_ = ioSomeSc.SetState(check.Critical)
		}
	}

	ioSomeSc.Output = fmt.Sprintf("Some - Avg10: %.2f, Avg60: %.2f, Avg300: %.2f", psiIo.Some.Avg10, psiIo.Some.Avg60, psiIo.Some.Avg300)
	ioCheck.AddSubcheck(ioSomeSc)

	return ioCheck
}

func checkPsiMemoryPressure(config *psiConfig) result.PartialResult {
	var memoryCheck result.PartialResult
	_ = memoryCheck.SetDefaultState(check.OK)
	memoryCheck.Output = "Memory"

	psiMemory, err := psi.ReadMemoryPressure()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			_ = memoryCheck.SetState(check.Unknown)
			memoryCheck.Output = "IO pressure file not found. Perhaps the PSI interface is not active on this system? It might be necessary to change the kernel config"

			return memoryCheck
		}

		check.ExitError(err)
	}

	memoryCheck.Perfdata = *psiMemory.Perfdata()

	//nolint:nestif
	if psiMemory.FullPresent {
		// Warn thresholds
		if config.WarningMemoryFullAvg10.IsSet {
			memoryCheck.Perfdata[psi.MemoryFullAvg10].Warn = &config.WarningMemoryFullAvg10.Th
		} else {
			memoryCheck.Perfdata[psi.MemoryFullAvg10].Warn = &config.WarningMemoryAvg.Th
		}

		if config.WarningMemoryFullAvg60.IsSet {
			memoryCheck.Perfdata[psi.MemoryFullAvg60].Warn = &config.WarningMemoryFullAvg60.Th
		} else {
			memoryCheck.Perfdata[psi.MemoryFullAvg60].Warn = &config.WarningMemoryAvg.Th
		}

		if config.WarningMemoryFullAvg300.IsSet {
			memoryCheck.Perfdata[psi.MemoryFullAvg300].Warn = &config.WarningMemoryFullAvg300.Th
		} else {
			memoryCheck.Perfdata[psi.MemoryFullAvg300].Warn = &config.WarningMemoryAvg.Th
		}

		// Critical thresholds
		if config.CriticalMemoryFullAvg10.IsSet {
			memoryCheck.Perfdata[psi.MemoryFullAvg10].Crit = &config.CriticalMemoryFullAvg10.Th
		} else {
			memoryCheck.Perfdata[psi.MemoryFullAvg10].Crit = &config.CriticalMemoryAvg.Th
		}

		if config.CriticalMemoryFullAvg60.IsSet {
			memoryCheck.Perfdata[psi.MemoryFullAvg60].Crit = &config.CriticalMemoryFullAvg60.Th
		} else {
			memoryCheck.Perfdata[psi.MemoryFullAvg60].Crit = &config.CriticalMemoryAvg.Th
		}

		if config.CriticalMemoryFullAvg300.IsSet {
			memoryCheck.Perfdata[psi.MemoryFullAvg300].Crit = &config.CriticalMemoryFullAvg300.Th
		} else {
			memoryCheck.Perfdata[psi.MemoryFullAvg300].Crit = &config.CriticalMemoryAvg.Th
		}

		memoryFullSc := result.PartialResult{}
		_ = memoryFullSc.SetDefaultState(check.OK)

		if memoryCheck.Perfdata[psi.MemoryFullAvg10].Warn.DoesViolate(psiMemory.Full.Avg10) ||
			memoryCheck.Perfdata[psi.MemoryFullAvg60].Warn.DoesViolate(psiMemory.Full.Avg60) ||
			memoryCheck.Perfdata[psi.MemoryFullAvg300].Warn.DoesViolate(psiMemory.Full.Avg300) {
			_ = memoryFullSc.SetState(check.Warning)
		}

		if memoryCheck.Perfdata[psi.MemoryFullAvg10].Crit.DoesViolate(psiMemory.Full.Avg10) ||
			memoryCheck.Perfdata[psi.MemoryFullAvg60].Crit.DoesViolate(psiMemory.Full.Avg60) ||
			memoryCheck.Perfdata[psi.MemoryFullAvg300].Crit.DoesViolate(psiMemory.Full.Avg300) {
			_ = memoryFullSc.SetState(check.Critical)
		}

		memoryFullSc.Output = fmt.Sprintf("Full - Avg10: %.2f, Avg60: %.2f, Avg300: %.2f", psiMemory.Full.Avg10, psiMemory.Full.Avg60, psiMemory.Full.Avg300)
		memoryCheck.AddSubcheck(memoryFullSc)
	}

	if config.WarningMemorySomeAvg10.IsSet {
		memoryCheck.Perfdata[psi.MemorySomeAvg10].Warn = &config.WarningMemorySomeAvg10.Th
	} else {
		memoryCheck.Perfdata[psi.MemorySomeAvg10].Warn = &config.WarningMemoryAvg.Th
	}

	if config.WarningMemorySomeAvg60.IsSet {
		memoryCheck.Perfdata[psi.MemorySomeAvg60].Warn = &config.WarningMemorySomeAvg60.Th
	} else {
		memoryCheck.Perfdata[psi.MemorySomeAvg60].Warn = &config.WarningMemoryAvg.Th
	}

	if config.WarningMemorySomeAvg300.IsSet {
		memoryCheck.Perfdata[psi.MemorySomeAvg300].Warn = &config.WarningMemorySomeAvg300.Th
	} else {
		memoryCheck.Perfdata[psi.MemorySomeAvg300].Warn = &config.WarningMemoryAvg.Th
	}

	// Critical thresholds
	if config.CriticalMemorySomeAvg10.IsSet {
		memoryCheck.Perfdata[psi.MemorySomeAvg10].Crit = &config.CriticalMemorySomeAvg10.Th
	} else {
		memoryCheck.Perfdata[psi.MemorySomeAvg10].Crit = &config.CriticalMemoryAvg.Th
	}

	if config.CriticalMemorySomeAvg60.IsSet {
		memoryCheck.Perfdata[psi.MemorySomeAvg60].Crit = &config.CriticalMemorySomeAvg60.Th
	} else {
		memoryCheck.Perfdata[psi.MemorySomeAvg60].Crit = &config.CriticalMemoryAvg.Th
	}

	if config.CriticalMemorySomeAvg300.IsSet {
		memoryCheck.Perfdata[psi.MemorySomeAvg300].Crit = &config.CriticalMemorySomeAvg300.Th
	} else {
		memoryCheck.Perfdata[psi.MemorySomeAvg300].Crit = &config.CriticalMemoryAvg.Th
	}

	memorySomeSc := result.PartialResult{}
	_ = memorySomeSc.SetDefaultState(check.OK)

	if (memoryCheck.GetStatus() != check.Critical) && (memoryCheck.GetStatus() != check.Warning) {
		if memoryCheck.Perfdata[psi.MemorySomeAvg10].Warn.DoesViolate(psiMemory.Some.Avg10) ||
			memoryCheck.Perfdata[psi.MemorySomeAvg60].Warn.DoesViolate(psiMemory.Some.Avg60) ||
			memoryCheck.Perfdata[psi.MemorySomeAvg300].Warn.DoesViolate(psiMemory.Some.Avg300) {
			_ = memorySomeSc.SetState(check.Warning)
		}
	}

	if memoryCheck.GetStatus() != check.Critical {
		if memoryCheck.Perfdata[psi.MemorySomeAvg10].Crit.DoesViolate(psiMemory.Some.Avg10) ||
			memoryCheck.Perfdata[psi.MemorySomeAvg60].Crit.DoesViolate(psiMemory.Some.Avg60) ||
			memoryCheck.Perfdata[psi.MemorySomeAvg300].Crit.DoesViolate(psiMemory.Some.Avg300) {
			_ = memorySomeSc.SetState(check.Critical)
		}
	}

	memorySomeSc.Output = fmt.Sprintf("Some - Avg10: %.2f, Avg60: %.2f, Avg300: %.2f", psiMemory.Some.Avg10, psiMemory.Some.Avg60, psiMemory.Some.Avg300)
	memoryCheck.AddSubcheck(memorySomeSc)

	return memoryCheck
}
