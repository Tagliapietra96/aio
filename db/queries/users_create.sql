-- File: users_create.sql
-- Created by: Matteo Tagliapietra 2024-09-01
-- Last modified: 2024-09-01
-- Purpose: Create a new user in the database.
INSERT INTO users (firstname, lastname, nickname, created_at, updated_at)
VALUES (?, ?, ?, ?, ?);