package database

import (
	"database/sql"
	"log"
	"strconv"
	"time"
)

type dayEntry struct {
	Time     time.Time
	Weight   sql.NullFloat64
	Calories sql.NullFloat64
}

const dateFormat = "02/01/2006"

func GetFinalRows(dbConn *sql.DB, numRows int) []dayEntry {
	sqlCntStmt :=
		"SELECT COUNT(date) FROM weight ORDER BY id DESC LIMIT " + strconv.Itoa(numRows) + ";"
	sqlStmt :=
		"SELECT date, weight_kg, calories_kcal FROM weight ORDER BY id DESC LIMIT " + strconv.Itoa(numRows) + ";"

	rows, err := dbConn.Query(sqlCntStmt)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var count int
	for rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			log.Fatal(err)
		}
	}

	if count > numRows {
		count = numRows
	}

	rows, err = dbConn.Query(sqlStmt)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	entrySlice := make([]dayEntry, 0, count)

	var date string
	var entry dayEntry

	for rows.Next() {

		err = rows.Scan(&date, &entry.Weight, &entry.Calories)

		if err != nil {
			log.Fatal(err)
		}

		parsedDate, _ := time.Parse(dateFormat, date)
		entry.Time = parsedDate
		entrySlice = append(entrySlice, entry)

	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	for i, j := 0, len(entrySlice)-1; i < j; i, j = i+1, j-1 {
		entrySlice[i], entrySlice[j] = entrySlice[j], entrySlice[i]
	}

	return entrySlice
}
