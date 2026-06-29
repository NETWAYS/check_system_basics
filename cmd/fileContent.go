package cmd

import (
	fileContent "github.com/NETWAYS/check_system_basics/internal/files"
	"github.com/NETWAYS/go-check"
	"github.com/spf13/cobra"
)

var FileContentConfig fileContent.FileContentconfig

var fileContentCmd = &cobra.Command{
	Use:     "fileContent",
	Short:   "Submodule to test for different properties on file content ",
	Example: ``,
	Run: func(_ *cobra.Command, _ []string) {
	},
}

func init() {
	rootCmd.AddCommand(fileContentCmd)
	fileContentCmd.DisableFlagsInUseLine = true

	fileContentFS := fileContentCmd.Flags()
	fileContentFS.StringArrayVar(&FileContentConfig.Paths, "paths", []string{}, "File paths to evaluate")
	fileContentFS.StringArrayVar(&FileContentConfig.OKPatterns, "ok-pattern", []string{}, "Regex pattern in file which are OK")
	fileContentFS.StringArrayVar(&FileContentConfig.WarningPatterns, "warning-pattern", []string{}, "Regex pattern in file which cause a WARNING")
	fileContentFS.StringArrayVar(&FileContentConfig.Paths, "critical-pattern", []string{}, "Regex pattern in file which cause a CRITICAL")
	fileContentFS.IntVar(&FileContentConfig.NotFoundStatus, "paths", check.OK, "Exit status if none of the patterns apply (OK (0), warning (1), critical (2), Uknown (3)) (default: OK)")
	fileContentFS.BoolVar(&FileContentConfig.Recursive, "recursive", true, "Recursively test all files in \"paths\"")
}
