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
	"strings"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/firmata"
)

type requester interface {
	PostJson(endpoint string, val interface{}) (*http.Response, error)
}

type tabletDispenser struct {
	motor    *gpio.StepperDriver
	tabletID api.TabletID
	rpm      uint
	step     int
}

func newTabletDispenser(motor *gpio.StepperDriver, tabletID api.TabletID, rpm uint, step int) *tabletDispenser {
	return &tabletDispenser{
		motor:    motor,
		tabletID: tabletID,
		rpm:      rpm,
		step:     step,
	}
}

// TODO: will have to add time and stop event for this and Close
func (t *tabletDispenser) Rotate() {
	log.Printf("started opening tablet dispenser for %s", t.tabletID)

	t.motor.SetSpeed(t.rpm)
	err := t.motor.Move(t.step)
	if err != nil {
		log.Fatalf("dispenser for tablet %s move error: %v", t.tabletID, err)
	}
	time.Sleep(time.Second)
	t.motor.SetSpeed(0)
	log.Printf("finished opening tablet dispenser for %s", t.tabletID)
}

func (t *tabletDispenser) DbgRotate() {
	t.motor.SetSpeed(t.rpm)
	for {
		err := t.motor.Move(t.step)
		if err != nil {
			log.Fatalf("dispenser for tablet %s move error: %v", t.tabletID, err)
		}
	}
}

type requestData struct {
	requester        requester
	deviceID         api.DeviceID
	tabletDispensers map[api.TabletID]*tabletDispenser
}

func main() {
	tty := flag.String("tty", "/dev/ttyMSM1", "arduino tty")
	in1 := flag.String("in1", "8", "in1")
	in2 := flag.String("in2", "9", "in2")
	in3 := flag.String("in3", "10", "in3")
	in4 := flag.String("in4", "11", "in4")
	tabletButtonPin := flag.String("tablet-button-pin", "12", `pin of "give me tablets!" button`)
	tabletButtonPollInterval := flag.Duration("tablet-button-poll-interval", 10*time.Millisecond, `poll interval of "give me tablets!" button`)
	stepsPerRev := flag.Uint("steps-per-rev", 2038, "steps per rev")
	step := flag.Int("step", 2038, "step")
	rpm := flag.Uint("rpm", 10, "rpm speed")
	heartbeetInterval := flag.Duration("heartbeat interval", 10*time.Second, "interval between heartbeats")
	server := flag.String("server", "130.193.56.206", "address of server to send data to")
	deviceID := flag.Int("device-id", 1337, "the (unique) id of the device")
	tabletID := flag.String("tablet-id", "0", "tablet id (type of tablets)")
	debugButton := flag.Bool("debug-button", false, "debug button events")
	flag.Parse()

	pins := [...]string{*in1, *in2, *in3, *in4}

	firmataAdaptor := firmata.NewAdaptor(*tty)
	motor := gpio.NewStepperDriver(firmataAdaptor, pins, gpio.StepperModes.SinglePhaseStepping, *stepsPerRev)
	tabletButton := gpio.NewButtonDriver(firmataAdaptor, *tabletButtonPin, *tabletButtonPollInterval)

	requester := httpRequester{
		server: *server,
	}
	tid := api.TabletID(*tabletID)
	rd := &requestData{
		requester: &requester,
		deviceID:  api.DeviceID(*deviceID),
		tabletDispensers: map[api.TabletID]*tabletDispenser{
			tid: newTabletDispenser(motor, tid, *rpm, *step),
		},
	}

	log.Printf("ready to start work: %+v", *rd)

	work := func() {
		go heartbeat(rd, *heartbeetInterval)
		go rd.tabletDispensers[tid].DbgRotate()

		err := tabletButton.Start()
		if err != nil {
			log.Fatalf("failed to start tablet button: %v", err)
		}

		tabletButtonEvents := tabletButton.Subscribe()
		lastButtonEventTime := time.Now()
		isDoubleEvent := func() bool {
			t := time.Now()
			ret := t.Sub(lastButtonEventTime) < 100*time.Millisecond
			lastButtonEventTime = t
			return ret
		}

		for event := range tabletButtonEvents {
			switch event.Name {
			case gpio.ButtonPush: // skipping, acting on push
				if isDoubleEvent() {
					continue
				}
				if *debugButton {
					log.Printf("button push event: %+v", event)
					continue
				}
				err = tabletButtonPush(rd)
				if err != nil {
					log.Fatalf("error processing button push: %v", err)
				}

			case gpio.ButtonRelease:
				if isDoubleEvent() {
					continue
				}
				if *debugButton {
					log.Printf("button release event: %+v", event)
					continue
				}
			case gpio.Error:
				if *debugButton {
					log.Printf("button error event: %+v", event)
					continue
				}
				err = event.Data.(error)
				log.Fatalf("error event from button: %v (%+v)", err, event)
			default:
				log.Fatalf("got unexpected event from button: %+v", event)
			}
		}

		log.Println("worker loop finished")
	}

	robot := gobot.NewRobot("bot",
		[]gobot.Connection{firmataAdaptor},
		[]gobot.Device{motor, tabletButton},
		work,
	)

	robot.Start()
}

