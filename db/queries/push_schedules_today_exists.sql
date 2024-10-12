-- File: push_schedules_today_exists.sql
-- Created by: Matteo Tagliapietra 2024-10-12
-- Last modified: 2024-10-12
-- Purpose: Get all push schedules for today.
SELECT EXISTS(
SELECT 1 
FROM push_schedules 
WHERE DATE(created_at) = DATE('now', 'localtime')
);