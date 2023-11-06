package load

import (
	"fmt"

	"github.com/NETWAYS/go-check/perfdata"
	"github.com/shirou/gopsutil/v3/load"
)

type Load struct {
	LoadAvg   *load.AvgStat
	CliConfig *LoadConfig
}

func GetActualLoadValues() (l *Load, err error) {
	l = &Load{}

	loadAvg, err := load.Avg()
	if err != nil {
		return nil, fmt.Errorf("could not evaluate load average: %w", err)
	}

	l.LoadAvg = loadAvg

	return l, err
}

func (l *Load) GetOutput() (output string) {
	output += fmt.Sprintf("load average: %.2f, %.2f, %.2f",
		l.LoadAvg.Load1,
		l.LoadAvg.Load5,
		l.LoadAvg.Load15)

	return output
}

func (l *Load) GetPerfData() perfdata.PerfdataList {
	perfList := perfdata.PerfdataList{
		{
			Label: "load1",
			Value: l.LoadAvg.Load1,
			Warn:  &l.CliConfig.Load1Th.Warn.Th,
			Crit:  &l.CliConfig.Load1Th.Crit.Th,
			Min:   nil,
			Max:   nil,
		},
		{
			Label: "load5",
			Value: l.LoadAvg.Load5,
			Warn:  &l.CliConfig.Load5Th.Warn.Th,
			Crit:  &l.CliConfig.Load5Th.Crit.Th,
			Min:   nil,
			Max:   nil,
		},
		{
			Label: "load15",
			Value: l.LoadAvg.Load15,
			Warn:  &l.CliConfig.Load15Th.Warn.Th,
			Crit:  &l.CliConfig.Load15Th.Crit.Th,
			Min:   nil,
			Max:   nil,
		},
	}

	return perfList
}
