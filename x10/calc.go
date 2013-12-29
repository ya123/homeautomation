package x10

import (
	"math"
)

var factor = 4.5

func StepToPercent(step int) (percent int) {
	switch step {
	case 1:
		return 1
	case 2:
		return 10
	default:
		result := factor*float64(step-2) + float64(10)
		return round(result)

	}
}

func min7(in int) int {
	if in <= 7 {
		return 7
	}
	return in
}

func PercentToStep(percent int) (step int) {
	if percent < 7 {
		//return 1
		return min7(1)
	}
	if percent < 12 {
		return min7(2)
	}

	result := float64(percent-9)/factor + float64(2)
	return min7(round(result))

}
func round(result float64) int {
	floor := math.Floor(result)
	if result >= floor+0.5 {
		return int(floor) + 1
	}
	return int(floor)
}
