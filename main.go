package main

import (
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var DbFile = "/var/lib/gopher-box/db"

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("GopherBox!\n"))
}

func EventsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "VARS: %v\n", vars)
}

func main() {
	db, err := bolt.Open(DbFile, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := mux.NewRouter()
	r.HandleFunc("/", IndexHandler)
	r.HandleFunc("/events", EventsHandler).Methods("GET", "POST")

	log.Fatal(http.ListenAndServe(":8000", r))
}
