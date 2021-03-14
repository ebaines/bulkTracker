package main

import (
	"math"
	"testing"
)

func TestSlidingAvgs(t *testing.T) {
	var values = []float64{82, 82, 82, 82, 82, 82, 82}
	avgs := slidingAvgs(values, 7)
	if avgs[0] != 82 {
		t.Error()
	}

	values = []float64{82, 82, 82, 82, 82, 82, 82, 81, 81, 83}
	correct := []float64{82, 81.9, 81.7, 81.9}
	avgs = slidingAvgs(values, 7)
	for i := 0; i < len(avgs); i++ {
		if math.Round(avgs[i]*10)/10 != correct[i] {
			t.Error()
		}
	}

}

func TestCalculateDiffs(t *testing.T) {
	var values = []float64{1, 1.2, 1.4, 1.6, 1.8, 2.0, 3.0, 5.0, 8.0}
	var correct = []float64{
		0.2, 0.2, 0.2, 0.2, 0.2, 1.0, 2.0, 3.0,
	}
	differences := calculateDifferences(values)

	for i := 0; i < len(differences); i++ {
		if math.Round(differences[i]*1000)/1000 != correct[i] {
			t.Error()
		}
	}
}

func TestCalculateSevenDayDifferences(t *testing.T) {
	var values = []float64{1, 1.2, 1.4, 1.6, 1.8, 2.0, 3.0, 5.0, 8.0}
	var correct = []float64{
		0, 0, 0, 0, 0, 0, 0, 4, 6.8,
	}
	differences := calculateSevenDayDifferences(values)
	for i := 0; i < len(differences); i++ {
		if math.Round(differences[i]*1000)/1000 != correct[i] {
			t.Error()
		}
	}
}

//func TestCalculateTDEE(t *testing.T) {
//	var caloriesConsumed = []float64{2750, 3000, 2500, 2250, 2000}
//	var weightDiff = []float64{0.250, 0.500, 0, -0.250, -0.5}
//
//	for i := 0; i < len(weightDiff); i++ {
//		fmt.Println(calculateTDEE(caloriesConsumed[i], weightDiff[i]))
//	}
//}
