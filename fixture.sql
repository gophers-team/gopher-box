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
INSERT INTO dispensing_schedule(id, plan_id, pill_id, amount, interval, created_at) VALUES
(1, 1, 1, 1, 1, CURRENT_TIMESTAMP);

------------------------------------------------------------------
-- Devices
INSERT INTO devices(id, plan_id, name, created_at) VALUES
(1, 1, 'Homer Simpson device', CURRENT_TIMESTAMP),
(2, 1, 'Rick Sanchez', CURRENT_TIMESTAMP),
(3, 1, 'Finn the human', CURRENT_TIMESTAMP),
(4, 1, 'Avocato', CURRENT_TIMESTAMP);
