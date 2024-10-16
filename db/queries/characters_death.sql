-- File: characters_death.sql
-- Purpose: Update the character's stats when on death.
UPDATE characters
SET hp = 100,
    max_hp = 100,
    pp = 50,
    max_pp = 50,
    karma = karma - 10,
    coins = 0,
    level = 1,
    xp = 0,
    next_level_xp = 50
WHERE id = 1;
