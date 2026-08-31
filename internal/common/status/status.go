package status

import (
	"fmt"
	"strconv"

	// "strings"

	"github.com/NETWAYS/go-check"
	// "github.com/spf13/pflag"
)

type Status struct {
	Status check.Status
	IsSet  bool
}

func (s *Status) String() string {
	return s.Status.String()
}

func (s *Status) Type() string {
	return "Monitoring Plugin Status"
}

func (s *Status) Set(input string) error {
	tmp, err := check.NewStatusFromString(input)
	if err == nil {
		s.Status = tmp
		s.IsSet = true

		return nil
	}

	intStatus, err := strconv.Atoi(input)
	if err != nil {
		return fmt.Errorf("failed to convert input \"%s\" to status", input)
	}

	tmp, err = check.NewStatus(intStatus)
	if err != nil {
		return fmt.Errorf("failed to convert input \"%s\" to status", input)
	}

	s.IsSet = true
	s.Status = tmp

	return nil
}
