package helpers

import "math"

func RoundDecimalPlaces(number float64, decimalPlaces int) float64 {
	rounded := math.Round(number*math.Pow(10.0, float64(decimalPlaces))) / math.Pow(10.0, float64(decimalPlaces))
	return rounded
}
