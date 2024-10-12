-- File: users_get.sql
-- Created by: Matteo Tagliapietra 2024-09-01
-- Last modified: 2024-09-01
-- Purpose: Get a user from the database.
SELECT id, 
created_at, 
lastedit_at, 
firstname, 
lastname, 
nickname, 
hp, 
pp, 
experience, 
level, 
budget, 
coins 
FROM users 
WHERE id = 1;