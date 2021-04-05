package api

import (
	"database/sql"
	"git.ebain.es/healthAndFitnessTracker/internal/helpers"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

const dateFormat = "02/01/2006"
const sqliteConnString = "file:/home/ebaines/Downloads/fitness.db"

type weightEntry struct {
	Time   int64   `json:"timestamp"`
	Weight float64 `json:"weight" binding:"required"`
}

type parsedWeightEntry struct {
	date   string
	weight float64
}

func AddWeight(c *gin.Context) {
	db, err := sql.Open("sqlite3", sqliteConnString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var json weightEntry

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var parsedEntry parsedWeightEntry

	if json.Time == 0{
		parsedEntry.date = time.Now().Format(dateFormat)
	} else{
		parsedEntry.date = time.Unix(json.Time, 0).Format(dateFormat)
	}
	parsedEntry.weight = helpers.RoundDecimalPlaces(json.Weight, 1)

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare("INSERT INTO weight(date, weight_kg) VALUES (?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(parsedEntry.date, parsedEntry.weight)
	if err != nil {
		log.Fatal(err)
	}

	tx.Commit()

	c.JSON(http.StatusCreated, gin.H{"status": "success"})
}

func DeleteWeight(c *gin.Context){
	db, err := sql.Open("sqlite3", sqliteConnString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	id := c.Param("id")

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare("DELETE FROM weight WHERE id=?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		log.Fatal(err)
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{"status": "success"})

}

func GetWeight(c *gin.Context){
	db, err := sql.Open("sqlite3", sqliteConnString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	id := c.Param("id")

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare("SELECT weight_kg FROM weight WHERE id=?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var weight float64

	for rows.Next(){
		err = rows.Scan(&weight)

		if err != nil {
			log.Fatal(err)
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "weight": weight})

}