package main

import (
	"log"
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

var selectDispensingSchedule = `
SELECT p.name, ds.amount FROM dispensing_schedule AS ds
INNER JOIN dispensing_plans AS dp ON ds.plan_id = dp.id
INNER JOIN devices AS d ON d.plan_id = dp.id
INNER JOIN pills AS p ON p.id = ds.pill_id
WHERE d.id = $1
`

func dispensingBegin(db *sqlx.DB, deviceID api.DeviceID) (operationID int64, pills map[api.TabletID]api.TabletAmount, err error) {
	tx := db.MustBegin()
	res := tx.MustExec(insertDeviceDispensingQuery, deviceID, "begin", time.Now())
	tx.Commit()
	operationID, err = res.LastInsertId()

	pills = make(map[api.TabletID]api.TabletAmount)
	tx = db.MustBegin()
	rows, err := tx.Queryx(selectDispensingSchedule, deviceID)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var name string
		var amount int
		err := rows.Scan(&name, &amount)
		if err != nil {
			log.Fatal(err)
		}
		pills[api.TabletID(name)] = api.TabletAmount(amount)
	}

	return operationID, pills, err
}
