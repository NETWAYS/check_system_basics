package cmd

import (
	"github.com/NETWAYS/check_system_basics/internal/files"
	"github.com/spf13/cobra"
)

var FilesConfig files.FilesConfig

var filesCmd = &cobra.Command{
	Use:     "files",
	Short:   "Submodule to reason about files and file attributes",
	Example: "",
	Run:     fileFunction,
}

func fileFunction(cmd *cobra.Command, args []string) {
}

func init() {
	rootCmd.AddCommand(filesCmd)

	filesFlagSet := filesCmd.Flags()

	filesFlagSet.StringVar(&FilesConfig.BasePath, "base-path", "", "The directory path which should be examined")
	filesFlagSet.StringArrayVar()
}
