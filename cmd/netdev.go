package cmd

import (
	"github.com/NETWAYS/check_system_basics/internal/netdev"
	"github.com/NETWAYS/go-check"
	"github.com/NETWAYS/go-check/perfdata"
	"github.com/NETWAYS/go-check/result"
	"github.com/spf13/cobra"
)

var NetdevConfig netdev.CheckConfig

var netdevCmd = &cobra.Command{
	Use:   "netdev",
	Short: "Submodule to detect, display and check against the interfaces on the local machine",
	Example: `./check_system_basics netdev
[WARNING] - states: warning=4 ok=2
\_ [WARNING] virbr0 is Down
\_ [WARNING] lo is Unknown
\_ [OK] enx00e04c6801bd is Up
\_ [OK] wlp170s0 is Up
\_ [WARNING] tun0 is Unknown
\_ [WARNING] docker0 is Down
|virbr0_rx_bytes=0 virbr0_rx_errors=0 virbr0_rx_dropped=0 virbr0_rx_packets=0 virbr0_tx_bytes=0 virbr0_tx_errors=0 virbr0_tx_dropped=0 virbr0_tx_packets=0 lo_rx_bytes=5991167 lo_rx_errors=0 lo_rx_dropped=0 lo_rx_packets=11264 lo_tx_bytes=5991167 lo_tx_errors=0 lo_tx_dropped=0 lo_tx_packets=11264 enx00e04c6801bd_rx_bytes=259382300 enx00e04c6801bd_rx_errors=0 enx00e04c6801bd_rx_dropped=21297 enx00e04c6801bd_rx_packets=711180 enx00e04c6801bd_tx_bytes=181029368 enx00e04c6801bd_tx_errors=0 enx00e04c6801bd_tx_dropped=0 enx00e04c6801bd_tx_packets=653583 wlp170s0_rx_bytes=4841533 wlp170s0_rx_errors=0 wlp170s0_rx_dropped=21110 wlp170s0_rx_packets=46034 wlp170s0_tx_bytes=160490 wlp170s0_tx_errors=0 wlp170s0_tx_dropped=0 wlp170s0_tx_packets=1453 tun0_rx_bytes=204101734 tun0_rx_errors=0 tun0_rx_dropped=0 tun0_rx_packets=349065 tun0_tx_bytes=121967676 tun0_tx_errors=0 tun0_tx_dropped=0 tun0_tx_packets=377347 docker0_rx_bytes=0 docker0_rx_errors=0 docker0_rx_dropped=0 docker0_rx_packets=0 docker0_tx_bytes=0 docker0_tx_errors=0 docker0_tx_dropped=0 docker0_tx_packets=0`,
	Run: NetdevCheck,
}

func init() {
	rootCmd.AddCommand(netdevCmd)

	fs := netdevCmd.Flags()

	fs.StringSliceVar(&NetdevConfig.Filters.IncludeInterfaceNames, "include-interface-name", nil,
		"Explicitly include only interfaces whose names match this regexp regex (may be repeated). E.g. 'eth', '^et.*'")
	fs.StringSliceVar(&NetdevConfig.Filters.ExcludeInterfaceNames,
		"exclude-interface-name",
		[]string{"^lo$"},
		"Ignore all interfaces where the interface name matches this regexp regex (may be repeated). E.g. 'eth', '^et.*'")
	fs.BoolVar(&NetdevConfig.DownIsCritical, "down-is-critical", false,
		"Setting this option will set the state to CRITICAL if an interface is DOWN")
	fs.BoolVar(&NetdevConfig.NotUpIsOK, "not-up-is-ok", false,
		"Setting this option will set the state to OK regardless of the actual state of an interface")
	fs.BoolVar(&NetdevConfig.UnknownIsOk, "unknown-is-ok", false,
		"Setting this option will set the state to OK if the interface is in a state of UNKNOWN")
}

func NetdevCheck(_ *cobra.Command, _ []string) {
	overall := result.Overall{}

	interfaces, err := netdev.GetAllInterfaces()

	if err != nil {
		check.ExitError(err)
	}

	interfaces, err = netdev.FilterInterfaces(&interfaces, &NetdevConfig.Filters)

	if err != nil {
		check.ExitError(err)
	}

	for i := range interfaces {
		sc := result.PartialResult{}
		_ = sc.SetDefaultState(check.OK)

		sc.Output = interfaces[i].Name + " is " + netdev.TranslateIfaceState(interfaces[i].Operstate)

		if !NetdevConfig.NotUpIsOK {
			switch interfaces[i].Operstate {
			case netdev.Up:
				_ = sc.SetState(check.OK)
			case netdev.Down:
				if NetdevConfig.DownIsCritical {
					_ = sc.SetState(check.Critical)
				} else {
					_ = sc.SetState(check.Warning)
				}
			case netdev.Unknown:
				if NetdevConfig.UnknownIsOk {
					_ = sc.SetState(check.OK)
				} else {
					_ = sc.SetState(check.Warning)
				}
			default:
				_ = sc.SetState(check.Warning)
			}
		}

		for j := range interfaces[i].Metrics {
			pd := perfdata.Perfdata{}
			pd.Label = interfaces[i].Name + "_" + netdev.GetIfaceStatNames()[j]
			pd.Value = interfaces[i].Metrics[j]

			sc.Perfdata.Add(&pd)
		}

		overall.AddSubcheck(sc)
	}

	check.ExitRaw(overall.GetStatus(), overall.GetOutput())
}
