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
	static := flag.String("static", "static", "static files dir")
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
		plans := []DispensingPlan{}
		err = db.Select(&plans, "SELECT * FROM dispensing_plans")
		if err != nil {
			log.Fatal(err)
		}
		for _, plan := range plans {
			fmt.Println("%v", plan)

			schedules := []DispensingSchedule{}
			err = db.Select(
				&schedules, "SELECT * FROM dispensing_schedule WHERE plan_id=$1 ORDER BY dispense_dow", plan.Id,
			)
			if err != nil {
				log.Fatal(err)
			}
			for _, schedule := range schedules {
				fmt.Println("%v", schedule)
			}
		}
	}
	if *server {
		r := mux.NewRouter()

		r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir(*static))))
		r.HandleFunc(api.DeviceDispenseEndpoint, DbHandler(db, dispenseHandler)).Methods(http.MethodPost)
		r.HandleFunc(api.DeviceStatusEndpoint, DbHandler(db, statusHandler)).Methods(http.MethodPost)
		r.HandleFunc(api.DeviceHeartbeatEndpoint, DbHandler(db, heartbeatHandler)).Methods(http.MethodPost)

		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), r))
	}
}