func heartbeat(rd *requestData, interval time.Duration) {
	heartbeat := api.DeviceHeartbeat{
		DeviceID: rd.deviceID,
	}
	t := time.NewTimer(0)
	defer t.Stop()
	for {
		<-t.C
		_, _ = rd.requester.PostJson(api.DeviceHeartbeatEndpoint, &heartbeat)
		t.Reset(interval)
	}
}

func tabletButtonPush(rd *requestData) error {
	s, err := status(rd)
	if err != nil {
		// TODO: it'll be nice to notify user that the server is down
		return err
	}

	resp := api.DeviceTabletDispenseRequest{
		DeviceID:    rd.deviceID,
		Fulfillment: make(map[api.TabletID]api.TabletAmount, len(s.Tablets)),
		OperationID: s.OperationID,
	}

	for t, amount := range s.Tablets {
		res := api.TabletAmount(0)
		if amount != 0 {
			res, err = dispenseTablet(rd, t, amount)
		}
		resp.Fulfillment[t] = res
	}

	return nil
}

func dispenseTablet(rd *requestData, tabletID api.TabletID, amount api.TabletAmount) (api.TabletAmount, error) {
	dispenser, ok := rd.tabletDispensers[tabletID]
	if !ok {
		return 0, fmt.Errorf("dispenser for tablet id %s not found", tabletID)
	}

	for i := api.TabletAmount(0); i < amount; i++ {
		dispenser.Rotate()
	}

	return amount, nil
}

func status(rd *requestData) (*api.DeviceTabletStatusResponse, error) {
	request := api.DeviceTabletStatusRequest{
		DeviceID: rd.deviceID,
	}
	resp, err := rd.requester.PostJson(api.DeviceStatusEndpoint, &request)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		errText := fmt.Sprintf("error reading status response: %v", err)
		log.Println(errText)
		return nil, errors.New(errText)
	}

	var status api.DeviceTabletStatusResponse
	err = json.Unmarshal(data, &status)
	if err != nil {
		return nil, fmt.Errorf("error parsing json status response: %v", err)
	}
	return &status, nil
}

type httpRequester struct {
	server string
}

func (h *httpRequester) PostJson(endpoint string, val interface{}) (*http.Response, error) {
	data, err := json.Marshal(val)
	if err != nil {
		return nil, fmt.Errorf("marshaling json object for %s failed: %v", endpoint, err)
	}
	url := fmt.Sprintf("http://%s%s", h.server, endpoint)
	body := bytes.NewBuffer(data)
	log.Printf("sending request to %s", url)
	resp, err := http.Post(url, "application/json", body)
	if err != nil {
		log.Printf("request to %s failed: %v", endpoint, err)
		return nil, fmt.Errorf("request to %s failed: %v", url, err)
	}

	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		data, readErr := ioutil.ReadAll(resp.Body)
		if readErr != nil {
			data = []byte{}
		}
		text := strings.TrimSpace(string(data))
		errText := fmt.Sprintf("request to %s failed with %d status code: %s", url, resp.StatusCode, text)
		log.Println(errText)
		return nil, errors.New(errText)
	}

	log.Printf("request to %s succeed", url)

	return resp, nil
}
