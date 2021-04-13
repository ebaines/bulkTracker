package api

import (
	"github.com/stretchr/testify/assert"
	"log"
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
	log.Print(w.Body.String())

	resetDb()
}
