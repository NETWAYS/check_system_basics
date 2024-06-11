package files

import (
	"io/fs"
	"regexp"
	"time"

	"github.com/NETWAYS/check_system_basics/internal/common/thresholds"
)

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
	FileSizeIncludeFilter thresholds.Thresholds
	// Exclude files where the size matches this threshold
	FileSizeExcludeFilter thresholds.Thresholds

	ModificationTimeOlderThanFilter   time.Time
	ModificationTimeYoungerThanFilter time.Time

	// Should the operation recursively descent into directories
	Recursive bool

	// Thresholds for the number of matching files
	WarningNumberOfFiles  thresholds.Thresholds
	CriticalNumberOfFiles thresholds.Thresholds
}
