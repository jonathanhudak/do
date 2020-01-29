package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

var app App

func setup() *App {
	// Initialize an in-memory database for full integration testing.
	app := &App{}
	app.Initialize("sqlite3", "test.db")
	return app
}

func teardown(app *App) {
	// Closing the connection discards the in-memory database.
	app.DB.Close()
}

func TestGetEntries(t *testing.T) {
	app := setup()
	defer teardown(app)

	server := httptest.NewServer(app.Router)
	defer server.Close()

	res, err := http.Get(fmt.Sprintf("%s/api/entries", server.URL))
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != 200 {
		t.Fatalf("Received non-200 response: %d\n", res.StatusCode)
	}
	result := []Entry{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		t.Fatal(err)
	}
}
