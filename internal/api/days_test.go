package api

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetDays(t *testing.T) {

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/days?sort=recent&count=3", nil)
	ROUTER.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "{\"data\":[{\"calories\":3000,\"date\":1581379200,\"weight\":83.8},{\"calories\":3000,\"date\":1581465600,\"weight\":83.9},{\"calories\":3000,\"date\":1581552000,\"weight\":84}],\"status\":\"success\"}", w.Body.String())

	resetDb()
}

func TestGetDaysDefaults(t *testing.T) {

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/days", nil)
	ROUTER.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "{\"data\":[{\"calories\":3100,\"date\":1580774400,\"weight\":83.6},{\"calories\":3150,\"date\":1580860800,\"weight\":83.8},{\"calories\":3150,\"date\":1580947200,\"weight\":83.9},{\"calories\":3150,\"date\":1581033600,\"weight\":83.8},{\"calories\":3100,\"date\":1581120000,\"weight\":83.3},{\"calories\":3100,\"date\":1581206400,\"weight\":83.3},{\"calories\":3000,\"date\":1581292800,\"weight\":83.3},{\"calories\":3000,\"date\":1581379200,\"weight\":83.8},{\"calories\":3000,\"date\":1581465600,\"weight\":83.9},{\"calories\":3000,\"date\":1581552000,\"weight\":84}],\"status\":\"success\"}", w.Body.String())

	resetDb()
}
