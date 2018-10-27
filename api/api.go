package api

var DeviceHeartbeatEndpoint = "/heartbeat"

type DeviceHeartbeat struct {
	DeviceID string `json:"device_id"`
}

var DeviceStatusEndpoint = "/status"

type DeviceTabletStatusRequest struct {
	DeviceID string `json:"device_id"`
}

type TabletID string
type TableAmount uint8

type DeviceTabletStatusResponse struct {
	Tablets map[TabletID]TableAmount `json:"tablets"`
}

var DeviceDispenseEndpoint = "/dispense"

type DeviceTabletDispenseRequest struct {
	DeviceID string `json:"device_id"`
	Fullfilled bool `json:"fullfilled"`
}