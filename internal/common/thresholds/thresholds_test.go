package thresholds

import (
	"testing"

	"github.com/NETWAYS/go-check"
	"github.com/stretchr/testify/assert"
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

	assert.Equal(t, tw, tw2)
}
