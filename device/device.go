package main

import (
	"log"
	"net/http"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Commands = []cli.Command{
		{
			Name:    "event",
			Aliases: []string{"e"},
			Usage:   "send event",
			Action: sendEvent(),
		},
		{
			Name:    "start",
			Aliases: []string{"s"},
			Usage:   "start device",
			Action: func(c *cli.Context) error {
				log.Println("Device is starting...")
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func sendEvent() error {
	resp, err := http.Get("http://130.193.56.206/dbtest")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Response: %v", resp)
	return nil
}
