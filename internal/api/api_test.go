package api

import (
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

const testSqliteConnString = "file:test.db"
const testSql = "test_db.sql"
const testDb = "test.db"

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

func dropDb(){
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

func deleteDbFile(){
	err := os.Remove(testDb)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist){
			log.Fatal(err)
		}
	}
}

func TestMain(m *testing.M) {
	deleteDbFile()

	db, err := sql.Open("sqlite3", testSqliteConnString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	DB = db
	DateFormat = "02/01/2006"

	populateDb()
	m.Run()
	deleteDbFile()
}

func TestAddDay(t *testing.T) {
	router := ConfigureRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/day/1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "{\"data\":{\"calories\":2900,\"time\":1579910400,\"weight\":82},\"status\":\"success\"}", w.Body.String())
}
