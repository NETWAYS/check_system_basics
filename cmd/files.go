package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/NETWAYS/check_system_basics/internal/common/thresholds"
	"github.com/NETWAYS/check_system_basics/internal/files"
	"github.com/NETWAYS/go-check"
	"github.com/NETWAYS/go-check/result"
	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"
)

var FilesConfigRaw files.FilesConfigRaw
var FilesConfig files.FilesConfig

var filesCmd = &cobra.Command{
	Use:     "files",
	Short:   "Submodule to reason about files and file attributes",
	Example: "",
	Run:     fileFunction,
}

func fileFunction(cmd *cobra.Command, args []string) {
	overall := result.Overall{}

	config, err := files.ParseRawConfig(FilesConfig, &FilesConfigRaw)
	if err != nil {
		check.ExitError(err)
	}

	rawFileList, err := files.GetFileList(config)
	if err != nil {
		check.ExitError(err)
	}

	// Filter files
	filteredFiles, err := files.FilterFiles(config, rawFileList)
	if err != nil {
		check.ExitError(err)
	}

	// Evaluate checks
	partialResults, err := evaluate(&config, filteredFiles)
	if err != nil {
		check.ExitError(err)
	}

	for idx := range partialResults {
		overall.AddSubcheck(partialResults[idx])
	}

	// exit
	check.ExitRaw(overall.GetStatus(), overall.GetOutput())
}

func evaluate(config *files.FilesConfig, fileList []os.FileInfo) ([]result.PartialResult, error) {

	returnVal := make([]result.PartialResult, 0)

	var totalSize int64

	for idx := range fileList {
		fileInfo := fileList[idx]

		totalSize += fileInfo.Size()
	}

	totalNumberOfFiles := len(fileList)
	// Number of matching files checks
	if config.CriticalNumberOfFiles.IsSet || config.WarningNumberOfFiles.IsSet {
		numberResult := result.PartialResult{}
		numberResult.SetDefaultState(check.OK)

		if config.WarningNumberOfFiles.IsSet && config.WarningNumberOfFiles.Th.DoesViolate(float64(totalNumberOfFiles)) {
			numberResult.SetState(check.Warning)
			numberResult.Output = fmt.Sprintf("Number of matching files (%d) violates warning threshold %v", totalNumberOfFiles, config.WarningNumberOfFiles.Th)
		}

		if config.CriticalNumberOfFiles.IsSet {
			numberResult.SetState(check.Critical)
			numberResult.Output = fmt.Sprintf("Number of matching files (%d) violates critical threshold %v", totalNumberOfFiles, config.CriticalNumberOfFiles.Th)
		}

		// TODO perfdata

		returnVal = append(returnVal, numberResult)

	}

	if config.CriticalTotalSize.IsSet || config.WarningTotalSize.IsSet {
		totalSizeResult := result.PartialResult{}
		totalSizeResult.SetDefaultState(check.OK)

		if config.WarningTotalSize.IsSet && config.WarningTotalSize.Th.DoesViolate(float64(totalSize)) {
			totalSizeResult.SetState(check.Warning)
			totalSizeResult.Output = fmt.Sprintf("Total size of all files (%s) violates Warning threshold (%v)", humanize.IBytes(uint64(totalSize)), config.WarningTotalSize.Th)
		}
		if config.CriticalTotalSize.IsSet && config.CriticalTotalSize.Th.DoesViolate(float64(totalSize)) {
			totalSizeResult.SetState(check.Critical)
			totalSizeResult.Output = fmt.Sprintf("Total size of all files (%s) violates Critical threshold (%v)", humanize.IBytes(uint64(totalSize)), config.CriticalTotalSize.Th)
		}

		// TODO perfdata
		returnVal = append(returnVal, totalSizeResult)
	}

	return returnVal, nil
}

func init() {
	rootCmd.AddCommand(filesCmd)

	filesFlagSet := filesCmd.Flags()

	// "Raw" config options, have to be parsed
	filesFlagSet.StringArrayVar(&FilesConfigRaw.FileNameIncludeRegex, "filenameIncludeRegex", []string{}, "Golang re expression for files to be included")
	filesFlagSet.StringArrayVar(&FilesConfigRaw.FileNameExcludeRegex, "filenameExcludeRegex", []string{}, "Golang re expression for files to be included")

	filesFlagSet.StringArrayVar(&FilesConfigRaw.FileTypeIncludeFilters, "includeFileType", []string{}, "File types to be included: may be \"directory\", \"symlink\", \"device\", \"fifo\" \"unix-socket\". By default all types are included")
	filesFlagSet.StringArrayVar(&FilesConfigRaw.FileTypeExcludeFilters, "excludeFileType", []string{}, "File types to be exclude: may be \"directory\", \"symlink\", \"device\", \"fifo\" \"unix-socket\". By default no types are excluded")

	filesFlagSet.DurationVar(&FilesConfigRaw.ModificationTimeOlderThanFilter, "modificationTimeOlderThan", 0*time.Hour, "Only files which were modified (or created) before that point in time are included. Examples: 1h (one hour), 1h10 (one hour + 10 minutes), 20s (twenty seconds)")

	filesFlagSet.DurationVar(&FilesConfigRaw.ModificationTimeYoungerThanFilter, "modificationTimeYoungerThan", 20*366*24*time.Hour, "Only files which were modified (or created) after that point in time are included. Examples: 1h (one hour), 1h10 (one hour + 10 minutes), 20s (twenty seconds)")

	// plain config options, will be taken directly
	filesFlagSet.StringVar(&FilesConfig.BasePath, "base-path", "", "The directory path which should be examined")

	filesFlagSet.BoolVar(&FilesConfig.Recursive, "recursive", false, "Whether to descent recursively into (non filtered) subdirectories")

	filesThresholds := []thresholds.ThresholdOption{
		{
			Th:          &FilesConfig.CriticalNumberOfFiles,
			FlagString:  "criticalCount",
			Description: "Critical threshold for number of matching files",
		},
		{
			Th:          &FilesConfig.WarningNumberOfFiles,
			FlagString:  "warningCount",
			Description: "Warning threshold for number of matching files",
		},
		{
			Th:          &FilesConfig.FileSizeIncludeFilter,
			FlagString:  "fileSizeIncludeFilter",
			Description: "Only include files where the size matches this thresholds (in bytes)",
		},
		{
			Th:          &FilesConfig.FileSizeExcludeFilter,
			FlagString:  "fileSizeExcludeFilter",
			Description: "Exclude files where the size matches this thresholds (in bytes)",
		},
		{
			Th:          &FilesConfig.CriticalTotalSize,
			FlagString:  "criticalTotalSize",
			Description: "Critical threshold for the total size of all matching files",
		},
		{
			Th:          &FilesConfig.WarningTotalSize,
			FlagString:  "warningTotalSize",
			Description: "Warning threshold for the total size of all matching files",
		},
	}

	thresholds.AddFlags(filesFlagSet, &filesThresholds)
}
