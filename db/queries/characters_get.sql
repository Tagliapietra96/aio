-- File: characters_get.sql
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