package main

import (
	"time"

	"github.com/gophers-team/gopher-box/api"
	"github.com/jmoiron/sqlx"
)

var saveHeartbeatQuery = `
INSERT INTO heartbeats (device_id, created_at)
VALUES ($1, $2)`

func saveHeartbeat(db *sqlx.DB, deviceID api.DeviceID) error {
	tx := db.MustBegin()
	tx.MustExec(saveHeartbeatQuery, deviceID, time.Now())
	tx.Commit()
	return nil
}

var insertDeviceDispensingQuery = `
INSERT INTO device_dispensings (device_id, status, created_at)
VALUES ($1, $2, $3)`

func dispensingBegin(db *sqlx.DB, deviceID api.DeviceID) (operationID int64, err error) {
	tx := db.MustBegin()
	res := tx.MustExec(insertDeviceDispensingQuery, deviceID, "begin", time.Now())
	tx.Commit()

	return res.LastInsertId()
}
