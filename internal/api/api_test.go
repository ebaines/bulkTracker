package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

const testSql = "test_db.sql"
const testDb = "test.db"
const testSqliteConnString = "file:" + testDb

type postRequestData struct {
	Time     int64   `json:"time,omitempty"`
	Weight   float64 `json:"weight,omitempty"`
	Calories int     `json:"calories,omitempty"`
}

var ROUTER *gin.Engine

func populateDb() {

	query, err := os.ReadFile(testSql)
	if err != nil {
		log.Fatal(err)
	}

	_, err = DB.Exec(string(query))
	if err != nil {
		log.Fatal(err)
	}
}

func dropDb() {
	// Truncate table
	_, err := DB.Exec("DROP TABLE weight;")
	if err != nil {
		log.Fatal(err)
	}
}

func resetDb() {
	dropDb()
	populateDb()
}

func deleteDbFile() {
	err := os.Remove(testDb)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			log.Fatal(err)
		}
	}
}

func dbRows() int {
	sqlCntStmt :=
		"SELECT COUNT(date) FROM weight ORDER BY id;"

	rows, err := DB.Query(sqlCntStmt)
	if err != nil {
		log.Print("error in dbRows")
		log.Fatal(err)
	}

	var count int

	for rows.Next() {
		rows.Scan(&count)
	}

	defer rows.Close()

	return count
}

func TestMain(m *testing.M) {
	deleteDbFile()

	db, err := sql.Open("sqlite3", testSqliteConnString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	router := ConfigureRouter()

	DB = db
	DateFormat = "2006-01-02"
	ROUTER = router
	populateDb()
	m.Run()
	deleteDbFile()
}

func TestGetDay(t *testing.T) {

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/day/1", nil)
	ROUTER.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "{\"data\":{\"calories\":2900,\"date\":1579910400,\"weight\":82},\"status\":\"success\"}", w.Body.String())
	assert.Equal(t, 20, dbRows())

	resetDb()
}

func TestGetNonExistentDay(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/day/1000", nil)
	ROUTER.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	resetDb()
}

func TestGetMalformedDay(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/day/2", nil)
	ROUTER.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	//assert.Equal(t, "{\"data\":{\"calories\":2900,\"date\":1579910400,\"weight\":82},\"status\":\"success\"}", w.Body.String())
	//assert.Equal(t, 428, dbRows())

	resetDb()
}

func TestAddDay(t *testing.T) {
	testData := []postRequestData{
		{
			Time:     1617267916,
			Weight:   80,
			Calories: 3000,
		},
		{
			Time:     1617354316,
			Calories: 3100,
		},
		{
			Time:   1617440716,
			Weight: 80.2,
		},
		{},
	}

	currentDay := time.Now().Format(DateFormat)

	correctData := []dayRecord{
		{
			date: "2021-04-01",
			weight: sql.NullFloat64{
				Valid:   true,
				Float64: 80,
			},
			calories: sql.NullInt64{
				Int64: 3000,
				Valid: true,
			},
		},
		{
			date: "2021-04-02",
			weight: sql.NullFloat64{
				Valid: false,
			},
			calories: sql.NullInt64{
				Int64: 3100,
				Valid: true,
			},
		},
		{
			date: "2021-04-03",
			weight: sql.NullFloat64{
				Valid:   true,
				Float64: 80.2,
			},
			calories: sql.NullInt64{
				Valid: false,
			},
		},
		{
			date:     currentDay,
			weight:   sql.NullFloat64{},
			calories: sql.NullInt64{},
		},
	}

	for i, postData := range testData {
		jsonData, err := json.Marshal(postData)
		if err != nil {
			log.Fatal(err)
		}

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/day", bytes.NewReader(jsonData))
		ROUTER.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		dbId := 21 + i
		assert.Equal(t, dbId, dbRows())
		storedData, err := dbGetDay(strconv.Itoa(dbId))
		if err != nil {
			t.Error()
		}
		assert.EqualValues(t, storedData, correctData[i])
	}

	resetDb()
}

func TestAddExistingDay(t *testing.T) {
	// Test adding when day already exists.
	day, err := time.Parse(DateFormat, "2020-01-25")
	if err != nil {
		t.Error(err)
	}

	jsonData, err := json.Marshal(postRequestData{Time: day.Unix()})
	if err != nil {
		log.Fatal(err)
	}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/day", bytes.NewReader(jsonData))
	ROUTER.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	resetDb()
}

func TestAddMalformedBody(t *testing.T){
	requestJson := "asodhflksahf"

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/day", strings.NewReader(requestJson))
	ROUTER.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	resetDb()
}

func TestReplaceExistingDay(t *testing.T) {
	day, err := time.Parse(DateFormat, "2020-01-25")
	if err != nil {
		t.Error(err)
	}

	jsonData, err := json.Marshal(postRequestData{Time: day.Unix()})
	if err != nil {
		log.Fatal(err)
	}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/day/1", bytes.NewReader(jsonData))
	ROUTER.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	resetDb()
}

func TestAddMalformedJSON(t *testing.T) {
	requestJson := "{\"malformed_json\": \"test\"}"

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/day", strings.NewReader(requestJson))
	ROUTER.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	resetDb()
}

func TestDeleteDay(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/day/1", nil)
	ROUTER.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	resetDb()
}
