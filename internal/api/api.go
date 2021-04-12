package api

import (
	"database/sql"
	"errors"
	"git.ebain.es/healthAndFitnessTracker/internal/helpers"
	"github.com/gin-gonic/gin"
	"github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"time"
)

var DB *sql.DB
var DateFormat string

type dayRecord struct {
	date     string
	weight   sql.NullFloat64
	calories sql.NullInt64
}

func parseDayJSON(json interface{}) (dayRecord, error) {
	var record dayRecord

	m := json.(map[string]interface{})

	record.date = time.Now().Format(DateFormat)

	for k, v := range m {
		switch k {
		case "time":
			timestamp := v.(float64)
			record.date = time.Unix(int64(timestamp), 0).Format(DateFormat)
		case "weight":
			weight := v.(float64)
			record.weight.Float64 = helpers.RoundDecimalPlaces(weight, 1)
			record.weight.Valid = true
		case "calories":
			calories := v.(float64)
			record.calories.Int64 = int64(calories)
			record.calories.Valid = true
		default:
			return record, errors.New("invalid JSON - extra key")
		}
	}

	return record, nil

}

func genDayJSON(record dayRecord) (map[string]interface{}, error) {
	ginJSON := make(map[string]interface{})

	parsedTime, err := time.Parse(DateFormat, record.date)
	if err != nil {
		return nil, err
	}
	timestamp := parsedTime.Unix()

	ginJSON["date"] = timestamp

	if record.weight.Valid{
		ginJSON["weight"] = record.weight.Float64
	}

	if record.calories.Valid{
		ginJSON["calories"] = record.calories.Int64
	}

	return ginJSON, nil

}

func handleSQLExecErr(c *gin.Context, err error) {
	log.Print(err)
	if sqliteErr, ok := err.(sqlite3.Error); ok {
		if sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			c.JSON(http.StatusBadRequest, gin.H{"status": "failure", "error": "Day already exists"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "failure", "error": "SQL error"})
		}
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failure", "error": "SQL error"})
	}
}

func AddDay(c *gin.Context) {
	var json interface{}

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failure", "error": err.Error()})
		return
	}

	record, err := parseDayJSON(json)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failure", "error": err.Error()})
		return
	}

	if c.Request.Method == "POST" {
		_, err := DB.Exec("INSERT INTO weight(date, weight_kg, calories_kcal) VALUES (?, ?, ?)", record.date, record.weight, record.calories)
		if err != nil {
			handleSQLExecErr(c, err)
			return
		}
		c.JSON(http.StatusCreated, gin.H{"status": "success"})
	} else if c.Request.Method == "PUT" {
		id := c.Param("id")

		_, err := DB.Exec("UPDATE weight SET date=?, weight_kg=?, calories_kcal=? WHERE id = ?", record.date, record.weight, record.calories, id)
		if err != nil {
			handleSQLExecErr(c, err)
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "success"})
	}
}

func DeleteDay(c *gin.Context) {
	id := c.Param("id")

	_, err := DB.Exec("DELETE FROM weight WHERE id=?", id)
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failure", "error": "SQL error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})

}

func dbGetDay(id string)(dayRecord, error){
	rows, err := DB.Query("SELECT date, weight_kg, calories_kcal FROM weight WHERE id=?", id)
	if err != nil {
		return dayRecord{}, err
	}
	defer rows.Close()

	var record dayRecord

	for rows.Next() {
		err = rows.Scan(&record.date, &record.weight, &record.calories)
		if err != nil {
			return dayRecord{}, err
		}
	}

	return record, nil
}

func GetDay(c *gin.Context) {
	id := c.Param("id")

	record, err := dbGetDay(id)
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failure", "error": "SQL error"})
		return
	}

	responseJSON, err := genDayJSON(record)
	if err != nil{
		log.Print(err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failure", "error": "Error generating JSON"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": responseJSON})

}
