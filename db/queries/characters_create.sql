-- File: characters_create.sql
-- Purpose: Create a new character in the database.
INSERT INTO characters (firstname, lastname, nickname, birthday, budget)
VALUES(?, ?, ?, ?, ?);