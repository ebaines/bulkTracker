package main

import (
	"database/sql"
	"git.ebain.es/healthAndFitnessTracker/internal/api"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
)

const dateFormat = "2006-01-02"
const sqliteConnString = "file:fitness.db"

var DB *sql.DB

func main() {
	db, err := sql.Open("sqlite3", sqliteConnString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	DB = db
	api.DB = db
	api.DateFormat = dateFormat

	r := api.ConfigureRouter()

	r.GET("/table", genTable)

	_ = r.Run()
}

func genTable(c *gin.Context) {
	table := processDatabase()
	c.Data(http.StatusOK,
		"text/html; charset=utf-8", []byte(table))
}
