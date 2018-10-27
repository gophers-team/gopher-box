package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/firmata"
)

func main() {
	tty := flag.String("tty", "/dev/ttyMSM1", "arduino tty")
	in1 := flag.String("in1", "8", "in1")
	in2 := flag.String("in2", "9", "in2")
	in3 := flag.String("in3", "10", "in3")
	in4 := flag.String("in4", "11", "in4")
	stepsPerRev := flag.Uint("steps-per-rev", 2038, "steps per rev")
	step := flag.Int("step", 2038, "step")
	rpm := flag.Uint("rpm", 10, "rpm speed")
	revRpm := flag.Uint("rev-rpm", 10, "rev rpm speed")
	flag.Parse()

	pins := [...]string{*in1, *in2, *in3, *in4}

	firmataAdaptor := firmata.NewAdaptor(*tty)
	motor := gpio.NewStepperDriver(firmataAdaptor, pins, gpio.StepperModes.SinglePhaseStepping, *stepsPerRev)

	work := func() {
		for {
			motor.SetSpeed(*rpm)
			err := motor.Move(*step)
			if err != nil {
				log.Println("move error: %v", err)
			}
			go sendEvent()
			time.Sleep(time.Second)
			motor.SetSpeed(*revRpm)
			motor.Move(-*step)
		}
	}

	robot := gobot.NewRobot("bot",
		[]gobot.Connection{firmataAdaptor},
		[]gobot.Device{motor},
		work,
	)

	robot.Start()
}

func sendEvent() error {
	resp, err := http.Get("http://130.193.56.206/dbtest")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Response: %v", resp)
	return nil
}
