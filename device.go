package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	resp, err := http.Get("http://130.193.56.206/dbtest")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Response: %v", resp)
}
