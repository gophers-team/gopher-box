package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/gophers-team/gopher-box/api"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/firmata"
)

type requester interface {
	PostJson(endpoint string, val interface{}) (*http.Response, error)
}

type requestData struct {
	requester requester
	deviceID  api.DeviceID
}

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
	heartbeetInterval := flag.Duration("heartbeet interval", 10*time.Second, "interval between heartbeets")
	server := flag.String("server", "130.193.56.206", "address of server to send data to")
	deviceID := flag.String("device-id", "gophers-device-1337", "the (unique) name of the device")
	flag.Parse()

	pins := [...]string{*in1, *in2, *in3, *in4}

	firmataAdaptor := firmata.NewAdaptor(*tty)
	motor := gpio.NewStepperDriver(firmataAdaptor, pins, gpio.StepperModes.SinglePhaseStepping, *stepsPerRev)

	requester := httpRequester{
		server: *server,
	}
	rd := &requestData{
		requester: &requester,
		deviceID: api.DeviceID(*deviceID),
	}

	work := func() {
		go heartbeet(rd, *heartbeetInterval)
		for {
			motor.SetSpeed(*rpm)
			err := motor.Move(*step)
			if err != nil {
				log.Printf("move error: %v", err)
			}
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

func heartbeet(rd *requestData, interval time.Duration) {
	heartbeat := api.DeviceHeartbeat{
		DeviceID: rd.deviceID,
	}
	t := time.NewTimer(0)
	defer t.Stop()
	for {
		<-t.C
		_, err := rd.requester.PostJson(api.DeviceHeartbeatEndpoint, &heartbeat)
		if err != nil {
			continue
		}
		t.Reset(interval)
	}
}

type httpRequester struct {
	server string
}

func (h *httpRequester) PostJson(endpoint string, val interface{}) (*http.Response, error) {
	data, err := json.Marshal(val)
	if err != nil {
		return nil, fmt.Errorf("marshaling json object for %s failed: %v", endpoint, err)
	}
	url := fmt.Sprintf("http://%s/%s", h.server, endpoint)
	body := bytes.NewBuffer(data)
	log.Printf("sending request to %s", url)
	resp, err := http.Post(url, "application/json", body)
	if err != nil {
		log.Printf("request to %s failed: %v", endpoint, err)
		return nil, fmt.Errorf("request to %s failed: %v", url, err)
	}

	if resp.StatusCode != 200 {
		data, readErr := ioutil.ReadAll(resp.Body)
		if readErr != nil {
			data = []byte{}
		}
		resp.Body.Close()
		text := string(data)
		errText := fmt.Sprintf("request to %s failed with %d status code: %s", url, resp.StatusCode, text)
		log.Println(errText)
		return nil, errors.New(errText)
	}

	log.Println("request to %s succeed", url)

	return resp, nil
}
