package regression

import (
	"errors"
	"fmt"
	"math"
	"sort"
)

type Coord struct {
	X float64
	Y float64
}

type coordDist struct {
	coord Coord
	dist  float64
}

type coordDistSlice []coordDist

func (s coordDistSlice) Len() int{
	return len(s)
}

func (s coordDistSlice) Swap(i, j int){
	s[i], s[j] = s[j], s[i]
}

func (s coordDistSlice) Less(i, j int) bool{
	return s[i].dist < s[j].dist
}

func CoordsToArrays(coords []Coord) ([]float64, []float64) {
	var xCoords []float64
	var yCoords []float64

	for i := 0; i < len(coords); i++ {
		xCoords = append(xCoords, coords[i].X)
		yCoords = append(yCoords, coords[i].Y)
	}
	return xCoords, yCoords
}

func findDist(a float64, b float64) float64 {
	return math.Abs(a - b)
}

func findMax(values []float64) float64 {
	max := 0.0
	for _, value := range values {
		if value > max {
			max = value
		}
	}
	return max
}

func findNearest(sortedCoords []Coord, count int, centreCoordIndex int) (coordDistSlice, error) {
	//check that they're sorted
	//COMPLETE
	//check that
	if centreCoordIndex < 0 || centreCoordIndex >= len(sortedCoords) {
		return nil, errors.New("findnearest: the centre coord index is out of bounds of the sortedcoords")
	}
	if count > len(sortedCoords) {
		return nil, errors.New("findnearest: cannot return more coords than input")
	}

	var distances coordDistSlice
	for i := 0; i < len(sortedCoords); i++ {
		distances = append(distances, coordDist{
			coord: sortedCoords[i],
			dist:  findDist(sortedCoords[centreCoordIndex].X, sortedCoords[i].X),
		},
		)
	}

	sort.Sort(coordDistSlice(distances))

	return distances[:count + 1], nil
}

func tricubeWeightFunction(sortedCoordDists coordDistSlice) []float64 {
	//https://uk.mathworks.com/help/curvefit/smoothing-data.html

	weights := make([]float64, len(sortedCoordDists))
	maxDist := sortedCoordDists[len(sortedCoordDists) - 1].dist

	for i := 0; i < len(sortedCoordDists); i++ {
		weights[i] = math.Pow(1-math.Pow(math.Abs(sortedCoordDists[i].dist/maxDist), 3), 3)
	}
	return weights
}

func weightedMean(values []float64, weights []float64) (float64, error) {
	if len(weights) != len(values) {
		return 0, errors.New("regression: weighted mean requires equal length weight and value slices")
	}

	var sumWeights float64
	for i := 0; i < len(weights); i++ {
		sumWeights = sumWeights + weights[i]
	}

	var sumWeightedValues float64
	for i := 0; i < len(values); i++ {
		sumWeightedValues = sumWeightedValues + values[i]*weights[i]
	}

	return sumWeightedValues / sumWeights, nil
}

func wLSRegression(coordinates coordDistSlice, weights []float64) (float64, float64, error) {
	if len(weights) != len(coordinates) {
		return 0, 0, errors.New("regression: wls regressions requires coordinate and weight slices of equal length")
	}

	var xCoords []float64
	var yCoords []float64

	for i := 0; i < len(coordinates); i++ {
		xCoords = append(xCoords, coordinates[i].coord.X)
		yCoords = append(yCoords, coordinates[i].coord.Y)
	}

	weightedMeanX, err := weightedMean(xCoords, weights)
	weightedMeanY, err := weightedMean(yCoords, weights)

	if err != nil {
		fmt.Println(err)
	}

	var sumNumerator float64
	var sumDenominator float64
	for i := 0; i < len(xCoords); i++ {
		sumNumerator = sumNumerator + weights[i]*(xCoords[i]-weightedMeanX)*(yCoords[i]-weightedMeanY)
		sumDenominator = sumDenominator + weights[i]*math.Pow(xCoords[i]-weightedMeanX, 2)
	}

	var slope = sumNumerator / sumDenominator
	var intercept = weightedMeanY - slope*weightedMeanX

	return slope, intercept, nil
}

func CalcLOESS(coordinates []Coord, nearestNeighboursCount int) []Coord {
	var loessPoints []Coord

	// For each coordinate, calculate WLS regression line, then evaluate at point of estimation.
	for i := 0; i < len(coordinates); i++ {
		var widthCoords coordDistSlice

		// Capture coordinates within the width
		widthCoords, err := findNearest(coordinates, nearestNeighboursCount, i)

		weights := tricubeWeightFunction(widthCoords)
		slope, intercept, err := wLSRegression(widthCoords, weights)
		if err != nil {
			fmt.Println(err)
		}

		//fmt.Println("Weights:", weights)
		//fmt.Println("Slope: ", slope)
		//fmt.Println("Intercept: ", intercept)

		estimatedValue := slope*coordinates[i].X + intercept
		//fmt.Println("\033[0;92m", estimatedValue, "\033[0m")
		loessPoints = append(loessPoints, Coord{
			X: coordinates[i].X,
			Y: estimatedValue,
		})
	}

	return loessPoints
}
