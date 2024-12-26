package util

import "math"

func MutiplyAndRound(a, b float64) float64 {
	result := a * b
	return math.Round(result*100) / 100
}
