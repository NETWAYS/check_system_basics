package filter

import (
	"regexp"
)

type Filterable interface {
	GetFilterableValue(uint) string
}

type Options struct {
	EmptyFilterNoMatch    bool
	MatchIncludedInResult bool
	RegexpMatching        bool
}

// Filters on Value with a list of filters
// @return, true if value should be included, false otherwise
func Filter[I Filterable](input []I, filterList *[]string, valueIdentifier uint, opts Options) ([]I, error) {
	if len(*filterList) == 0 {
		if opts.EmptyFilterNoMatch {
			return []I{}, nil
		}

		return input, nil
	}

	result := make([]I, 0)

	for i := range input {
		isMatch := false

		for j := range *filterList {
			var match bool

			var err error

			if opts.RegexpMatching {
				match, err = regexp.MatchString(
					(*filterList)[j],
					input[i].GetFilterableValue(valueIdentifier),
				)

				if err != nil {
					return []I{}, err
				}
			} else {
				match = (*filterList)[j] == input[i].GetFilterableValue(valueIdentifier)
			}

			if match {
				isMatch = true
				break
			}
		}

		if (isMatch && opts.MatchIncludedInResult) ||
			(!isMatch && !opts.MatchIncludedInResult) {
			result = append(result, input[i])
		}
	}

	return result, nil
}
