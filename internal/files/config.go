package files

import (
	"errors"
	"io/fs"
	"regexp"
	"strings"
	"time"

	"github.com/NETWAYS/check_system_basics/internal/common/thresholds"
)

// "Raw" config from command line
type FilesConfigRaw struct {
	// Strings have to be compiled to regexes
	FileNameIncludeRegex []string
	FileNameExcludeRegex []string

	// translate this to filemode bits
	FileTypeIncludeFilters []string
	FileTypeExcludeFilters []string

	// for modification time we have to get a date (time.time), but it
	// probably more useful to specify a time interval relative to Now()
	ModificationTimeOlderThanFilter   time.Duration
	ModificationTimeYoungerThanFilter time.Duration
}

// "processed" config after parsing and sanitysing
type FilesConfig struct {
	// The path (a directory) in which all operations and checks take place
	BasePath string

	// Include files where the name matches this regex in checks
	FileNameIncludeRegex []regexp.Regexp
	// Exclude files where the name matches this regex from checks
	FileNameExcludeRegex []regexp.Regexp

	// Include files of the specific type
	FileTypeIncludeFilters []fs.FileMode
	// Excludefiles of the specific type
	FileTypeExcludeFilters []fs.FileMode

	// Include files where the size matches this threshold
	FileSizeIncludeFilter thresholds.ThresholdWrapper
	// Exclude files where the size matches this threshold
	FileSizeExcludeFilter thresholds.ThresholdWrapper

	ModificationTimeOlderThanFilter   time.Time
	ModificationTimeYoungerThanFilter time.Time

	// Should the operation recursively descent into directories
	Recursive bool

	// Thresholds for the number of matching files
	WarningNumberOfFiles  thresholds.ThresholdWrapper
	CriticalNumberOfFiles thresholds.ThresholdWrapper

	// Thresholds for total size in file system in bytes
	WarningTotalSize  thresholds.ThresholdWrapper
	CriticalTotalSize thresholds.ThresholdWrapper
}

func ParseRawConfig(result FilesConfig, raw *FilesConfigRaw) (FilesConfig, error) {

	// File name regexes
	//  include
	err := regexParser(&raw.FileNameIncludeRegex, &result.FileNameIncludeRegex)
	if err != nil {
		return result, err
	}

	// exclude
	err = regexParser(&raw.FileNameExcludeRegex, &result.FileNameExcludeRegex)
	if err != nil {
		return result, err
	}

	// file type
	//  include
	err = fileTypeParser(&raw.FileTypeIncludeFilters, &result.FileTypeIncludeFilters)
	if err != nil {
		return result, err
	}
	//  exclude
	err = fileTypeParser(&raw.FileTypeExcludeFilters, &result.FileTypeExcludeFilters)
	if err != nil {
		return result, err
	}

	// modification time
	result.ModificationTimeOlderThanFilter = modificationTimeParser(raw.ModificationTimeOlderThanFilter)
	result.ModificationTimeYoungerThanFilter = modificationTimeParser(raw.ModificationTimeYoungerThanFilter)

	return result, nil
}

func modificationTimeParser(input time.Duration) time.Time {
	// input (a duration) is interpreted as a time before Now()
	// therefor 20h (20 hours) at 21:00 o'clock ->  01:00 o'clock

	return time.Now().Add(-input)
}

func regexParser(input *[]string, output *[]regexp.Regexp) error {
	for idx := range *input {
		rgx, err := regexp.Compile((*input)[idx])
		if err != nil {
			return err
		}

		(*output) = append((*output), *rgx)
	}

	return nil
}

func fileTypeParser(input *[]string, output *[]fs.FileMode) error {
	for idx := range *input {
		switch strings.ToLower((*input)[idx]) {
		case "d", "dir", "directory":
			*output = append(*output, fs.ModeDir)
		case "l", "symlink":
			*output = append(*output, fs.ModeSymlink)
		case "p", "fifo":
			*output = append(*output, fs.ModeNamedPipe)
		case "socket":
			*output = append(*output, fs.ModeSocket)
		case "sticky":
			*output = append(*output, fs.ModeSticky)
		case "setuid":
			*output = append(*output, fs.ModeSetuid)
		case "setgid":
			*output = append(*output, fs.ModeSetgid)
		default:
			err := errors.New("Unknown file type parameter: " + (*input)[idx])
			return err
		}
	}
	return nil
}
