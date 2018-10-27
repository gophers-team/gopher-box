package main

import (
	"log"
	"time"

	"github.com/gophers-team/gopher-box/api"
	"github.com/jmoiron/sqlx"
)

func saveHeartbeat(db *sqlx.DB, deviceID api.DeviceID) error {
	tx := db.MustBegin()
	tx.MustExec(`
		INSERT INTO heartbeats (device_id, created_at)
		VALUES ($1, $2)
	`, deviceID, time.Now())
	return tx.Commit()
}

func getPills(db *sqlx.DB, deviceID api.DeviceID) (pills map[api.TabletID]api.TabletAmount, err error) {
	pills = make(map[api.TabletID]api.TabletAmount)
	tx := db.MustBegin()
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
	res := tx.MustExec(`
		INSERT INTO device_dispensings (device_id, status, created_at)
		VALUES ($1, $2, $3)
	`, deviceID, "begin", time.Now())
	tx.Commit()
	operationID, err = res.LastInsertId()

	pills, err = getPills(db, deviceID)
	return operationID, pills, err
}

func getDeviceInfos(db *sqlx.DB) []api.DeviceInfo {
	var infos []api.DeviceInfo
	tx := db.MustBegin()
	rows, err := tx.Queryx(
		`SELECT device_id, MAX(created_at) FROM heartbeats GROUP BY device_id;`,
	)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var deviceInfo api.DeviceInfo
		var createdAt []uint8
		err := rows.Scan(&deviceInfo.DeviceID, &createdAt)
		if err != nil {
			log.Fatal(err)
		}
		deviceInfo.Status = api.DeviceStatusOnline
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
