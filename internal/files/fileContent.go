package fileContent

import (
	"github.com/NETWAYS/check_system_basics/internal/common/regexp"
	"github.com/NETWAYS/check_system_basics/internal/common/status"
	"github.com/NETWAYS/check_system_basics/internal/common/thresholds"
)

type FileContentconfig struct {
	Paths                []string
	OKPatterns           regexp.SBRegexList
	WarningPatterns      regexp.SBRegexList
	CriticalPatterns     regexp.SBRegexList
	NotFoundStatus       status.Status
	Recursive            bool
	MetricPattern        regexp.SBRegex
	MetricLabel          string
	MetricThresholds     thresholds.Thresholds
	MetricNotFoundStatus status.Status
}
