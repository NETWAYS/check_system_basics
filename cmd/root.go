package cmd

import (
	"fmt"
	"os"

	intConfig "github.com/NETWAYS/check_system_basics/internal/common/config"
	"github.com/NETWAYS/go-check"
	"github.com/spf13/cobra"
)

var Timeout = 30
var debug = false

var (
	version string
)

var rootCmd = &cobra.Command{
	Use:     "check_system_basics",
	Short:   "Icinga check plugin to check various Linux metrics",
	Version: version,
	PersistentPreRun: func(_ *cobra.Command, _ []string) {
		go check.HandleTimeout(Timeout)
	},
	Run: RunFunction,
}

func Execute() {
	defer check.CatchPanic()

	if err := rootCmd.Execute(); err != nil {
		check.ExitError(err)
	}
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.DisableAutoGenTag = true

	rootCmd.SetHelpCommand(&cobra.Command{
		Use:    "no-help",
		Hidden: true,
	})

	pfs := rootCmd.PersistentFlags()

	pfs.IntVarP(&Timeout, "timeout", "t", Timeout,
		"Timeout for the check")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Enable debug output")

	rootCmd.Flags().SortFlags = false
	pfs.SortFlags = false

	flagSet := rootCmd.Flags()
	flagSet.Bool("dump-icinga2-config", false, "Dump icinga2 config for this plugin")
}

func RunFunction(cmd *cobra.Command, args []string) {
	flagSet := cmd.Flags()

	dumpConfig, err := flagSet.GetBool("dump-icinga2-config")
	if err != nil {
		check.ExitError(err)
	}

	if dumpConfig {
		ConfigDump(cmd, cmd.CommandPath())
		os.Exit(check.OK)
	}

	Help(cmd, args)
}

func Help(cmd *cobra.Command, _ []string) {
	_ = cmd.Usage()

	os.Exit(check.Unknown)
}

func ConfigDump(cmd *cobra.Command, executableName string) {
	result := intConfig.GenerateIcinga2Config(cmd, "system_basics", executableName, true)

	fmt.Println(result)
}
