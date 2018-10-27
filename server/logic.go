package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gophers-team/gopher-box/api"
	"github.com/jmoiron/sqlx"
)

func saveHeartbeat(db *sqlx.DB, deviceID api.DeviceID) {
	tx := db.MustBegin()
	defer tx.Commit()
	tx.MustExec(`
		INSERT INTO heartbeats (device_id, created_at)
		VALUES ($1, $2)
	`, deviceID, time.Now())
}

func getPills(db *sqlx.DB, deviceID api.DeviceID) (pills map[api.TabletID]api.TabletAmount, err error) {
	pills = make(map[api.TabletID]api.TabletAmount)
	tx := db.MustBegin()
	defer tx.Commit()

	rows, err := tx.Queryx(`
		SELECT p.name, ds.amount FROM dispensing_schedule AS ds
		INNER JOIN dispensing_plans AS dp ON ds.plan_id = dp.id
		INNER JOIN devices AS d ON d.plan_id = dp.id
		INNER JOIN pills AS p ON p.id = ds.pill_id
		WHERE d.id = $1
	`, deviceID)
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
	return pills, nil
}

func dispensingBegin(db *sqlx.DB, deviceID api.DeviceID) (operationID int64, pills map[api.TabletID]api.TabletAmount, err error) {
	tx := db.MustBegin()
	defer tx.Commit()

	res := tx.MustExec(`
		INSERT INTO device_dispensings (device_id, status, created_at)
		VALUES ($1, $2, $3)
	`, deviceID, "begin", time.Now())
	operationID, err = res.LastInsertId()

	pills, err = getPills(db, deviceID)
	return operationID, pills, err
}

func getDeviceInfos(db *sqlx.DB) []api.DeviceInfo {
	var infos []api.DeviceInfo
	tx := db.MustBegin()
	defer tx.Commit()

	rows, err := tx.Queryx(`
		SELECT h.device_id, d.name, MAX(h.created_at)
		FROM heartbeats as h
		JOIN devices as d ON h.device_id = d.id
		GROUP BY h.device_id, d.name;`,
	)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var deviceInfo api.DeviceInfo
		var createdAt time.Time
		err := rows.Scan(&deviceInfo.DeviceID, &deviceInfo.Name, &createdAt)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(createdAt)
		deviceInfo.Status = api.DeviceStatusOnline
		if time.Since(createdAt) >= 5 * time.Second{
			deviceInfo.Status = api.DeviceStatusOffline
		}
		infos = append(infos, deviceInfo)
	}
	return infos
}

func dispensingEnd(db *sqlx.DB, operationID api.OperationID) (err error) {
	tx := db.MustBegin()

	tx.MustExec(`
		UPDATE device_dispensings
		SET status = $1
		WHERE operation_id = $3
	`, operationID, "finished")
	err = tx.Commit()

	return err
}
