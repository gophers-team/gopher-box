------------------------------------------------------------------
-- Pills
INSERT INTO pills VALUES
(1, 'red', CURRENT_TIMESTAMP),
(2, 'green', CURRENT_TIMESTAMP),
(3, 'blue', CURRENT_TIMESTAMP);

------------------------------------------------------------------
-- Parkinson plan
INSERT INTO dispensing_plans
VALUES(1, 'parkinson', CURRENT_TIMESTAMP);

-- id, plan_id, pill_id, amount, dow, duration (mins), time, created
INSERT INTO dispensing_schedule VALUES
(1, 1, 1, 3, 1, "2018-01-01 10:00:00", 15, CURRENT_TIMESTAMP),
(2, 1, 1, 3, 2, "2018-01-01 10:00:00", 15, CURRENT_TIMESTAMP),
(3, 1, 1, 3, 3, "2018-01-01 10:00:00", 15, CURRENT_TIMESTAMP),
(4, 1, 1, 3, 4, "2018-01-01 10:00:00", 15, CURRENT_TIMESTAMP),
(5, 1, 1, 3, 5, "2018-01-01 10:00:00", 15, CURRENT_TIMESTAMP),
(6, 1, 1, 3, 6, "2018-01-01 10:00:00", 15, CURRENT_TIMESTAMP),
(7, 1, 1, 3, 7, "2018-01-01 10:00:00", 15, CURRENT_TIMESTAMP);

------------------------------------------------------------------
-- Devices
INSERT INTO devices VALUES
(1, 1, "Gomer Simpson's device", CURRENT_TIMESTAMP),
(2, 1, "Rick Sanchez", CURRENT_TIMESTAMP),
(3, 1, "Finn the human", CURRENT_TIMESTAMP),
(4, 1, "Avocato", CURRENT_TIMESTAMP),
