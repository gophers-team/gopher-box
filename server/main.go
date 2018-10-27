package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func DbHandler(db *sqlx.DB, handler func(db *sqlx.DB, w http.ResponseWriter, r * http.Request)) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r * http.Request) {
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

func main() {
	db, err := InitDb()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := mux.NewRouter()
	r.HandleFunc("/", IndexHandler)
	r.HandleFunc("/events", DbHandler(db, EventsHandler)).Methods(http.MethodGet, http.MethodPost)

	log.Fatal(http.ListenAndServe(":8000", r))
}
