-- File: characters_daily_login.sql
-- Purpose: Increments character stats and create a new daily login in the database.

-- Increment character stats
UPDATE characters
SET hp = CASE
    WHEN hp < 0.75 * max_hp
        THEN hp + (0.25 * max_hp)
        ELSE max_hp
    END,
    pp = max_pp
WHERE id = 1;

-- Create a new daily login
INSERT INTO daily_logins DEFAULT VALUES;