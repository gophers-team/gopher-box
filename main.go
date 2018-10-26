package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("GopherBox!\n"))
}

func EventsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "VARS: %v\n", vars)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", IndexHandler)
	r.HandleFunc("/events", EventsHandler).Methods("GET", "POST")

	log.Fatal(http.ListenAndServe(":8000", r))
}
