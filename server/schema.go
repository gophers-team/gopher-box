package main

var schema = `
DROP TABLE IF EXISTS pills;
DROP TABLE IF EXISTS devices;
DROP TABLE IF EXISTS heartbeats;
DROP TABLE IF EXISTS device_dispensings;
DROP TABLE IF EXISTS dispensing_plans;
DROP TABLE IF EXISTS dispensing_schedule;

CREATE TABLE dispensing_plans (
    id          SERIAL PRIMARY KEY,
	name 		VARCHAR(255) NOT NULL,
    created_at  TIMESTAMP
);

CREATE TABLE devices (
    id          SERIAL PRIMARY KEY,
	plan_id    	INTEGER REFERENCES dispensing_plans(id),
	name 		VARCHAR(255) NOT NULL,
    created_at  TIMESTAMP
);

CREATE TABLE pills (
    id          SERIAL PRIMARY KEY,
	name 		VARCHAR(255) NOT NULL,
    created_at  TIMESTAMP
);

CREATE TABLE heartbeats (
    device_id   INTEGER NOT NULL,
    created_at  TIMESTAMP NOT NULL
);

CREATE TABLE dispensing_schedule (
	id          		SERIAL PRIMARY KEY,
    plan_id     		INTEGER REFERENCES dispensing_plans(id),
	pill_id				INTEGER REFERENCES pills(id),
	amount				INTEGER NOT NULL,
	dispense_dow		INTEGER NOT NULL,
	dispense_time   	TIMESTAMP NOT NULL,
	schedule_duration 	INTEGER NOT NULL,
	created_at  		TIMESTAMP
);

CREATE TABLE device_dispensings (
	id          	SERIAL PRIMARY KEY,
    device_id     	INTEGER REFERENCES devices(id),
	schedule_id		INTEGER  REFERENCES dispensing_schedule(id),
	pills_dispensed	INTEGER,
	status 			VARCHAR(255) NOT NULL,
    created_at  	TIMESTAMP
);
`
