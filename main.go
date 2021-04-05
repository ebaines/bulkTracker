package main

import (
	"bytes"
	"database/sql"
	api "git.ebain.es/healthAndFitnessTracker/internal/api"
	database "git.ebain.es/healthAndFitnessTracker/internal/database"
	regression "git.ebain.es/healthAndFitnessTracker/internal/regression"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/wcharczuk/go-chart"
)

const dateFormat = "02/01/2006"
const sqliteConnString = "file:/home/ebaines/Downloads/fitness.db"


func main() {
	r := gin.Default()

	r.GET("/table", genTable)

	apiRouter := r.Group("/api")
	{
		apiRouter.GET("/weight/:id", api.GetWeight)
		apiRouter.POST("/weight", api.AddWeight)
		apiRouter.DELETE("/weight/:id", api.DeleteWeight)
	}
	_ = r.Run()
}

func genTable(c *gin.Context) {
	table := processDatabase()
	c.Data(http.StatusOK,
		"text/html; charset=utf-8", []byte(table))
}

func processDatabase() string {
	db, err := sql.Open("sqlite3", sqliteConnString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var dateRange = 1000
	records := database.GetFinalRows(db, dateRange)

	var dates = make([]time.Time, 0, 1000)
	var weightDates = make([]time.Time, 0, 1000)
	var calorieDates = make([]time.Time, 0, 1000)
	var weights, calories []float64

	// Only plot the datapoint if the weight/calorie isn't NULL in the table.
	for _, record := range records {
		dates = append(dates, record.Time)
		if record.Weight.Valid {
			weightDates = append(weightDates, record.Time)
			weights = append(weights, record.Weight.Float64)
		}
		if record.Calories.Valid {
			calorieDates = append(calorieDates, record.Time)
			calories = append(calories, record.Calories.Float64)
		}
	}

	//Calculate smoothed line for weights.
	_, loessWeights := regression.CoordsToArrays(loessSmoothTimeSeries(dates, weightDates, weights, 0.2))

	//Calculate smoothed line for calories
	_, loessCalories := regression.CoordsToArrays(loessSmoothTimeSeries(dates, calorieDates, calories, 0.2))

	//Calculate weight change per day and smooth.
	dayWeightDelta := calculateDayDifferences(weights, 1)
	_, loessDayWeightDelta := regression.CoordsToArrays(loessSmoothTimeSeries(dates, weightDates, dayWeightDelta, 0.4))

	//Calculate calories consumed per kg of bodyweight each day and smooth.
	//var caloriesPerKg []float64
	//for i, weight := range weights {
	//	caloriesPerKg = append(caloriesPerKg, calories[i]/weight)
	//}
	//_, loessCaloriesPerKg := regression.CoordsToArrays(loessSmoothTimeSeries(dates, calorieDates, caloriesPerKg, 0.3))

	//Calculate current TDEE
	calorieSlidingAverage := slidingAvgs(loessCalories, 14)
	weightDeltaSlidingAverage := slidingAvgs(loessDayWeightDelta, 14)
	var tdee []float64
	for i, calorieAverage := range calorieSlidingAverage {
		//tdee = append(tdee, calorieAverage-weightDeltaSlidingAverage[i]*7000)
		tdee = append(tdee, calculateTDEE(calorieAverage, weightDeltaSlidingAverage[i]))
	}

	differences := calculateDayDifferences(loessWeights, 7)
	bigDifferences := calculateDayDifferences(loessWeights, 28)
	differencesCals := calculateDayDifferences(loessCalories, 7)

	var t HtmlTable
	//t.setHeaders([]string{"Date", "Calories", "Day's Weight", "Rolling Weight", "Rolling Smoothed Calories", "1 Day ΔM", "7 Day ΔM", "28 Day ΔM", "7 Day ΔKCal", "TDEE"})
	t.setHeaders([]string{"Date", "Rolling Weight", "Rolling Smoothed Calories", "1 Day ΔM", "7 Day ΔM", "28 Day ΔM", "7 Day ΔKCal", "TDEE"})
	for i := 0; i < len(differences); i++ {
		t.addRow([]string{
			dates[i].Format(dateFormat),
			//strconv.FormatFloat(calories[i], 'f', 2, 64),
			//strconv.FormatFloat(weights[i], 'f', 2, 64),
			strconv.FormatFloat(loessWeights[i], 'f', 2, 64),
			strconv.FormatFloat(loessCalories[i], 'f', 2, 64),
			strconv.FormatFloat(loessDayWeightDelta[i], 'f', 2, 64),
			strconv.FormatFloat(differences[i], 'f', 2, 64),
			strconv.FormatFloat(bigDifferences[i], 'f', 2, 64),
			strconv.FormatFloat(differencesCals[i], 'f', 2, 64),
			strconv.FormatFloat(bufferStart(tdee, 27, 0.0, i), 'f', 2, 64),
		})
	}
	renderedTable := t.render()

	graph := chart.Chart{
		XAxis: chart.XAxis{
			Name: "Date",
		},
		YAxis: chart.YAxis{
			Name: "Weight /kg",
		},
		YAxisSecondary: chart.YAxis{
			Name: "Calories /kcal",
		},
		Series: []chart.Series{
			chart.TimeSeries{
				Name:    "Smoothed Daily Weight",
				XValues: weightDates,
				YValues: loessWeights,
			},
			chart.TimeSeries{
				Name:    "Smoothed Daily Calories",
				YAxis:   chart.YAxisSecondary,
				XValues: calorieDates[7:],
				YValues: loessCalories[7:],
			},
		},
	}
	graph.Elements = []chart.Renderable{
		chart.LegendThin(&graph),
	}

	//graph2 := chart.Chart{
	//	XAxis: chart.XAxis{
	//		Name: "Date",
	//	},
	//	YAxis: chart.YAxis{
	//		Name: "Weight Gain /kg/day",
	//	},
	//	YAxisSecondary: chart.YAxis{
	//		Name: "Calories per Kg /kcal",
	//	},
	//	Series: []chart.Series{
	//		chart.TimeSeries{
	//			Name:    "Smoothed Weight Gain",
	//			XValues: weightDates,
	//			YValues: loessDayWeightDelta,
	//		},
	//		chart.TimeSeries{
	//			Name:    "Smoothed Daily Calories",
	//			YAxis:   chart.YAxisSecondary,
	//			XValues: calorieDates,
	//			YValues: loessCaloriesPerKg,
	//		},
	//	},
	//}

	//graph2.Elements = []chart.Renderable{
	//	chart.LegendThin(&graph2),
	//}

	graph3 := chart.Chart{
		XAxis: chart.XAxis{
			Name: "Date",
		},
		YAxis: chart.YAxis{
			Name: "Calories /kcal",
		},
		YAxisSecondary: chart.YAxis{
			Name: "Weight /kg",
		},
		Series: []chart.Series{
			//chart.TimeSeries{
			//	Name:    "TDEE per KG",
			//	XValues: dates[28:],
			//	YValues: tdeePerKg,
			//},
			chart.TimeSeries{
				Name:    "TDEE",
				XValues: weightDates[28:],
				YValues: tdee,
			},
			chart.TimeSeries{
				Name:    "Weight /kg",
				YAxis:   chart.YAxisSecondary,
				XValues: weightDates,
				YValues: loessWeights,
			},
		},
	}

	graph3.Elements = []chart.Renderable{
		chart.LegendThin(&graph3),
	}

	buffer := bytes.NewBuffer([]byte{})
	err = graph.Render(chart.PNG, buffer)
	//buffer2 := bytes.NewBuffer([]byte{})
	//err = graph2.Render(chart.PNG, buffer2)
	buffer3 := bytes.NewBuffer([]byte{})
	err = graph3.Render(chart.PNG, buffer3)

	//f, err := os.Create("./output.png")
	err = ioutil.WriteFile("output.png", buffer.Bytes(), 0644)
	//f, err := os.Create("./output.png")
	//err = ioutil.WriteFile("output2.png", buffer2.Bytes(), 0644)
	err = ioutil.WriteFile("output3.png", buffer3.Bytes(), 0644)

	return renderedTable
}

func slidingAvgs(dayValues []float64, width int) []float64 {
	var weekAvg float64
	var avgSlice = make([]float64, 0)

	for i := width - 1; i < len(dayValues); i++ {
		weekAvg = 0
		for _, dayValue := range dayValues[i-(width-1) : i+1] {
			weekAvg += dayValue
		}

		weekAvg = weekAvg / float64(width)

		avgSlice = append(avgSlice, weekAvg)
	}
	return avgSlice
}

func calculateDifferences(values []float64) []float64 {
	var diffSlice = make([]float64, 0)

	for i := 1; i < len(values); i++ {
		diffSlice = append(diffSlice, values[i]-values[i-1])
	}
	return diffSlice
}

func calculateDayDifferences(values []float64, days int) []float64 {
	var diffSlice = make([]float64, 0)

	for i := 0; i < len(values); i++ {
		if i-days < 0 {
			diffSlice = append(diffSlice, 0)
		} else {
			diffSlice = append(diffSlice, values[i]-values[i-days])
		}
	}
	return diffSlice
}

func calculateTDEE(avgDayCalories float64, avgDayWeightDiff float64) float64 {
	// Approx 3500kcal = 450g fat
	fatCaloriesDiff := (avgDayWeightDiff * 3500) / 0.450
	tdee := avgDayCalories - fatCaloriesDiff
	return tdee
}

func loessSmoothTimeSeries(datesToEstimate []time.Time, dates []time.Time, yPoints []float64, bandwidth float64) []regression.Coord {
	// Calculate smoothed line for weights.
	var xPointsToEstimate = make([]float64, 0, len(datesToEstimate))
	for xPoint, _ := range datesToEstimate {
		xPointsToEstimate = append(xPointsToEstimate, float64(xPoint))
	}
	var xPoints = make([]float64, 0, len(dates))
	for xPoint, _ := range dates {
		xPoints = append(xPoints, float64(xPoint))
	}

	var coordinates []regression.Coord
	for i := 0; i < len(xPoints); i++ {
		coordinates = append(coordinates, regression.Coord{
			X: xPoints[i],
			Y: yPoints[i],
		})
	}

	loessCoords, err := regression.CalcLOESS(xPointsToEstimate, coordinates, bandwidth)
	if err != nil {
		log.Fatal(err)
	}

	return loessCoords

}

func bufferStart(series []float64, start int, buffer float64, i int) float64 {
	if i < start {
		return buffer
	} else {
		return series[i-start]
	}
}
