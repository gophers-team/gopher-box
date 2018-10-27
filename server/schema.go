package main

var schema = `
DROP TABLE IF EXISTS devices;
DROP TABLE IF EXISTS device_dispensings;
DROP TABLE IF EXISTS pills;
DROP TABLE IF EXISTS dispensing_plans;
DROP TABLE IF EXISTS dispensing_schedule;

CREATE TABLE devices (
    id          INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
	plan_id    	INTEGER NOT NULL,
	name 		VARCHAR(255) NOT NULL,
    created_at  DATETIME,
	FOREIGN KEY (plan_id) REFERENCES plans(id)
);

CREATE TABLE device_dispensings (
	id          	INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    device_id     	INTEGER NOT NULL,
    operation_id	INTEGER NOT NULL,
	schedule_id		INTEGER,
	pills_dispensed	INTEGER,
    created_at  	DATETIME,
	FOREIGN KEY (device_id) REFERENCES devices(id)
	FOREIGN KEY (schedule_id) REFERENCES dispensing_schedule(id)
);

CREATE TABLE pills (
    id          INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
	name 		VARCHAR(255) NOT NULL,
    created_at  DATETIME
);

CREATE TABLE dispensing_plans (
    id          INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
	name 		VARCHAR(255) NOT NULL,
    created_at  DATETIME
);

CREATE TABLE dispensing_schedule (
	id          		INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    plan_id     		INTEGER NOT NULL,
	pill_id				INTEGER NOT NULL,
	amount				INTEGER NOT NULL,
	dispense_dow		INTEGER NOT NULL,
	dispense_time   	DATETIME NOT NULL,
	schedule_duration 	INTEGER NOT NULL,
	created_at  		DATETIME,
	FOREIGN KEY (plan_id) REFERENCES plans(id)
	FOREIGN KEY (pill_id) REFERENCES pills(id)
);

`
