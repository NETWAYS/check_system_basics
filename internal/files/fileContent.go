package fileContent

import ()

type FileContentconfig struct {
	Paths            []string
	OKPatterns       []string
	WarningPatterns  []string
	CriticalPatterns []string
	NotFoundStatus   int
	Recursive        bool
}
