package sensors

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/NETWAYS/go-check"
	"github.com/NETWAYS/go-check/perfdata"
)

/*
 * References: http://blog.foool.net/wp-content/uploads/linuxdocs/hwmon.pdf
 * https://www.kernel.org/doc/html/latest/hwmon/hwmon-kernel-api.html
 */

type Sensor struct {
	Name     string
	Path     string
	Alarm    bool
	Perfdata perfdata.Perfdata
}

type Device struct {
	Name    string
	Sensors []Sensor
}

const (
	inputFileSuffix         string = "_input"
	critThresholdFileSuffix string = "_crit"
	lowestValueFileSuffix   string = "_lowest"
	highestValueFileSuffix  string = "_highest"
	maxValueFileSuffix      string = "_max"
)

func (d *Device) String() string {
	result := d.Name
	result += ": "

	for idx := range d.Sensors {
		result += d.Sensors[idx].String()
		if idx < len(d.Sensors) {
			result += ";"
		}
	}

	return result
}

func (s *Sensor) String() string {
	return fmt.Sprintf("%s - %v%s", s.Name, s.Perfdata.Value, s.Perfdata.Uom)
}

func GetDefaultDevices() ([]Device, error) {
	return GetDevices("/sys/class/hwmon")
}

func GetDevices(hwmonPath string) ([]Device, error) {
	/*
	 * hwmonPath should look like:
	 * hwmon0 hwmon1 hwmon2 ...
	 */
	files, err := filepath.Glob(hwmonPath + "/*")

	if err != nil {
		return []Device{}, err
	}

	devices := make([]Device, len(files))

	/* Iterate through the hwmon* files, detect the sensors and ingest them */
	for i, file := range files {
		/*
		 * Files with details of the device (name, etc)
		 * and the sensors (tempN_label, tempN_input, ...)
		 */
		subfiles, err := filepath.Glob(file + "/*")
		if err != nil {
			return devices, err
		}

		for _, subfile := range subfiles {
			if !strings.Contains(subfile, "name") {
				continue
			}

			bytes, err := os.ReadFile(subfile)
			if err != nil {
				return devices, err
			}

			devices[i].Name = strings.TrimSpace(string(bytes))

			devices[i].Sensors, err = readSensorData(file)
			if err != nil {
				return devices, err
			}
		}
	}

	return devices, nil
}

func filterBasePathList(files []string, prefix string) []string {
	result := make([]string, 0)

	for _, file := range files {
		if strings.Contains(path.Base(file), prefix) {
			result = append(result, file)
		}
	}

	return result
}

// nolint: funlen
func readSensorData(devicePath string) ([]Sensor, error) {
	// Possible sensor types
	var m = []string{
		"power",
		"pwm",
		"temp",
		"freq",
		"fan",
		"in",
		"curr",
		"energy",
		"humidity",
		"intrusion",
	}

	sensors := make([]Sensor, 0)

	// Get the files for the device
	dirEntryList, err := filepath.Glob(devicePath + "/*")
	if err != nil {
		return sensors, err
	}

	fileList := make([]string, 0)

	for idx := range dirEntryList {
		fileInfo, erro := os.Stat(dirEntryList[idx])
		if erro != nil {
			return sensors, erro
		}

		if !fileInfo.IsDir() {
			fileList = append(fileList, dirEntryList[idx])
		}
	}

	// Read the device name
	deviceName, err := readStringFromFile(devicePath + "/name")
	if err != nil {
		deviceName = filepath.Base(devicePath)
	}

	// Iterate over the files in the device directory and filter
	// the sensors out
	for _, key := range m {
		sensorFiles := filterBasePathList(fileList, key)
		if len(sensorFiles) == 0 {
			continue
		}

		// Ok, so we have the files for sensor type *key* now
		var idx int
		if key == "in" || key == "intrusion" {
			// voltage and intrusion detection start with 0
			idx = 0
		} else {
			idx = 1
		}

		for ; ; idx++ {
			tmp := filterBasePathList(sensorFiles, strconv.Itoa(idx))

			if len(tmp) == 0 {
				// No files found, so no more sensors
				break
			}
			// Now we should have all files for that specific sensor

			switch key {
			case "in":
				sensor, err := readVoltageSensor(deviceName, devicePath, idx)
				if err != nil {
					continue
				}

				sensors = append(sensors, sensor)

			case "fan":
				sensor, err := readFanSensor(devicePath, idx)
				if err != nil {
					continue
				}

				sensors = append(sensors, sensor)

			case "pwm":
				sensor, err := readPwmSensor(devicePath, idx)
				if err != nil {
					continue
				}

				sensors = append(sensors, sensor)

			case "temp":
				sensor, err := readTempSensor(deviceName, devicePath, idx)
				if err != nil {
					continue
				}

				sensors = append(sensors, sensor)

			case "curr":
				sensor, err := readCurrSensor(devicePath, idx)
				if err != nil {
					continue
				}

				sensors = append(sensors, sensor)

			case "power":
				sensor, err := readPowerSensor(deviceName, devicePath, idx)
				if err != nil {
					continue
				}

				sensors = append(sensors, sensor)

			case "energy":
				sensor, err := readEnergySensor(devicePath, idx)
				if err != nil {
					continue
				}

				sensors = append(sensors, sensor)

			case "humidity":
				sensor, err := readHumiditySensor(devicePath, idx)
				if err != nil {
					continue
				}

				sensors = append(sensors, sensor)

			default:
				continue
			}
		}
	}

	return sensors, nil
}

