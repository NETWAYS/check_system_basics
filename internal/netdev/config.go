package netdev

import (
	"net"

	"github.com/NETWAYS/check_system_basics/internal/common/thresholds"
)

type CheckConfig struct {
	WarningTotalCountOfInterfaces  thresholds.ThresholdWrapper
	CriticalTotalCountOfInterfaces thresholds.ThresholdWrapper

	DownIsCritical bool
	NotUpIsOK      bool

	Filters Filter
}

type Filter struct {
	IncludeInterfaceNames []string
	ExcludeInterfaceNames []string

	IncludeIPRange []net.IP
	ExcludeIPRange []net.IP
}
