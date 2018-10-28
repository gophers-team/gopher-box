package api

type DeviceID int
type EventType string
type OperationID int64
type DeviceStatus string

const (
	DeviceStatusUnspecified DeviceStatus = "unspecified"
	DeviceStatusOnline                   = "online"
	DeviceStatusOffline                  = "offline"
)

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
	Tablets     map[TabletID]TabletAmount `json:"tablets"`
	OperationID OperationID               `json:"operation_id"`
}

var DeviceDispenseEndpoint = "/dispense"

type DeviceTabletDispenseRequest struct {
	DeviceID    DeviceID `json:"device_id"`
	Fulfillment map[TabletID]TabletAmount
	OperationID OperationID `json:"operation_id"`
}

var DeviceEndpoint = "/device"

type DeviceInfo struct {
	DeviceID DeviceID     `json:"device_id"`
	Name     string       `json:"name"`
	Status   DeviceStatus `json:"status"`
	Info     string       `json:"info"`
}

type DeviceListResponse []DeviceInfo
