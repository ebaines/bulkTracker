package api

import (
	"database/sql"
	"encoding/json"
	"git.ebain.es/healthAndFitnessTracker/internal/regression"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"time"
)

type smoothedDayRecord struct {
	date             string
	weight           sql.NullFloat64
	calories         sql.NullInt64
	smoothedWeight   float64
	smoothedCalories float64
}

func (record smoothedDayRecord) MarshalJSON() ([]byte, error) {
	ginJSON := make(map[string]interface{})

	ginJSON["date"] = record.date
	ginJSON["smoothedWeight"] = record.smoothedWeight
	ginJSON["smoothedCalories"] = record.smoothedCalories

	if record.weight.Valid {
		ginJSON["weight"] = record.weight.Float64
	}

	if record.calories.Valid {
		ginJSON["calories"] = record.calories.Int64
	}

	marshalledJson, err := json.Marshal(ginJSON)
	if err != nil {
		return []byte{}, err
	}

	return marshalledJson, nil
}

//-------------
//GET
//-------------

func loessSmoothTimeSeries(datesToEstimate []time.Time, dates []time.Time, yPoints []float64, bandwidth float64) ([]regression.Coord, error) {
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
		log.Print(err)
		return []regression.Coord{}, err
	}

	return loessCoords, nil
}

func dbGetDays(count string) ([]dayRecord, error) {
	requestedDays, err := strconv.Atoi(count)
	if err != nil {
		return []dayRecord{}, err
	}

	rows, err := DB.Query("SELECT * FROM (SELECT  date, weight_kg, calories_kcal FROM weight ORDER BY date DESC LIMIT ?) ORDER BY date ASC;", requestedDays)
	if err != nil {
		return []dayRecord{}, err
	}
	defer rows.Close()

	dayRecords := make([]dayRecord, 0, requestedDays)
	var record dayRecord

	for rows.Next() {
		err = rows.Scan(&record.date, &record.weight, &record.calories)
		if err != nil {
			return []dayRecord{}, err
		}

		dayRecords = append(dayRecords, record)
	}

	return dayRecords, nil
}

func dbGetDayRange(startDate string, endDate string) ([]dayRecord, error) {
	rows, err := DB.Query("SELECT  date, weight_kg, calories_kcal FROM weight WHERE date BETWEEN ? AND ? ORDER BY date;", startDate, endDate)
	if err != nil {
		return []dayRecord{}, err
	}
	defer rows.Close()

	dayRecords := make([]dayRecord, 0)
	var record dayRecord

	for rows.Next() {
		err = rows.Scan(&record.date, &record.weight, &record.calories)
		if err != nil {
			return []dayRecord{}, err
		}

		dayRecords = append(dayRecords, record)

	}

	return dayRecords, nil
}

func GetDays(c *gin.Context) {
	filterType := c.DefaultQuery("filterType", "recent")
	dataType := c.Query("dataType")
	count := c.DefaultQuery("count", "3")
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")

	// Scope variables outside the switch statement.
	var records []dayRecord
	var err error

	switch filterType {
	case "recent":
		var err error

		// Get the most recent $count records.
		records, err = dbGetDays(count)
		if err != nil {
			log.Print(err)
			c.JSON(http.StatusInternalServerError, gin.H{"status": "failure", "error": "SQL error"})
			return
		}
	case "dateRange":
		// Verify that the dates are in a valid format.
		dates := []string{startDate, endDate}

		for _, date := range dates {
			_, err := time.Parse(DateFormat, date)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"status": "failure", "error": "Invalid/malformed date range specified."})
				return
			}
		}

		// Get records within the date range.
		records, err = dbGetDayRange(startDate, endDate)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "failure", "error": "SQL error"})
			return
		}
	case "":
		// Default is no filter type is specified.
		records, err = dbGetDays(count)
		if err != nil {
			log.Print(err)
			c.JSON(http.StatusInternalServerError, gin.H{"status": "failure", "error": "SQL error"})
			return
		}
	default:
		// Default if the filters make no sense.
		c.JSON(http.StatusBadRequest, gin.H{"status": "failure", "error": "Bad filters"})
		return
	}

	var returnData []smoothedDayRecord

	// If requested, smooth the data.
	if dataType == "smoothed" {
		log.Print("smoothing")
		recordLength := len(records)

		datesToEstimate := make([]time.Time, 0, recordLength)
		weightDates := make([]time.Time, 0, recordLength)
		calorieDates := make([]time.Time, 0, recordLength)

		weights := make([]float64, 0, recordLength)
		calories := make([]float64, 0, recordLength)

		for _, record := range records {
			parsedTime, err := time.Parse(DateFormat, record.date)
			if err != nil {
				log.Print(err)
				c.JSON(http.StatusInternalServerError, gin.H{"status": "failure", "error": "SQL error"})
				return
			}
			datesToEstimate = append(datesToEstimate, parsedTime)
			if record.weight.Valid {
				weightDates = append(weightDates, parsedTime)
				weights = append(weights, record.weight.Float64)
			}
			if record.calories.Valid {
				calorieDates = append(calorieDates, parsedTime)
				calories = append(calories, float64(record.calories.Int64))
			}

		}

		smoothedWeights, err := loessSmoothTimeSeries(datesToEstimate, weightDates, weights, 0.3)
		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"status": "failure", "error": "Error smoothing data"})
		}
		smoothedCalories, err := loessSmoothTimeSeries(datesToEstimate, calorieDates, calories, 0.3)
		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"status": "failure", "error": "Error smoothing data"})
		}
		
		returnData = make([]smoothedDayRecord, 0, len(datesToEstimate))

		for i := 0; i < len(datesToEstimate); i++ {
			currentDate := datesToEstimate[i].Format(DateFormat)
			data := smoothedDayRecord{
				date:             currentDate,
				smoothedWeight:   smoothedWeights[i].Y,
				smoothedCalories: smoothedCalories[i].Y,
				weight:           records[i].weight,
				calories:         records[i].calories,
			}

			returnData = append(returnData, data)
		}
	}

	if dataType == "smoothed" {
		c.JSON(http.StatusOK, gin.H{"status": "success", "data": returnData})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": "success", "data": records})
	}
}
