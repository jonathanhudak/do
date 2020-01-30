package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

type App struct {
	Router *mux.Router
	DB     *gorm.DB
}

type Entry struct {
	gorm.Model
	Title string `json:title`
}

func (a *App) getEntries(w http.ResponseWriter, req *http.Request) {
	var entry []Entry
	a.DB.Find(&entry)
	js, err := json.Marshal(entry)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (a *App) getEntry(w http.ResponseWriter, req *http.Request) {
	var entry Entry
	vars := mux.Vars(req)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	a.DB.Find(&entry, id)
	js, err := json.Marshal(entry)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (a *App) createEntry(w http.ResponseWriter, req *http.Request) {
	var entry Entry
	db := a.DB
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
}

func (a *App) updateEntry(w http.ResponseWriter, req *http.Request) {
	db := a.DB
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
}

func (a *App) deleteEntry(w http.ResponseWriter, req *http.Request) {
	db := a.DB
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
}
