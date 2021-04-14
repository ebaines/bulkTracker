package api

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"time"
)

//-------------
//GET
//-------------

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
	filterType := c.Query("filterType")
	count := c.Query("count")
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")

	// Set defaults
	if filterType == "" {
		filterType = "recent"
	}

	if count == "" {
		count = "3"
	}

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

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": records})
}
