package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gophers-team/gopher-box/api"
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

var heartbeatQuery = `
INSERT INTO heartbeats (
	device_id,
	created_at
)
VALUES ($1, $2)`

func heartbeatHandler(db *sqlx.DB, w http.ResponseWriter, r *http.Request) {
	var h api.DeviceHeartbeat
	if err := json.NewDecoder(r.Body).Decode(&h); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to unmarshal /heartbeat request"))
		return
	}

	tx := db.MustBegin()
	tx.MustExec(heartbeatQuery, h.DeviceID, time.Now())
	tx.Commit()

	w.WriteHeader(http.StatusOK)
}

func dispenseHandler(db *sqlx.DB, w http.ResponseWriter, r *http.Request) {
	var t api.DeviceTabletDispenseRequest
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to unmarshal /dispense request"))
		return
	}

	w.WriteHeader(http.StatusOK)
}

func statusHandler(db *sqlx.DB, w http.ResponseWriter, r *http.Request) {
	var s api.DeviceTabletStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to unmarshal /status request"))
		return
	}

	fmt.Println(s)
	w.WriteHeader(http.StatusOK)
}
