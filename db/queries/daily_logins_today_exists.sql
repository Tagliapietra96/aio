-- File: daily_logins_today_exists.sql
-- Created by: Matteo Tagliapietra 2024-10-15
-- Last modified: 2024-10-15
-- Purpose: Check if there are any daily logins in the database for today.
SELECT EXISTS(
SELECT 1 
FROM daily_logins 
WHERE DATE(created_at) = DATE('now', 'localtime')
);