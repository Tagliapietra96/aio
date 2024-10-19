-- File: characters_exists.sql
-- Purpose: Check if there are nay characters in the database.
SELECT EXISTS(SELECT 1 FROM characters);