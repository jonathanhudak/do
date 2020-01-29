package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type App struct {
	Router *mux.Router
	DB     *gorm.DB
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

func (a *App) Initialize(driverType, connectionStr string) {
	var err error
	a.DB, err = gorm.Open(driverType, connectionStr)
	if err != nil {
		fmt.Print(err)
		panic("failed to connect database")
	}

	a.DB.AutoMigrate(&Entry{})

	r := mux.NewRouter()

	r.HandleFunc("/api/entries", a.getEntries).Methods(http.MethodGet)

	a.Router = r
}

func (a *App) Run(addr string) {
	srv := &http.Server{
		Handler:      a.Router,
		Addr:         fmt.Sprintf("127.0.0.1%s", addr),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
