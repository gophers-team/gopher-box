package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gophers-team/gopher-box/api"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

func DbHandler(db *sqlx.DB, handler func(db *sqlx.DB, w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		handler(db, w, r)
	}
	return http.HandlerFunc(fn)
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("GopherBox!\n"))
}

func EventsHandler(db *sqlx.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		tx := db.MustBegin()
		tx.MustExec("INSERT INTO events (type, created_at) VALUES ($1, $2)", "Heartbeat", time.Now())
		tx.Commit()
	}
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "VARS: %v\n", vars)
}

func deviceDispenseHandler(db *sqlx.DB, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Failed to read device /dispense request")
		return
	}
	var t api.DeviceTabletDispenseRequest
	err = json.Unmarshal(data, &t)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Failed to read device /dispense request"))
	}
}
