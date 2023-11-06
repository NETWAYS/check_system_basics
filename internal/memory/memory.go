package memory

import (
	"fmt"

	"github.com/shirou/gopsutil/v3/mem"
)

type Mem struct {
	VirtMem *mem.VirtualMemoryStat
	// MemUsedPercentage      float64 // already in VirtMem
	MemAvailablePercentage float64

	SwapInfo *mem.SwapMemoryStat
}

func LoadMemStat() (m *Mem, err error) {
	m = &Mem{}

	virtMem, err := mem.VirtualMemory()
	if err != nil {
		return nil, fmt.Errorf("could not load virtual memory stat: %w", err)
	}

	m.VirtMem = virtMem

	m.MemAvailablePercentage = float64(virtMem.Available) / float64(virtMem.Total/100)

	swapMem, err := mem.SwapMemory()
	if err != nil {
		return nil, fmt.Errorf("could not load swap memory stat: %w", err)
	}

	m.SwapInfo = swapMem

	return
}
