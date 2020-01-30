package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"gotest.tools/assert"
)

var app App

func setup() *App {
	// Initialize an in-memory database for full integration testing.
	app := &App{}
	app.Initialize("sqlite3", ":memory:")
	return app
}

func teardown(app *App) {
	// Closing the connection discards the in-memory database.
	app.DB.Close()
}

func putRequest(url string, data io.Reader) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPut, url, data)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	if err != nil {
		log.Fatal(err)
	}
	res, err := client.Do(req)
	return res, err
}

func deleteRequest(url string, data io.Reader) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodDelete, url, data)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	if err != nil {
		log.Fatal(err)
	}
	res, err := client.Do(req)
	return res, err
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

func TestGetEntry(t *testing.T) {
	app := setup()
	defer teardown(app)

	server := httptest.NewServer(app.Router)
	defer server.Close()

	res, err := http.Get(fmt.Sprintf("%s/api/entry/1", server.URL))
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != 200 {
		t.Fatalf("Received non-200 response: %d\n", res.StatusCode)
	}
	result := Entry{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		t.Fatal(err)
	}
}

func TestCreateEntry(t *testing.T) {
	app := setup()
	defer teardown(app)

	var initialCount int
	app.DB.Table("entries").Count(&initialCount)

	server := httptest.NewServer(app.Router)
	defer server.Close()

	entry := Entry{
		Title: "Testing entry creation",
	}
	requestBody, _ := json.Marshal(entry)
	res, err := http.Post(fmt.Sprintf("%s/api/entry", server.URL), "application/json", bytes.NewBuffer(requestBody))

	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != 200 {
		t.Fatalf("Received non-200 response: %d\n", res.StatusCode)
	}
	result := Entry{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		t.Fatal(err)
	}

	var newCount int
	app.DB.Table("entries").Count(&newCount)
	assert.Equal(t, newCount, initialCount+1)
}

func TestUpdateEntry(t *testing.T) {
	app := setup()
	defer teardown(app)

	// Given I have an entry

	server := httptest.NewServer(app.Router)
	defer server.Close()

	entry := Entry{
		Title: "Testing entry update foo",
	}
	requestBody, _ := json.Marshal(entry)
	res, err := http.Post(fmt.Sprintf("%s/api/entry", server.URL), "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != 200 {
		t.Fatalf("Received non-200 response: %d\n", res.StatusCode)
	}
	createResult := Entry{}
	if err := json.NewDecoder(res.Body).Decode(&createResult); err != nil {
		t.Fatal(err)
	}

	// When I update the entry
	var initialCount int
	app.DB.Table("entries").Count(&initialCount)

	titleUpdate := "Testing entry update"
	entryUpdate := Entry{
		Title: titleUpdate,
	}
	putRequestBody, _ := json.Marshal(entryUpdate)
	putRes, err := putRequest(fmt.Sprintf("%s/api/entry/%d", server.URL, createResult.ID), bytes.NewBuffer(putRequestBody))

	if putRes.StatusCode != 200 {
		t.Fatalf("Received non-200 response: %d\n", putRes.StatusCode)
	}

	result := Entry{}
	if err := json.NewDecoder(putRes.Body).Decode(&result); err != nil {
		t.Fatal(err)
	}

	var nextCount int
	app.DB.Table("entries").Count(&nextCount)
	assert.Equal(t, nextCount, initialCount)
	assert.Equal(t, result.Title, titleUpdate)
}

func TestDeleteEntry(t *testing.T) {
	app := setup()

	// Given I have an entry
	server := httptest.NewServer(app.Router)
	defer server.Close()

	entry := Entry{
		Title: "Testing entry deletion",
	}
	requestBody, _ := json.Marshal(entry)
	res, err := http.Post(fmt.Sprintf("%s/api/entry", server.URL), "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != 200 {
		t.Fatalf("Received non-200 response: %d\n", res.StatusCode)
	}
	createResult := Entry{}
	if err := json.NewDecoder(res.Body).Decode(&createResult); err != nil {
		t.Fatal(err)
	}

	// When I delete the entry
	var initialCount int
	app.DB.Table("entries").Where("deleted_at", "NULL").Count(&initialCount)
	assert.Equal(t, initialCount, 0)

	delRes, err := deleteRequest(fmt.Sprintf("%s/api/entry/%d", server.URL, createResult.ID), nil)

	if delRes.StatusCode != 200 {
		t.Fatalf("Received non-200 response: %d\n", delRes.StatusCode)
	}

	var nextCount int
	app.DB.Table("entries").Where("deleted_at", "NULL").Count(&nextCount)
	assert.Equal(t, nextCount, 1)
}
