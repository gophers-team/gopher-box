package main

import (
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var DbFile = "/var/lib/gopher-box/db"

func DbHandler(db *bolt.DB, handler func(db *bolt.DB, w http.ResponseWriter, r * http.Request)) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r * http.Request) {
		handler(db, w, r)
	}
	return http.HandlerFunc(fn)
}

func InitDb() (*bolt.DB, error) {
	db, err := bolt.Open(DbFile, 0600, nil)
	if err != nil {
		return nil, err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("DB"))
		return err
	})

	return db, err
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("GopherBox!\n"))
}

func EventsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "VARS: %v\n", vars)
}

func DbTestHandler(db *bolt.DB, w http.ResponseWriter, r * http.Request) {
	db.Update(func(tx *bolt.Tx) error {
		tx.Bucket([]byte("DB")).Put([]byte("answer"), []byte("42\n"))
		return nil
	})

	var v []byte
	db.View(func(tx *bolt.Tx) error {
		v = tx.Bucket([]byte("DB")).Get([]byte("answer"))
		return nil
	})
	w.Write(v)
}

func main() {
	db, err := InitDb()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := mux.NewRouter()
	r.HandleFunc("/", IndexHandler)
	r.HandleFunc("/events", EventsHandler).Methods("GET", "POST")
	r.HandleFunc("/dbtest", DbHandler(db, DbTestHandler))

	log.Fatal(http.ListenAndServe(":80", r))
}
