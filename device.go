package main

import (
	"fmt"
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
			Action: func(c *cli.Context) error {
				resp, err := http.Get("http://130.193.56.206/dbtest")
				if err != nil {
					log.Fatal(err)
				}
				fmt.Printf("Response: %v", resp)
				return nil
			},
		},
		{
			Name:    "start",
			Aliases: []string{"s"},
			Usage:   "start device",
			Action: func(c *cli.Context) error {
				fmt.Println("Device is starting...")
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
