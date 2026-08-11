package fileContent

import (
	"github.com/NETWAYS/check_system_basics/internal/common/regexp"
	"github.com/NETWAYS/check_system_basics/internal/common/status"
)

type FileContentconfig struct {
	Paths            []string
	OKPatterns       regexp.SBRegexList
	WarningPatterns  regexp.SBRegexList
	CriticalPatterns regexp.SBRegexList
	NotFoundStatus   status.Status
	Recursive        bool
}
