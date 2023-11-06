package cmd

import (
	"fmt"

	"github.com/NETWAYS/check_system_basics/internal/sensors"
	"github.com/NETWAYS/go-check"
	"github.com/NETWAYS/go-check/result"

	"github.com/spf13/cobra"
)

var sensorsCmd = &cobra.Command{
	Use:   "sensors",
	Short: "Submodule to read the hardware sensors known to linux and check if they exceed the internal threshold",
	Long: `This plugin tried to read all sensors from /sys/class/hwmon and read the values and
thresholds respecting the sensor type and the respective specialities`,
	Example: `./check_system_basics sensors
[OK] - states: ok=6
\_ [OK] acpitz
		\_ [OK] acpitz_temp1: Ok - 46C
\_ [OK] BAT1
		\_ [OK] BAT1_in0: Ok - 17.363V
		\_ [OK] BAT1_curr1: Ok - 0A
\_ [OK] nvme
		\_ [OK] Composite: Ok - 40C
\_ [OK] ACAD
\_ [OK] coretemp
		\_ [OK] Package id 0: Ok - 46C
		\_ [OK] Core 0: Ok - 44C
		\_ [OK] Core 1: Ok - 42C
		\_ [OK] Core 2: Ok - 43C
		\_ [OK] Core 3: Ok - 45C
\_ [OK] iwlwifi_1
		\_ [OK] iwlwifi_1_temp1: Ok - 51C
|acpitz_temp1=46C;;~:210 BAT1_in0=17.363V BAT1_curr1=0A Composite=40C;~:83;-5:87  'Package id 0'=46C;~:100;~:100 'Core 0'=44C;~:100;~:100 'Core 1'=42C;~:100;~:100 'Core 2'=43C;~:100;~:100 'Core 3'=45C;~:100;~:100 iwlwifi_1_temp1=51C`,
	Run: func(cmd *cobra.Command, args []string) {
		devices, err := sensors.GetDefaultDevices()
		if err != nil {
			check.ExitError(err)
		}

		var overall result.Overall

		if len(devices) == 0 {
			overall.Add(check.Unknown, "No devices found")
			check.ExitRaw(overall.GetStatus(), overall.GetOutput())
		}
		var (
			alarms uint = 0
		)

		for _, device := range devices {
			var devicePartial result.PartialResult
			_ = devicePartial.SetDefaultState(check.OK)
			devicePartial.Output = device.Name
			for idx, sensor := range device.Sensors {
				var ssc result.PartialResult
				_ = ssc.SetDefaultState(check.OK)
				ssc.Perfdata.Add(&(device.Sensors[idx]).Perfdata)
				if sensor.Alarm {
					ssc.Output = "Alarm!"
					_ = ssc.SetState(check.Critical)
					alarms++
				} else {
					ssc.Output = "Ok"
					_ = ssc.SetState(check.OK)
				}

				// Add perfdata label (sensor name) to ouptput to make it more descriptive
				if len(ssc.Perfdata) == 1 {
					ssc.Output = fmt.Sprintf("%s: %s - %v%s", ssc.Perfdata[0].Label, ssc.Output, ssc.Perfdata[0].Value, ssc.Perfdata[0].Uom)
				}

				devicePartial.AddSubcheck(ssc)
			}

			overall.AddSubcheck(devicePartial)
		}

		check.ExitRaw(overall.GetStatus(), overall.GetOutput())
	},
}

func init() {
	rootCmd.AddCommand(sensorsCmd)
}
