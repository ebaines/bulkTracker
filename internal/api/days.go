package api

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

//-------------
//GET
//-------------

func dbGetDays(count string) ([]dayRecord, error) {
	requestedDays, err := strconv.Atoi(count)
	if err != nil {
		return []dayRecord{}, err
	}

	rows, err := DB.Query("SELECT * FROM (SELECT  date, weight_kg, calories_kcal FROM WEIGHT ORDER BY date DESC LIMIT ?) ORDER BY date ASC;", requestedDays)
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

func GetDays(c *gin.Context) {
	count := c.Query("count")
	sorting := c.Query("sort")

	// Set defaults
	if count == "" {
		count = "10"
	}
	
	if sorting == ""{
		sorting = "recent"
	}

	var records []dayRecord

	switch sorting {
	case "recent":
		var err error
		records, err = dbGetDays(count)
		if err != nil {
			log.Print(err)
			c.JSON(http.StatusInternalServerError, gin.H{"status": "failure", "error": "SQL error"})
			return
		}
	case "":
		var err error
		records, err = dbGetDays(count)
		if err != nil {
			log.Print(err)
			c.JSON(http.StatusInternalServerError, gin.H{"status": "failure", "error": "SQL error"})
			return
		}
	default:
		c.JSON(http.StatusBadRequest, gin.H{"status": "failure", "error": "Bad filters"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": records})
}
