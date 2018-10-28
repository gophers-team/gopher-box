------------------------------------------------------------------
-- Pills
INSERT INTO pills(id, name, created_at) VALUES
(1, 'red', CURRENT_TIMESTAMP),
(2, 'green', CURRENT_TIMESTAMP),
(3, 'blue', CURRENT_TIMESTAMP);

------------------------------------------------------------------
-- Parkinson plan
INSERT INTO dispensing_plans(id, name, created_at)
VALUES(1, 'parkinson', CURRENT_TIMESTAMP);

-- id, plan_id, pill_id, amount, dow, duration (mins), time, created
INSERT INTO dispensing_schedule(id, plan_id, pill_id, amount, dispense_dow, dispense_time, schedule_duration, created_at) VALUES
(1, 1, 1, 3, 1, CURRENT_TIMESTAMP, 15, CURRENT_TIMESTAMP),
(2, 1, 1, 3, 2, CURRENT_TIMESTAMP, 15, CURRENT_TIMESTAMP),
(3, 1, 1, 3, 3, CURRENT_TIMESTAMP, 15, CURRENT_TIMESTAMP),
(4, 1, 1, 3, 4, CURRENT_TIMESTAMP, 15, CURRENT_TIMESTAMP),
(5, 1, 1, 3, 5, CURRENT_TIMESTAMP, 15, CURRENT_TIMESTAMP),
(6, 1, 1, 3, 6, CURRENT_TIMESTAMP, 15, CURRENT_TIMESTAMP),
(7, 1, 1, 3, 7, CURRENT_TIMESTAMP, 15, CURRENT_TIMESTAMP);

------------------------------------------------------------------
-- Devices
INSERT INTO devices(id, plan_id, name, created_at) VALUES
(1, 1, 'Homer Simpson device', CURRENT_TIMESTAMP),
(2, 1, 'Rick Sanchez', CURRENT_TIMESTAMP),
(3, 1, 'Finn the human', CURRENT_TIMESTAMP),
(4, 1, 'Avocato', CURRENT_TIMESTAMP);
