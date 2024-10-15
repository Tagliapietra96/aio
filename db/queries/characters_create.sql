-- File: characters_create.sql
-- Created by: Matteo Tagliapietra 2024-10-15
-- Last modified: 2024-10-15
-- Purpose: Create a new character in the database.
INSERT INTO characters (firstname, lastname, nickname, birthday, budget)
VALUES(?, ?, ?, ?, ?);