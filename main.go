package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Entry struct {
	gorm.Model
	Title string
}

func main() {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	// Migrate the schema
	db.AutoMigrate(&Entry{})

	r := mux.NewRouter()

	// List Entries
	r.HandleFunc("/api/entries", func(w http.ResponseWriter, req *http.Request) {
		var entry []Entry
		db.Find(&entry)
		js, err := json.Marshal(entry)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}).Methods(http.MethodGet)

	// Get Entry
	r.HandleFunc("/api/entry/{id}", func(w http.ResponseWriter, req *http.Request) {
		var entry Entry
		vars := mux.Vars(req)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}
		db.Find(&entry, id)
		js, err := json.Marshal(entry)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}).Methods(http.MethodGet)

	// Create Entry
	r.HandleFunc("/api/entry", func(w http.ResponseWriter, req *http.Request) {
		var entry Entry
		err := json.NewDecoder(req.Body).Decode(&entry)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}
		db.Create(&entry)
		db.Take(&entry)
		js, err := json.Marshal(entry)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}).Methods(http.MethodPost)

	// Update Entry
	r.HandleFunc("/api/entry/{id}", func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}
		var entry, entryUpdate Entry
		db.Find(&entry, id)

		err = json.NewDecoder(req.Body).Decode(&entryUpdate)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}
		db.Model(&entry).Update("title", entryUpdate.Title)
		db.Take(&entry)
		js, err := json.Marshal(entry)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}).Methods(http.MethodPut)

	// Delete Entry
	r.HandleFunc("/api/entry/{id}", func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}
		var entry Entry
		db.Find(&entry, id)
		db.Delete(&entry)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}).Methods(http.MethodDelete)

	http.Handle("/", r)

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
