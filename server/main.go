package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gophers-team/gopher-box/api"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	initDb := flag.Bool("init-db", false, "initialize database")
	dbFile := flag.String("db-file", "/var/lib/gopher-box/events.db", "database file path")

	server := flag.Bool("server", false, "start server")
	port := flag.Int("port", 80, "server's port")

	showPlans := flag.Bool("show-plans", false, "show plans")

	flag.Parse()

	db, err := InitDb(*dbFile)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if *initDb {
		db.MustExec(schema)
	}
	if *showPlans {
		plans := []Plan{}
		err = db.Select(&plans, "SELECT * FROM dispensing_plans")
		if err != nil {
			log.Fatal(err)
		}
		for _, plan := range plans {
			fmt.Println("%v", plan)
		}
	}
	if *server {
		r := mux.NewRouter()
		r.HandleFunc("/", IndexHandler)
		r.HandleFunc("/events", DbHandler(db, EventsHandler)).Methods(http.MethodGet, http.MethodPost)
		r.HandleFunc(api.DeviceDispenseEndpoint, DbHandler(db, deviceDispenseHandler)).Methods(http.MethodPost)

		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), r))
	}
}
