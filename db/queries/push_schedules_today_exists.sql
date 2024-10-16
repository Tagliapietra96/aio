-- File: push_schedules_today_exists.sql
-- Purpose: Get all push schedules for today.
SELECT EXISTS(
SELECT 1 
FROM push_schedules 
WHERE DATE(created_at) = DATE('now', 'localtime')
);