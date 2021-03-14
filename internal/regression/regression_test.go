package regression

import (
	"math"
	"testing"
)

func TestCoordsToArrays(t *testing.T) {
	coordinates := []Coord{
		{
			X: 0.5578196,
			Y: 18.63654,
		},
		{
			X: 2.0217271,
			Y: 103.49646,
		},
		{
			X: 2.5773252,
			Y: 150.35391,
		},
		{
			X: 3.4140288,
			Y: 190.51031,
		},
		{
			X: 4.3014084,
			Y: 208.70115,
		},
		{
			X: 4.7448394,
			Y: 213.71135,
		},
		{
			X: 5.1073781,
			Y: 228.49353,
		},
	}

	xPoints, yPoints := CoordsToArrays(coordinates)
	correctX := []float64{0.5578196, 2.0217271, 2.5773252, 3.4140288, 4.3014084, 4.7448394, 5.1073781}
	correctY := []float64{18.63654, 103.49646, 150.35391, 190.51031, 208.70115, 213.71135, 228.49353}
	for i := 0; i < len(xPoints); i++ {
		if xPoints[i] != correctX[i] {
			t.Errorf("CoordsToArrays returned %v not %v", xPoints[i], correctX[i] )
		}
	}
	for i := 0; i < len(xPoints); i++ {
		if yPoints[i] != correctY[i] {
			t.Errorf("CoordsToArrays returned %v not %v", yPoints[i], correctY[i] )
		}
	}
}

func TestFindDistance(t *testing.T) {
	coordinates := []Coord{
		{
			X: 0.5578196,
			Y: 18.63654,
		},
		{
			X: 2.0217271,
			Y: 103.49646,
		},
		{
			X: 2.5773252,
			Y: 150.35391,
		},
		{
			X: 3.4140288,
			Y: 190.51031,
		},
		{
			X: 4.3014084,
			Y: 208.70115,
		},
		{
			X: 4.7448394,
			Y: 213.71135,
		},
		{
			X: 5.1073781,
			Y: 228.49353,
		},
	}

	correctDistances := []float64{
		0.000000,
		1.463908,
		2.019506,
		2.856209,
		3.743589,
		4.187020,
		4.549559,
	}

	xPoints, _ := CoordsToArrays(coordinates)
	for i := 0; i < len(xPoints); i++ {
		distance := findDist(xPoints[0], xPoints[i])
		if math.Round(distance*1000000)/1000000 != correctDistances[i] {
			t.Errorf("findDist returned %v not %v", math.Round(distance*1000000)/1000000, correctDistances[i])
		}
	}

}

func TestFindMax(t *testing.T) {
	values := []float64{82, 82, 82, 82, 82, 82, 82, 81, 81, 83}
	max := findMax(values)

	if max != 83 {
		t.Error()
	}

	values = []float64{-2, -1, 0, 1, 2}
	max = findMax(values)

	if max != 2 {
		t.Error()
	}

	values = []float64{0, 2, -1, 1, -2}
	max = findMax(values)

	if max != 2 {
		t.Error()
	}

	values = []float64{0, -3, -1, -4, -2}
	max = findMax(values)

	if max != 0 {
		t.Error()
	}
}

//func TestFindNearest(t *testing.T){
//	t.Error()
//}

//
//func TestTricubeWeightFunction(t *testing.T) {
//	values := []float64{1, 2, 3, 4, 5}
//	correct := []float64{0, 0.6699, 1, 0.6699, 0}
//	weights := tricubeWeightFunction(values, 3)
//	for i := 0; i < len(weights); i++ {
//		if math.Round(weights[i]*10000)/10000 != correct[i] {
//			t.Error()
//		}
//	}
//
//	fmt.Println("----------------------")
//	coordinates := []Coord{
//		{
//			X: 0.5578196,
//			Y: 18.63654,
//		},
//		{
//			X: 2.0217271,
//			Y: 103.49646,
//		},
//		{
//			X: 2.5773252,
//			Y: 150.35391,
//		},
//		{
//			X: 3.4140288,
//			Y: 190.51031,
//		},
//		{
//			X: 4.3014084,
//			Y: 208.70115,
//		},
//		{
//			X: 4.7448394,
//			Y: 213.71135,
//		},
//		{
//			X: 5.1073781,
//			Y: 228.49353,
//		},
//	}
//	var xCoords []float64
//
//	for i := 0; i < len(coordinates); i++ {
//		xCoords = append(xCoords, coordinates[i].X)
//	}
//	weights = tricubeWeightFunction(xCoords, xCoords[0])
//	fmt.Println("Weights")
//	fmt.Println(weights)
//
//}

func TestCalcLOESS(t *testing.T) {

	values := []Coord{
		{
			X: 0.5578196,
			Y: 18.63654,
		},
		{
			X: 2.0217271,
			Y: 103.49646,
		},
		{
			X: 2.5773252,
			Y: 150.35391,
		},
		{
			X: 3.4140288,
			Y: 190.51031,
		},
		{
			X: 4.3014084,
			Y: 208.70115,
		},
		{
			X: 4.7448394,
			Y: 213.71135,
		},
		{
			X: 5.1073781,
			Y: 228.49353,
		},
		{
			X: 6.5411662,
			Y: 233.55387,
		},
		{
			X: 6.7216176,
			Y: 234.55054,
		},
		{
			X: 7.2600583,
			Y: 223.89225,
		},
		{
			X: 8.1335874,
			Y: 227.68339,
		},
		{
			X: 9.1224379,
			Y: 223.91982,
		},
		{X: 11.9296663, Y: 168.01999},
		{X: 12.3797674, Y: 164.95750},
		{X: 13.2728619, Y: 152.61107},
		{X: 14.2767453, Y: 160.78742},
		{X: 15.3731026, Y: 168.55567},
		{X: 15.6476637, Y: 152.42658},
		{X: 18.5605355, Y: 221.70702},
		{X: 18.5866354, Y: 222.69040},
		{X: 18.7572812, Y: 243.18828},
	}

	answers := []float64{
		20.593, 107.160, 139.767, 174.263, 207.233, 216.662, 220.544, 229.861, 229.835, 229.430, 226.604, 220.390, 172.348, 163.842, 161.849, 160.335, 160.192, 161.056, 227.340, 227.899, 231.559,
	}

	loessPoints := CalcLOESS(values, 6)
	_, yPoints := CoordsToArrays(loessPoints)

	for i := 0; i < len(yPoints); i++ {
		if math.Round(yPoints[i]*1000)/1000 != answers[i] {
			t.Errorf("Loess returned %v not %v", math.Round((yPoints[i]*1000))/1000, answers[i])
		}
	}

}
