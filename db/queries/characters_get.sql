-- File: characters_get.sql
-- Created by: Matteo Tagliapietra 2024-10-15
-- Last modified: 2024-10-15
-- Purpose: Get the character from the database.
SELECT
firstname,
lastname,
nickname,
birthday,
budget,
balance,
coins,
xp,
next_level_xp,
level,
pp,
max_pp,
hp,
max_hp,
karma,
created_at,
updated_at
FROM characters
WHERE id = 1;