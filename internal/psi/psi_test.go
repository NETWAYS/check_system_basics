package psi

import (
	"reflect"
	"testing"

	"github.com/NETWAYS/go-check/perfdata"
)

func TestPressureValueString(t *testing.T) {
	pv := PressureValue{
		Avg10:  1.5,
		Avg60:  2.5,
		Avg300: 3.5,
	}

	// Expected
	expected := perfdata.PerfdataList{}
	expected.Add(&perfdata.Perfdata{Label: "preavg10", Value: 1.5, Min: 0, Max: 100, Uom: "%"})
	expected.Add(&perfdata.Perfdata{Label: "preavg60", Value: 2.5, Min: 0, Max: 100, Uom: "%"})
	expected.Add(&perfdata.Perfdata{Label: "preavg300", Value: 3.5, Min: 0, Max: 100, Uom: "%"})
	expected.Add(&perfdata.Perfdata{Label: "pretotal", Value: uint64(0), Min: 0, Uom: "c"})

	if !reflect.DeepEqual(&expected, pv.Perfdata("pre")) {
		t.Fatalf("expected %v, got %v", &expected, pv.Perfdata("pre"))
	}
}

func TestPressureElementString(t *testing.T) {

	pecpu := PressureElement{
		Some:        PressureValue{Avg10: 0.1, Avg60: 0.6, Avg300: 0.3, Total: 0},
		Full:        PressureValue{Avg10: 0.1, Avg60: 0.6, Avg300: 0.3, Total: 0},
		FullPresent: true,
	}

	// Expected
	pecpuexpected := perfdata.PerfdataList{}
	pecpuexpected.Add(&perfdata.Perfdata{Label: "cpu-some-avg10", Value: float64(0.1), Min: 0, Max: 100, Uom: "%"})
	pecpuexpected.Add(&perfdata.Perfdata{Label: "cpu-some-avg60", Value: float64(0.6), Min: 0, Max: 100, Uom: "%"})
	pecpuexpected.Add(&perfdata.Perfdata{Label: "cpu-some-avg300", Value: float64(0.3), Min: 0, Max: 100, Uom: "%"})
	pecpuexpected.Add(&perfdata.Perfdata{Label: "cpu-some-total", Value: uint64(0), Min: 0, Uom: "c"})
	pecpuexpected.Add(&perfdata.Perfdata{Label: "cpu-full-avg10", Value: float64(0.1), Min: 0, Max: 100, Uom: "%"})
	pecpuexpected.Add(&perfdata.Perfdata{Label: "cpu-full-avg60", Value: float64(0.6), Min: 0, Max: 100, Uom: "%"})
	pecpuexpected.Add(&perfdata.Perfdata{Label: "cpu-full-avg300", Value: float64(0.3), Min: 0, Max: 100, Uom: "%"})
	pecpuexpected.Add(&perfdata.Perfdata{Label: "cpu-full-total", Value: uint64(0), Min: 0, Uom: "c"})

	if !reflect.DeepEqual(&pecpuexpected, pecpu.Perfdata()) {
		t.Fatalf("expected %v, got %v", &pecpuexpected, pecpu.Perfdata())
	}

	peio := PressureElement{
		Some:        PressureValue{Avg10: 0.1, Avg60: 0.6, Avg300: 0.3, Total: 0},
		Full:        PressureValue{Avg10: 0.1, Avg60: 0.6, Avg300: 0.3, Total: 0},
		FullPresent: false,
		Type:        2,
	}

	// Expected
	peioexpected := perfdata.PerfdataList{}
	peioexpected.Add(&perfdata.Perfdata{Label: "io-some-avg10", Value: float64(0.1), Min: 0, Max: 100, Uom: "%"})
	peioexpected.Add(&perfdata.Perfdata{Label: "io-some-avg60", Value: float64(0.6), Min: 0, Max: 100, Uom: "%"})
	peioexpected.Add(&perfdata.Perfdata{Label: "io-some-avg300", Value: float64(0.3), Min: 0, Max: 100, Uom: "%"})
	peioexpected.Add(&perfdata.Perfdata{Label: "io-some-total", Value: uint64(0), Min: 0, Uom: "c"})

	if !reflect.DeepEqual(&peioexpected, peio.Perfdata()) {
		t.Fatalf("expected %v, got %v", &peioexpected, peio.Perfdata())
	}
}

func TestReadPressureFile(t *testing.T) {
	cpuPressure, err := readPressureFile("testdata/cpu")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expectedResult := PressureElement{
		Some:        PressureValue{Avg10: 0.0, Avg60: 0.08, Avg300: 0.05, Total: 3391622},
		Full:        PressureValue{Avg10: 0.0, Avg60: 0.00, Avg300: 0.00, Total: 0},
		FullPresent: true,
	}

	if !reflect.DeepEqual(&expectedResult, cpuPressure) {
		t.Fatalf("expected %v, got %v", &expectedResult, cpuPressure)
	}
}
