package main

import (
	"encoding/json"
	"log"
	"net/http"

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

func heartbeatHandler(db *sqlx.DB, w http.ResponseWriter, r *http.Request) {
	var h api.DeviceHeartbeat
	if err := json.NewDecoder(r.Body).Decode(&h); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to unmarshal /heartbeat request"))
		return
	}
	saveHeartbeat(db, h.DeviceID)
	w.WriteHeader(http.StatusOK)
}

func dispenseHandler(db *sqlx.DB, w http.ResponseWriter, r *http.Request) {
	var t api.DeviceTabletDispenseRequest
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to unmarshal /dispense request"))
		return
	}
	err := dispensingEnd(db, t.OperationID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	return
}

func statusHandler(db *sqlx.DB, w http.ResponseWriter, r *http.Request) {
	var s api.DeviceTabletStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to unmarshal /status request"))
		return
	}
	operationID, pills, err := dispensingBegin(db, s.DeviceID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed begin dispensing"))
		return
	}

	resp := api.DeviceTabletStatusResponse{
		OperationID: api.OperationID(operationID),
		Tablets:     pills,
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&resp); err != nil {
		log.Panic(err)
	}
}

func deviceListHandler(db *sqlx.DB, w http.ResponseWriter, r *http.Request) {
	infos := getDeviceInfos(db)
	resp := api.DeviceListResponse(infos)
	if err := json.NewEncoder(w).Encode(&resp); err != nil {
		log.Panic(err)
	}
}
