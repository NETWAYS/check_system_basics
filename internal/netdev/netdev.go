package netdev

import (
	"os"
	"strconv"
	"strings"

	"github.com/NETWAYS/check_system_basics/internal/common/filter"
)

// types and constants
// dormant = 5
const (
	Up             = 0
	Testing        = 1
	Lowerlayerdown = 2
	Down           = 3
	Unknown        = 4
)

const (
	netDevicePath = "/sys/class/net"
)

func TranslateIfaceState(state uint) string {
	switch state {
	case Up:
		return "Up"
	case Testing:
		return "Testing"
	case Lowerlayerdown:
		return "Lowerlayerdown"
	case Down:
		return "Down"
	case Unknown:
		return "Unknown"
	default:
		return ""
	}
}

// Constants and the string array MUST be kept in sync!
const (
	rx_bytes int = iota
	rx_errs
	rx_drop
	rx_packets
	/*
		rx_fifo uint
		rx_frame uint
		rx_compressed uint
		rx_multicast uint
	*/

	tx_bytes
	tx_errs
	tx_drop
	tx_packets
	/*
		tx_fifo uint
		tx_frame uint
		tx_compressed uint
		tx_multicast uint
	*/
	metricLength
)

func GetIfaceStatNames() []string {
	return []string{
		"rx_bytes",
		"rx_errors",
		"rx_dropped",
		"rx_packets",
		/*
			"rx_fifo",
			"rx_frame"
			"rx_compressed"
			"rx_multicast"
		*/

		"tx_bytes",
		"tx_errors",
		"tx_dropped",
		"tx_packets",
		/*
			"tx_fifo",
			"tx_frame",
			"tx_compressed",
			"tx_multicast",
		*/
	}
}

type statistics [metricLength]uint64

type IfaceData struct {
	Name      string
	Operstate uint
	Metrics   statistics
}

const (
	IfaceDataName = iota
)

func (iface IfaceData) GetFilterableValue(ident uint) string {
	switch ident {
	case IfaceDataName:
		return iface.Name
	default:
		return ""
	}
}

func GetAllInterfaces() ([]IfaceData, error) {
	interfaces, err := listInterfaces()

	if err != nil {
		return []IfaceData{}, err
	}

	if len(interfaces) == 0 {
		return []IfaceData{}, nil
	}

	result := make([]IfaceData, len(interfaces))

	for i := range interfaces {
		result[i].Name = interfaces[i]

		err = getInterfaceState(&result[i])

		if err != nil {
			return result, err
		}

		err = getInfacesStatistics(&result[i])
		if err != nil {
			return result, err
		}
	}

	return result, nil
}

func listInterfaces() ([]string, error) {
	file, err := os.Open(netDevicePath)
	if err != nil {
		return []string{}, err
	}

	devices, err := file.Readdirnames(0)
	if err != nil {
		return []string{}, err
	}

	result := make([]string, 0, len(devices))
	for idx := range devices {
		fileInfo, err := os.Stat(netDevicePath + devices[idx])
		if err != nil {
			// Could not stat file there, not sure if this can be handled usefully. Just die for now.
			return []string{}, err
		}

		if fileInfo.Mode().IsDir() {
			result = append(result, devices[idx])
		}
	}

	return result, nil
}

// getInterfaceState receives the name of an interfaces and returns
// an integer result code representing the state of the interface
// @result = 0 => Interface is up
// @result = 2 => Interface is down
// @result = 3 => Interface is unknown or state of the interface is unknown for some reason
func getInterfaceState(data *IfaceData) error {
	basePath := netDevicePath + data.Name

	bytes, err := os.ReadFile(basePath + "/operstate")
	if err != nil {
		return err
	}

	switch strings.Trim(string(bytes), " \n") {
	case "up":
		data.Operstate = Up
		return nil
	case "testing":
		data.Operstate = Testing
		return nil
	case "down":
		data.Operstate = Down
		return nil
	case "lowerlayerdown":
		data.Operstate = Lowerlayerdown
		return nil
	default:
		data.Operstate = Unknown
		return nil
	}
}

// Get interfaces statistics
// @result: ifaceStats, err
func getInfacesStatistics(data *IfaceData) error {
	basePath := netDevicePath + data.Name + "/statistics"

	var val uint64

	for idx, stat := range GetIfaceStatNames() {
		numberBytes, err := os.ReadFile(basePath + "/" + stat)

		if err != nil {
			return err
		}

		numberString := string(numberBytes)
		val, err = strconv.ParseUint(numberString[:len(numberString)-1], 10, 64)

		if err != nil {
			return err
		}

		data.Metrics[idx] = val
	}

	return nil
}

func FilterInterfaces(interfaces *[]IfaceData, filters *Filter) ([]IfaceData, error) {
	foo, err := filter.Filter(*interfaces,
		&filters.IncludeInterfaceNames,
		IfaceDataName,
		filter.Options{
			MatchIncludedInResult: true,
			RegexpMatching:        true,
			EmptyFilterNoMatch:    false,
		},
	)

	if err != nil {
		return []IfaceData{}, err
	}

	foo, err = filter.Filter(foo,
		&filters.ExcludeInterfaceNames,
		IfaceDataName,
		filter.Options{
			MatchIncludedInResult: false,
			RegexpMatching:        true,
			EmptyFilterNoMatch:    false,
		},
	)

	if err != nil {
		return []IfaceData{}, err
	}

	return foo, nil
}
