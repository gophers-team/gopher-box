package main

import (
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

func getPills(db *sqlx.DB, deviceID api.DeviceID) (pills map[api.TabletID]api.TabletAmount, schedule_id int, err error) {
	pills = make(map[api.TabletID]api.TabletAmount)
	tx := db.MustBegin()
	defer tx.Commit()

	rows, err := tx.Queryx(`
		SELECT p.name, ds.id, ds.amount, ds.interval, COALESCE(MAX(dd.changed_at), '2018-10-10 01:16:30.3404') FROM dispensing_schedule AS ds
		INNER JOIN dispensing_plans AS dp ON ds.plan_id = dp.id
		INNER JOIN devices AS d ON d.plan_id = dp.id
		INNER JOIN pills AS p ON p.id = ds.pill_id
		LEFT JOIN device_dispensings as dd ON ds.id = dd.schedule_id
		WHERE d.id = $1 AND (dd.status = $2 OR dd.status IS NULL)
		GROUP BY p.name, ds.id, ds.amount, ds.interval;
	`, deviceID, DispensingStatusFinished)
	if err != nil {
		log.Fatal(err)
	}

	// var schedule_id int // :DDDDDDDD
	var name string
	var interval int
	var amount   int
	var changedAt time.Time
	// :DDDDDDDDD
	for rows.Next() {
		err := rows.Scan(&name, &schedule_id, &amount, &interval, &changedAt)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("pills", name, amount, changedAt)
	}
	log.Println("diff", time.Since(changedAt).Minutes())
	if time.Since(changedAt).Minutes() < float64(interval) {
		amount = 0
	}
	pills[api.TabletID(name)] = api.TabletAmount(amount)
	return pills, schedule_id, nil
}

func dispensingBegin(db *sqlx.DB, deviceID api.DeviceID) (operationID int64, pills map[api.TabletID]api.TabletAmount, err error) {
	tx := db.MustBegin()
	defer tx.Commit()

	pills, scheduleId, err := getPills(db, deviceID)

	now := time.Now()
	row := tx.QueryRow(`
		INSERT INTO device_dispensings (device_id, status, created_at, changed_at, schedule_id)
		VALUES ($1, $2, $3, $3, $4)
        RETURNING id
	`, deviceID, DispensingStatusBegin, now, scheduleId)

	err = row.Scan(&operationID)
	if err != nil {
		log.Fatal(err)
	}

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
		GROUP BY h.device_id, d.name
		ORDER BY h.device_id`,
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
		deviceInfo.Status = api.DeviceStatusOnline
		if time.Since(createdAt) >= 5 * time.Second{
			deviceInfo.Status = api.DeviceStatusOffline
		}
		deviceInfo.Info = "info"
		infos = append(infos, deviceInfo)
	}
	return infos
}

func dispensingEnd(db *sqlx.DB, operationID api.OperationID, status DispensingStatus) (err error) {
	tx := db.MustBegin()

	tx.MustExec(`
		UPDATE device_dispensings
		SET status = $1, changed_at = $2
		WHERE id = $3
	`, status, time.Now(), operationID)
	err = tx.Commit()

	return err
}
