--------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------

-- File Name: tables.sql
-- Tables creation script
-- In this file we define the tables and triggers for the database
-- We use SQLite3 as the database engine

--------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------

--
-- characters table
--

-- the characters table is used to store the user character information
-- the character represents the user in the application, making the experience more engaging
CREATE TABLE IF NOT EXISTS characters (
    id INTEGER PRIMARY KEY AUTOINCREMENT, -- unique identifier for the character
    firstname TEXT NOT NULL, -- character's first name
    lastname TEXT NOT NULL, -- character's last name
    nickname TEXT NOT NULL UNIQUE, -- character's nickname, must be unique
    birthday TEXT NOT NULL, -- character's birthday, used for birthday greetings
    budget REAL NOT NULL DEFAULT 0.0, -- character's budget for financial management
    balance REAL NOT NULL DEFAULT 0.0, -- character's balance for financial management
    coins INTEGER NOT NULL DEFAULT 0, -- character's coins, used for rewards
    xp INTEGER NOT NULL DEFAULT 0, -- character's experience points
    next_level_xp INTEGER NOT NULL DEFAULT 50, -- experience points needed for next level
    level INTEGER NOT NULL DEFAULT 1, -- character's level
    pp INTEGER NOT NULL DEFAULT 50, -- character's current power points
    max_pp INTEGER NOT NULL DEFAULT 50, -- character's maximum power points
    hp INTEGER NOT NULL DEFAULT 100, -- character's current health points
    max_hp INTEGER NOT NULL DEFAULT 100, -- character's maximum health points
    karma INTEGER NOT NULL DEFAULT 0, -- character's karma, an indicator of performance
    created_at TEXT NOT NULL DEFAULT (datetime('now', 'localtime')), -- record creation timestamp
    updated_at TEXT NOT NULL DEFAULT (datetime('now', 'localtime')) -- record update timestamp
);

-- characters table indexes
CREATE INDEX IF NOT EXISTS characters_id_index ON characters (id);
CREATE INDEX IF NOT EXISTS characters_firstname_index ON characters (firstname);
CREATE INDEX IF NOT EXISTS characters_lastname_index ON characters (lastname);
CREATE INDEX IF NOT EXISTS characters_nickname_index ON characters (nickname);
CREATE INDEX IF NOT EXISTS characters_birthday_index ON characters (birthday);
CREATE INDEX IF NOT EXISTS characters_budget_index ON characters (budget);
CREATE INDEX IF NOT EXISTS characters_balance_index ON characters (balance);
CREATE INDEX IF NOT EXISTS characters_coins_index ON characters (coins);
CREATE INDEX IF NOT EXISTS characters_xp_index ON characters (xp);
CREATE INDEX IF NOT EXISTS characters_next_level_xp_index ON characters (next_level_xp);
CREATE INDEX IF NOT EXISTS characters_level_index ON characters (level);
CREATE INDEX IF NOT EXISTS characters_pp_index ON characters (pp);
CREATE INDEX IF NOT EXISTS characters_max_pp_index ON characters (max_pp);
CREATE INDEX IF NOT EXISTS characters_hp_index ON characters (hp);
CREATE INDEX IF NOT EXISTS characters_max_hp_index ON characters (max_hp);
CREATE INDEX IF NOT EXISTS characters_karma_index ON characters (karma);
CREATE INDEX IF NOT EXISTS characters_created_at_index ON characters (created_at);
CREATE INDEX IF NOT EXISTS characters_updated_at_index ON characters (updated_at);

--------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------

--
-- daily_logins table
--

-- the daily_logins table is used to store the user dayly logins
-- every day the user logs in the app, a dayly login is inserted in the daily_logins table
-- every dayly login insert if the user restores the pp to the max values, and if the user has less than 75% of the max hp, the user restores the 25% of the max hp
-- the created_at field is unique for prevent duplicate dayly logins, or conflicts in the transactions
CREATE TABLE IF NOT EXISTS daily_logins (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at TEXT NOT NULL DEFAULT (datetime('now', 'localtime')) UNIQUE
);

-- daily_logins table indexes
CREATE INDEX IF NOT EXISTS daily_logins_id_index ON daily_logins (id);
CREATE INDEX IF NOT EXISTS daily_logins_created_at_index ON daily_logins (created_at);

--------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------
