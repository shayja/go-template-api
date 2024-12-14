package app_test

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/shayja/go-template-api/cmd/app"
	"github.com/stretchr/testify/assert"
)

func setupMockDB() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic("Failed to create sqlmock")
	}
	return db, mock
}

func TestRoutes_ProductEndpoints(t *testing.T) {
	a := &app.App{}
	db, _ := setupMockDB()
	defer db.Close()

	a.DB = db
	a.Routes()

	assert.NotNil(t, a.Router, "Router should not be nil after calling Routes")
	assert.NotNil(t, a.DB, "DB should not be nil after calling Routes")

	server := httptest.NewServer(a.Router)
	defer server.Close()

	// Test GET /product with mock JWT middleware bypass
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/product", server.URL), nil)
	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err, "Expected no error for the product GET request")
	assert.NotNil(t, resp, "Expected a response for the product GET request")
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode, "Expected 401 Unauthorized when JWT is missing")
}
