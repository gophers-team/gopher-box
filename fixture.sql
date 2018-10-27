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

-- id, plan_id, pill_id, amount, dow, time, created
INSERT INTO dispense_schedule VALUES
(1, 1, 1, 3, 1, "10:00", CURRENT_TIMESTAMP),
(2, 1, 1, 3, 2, "10:00", CURRENT_TIMESTAMP),
(3, 1, 1, 3, 3, "10:00", CURRENT_TIMESTAMP),
(4, 1, 1, 3, 4, "10:00", CURRENT_TIMESTAMP),
(5, 1, 1, 3, 5, "10:00", CURRENT_TIMESTAMP),
(6, 1, 1, 3, 6, "10:00", CURRENT_TIMESTAMP),
(7, 1, 1, 3, 7, "10:00", CURRENT_TIMESTAMP);
