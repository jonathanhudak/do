package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

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
	r.HandleFunc("/api/entry/{id}", a.getEntry).Methods(http.MethodGet)
	r.HandleFunc("/api/entry", a.createEntry).Methods(http.MethodPost)
	r.HandleFunc("/api/entry/{id}", a.updateEntry).Methods(http.MethodPut)
	r.HandleFunc("/api/entry/{id}", a.deleteEntry).Methods(http.MethodDelete)

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
