package main

var schema = `
DROP TABLE IF EXISTS events;
DROP TABLE IF EXISTS pills;
DROP TABLE IF EXISTS dispensing_plans;
DROP TABLE IF EXISTS dispense_schedule;
DROP TABLE IF EXISTS operations;

CREATE TABLE devices (
    id          INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
	name 		VARCHAR(255) NOT NULL,
    created_at  DATETIME
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

CREATE TABLE dispense_schedule (
    id          	INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    plan_id     	INTEGER NOT NULL,
	pill_id			INTEGER NOT NULL,
	amount			INTEGER NOT NULL,
	dispence_dow	INTEGER NOT NULL,
	dispence_time   DATETIME NOT NULL,
	created_at  	DATETIME,
	FOREIGN KEY (plan_id) REFERENCES plans(id)
	FOREIGN KEY (pill_id) REFERENCES pills(id)
);

CREATE TABLE events (
    id          INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
	device_id	VARCHAR(127) NOT NULL,
    event_type  UNSIGNED INT8 NOT NULL,
    description VARCHAR(255),
    timestamp   DATETIME,
    created_at  DATETIME
);

CREATE TABLE operations (
	id        INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
	device_id VARCHAR(127) NOT NULL
);
`
