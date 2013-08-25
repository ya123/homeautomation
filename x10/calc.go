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
func PercentToStep(percent int) (step int) {
	if percent < 7 {
		return 1
	}
	if percent < 12 {
		return 2
	}

	result := float64(percent-9)/factor + float64(2)
	return round(result)

}
func round(result float64) int {
	floor := math.Floor(result)
	if result >= floor+0.5 {
		return int(floor) + 1
	}
	return int(floor)
}
