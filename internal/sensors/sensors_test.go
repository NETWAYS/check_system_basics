package sensors

import (
	"testing"

	"github.com/NETWAYS/go-check"
	"github.com/NETWAYS/go-check/perfdata"
)

func TestSensorAndDeviceString(t *testing.T) {
	s := Sensor{
		Name:  "testname",
		Path:  "testpath",
		Alarm: false,
		Perfdata: perfdata.Perfdata{
			Label: "test",
			Value: 10.0,
			Uom:   "%",
			Warn:  &check.Threshold{Upper: 80},
			Crit:  &check.Threshold{Upper: 90},
			Min:   0,
			Max:   100,
		},
	}

	if "testname - 10%" != s.String() {
		t.Fatalf("expected %v, got %v", "testname - 10%", s.String())
	}

	d := Device{
		Name:    "test",
		Sensors: []Sensor{s},
	}

	expected := "test: testname - 10%;"
	if expected != d.String() {
		t.Fatalf("expected %v, got %v", expected, d.String())
	}
}

func TestReadSensorDataSingle(t *testing.T) {
	sensors, err := readSensorData("testdata/01")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := "tempSensor=20C"
	if expected != sensors[0].Perfdata.String() {
		t.Fatalf("expected %v, got %v", expected, sensors[0].Perfdata.String())
	}
}

func TestReadSensorDataMultiple(t *testing.T) {
	sensors, err := readSensorData("testdata/02")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := "tempSensor=20C"
	if expected != sensors[0].Perfdata.String() {
		t.Fatalf("expected %v, got %v", expected, sensors[0].Perfdata.String())
	}

	expected = "testDevice2_temp2=20C;~:30;~:40"
	if expected != sensors[1].Perfdata.String() {
		t.Fatalf("expected %v, got %v", expected, sensors[1].Perfdata.String())
	}

	expected = "tempSensor3=20C;~:30;~:60"
	if expected != sensors[2].Perfdata.String() {
		t.Fatalf("expected %v, got %v", expected, sensors[2].Perfdata.String())
	}
}
