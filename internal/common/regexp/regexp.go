package regexp

import (
	"regexp"
	"strings"
)

type SBRegex struct {
	Regex regexp.Regexp
	IsSet bool
}

type SBRegexList []SBRegex

func (s *SBRegex) String() string {
	return s.Regex.String()
}

func (s *SBRegex) Type() string {
	return "Golang re regular Expression"
}

func (s *SBRegex) Set(input string) error {
	re, err := regexp.Compile(input)
	if err != nil {
		return err
	}

	s.Regex = *re
	s.IsSet = true

	return nil
}

func (s *SBRegexList) String() string {
	builder := strings.Builder{}

	length := len(*s)
	for index, entry := range *s {
		if index == length-1 {
			builder.WriteString(entry.String())
		} else {
			builder.WriteString(entry.String() + ", ")
		}
	}

	return builder.String()
}

func (s *SBRegexList) Type() string {
	return "Golang re regular Expression list"
}

func (s *SBRegexList) Set(input string) error {
	re, err := regexp.Compile(input)
	if err != nil {
		return err
	}

	newFoo := SBRegex{
		Regex: *re,
		IsSet: true,
	}

	*s = append(*s, newFoo)

	return nil
}