func readHumiditySensor(devicePath string, index int) (Sensor, error) {
	basePath := devicePath + "/humidity" + strconv.Itoa(index)

	var sensor Sensor

	// Look for label
	label := readSensorLabel(basePath)
	sensor.Name = label
	sensor.Perfdata.Label = label

	// Look for input (the actual value)
	value, err := readIntFromFile(basePath + inputFileSuffix)

	if err != nil {
		return sensor, err
	}

	sensor.Perfdata.Value = value // ???

	return sensor, nil
}

func readEnergySensor(devicePath string, index int) (Sensor, error) {
	basePath := devicePath + "/energy" + strconv.Itoa(index)

	var sensor Sensor

	// Look for label
	label := readSensorLabel(basePath)
	sensor.Name = label
	sensor.Perfdata.Label = label

	// Look for input (the actual value)
	value, err := readIntFromFile(basePath + inputFileSuffix)
	if err != nil {
		return sensor, err
	}

	sensor.Perfdata.Value = float64(value) / 1000 // micro Joule
	sensor.Perfdata.Uom = "Ws"

	return sensor, nil
}

func readCurrSensor(devicePath string, index int) (Sensor, error) {
	basePath := devicePath + "/curr" + strconv.Itoa(index)

	var sensor Sensor

	// Look for label
	label := readSensorLabel(basePath)
	sensor.Name = label
	sensor.Perfdata.Label = label

	// Look for input (the actual value)
	value, err := readIntFromFile(basePath + inputFileSuffix)
	if err != nil {
		return sensor, err
	}

	sensor.Perfdata.Value = float64(value) / 1000 // Milli Ampere
	sensor.Perfdata.Uom = "A"

	// == Min
	// Is there a currN_lowest file? Use it for max value
	value, err = readIntFromFile(basePath + lowestValueFileSuffix)
	if err == nil {
		sensor.Perfdata.Min = float64(value) / 1000
	}

	// == Max
	// Is there a currN_highest file? Use it for max value
	value, err = readIntFromFile(basePath + highestValueFileSuffix)
	if err == nil {
		sensor.Perfdata.Max = float64(value) / 1000
	}

	// Crit threshold
	tmp := check.Threshold{
		Lower:  check.NegInf,
		Upper:  check.PosInf,
		Inside: false,
	}
	critPresent := false
	// Is there a currN_lcrit file? If yes, use that as lower critical
	value, err = readIntFromFile(basePath + "_lcrit")
	if err == nil {
		tmp.Lower = float64(value) / 1000
		critPresent = true
	}
	// Is there a currN_crit file? If yes, use that as upper critical
	value, err = readIntFromFile(basePath + critThresholdFileSuffix)
	if err == nil {
		tmp.Upper = float64(value) / 1000
		critPresent = true
	}

	if critPresent {
		sensor.Perfdata.Crit = &tmp
	}

	// Alarm
	sensor.Alarm = readSensorAlarm(basePath)

	return sensor, nil
}

func readPwmSensor(devicePath string, index int) (Sensor, error) {
	basePath := devicePath + "/pwm" + strconv.Itoa(index)

	var sensor Sensor

	// Look for label
	label := readSensorLabel(basePath)
	sensor.Name = label
	sensor.Perfdata.Label = label

	// Look for input (the actual value)
	value, err := readIntFromFile(basePath + "_freq")
	if err != nil {
		return sensor, err
	}

	sensor.Perfdata.Value = float64(value) // Hertz

	return sensor, nil
}

