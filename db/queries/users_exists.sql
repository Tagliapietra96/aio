-- File: users_exists.sql
-- Created by: Matteo Tagliapietra 2024-09-01
-- Last modified: 2024-09-01
-- Purpose: Check if a user exists in the database.
SELECT EXISTS(SELECT 1 FROM users WHERE id = 1);