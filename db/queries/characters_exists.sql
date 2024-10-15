-- File: characters_exists.sql
-- Created by: Matteo Tagliapietra 2024-10-15
-- Last modified: 2024-10-15
-- Purpose: Check if there are nay characters in the database.
SELECT EXISTS(SELECT 1 FROM characters);