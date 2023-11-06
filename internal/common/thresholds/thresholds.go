package thresholds

import (
	"strings"

	"github.com/NETWAYS/go-check"
	"github.com/spf13/pflag"
)

type ThresholdWrapper struct {
	Th    check.Threshold
	IsSet bool
}

type Thresholds struct {
	Warn ThresholdWrapper
	Crit ThresholdWrapper
}

func (t *ThresholdWrapper) Set(foo string) error {
	tmp, err := check.ParseThreshold(foo)
	if err != nil {
		return err
	}

	t.Th.Inside = tmp.Inside
	t.Th.Lower = tmp.Lower
	t.Th.Upper = tmp.Upper
	t.IsSet = true

	return nil
}

func (t *ThresholdWrapper) String() string {
	return t.Th.String()
}

func (t *ThresholdWrapper) Type() string {
	return "Range_Expression"
}

type ThresholdOption struct {
	Th          *ThresholdWrapper
	FlagString  string
	Description string
	Default     ThresholdWrapper
}

func AddFlags(flagP *pflag.FlagSet, ths *[]ThresholdOption) {
	for i := range *ths {
		desc := strings.Builder{}
		desc.WriteString((*ths)[i].Description)

		if (*ths)[i].Default.IsSet {
			desc.WriteString(" (Default: " + (*ths)[i].Default.Th.String() + ")")
			(*ths)[i].Th.IsSet = (*ths)[i].Default.IsSet
			(*ths)[i].Th.Th = (*ths)[i].Default.Th
		}

		flagP.Var((*ths)[i].Th, (*ths)[i].FlagString, desc.String())
	}
}
