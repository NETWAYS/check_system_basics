package sensors

import (
	"testing"

	"github.com/NETWAYS/go-check"
	"github.com/NETWAYS/go-check/perfdata"
	"github.com/stretchr/testify/assert"
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

	assert.Equal(t, "testname - 10%", s.String())

	d := Device{
		Name:    "test",
		Sensors: []Sensor{s},
	}

	assert.Equal(t, "test: testname - 10%;", d.String())
}

func TestReadSensorDataSingle(t *testing.T) {
	sensors, err := readSensorData("testdata/01")

	assert.Equal(t, nil, err)

	assert.Equal(t, "tempSensor=20C", sensors[0].Perfdata.String())
}

func TestReadSensorDataMultiple(t *testing.T) {
	sensors, err := readSensorData("testdata/02")
	assert.Equal(t, nil, err)

	assert.Equal(t, "tempSensor=20C", sensors[0].Perfdata.String())
	assert.Equal(t, "testDevice2_temp2=20C;~:30;~:40", sensors[1].Perfdata.String())
	assert.Equal(t, "tempSensor3=20C;~:30;~:60", sensors[2].Perfdata.String())
}