func readFanSensor(devicePath string, index int) (Sensor, error) {
	basePath := devicePath + "/fan" + strconv.Itoa(index)

	var sensor Sensor

	// Look for label
	label := readSensorLabel(basePath)
	sensor.Name = label
	sensor.Perfdata.Label = label

	// Look for input (the actual value)
	value, err := readIntFromFile(basePath + inputFileSuffix)
	if err != nil {
		return sensor, err
	}

	sensor.Perfdata.Value = float64(value) // Rounds per Minute

	// == Min
	sensor.Perfdata.Min = 0

	// == Max
	// Is there a tempN_highest file? Use it for max value
	value, err = readIntFromFile(basePath + maxValueFileSuffix)
	if err == nil {
		sensor.Perfdata.Max = float64(value)
	}

	sensor.Alarm = readSensorAlarm(basePath)

	return sensor, nil
}

func readVoltageSensor(_, devicePath string, index int) (Sensor, error) {
	basePath := devicePath + "/in" + strconv.Itoa(index)

	var sensor Sensor

	// Look for label
	label := readSensorLabel(basePath)
	sensor.Name = label
	sensor.Perfdata.Label = label

	// Look for input (the actual value)
	value, err := readIntFromFile(basePath + inputFileSuffix)
	if err != nil {
		return sensor, err
	}

	sensor.Perfdata.Value = float64(value) / 1000 // milli Volt to Volt
	sensor.Perfdata.Uom = "V"

	// == Warn thresholds
	tmpWarn := check.Threshold{
		Lower:  check.NegInf,
		Upper:  check.PosInf,
		Inside: false,
	}
	warnPresent := false
	// Is there a inN_max file? If yes, use that as warning
	value, err = readIntFromFile(basePath + maxValueFileSuffix)
	if err == nil {
		tmpWarn.Upper = float64(value) / 1000
		warnPresent = true
	}

	if warnPresent {
		sensor.Perfdata.Warn = &tmpWarn
	}

	// == Crit threshold
	tmpCrit := check.Threshold{
		Lower:  check.NegInf,
		Upper:  check.PosInf,
		Inside: false,
	}
	critPresent := false
	// Is there a inN_crit file? If yes, use that as critical
	value, err = readIntFromFile(basePath + critThresholdFileSuffix)
	if err == nil {
		tmpCrit.Upper = float64(value) / 1000
		critPresent = true
	}

	// Is there a inN_min file? If yes, use that as lower boundary
	value, err = readIntFromFile(basePath + "_min")
	if err == nil {
		tmpCrit.Lower = float64(value) / 1000
		critPresent = true
	}

	if critPresent {
		sensor.Perfdata.Crit = &tmpCrit
	}

	// == Min
	value, err = readIntFromFile(basePath + lowestValueFileSuffix)
	if err == nil {
		sensor.Perfdata.Min = float64(value) / 1000
	}

	// == Max
	// Is there a tempN_highest file? Use it for max value
	value, err = readIntFromFile(basePath + highestValueFileSuffix)
	if err == nil {
		sensor.Perfdata.Max = float64(value) / 1000
	}

	sensor.Alarm = readSensorAlarm(basePath)

	return sensor, nil
}

func readPowerSensor(_, devicePath string, index int) (Sensor, error) {
	basePath := devicePath + "/power" + strconv.Itoa(index)

	var sensor Sensor

	label := readSensorLabel(basePath)
	sensor.Name = label
	sensor.Perfdata.Label = label

	// Look for input (the actual value)
	value, err := readIntFromFile(basePath + inputFileSuffix)

	if err != nil {
		return sensor, err
	}

	sensor.Perfdata.Value = value // micro Watt
	sensor.Perfdata.Uom = "uW"

	// == Min
	value, err = readIntFromFile(basePath + "_input_lowest")
	if err == nil {
		sensor.Perfdata.Min = value
	}

	// == Max
	// Is there a tempN_highest file? Use it for max value
	value, err = readIntFromFile(basePath + "_input_highest")
	if err == nil {
		sensor.Perfdata.Max = value
	}

	// == Warn threshold
	tmpWarn := check.Threshold{
		Lower:  0,
		Upper:  check.PosInf,
		Inside: false,
	}
	warnPresent := false
	// Is there a powerN_cap file? If yes, use that as warning
	value, err = readIntFromFile(basePath + "_cap")
	if err == nil {
		tmpWarn.Upper = float64(value)
		warnPresent = true
	}

	if warnPresent {
		sensor.Perfdata.Warn = &tmpWarn
	}

	// == Crit threshold
	tmpCrit := check.Threshold{
		Lower:  0,
		Upper:  check.PosInf,
		Inside: false,
	}
	critPresent := false
	// Is there a powerN_crit file? If yes, use that as critical
	value, err = readIntFromFile(basePath + critThresholdFileSuffix)
	if err == nil {
		tmpCrit.Upper = float64(value)
		critPresent = true
	}

	if critPresent {
		sensor.Perfdata.Crit = &tmpCrit
	}

	return sensor, nil
}

