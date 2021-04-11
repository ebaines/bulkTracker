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

	tx, err := DB.Begin()
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failure", "error": "SQL error"})
		return
	}

	if c.Request.Method == "POST" {
		stmt, err := tx.Prepare("INSERT INTO weight(date, weight_kg, calories_kcal) VALUES (?, ?, ?)")
		if err != nil {
			log.Print(err)
			c.JSON(http.StatusInternalServerError, gin.H{"status": "failure", "error": "SQL error"})
			return
		}
		defer stmt.Close()

		_, err = stmt.Exec(record.date, record.weight, record.calories)
		if err != nil {
			handleSQLExecErr(c, err)
			err = tx.Rollback()
			if err != nil {
				log.Print(err)
			}
			return
		}
	} else if c.Request.Method == "PUT" {
		id := c.Param("id")

		stmt, err := tx.Prepare("UPDATE weight SET date=?, weight_kg=?, calories_kcal=? WHERE id = ?")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()

		_, err = stmt.Exec(record.date, record.weight, record.calories, id)
		if err != nil {
			handleSQLExecErr(c, err)
			err = tx.Rollback()
			if err != nil {
				log.Print(err)
			}
			return
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failure", "error": "SQL error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "success"})
}

func DeleteDay(c *gin.Context) {
	id := c.Param("id")

	tx, err := DB.Begin()
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failure", "error": "SQL error"})
		return
	}

	stmt, err := tx.Prepare("DELETE FROM weight WHERE id=?")
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failure", "error": "SQL error"})
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failure", "error": "SQL error"})
		return
	}

	err = tx.Commit()
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failure", "error": "SQL error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})

}

func GetDay(c *gin.Context) {
	id := c.Param("id")

	tx, err := DB.Begin()
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failure", "error": "SQL error"})
		return
	}

	stmt, err := tx.Prepare("SELECT date, weight_kg, calories_kcal FROM weight WHERE id=?")
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failure", "error": "SQL error"})
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(id)
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failure", "error": "SQL error"})
		return
	}
	defer rows.Close()

	var date string
	var weight float64
	var calories float64
	var timestamp int64

	for rows.Next() {
		err = rows.Scan(&date, &weight, &calories)
		if err != nil {
			log.Print(err)
			c.JSON(http.StatusInternalServerError, gin.H{"status": "failure", "error": "SQL error"})
			return
		}

		parsedTime, err := time.Parse(DateFormat, date)
		if err != nil {
			log.Print(err)
			c.JSON(http.StatusInternalServerError, gin.H{"status": "failure", "error": "Error parsing data from database."})
			return
		}
		timestamp = parsedTime.Unix()
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"time": timestamp, "weight": weight, "calories": calories}})

}
