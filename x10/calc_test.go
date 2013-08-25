package x10

import (
	"testing"
)

var (
	stepToPercent = map[int]int{
		1:  1,
		2:  10,
		3:  15,
		4:  19,
		5:  24,
		6:  28,
		8:  37,
		16: 73,
		20: 91,
		22: 100,
	}
)
var (
	percentToStep = map[int]int{
		1:   1,
		4:   1,
		7:   2,
		9:   2,
		12:  3,
		15:  3,
		16:  4,
		22:  5,
		34:  8,
		43:  10,
		50:  11,
		100: 22,
	}
)

func TestStepToPercent(t *testing.T) {
	for in, expected := range stepToPercent {
		if StepToPercent(in) != expected {
			t.Errorf("falscher wert für step = %v: %v", in, StepToPercent(in))
		}
	}

}

func TestPercentToStep(t *testing.T) {
	for in, expected := range percentToStep {
		if PercentToStep(in) != expected {
			t.Errorf("falscher wert für percent = %v: %v", in, PercentToStep(in))
		}
	}

}
