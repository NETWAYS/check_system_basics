package psi

import (
	"bufio"
	"os"
	"strconv"
	"strings"

	"github.com/NETWAYS/go-check/perfdata"
)

const (
	CPUSomeAvg10 uint = iota
	CPUSomeAvg60
	CPUSomeAvg300
	CPUFullAvg10
	CPUFullAvg60
	CPUFullAvg300
)

const (
	IoSomeAvg10 uint = iota
	IoSomeAvg60
	IoSomeAvg300
	IoFullAvg10
	IoFullAvg60
	IoFullAvg300
)

const (
	MemorySomeAvg10 uint = iota
	MemorySomeAvg60
	MemorySomeAvg300
	MemoryFullAvg10
	MemoryFullAvg60
	MemoryFullAvg300
)

type PressureValue struct {
	Avg10  float64
	Avg60  float64
	Avg300 float64
	Total  uint64
}

type PressureElement struct {
	Some        PressureValue
	Full        PressureValue
	FullPresent bool
	Type        PressureType
}

type PressureType uint

const (
	cpu PressureType = iota
	memory
	io
)

func (p *PressureValue) Perfdata(prefix string) *perfdata.PerfdataList {
	var ret perfdata.PerfdataList

	avg10 := perfdata.Perfdata{}
	avg10.Label = prefix + "avg10"
	avg10.Value = p.Avg10
	avg10.Uom = "%"
	avg10.Min = 0
	avg10.Max = 100
	ret.Add(&avg10)

	avg60 := perfdata.Perfdata{}
	avg60.Label = prefix + "avg60"
	avg60.Value = p.Avg60
	avg60.Uom = "%"
	avg60.Min = 0
	avg60.Max = 100
	ret.Add(&avg60)

	avg300 := perfdata.Perfdata{}
	avg300.Label = prefix + "avg300"
	avg300.Value = p.Avg300
	avg300.Uom = "%"
	avg300.Min = 0
	avg300.Max = 100
	ret.Add(&avg300)

	total := perfdata.Perfdata{}
	total.Label = prefix + "total"
	total.Value = p.Total
	total.Min = 0
	total.Uom = "c"
	ret.Add(&total)

	return &ret
}

func (p *PressureElement) Perfdata() *perfdata.PerfdataList {
	switch p.Type {
	case cpu:
		tmp := *p.Some.Perfdata("cpu-some-")
		if p.FullPresent {
			tmp = append(tmp, *p.Full.Perfdata("cpu-full-")...)
		}

		return &tmp
	case io:
		tmp := *p.Some.Perfdata("io-some-")
		if p.FullPresent {
			tmp = append(tmp, *p.Full.Perfdata("io-full-")...)
		}

		return &tmp
	case memory:
		tmp := *p.Some.Perfdata("memory-some-")
		if p.FullPresent {
			tmp = append(tmp, *p.Full.Perfdata("memory-full-")...)
		}

		return &tmp
	default:
		return nil
	}
}

func parsePressureValue(val string) (PressureValue, error) {
	tmp := strings.Split(val, " ")

	var result PressureValue

	tmpString := strings.Split(tmp[1], "=")

	tmpFloat, err := strconv.ParseFloat(tmpString[1], 64)
	if err != nil {
		return PressureValue{}, err
	}

	result.Avg10 = tmpFloat

	tmpString = strings.Split(tmp[2], "=")

	tmpFloat, err = strconv.ParseFloat(tmpString[1], 64)
	if err != nil {
		return PressureValue{}, err
	}

	result.Avg60 = tmpFloat

	tmpString = strings.Split(tmp[3], "=")

	tmpFloat, err = strconv.ParseFloat(tmpString[1], 64)
	if err != nil {
		return PressureValue{}, err
	}

	result.Avg300 = tmpFloat

	tmpString = strings.Split(tmp[4], "=")

	tmpTotal, err := strconv.ParseUint(tmpString[1], 10, 64)
	if err != nil {
		return PressureValue{}, err
	}

	result.Total = tmpTotal

	return result, nil
}

func readPressureFile(pressurePath string) (*PressureElement, error) {
	readFile, err := os.Open(pressurePath)

	if err != nil {
		return nil, err
	}

	defer readFile.Close()

	fileScanner := bufio.NewScanner(readFile)

	fileScanner.Split(bufio.ScanLines)

	lines := make([]string, 0)
	for fileScanner.Scan() {
		lines = append(lines, fileScanner.Text())
	}

	var result PressureElement

	// TODO Is this correct?
	// nolint: gosec
	tmpPval, err := parsePressureValue(lines[0])
	if err != nil {
		return nil, err
	}

	result.Some = tmpPval

	// TODO Is this correct?
	// nolint: gosec
	if len(lines) > 1 {
		tmpPval, err = parsePressureValue(lines[1])
		if err != nil {
			return nil, err
		}

		result.Full = tmpPval
		result.FullPresent = true
	} else {
		result.FullPresent = false
	}

	return &result, nil
}

func ReadCPUPressure() (*PressureElement, error) {
	tmp, err := readPressureFile("/proc/pressure/cpu")
	if err != nil {
		return nil, err
	}

	tmp.Type = cpu

	return tmp, nil
}
func ReadIoPressure() (*PressureElement, error) {
	tmp, err := readPressureFile("/proc/pressure/io")
	if err != nil {
		return nil, err
	}

	tmp.Type = io

	return tmp, nil
}
func ReadMemoryPressure() (*PressureElement, error) {
	tmp, err := readPressureFile("/proc/pressure/memory")
	if err != nil {
		return nil, err
	}

	tmp.Type = memory

	return tmp, nil
}
