package thresholds

import (
	"reflect"
	"testing"

	"github.com/NETWAYS/go-check"
)

func TestThresholdSet(t *testing.T) {
	tw := ThresholdWrapper{
		Th: check.Threshold{
			Inside: true,
			Lower:  30,
			Upper:  50,
		},
		IsSet: true,
	}

	tw2 := ThresholdWrapper{}

	tw2.Set(tw.String())

	if !reflect.DeepEqual(tw, tw2) {
		t.Fatalf("expected %v, got %v", tw, tw2)
	}
}
