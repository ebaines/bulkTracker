package api

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetDays(t *testing.T) {

	correctResponseData := map[string]interface{}{
		"status": "success",
		"data": []interface{}{
			map[string]interface{}{
				"date":     "2020-02-09",
				"calories": 3100.0,
				"weight":   83.3,
			},
			map[string]interface{}{
				"date":     "2020-02-10",
				"calories": 3000.0,
				"weight":   83.3,
			},
			map[string]interface{}{
				"date":     "2020-02-11",
				"calories": 3000.0,
				"weight":   83.8,
			},
			map[string]interface{}{
				"date":     "2020-02-12",
				"calories": 3000.0,
				"weight":   83.9,
			},
			map[string]interface{}{
				"date":     "2020-02-13",
				"calories": 3000.0,
				"weight":   84.0,
			},
		},
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/days?filterType=recent&count=5", nil)
	ROUTER.ServeHTTP(w, req)

	parsedJson, err := parseTestJSON(w.Body.String())
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, correctResponseData, parsedJson)

	resetDb()
}

func TestGetDaysDefaults(t *testing.T) {
	correctResponseData := map[string]interface{}{
		"status": "success",
		"data": []interface{}{
			map[string]interface{}{
				"date":     "2020-02-11",
				"calories": 3000.0,
				"weight":   83.8,
			},
			map[string]interface{}{
				"date":     "2020-02-12",
				"calories": 3000.0,
				"weight":   83.9,
			},
			map[string]interface{}{
				"date":     "2020-02-13",
				"calories": 3000.0,
				"weight":   84.0,
			},
		},
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/days", nil)
	ROUTER.ServeHTTP(w, req)

	parsedJson, err := parseTestJSON(w.Body.String())
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, correctResponseData, parsedJson)
	resetDb()
}

func TestGetDaysDateRange(t *testing.T) {
	correctResponseData := map[string]interface{}{
		"status": "success",
		"data": []interface{}{
			map[string]interface{}{
				"date":     "2020-02-11",
				"calories": 3000.0,
				"weight":   83.8,
			},
			map[string]interface{}{
				"date":     "2020-02-12",
				"calories": 3000.0,
				"weight":   83.9,
			},
			map[string]interface{}{
				"date":     "2020-02-13",
				"calories": 3000.0,
				"weight":   84.0,
			},
		},
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/days?filterType=dateRange&startDate=2020-02-11&endDate=2020-02-13", nil)
	ROUTER.ServeHTTP(w, req)
	
	parsedJson, err := parseTestJSON(w.Body.String())
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, correctResponseData, parsedJson)

	resetDb()
}

func TestGetDaysSmoothedDateRange(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/days?filterType=dateRange&dataType=smooth&startDate=2020-01-25&endDate=2020-02-13", nil)
	ROUTER.ServeHTTP(w, req)

	parsedJson, err := parseTestJSON(w.Body.String())
	if err != nil {
		t.Error(err)
	}
	
	assertedJson := parsedJson.(map[string]interface{})
	data := assertedJson["data"].([]interface{})
	
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, len(data), 20)

	resetDb()
}

func TestGetDaysInvalidDateRange(t *testing.T) {
	correctResponseData := map[string]interface{}{
		"status": "failure",
		"error":  "Invalid/malformed date range specified.",
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/days?filterType=dateRange&startDate=2020-03-49&endDate=2020-01-30", nil)
	ROUTER.ServeHTTP(w, req)

	parsedJson, err := parseTestJSON(w.Body.String())
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, correctResponseData, parsedJson)

	resetDb()
}