func readTempSensor(_, devicePath string, index int) (Sensor, error) {
	basePath := devicePath + "/temp" + strconv.Itoa(index)

	var sensor Sensor
	// Look for label
	label := readSensorLabel(basePath)
	sensor.Name = label
	sensor.Perfdata.Label = label

	// Look for input (the actual value)
	value, err := readIntFromFile(basePath + inputFileSuffix)

	if err != nil {
		return sensor, err
	}

	sensor.Perfdata.Value = value / 1000 // milli celsius to celsius
	sensor.Perfdata.Uom = "C"

	// == Warn thresholds
	tmpWarn := check.Threshold{
		Lower:  check.NegInf,
		Upper:  check.PosInf,
		Inside: false,
	}
	warnPresent := false
	// Is there a tempN_max file? If yes, use that as warning
	value, err = readIntFromFile(basePath + maxValueFileSuffix)

	if err == nil {
		tmpWarn.Upper = float64(value / 1000)
		warnPresent = true
	}

	if warnPresent {
		sensor.Perfdata.Warn = &tmpWarn
	}

	// Crit threshold
	tmpCrit := check.Threshold{
		Lower:  check.NegInf,
		Upper:  check.PosInf,
		Inside: false,
	}

	critPresent := false
	// Is there a tempN_crit file? If yes, use that as critical
	value, err = readIntFromFile(basePath + critThresholdFileSuffix)

	if err == nil {
		tmpCrit.Upper = float64(value / 1000)
		critPresent = true
	}

	// Is there a tempN_emergency file? If yes, use that instead of crit
	value, err = readIntFromFile(basePath + "_emergency")

	if err == nil {
		tmpCrit.Upper = float64(value / 1000)
		critPresent = true
	}

	// Is there a tempN_min file? If yes, use that as lower boundary
	value, err = readIntFromFile(basePath + "_min")

	if err == nil {
		tmpCrit.Lower = float64(value / 1000)
		critPresent = true
	}

	if critPresent {
		sensor.Perfdata.Crit = &tmpCrit
	}

	// == Min
	// Is there a tempN_lowest file? Use it for min value
	value, err = readIntFromFile(basePath + lowestValueFileSuffix)

	if err == nil {
		sensor.Perfdata.Min = value
	}

	// == Max
	// Is there a tempN_highest file? Use it for max value
	value, err = readIntFromFile(basePath + highestValueFileSuffix)

	if err == nil {
		sensor.Perfdata.Max = value
	}

	sensor.Alarm = readSensorAlarm(basePath)

	return sensor, nil
}

func readStringFromFile(fp string) (string, error) {
	tmp, err := os.ReadFile(fp)

	if err != nil {
		return "", err
	}

	return strings.TrimRight(string(tmp), "\n"), nil
}

func readIntFromFile(fp string) (int, error) {
	tmp, err := os.ReadFile(fp)

	if err != nil {
		return 0, err
	}

	value, err := strconv.ParseInt(strings.TrimSpace(string(tmp)), 10, 64)

	if err != nil {
		return 0, err
	}

	return int(value), nil
}

func readBoolFromFile(fp string) (bool, error) {
	tmp, err := os.ReadFile(fp)

	if err != nil {
		return false, err
	}

	value, err := strconv.ParseInt(strings.TrimSpace(string(tmp)), 10, 64)

	if err != nil {
		return false, err
	}

	if value == 0 {
		return false, nil
	}

	return true, nil
}

// @param:
// sensorBasePath: something like /sys/class/hwmon/hwmon3/in2
func readSensorLabel(sensorBasePath string) string {
	// See if Sensor is alarmed, use that for the status
	label, err := readStringFromFile(sensorBasePath + "_label")

	if err == nil {
		return label
	}

	// Try the "name" of the device and add the sensor type
	devicePath := path.Dir(sensorBasePath)
	sensorBaseName := path.Base(sensorBasePath)
	label, err = readStringFromFile(devicePath + "/name")

	if err == nil {
		return label + "_" + sensorBaseName
	}

	rescueName := path.Base(devicePath)

	return rescueName + "_" + sensorBaseName
}

// @param:
// sensorBasePath: something like /sys/class/hwmon/hwmon3/in2
func readSensorAlarm(sensorBasePath string) bool {
	// See if Sensor is alarmed, use that for the status
	alarmFiles, err := filepath.Glob(sensorBasePath + "*_alarm")

	if err != nil {
		return false
	}

	for _, alarmPile := range alarmFiles {
		alarm, err := readBoolFromFile(alarmPile)
		if err == nil && alarm {
			return true
		}
	}

	return false
}
