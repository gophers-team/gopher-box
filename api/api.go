package api

type DeviceID    string
type EventType   string
type OperationID string

var DeviceHeartbeatEndpoint = "/heartbeat"

type DeviceHeartbeat struct {
	DeviceID DeviceID `json:"device_id"`
}

var DeviceStatusEndpoint = "/status"

type DeviceTabletStatusRequest struct {
	DeviceID DeviceID `json:"device_id"`
}

type TabletID string
type TabletAmount uint8

type DeviceTabletStatusResponse struct {
	Tablets map[TabletID]TabletAmount `json:"tablets"`
	OperationID OperationID           `json:"operation_id"`
}

var DeviceDispenseEndpoint = "/dispense"

type DeviceTabletDispenseRequest struct {
	DeviceID    DeviceID    `json:"device_id"`
	Fulfillment map[TabletID]TabletAmount
	OperationID OperationID `json:"operation_id"`
}
