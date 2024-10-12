--------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------

-- File Name: tables.sql
-- Created by: Matteo Tagliapietra 2024-09-01
-- Last modified: 2024-10-12

-- Tables creation script
-- In this file we define the tables and triggers for the database
-- We use SQLite3 as the database engine

--------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------

--
-- push_schedules table
--

-- the push_schedules table is used to store the user push schedules
-- the push schedules are used to store the user push notifications
-- the notifications are used to log every time the app push the db to the remote
-- the created_at field is unique for prevent duplicate push schedules, or conflicts in the transactions
CREATE TABLE IF NOT EXISTS push_schedules (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at TEXT NOT NULL DEFAULT (datetime('now', 'localtime')) UNIQUE
);

-- push_schedules table indexes
CREATE INDEX IF NOT EXISTS push_schedules_id_index ON push_schedules (id);
CREATE INDEX IF NOT EXISTS push_schedules_created_at_index ON push_schedules (created_at);

--------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------

--
-- users table
--

-- from now only one user is allowed, the user with id = 1. This user is the main user and is used to store the user stats
-- if there is no user with id = 1, the app at the first run will insert the user with id = 1 asking the user to fill the user data
-- this entity is the main entity of the app, and all the other entities are related to this entity (DIRECTLY OR INDIRECTLY)
-- the created_at and nickname fields are unique, in vision of future implementation of multiple users
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at TEXT NOT NULL DEFAULT (datetime('now', 'localtime')) UNIQUE,
    updated_at TEXT NOT NULL DEFAULT (datetime('now', 'localtime')),
    firstname TEXT NOT NULL, -- user first name
    lastname TEXT NOT NULL, -- user last name
    nickname TEXT NOT NULL UNIQUE, -- user nickname, must be unique
    max_hp INTEGER NOT NULL DEFAULT 100, -- user max health points, they are increased by 10 after level up
    hp INTEGER NOT NULL DEFAULT 100, -- user health points they are increased by 10 after level up and descreased by bad events, if hp <= 0 user is dead and skills and user xp are reset
    max_pp INTEGER NOT NULL DEFAULT 50, -- user max productivity points, they are increased by 5 after level up
    pp INTEGER NOT NULL DEFAULT 50, -- user productivity points they are increased by 5 after level up and descreased by missions and rituals, if pp <= 0 user is burned out and skills and not able to do missions and rituals
    xp INTEGER NOT NULL DEFAULT 0, -- user experience points they are increased by missions, rituals and goals, if xp >= next_level_xp user level up
    karma INTEGER NOT NULL DEFAULT 0, -- user karma, it is increased by good events and decreased by bad events
    next_level_xp INTEGER NOT NULL DEFAULT 50, -- user next level xp, it is increased by 50 after level up
    level INTEGER NOT NULL DEFAULT 1, -- user level, it is increased by 1 after level up
    budget REAL NOT NULL DEFAULT 0.0, -- user monthly budget, it is used to check if user can afford expenses
    balance REAL NOT NULL DEFAULT 0.0, -- user balance, it is updated after account insert, update and delete, is the sum of all user accounts balance
    coins INTEGER NOT NULL DEFAULT 0 -- user coins, they are increased by goals and can be used to buy hp, or used to get the wishlist items
);

--
-- users table indexes
--

CREATE INDEX IF NOT EXISTS users_id_index ON users (id);

--
-- users table triggers
--

-- update user updateded_at after user update
CREATE TRIGGER IF NOT EXISTS update_user_updated_at_after_user_update
AFTER UPDATE ON users
FOR EACH ROW
BEGIN
    UPDATE users
    SET updated_at = datetime('now', 'localtime')
    WHERE id = NEW.id;
END;

--------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------

--
-- dayly_logins table
--

-- the dayly_logins table is used to store the user dayly logins
-- every day the user logs in the app, a dayly login is inserted in the dayly_logins table
-- every dayly login insert if the user restores the pp to the max values, and if the user has less than 75% of the max hp, the user restores the 25% of the max hp
-- the created_at field is unique for prevent duplicate dayly logins, or conflicts in the transactions
CREATE TABLE IF NOT EXISTS dayly_logins (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at TEXT NOT NULL DEFAULT (datetime('now', 'localtime')) UNIQUE,
    updated_at TEXT NOT NULL DEFAULT (datetime('now', 'localtime')),
    user_id INTEGER NOT NULL DEFAULT 1, -- user id, the user that owns the dayly login
    FOREIGN KEY (user_id) REFERENCES users(id)
);

--
-- dayly_logins table indexes
--

CREATE INDEX IF NOT EXISTS dayly_logins_id_index ON dayly_logins (id);
CREATE INDEX IF NOT EXISTS dayly_logins_created_at_index ON dayly_logins (created_at);
CREATE INDEX IF NOT EXISTS dayly_logins_user_id_index ON dayly_logins (user_id);

--
-- dayly_logins table triggers
--

-- update dayly_login updateded_at after dayly_login update
CREATE TRIGGER IF NOT EXISTS update_dayly_login_updated_at_after_dayly_login_update
AFTER UPDATE ON dayly_logins
FOR EACH ROW
BEGIN
    UPDATE dayly_logins
    SET updated_at = datetime('now', 'localtime')
    WHERE id = NEW.id;
END;

-- restore user hp and pp after dayly_login insert
CREATE TRIGGER IF NOT EXISTS restore_user_hp_and_pp_after_dayly_login_insert
AFTER INSERT ON dayly_logins
FOR EACH ROW
BEGIN
    -- restore user hp and pp
    UPDATE users
    SET 
        hp = CASE
            WHEN hp < 0.75 * max_hp THEN hp + (0.25 * max_hp)
            ELSE max_hp
        END,
        pp = max_pp
    WHERE id = NEW.user_id;
END;

--------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------

--
-- death_logs table
--

-- the death_logs table is used to store the user death logs
-- every time the user finishes the hp, a death log is inserted in the death_logs table
-- when a death log is inserted, the user xp, level, hp, pp and coins are reset and the user get -10 karma
-- the created_at field is unique for prevent duplicate death logs, or conflicts in the transactions
CREATE TABLE IF NOT EXISTS death_logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at TEXT NOT NULL DEFAULT (datetime('now', 'localtime')) UNIQUE,
    updated_at TEXT NOT NULL DEFAULT (datetime('now', 'localtime')),
    user_id INTEGER NOT NULL DEFAULT 1, -- user id, the user that owns the death log
    FOREIGN KEY (user_id) REFERENCES users(id)
);

--
-- death_logs table indexes
--

CREATE INDEX IF NOT EXISTS death_logs_id_index ON death_logs (id);
CREATE INDEX IF NOT EXISTS death_logs_created_at_index ON death_logs (created_at);
CREATE INDEX IF NOT EXISTS death_logs_user_id_index ON death_logs (user_id);

--
-- death_logs table triggers
--

-- update death_log updateded_at after death_log update
CREATE TRIGGER IF NOT EXISTS update_death_log_updated_at_after_death_log_update
AFTER UPDATE ON death_logs
FOR EACH ROW
BEGIN
    UPDATE death_logs
    SET updated_at = datetime('now', 'localtime')
    WHERE id = NEW.id;
END;

-- reset user xp, level, hp, pp, skills and coins after death_log insert
CREATE TRIGGER IF NOT EXISTS reset_user_xp_level_hp_pp_skills_and_coins_after_death_log_insert
AFTER INSERT ON death_logs
FOR EACH ROW
BEGIN
    -- reset user xp, level, hp, pp, skills and coins
    UPDATE users
    SET 
        xp = 0,
        next_level_xp = 50,
        level = 1,
        max_hp = 100,
        hp = 100,
        max_pp = 50,
        pp = 50,
        karma = karma - 10,
        coins = 0
    WHERE id = NEW.user_id;
END;

--
-- logs table
--

-- the logs table is used to store the logs of the app, the logs are used to store the app events
-- every time an event occurs in the app, a log is inserted in the logs table
-- the karma field is used to store the log karma, it can be 0 for neutral events, 1 for good events, -1 for bad events
-- for example if a user completes a mission, a log is inserted in the logs table with the message "You earned 10 xp for completing the mission 'Task name'"
CREATE TABLE IF NOT EXISTS logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at TEXT NOT NULL DEFAULT (datetime('now', 'localtime')) UNIQUE,
    updated_at TEXT NOT NULL DEFAULT (datetime('now', 'localtime')),
    message TEXT NOT NULL, -- log message, it can be every event message possible
    karma INTEGER NOT NULL DEFAULT 0, -- log karma, 0 for neutral events, 1 for good events, -1 for bad events
    user_id INTEGER NOT NULL DEFAULT 1, -- user id, the user that generated the log
    FOREIGN KEY (user_id) REFERENCES users(id)
);

--
-- logs table indexes
--

CREATE INDEX IF NOT EXISTS logs_id_index ON logs (id);
CREATE INDEX IF NOT EXISTS logs_created_at_index ON logs (created_at);
CREATE INDEX IF NOT EXISTS logs_user_id_index ON logs (user_id);

--
-- logs table triggers
--

-- update log updateded_at after log update
CREATE TRIGGER IF NOT EXISTS update_log_updated_at_after_log_update
AFTER UPDATE ON logs
FOR EACH ROW
BEGIN
    UPDATE logs
    SET updated_at = datetime('now', 'localtime')
    WHERE id = NEW.id;
END;

--------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------

--
-- accounts table
--

-- the accounts table is used to store the user accounts, the accounts are used to store the user money
-- the user can have multiple accounts, for example a cash account, a bank account, a credit card account, etc.
-- the currency field is used to store the account currency, every transaction related to the account gets the account currency
-- name and created_at fields are unique for prevent duplicate accounts, or conflicts in the transactions
CREATE TABLE IF NOT EXISTS accounts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at TEXT NOT NULL DEFAULT (datetime('now', 'localtime')) UNIQUE,
    updated_at TEXT NOT NULL DEFAULT (datetime('now', 'localtime')),
    name TEXT NOT NULL UNIQUE, -- account name, must be unique
    description TEXT, -- account description
    currency TEXT NOT NULL DEFAULT 'EUR', -- account currency, it is used to store the account currency
    balance REAL NOT NULL DEFAULT 0.0, -- account balance, it will be affected by transferts, expenses and incomes
    user_id INTEGER NOT NULL DEFAULT 1, -- user id, the user that owns the account
    FOREIGN KEY (user_id) REFERENCES users(id)
);

--
-- accounts table indexes
--

CREATE INDEX IF NOT EXISTS accounts_id_index ON accounts (id);
CREATE INDEX IF NOT EXISTS accounts_name_index ON accounts (name);
CREATE INDEX IF NOT EXISTS accounts_user_id_index ON accounts (user_id);

--
-- accounts table triggers
--

-- update account updateded_at after account update
CREATE TRIGGER IF NOT EXISTS update_account_updated_at_after_account_update
AFTER UPDATE ON accounts
FOR EACH ROW
BEGIN
    UPDATE accounts
    SET updated_at = datetime('now', 'localtime')
    WHERE id = NEW.id;
END;

-- update user balance after account insert
CREATE TRIGGER IF NOT EXISTS update_user_balance_after_account_insert
AFTER INSERT ON accounts
FOR EACH ROW
BEGIN
    UPDATE users
    SET balance = balance + NEW.balance
    WHERE id = NEW.user_id;
END;

-- update user balance after account update balance
CREATE TRIGGER IF NOT EXISTS update_user_balance_after_account_update
AFTER UPDATE ON accounts
FOR EACH ROW
WHEN NEW.balance != OLD.balance
BEGIN
    UPDATE users
    SET balance = balance + (NEW.balance - OLD.balance)
    WHERE id = NEW.user_id;
END;

-- update user balance after account delete
CREATE TRIGGER IF NOT EXISTS update_user_balance_after_account_delete
AFTER DELETE ON accounts
FOR EACH ROW
BEGIN
    UPDATE users
    SET balance = balance - OLD.balance
    WHERE id = OLD.user_id;
END;

--------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------

--
-- transfers table
--

-- the transfers table is used to store the user transfers, the transfers are used to store the user money transfers between accounts
-- every time a user transfer money from an account to another account, the balance of the accounts are updated
-- the table has triggers to update the account balance after transfer insert, update and delete
-- the created_at field and date field are unique for prevent duplicate transfers, or conflicts in the transactions
CREATE TABLE IF NOT EXISTS transfers (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at TEXT NOT NULL DEFAULT (datetime('now', 'localtime')) UNIQUE,
    updated_at TEXT NOT NULL DEFAULT (datetime('now', 'localtime')),
    date TEXT NOT NULL DEFAULT (datetime('now', 'localtime')) UNIQUE, -- transfer date
    amount REAL NOT NULL DEFAULT 0.0, -- transfer amount
    from_account_id INTEGER NOT NULL, -- from account id, the account that sends the money
    to_account_id INTEGER NOT NULL, -- to account id, the account that receives the money
    FOREIGN KEY (from_account_id) REFERENCES accounts(id),
    FOREIGN KEY (to_account_id) REFERENCES accounts(id)
);

--
-- transfers table indexes
--

CREATE INDEX IF NOT EXISTS transfers_id_index ON transfers (id);
CREATE INDEX IF NOT EXISTS transfers_date_index ON transfers (date);
CREATE INDEX IF NOT EXISTS transfers_amount_index ON transfers (amount);
CREATE INDEX IF NOT EXISTS transfers_from_account_id_index ON transfers (from_account_id);
CREATE INDEX IF NOT EXISTS transfers_to_account_id_index ON transfers (to_account_id);

--
-- transfers table triggers
--

-- update transfer updateded_at after transfer update
CREATE TRIGGER IF NOT EXISTS update_transfer_updated_at_after_transfer_update
AFTER UPDATE ON transfers
FOR EACH ROW
BEGIN
    UPDATE transfers
    SET updated_at = datetime('now', 'localtime')
    WHERE id = NEW.id;
END;

-- update from_account and to_account balance after transfer insert
CREATE TRIGGER IF NOT EXISTS update_from_account_and_to_account_balance_after_transfer_insert
AFTER INSERT ON transfers
FOR EACH ROW
BEGIN
    -- update from_account balance
    UPDATE accounts
    SET balance = balance - NEW.amount
    WHERE id = NEW.from_account_id;

    -- update to_account balance
    UPDATE accounts
    SET balance = balance + NEW.amount
    WHERE id = NEW.to_account_id;

    -- update from_account user balance
    UPDATE users
    SET balance = balance - NEW.amount
    WHERE id = (SELECT user_id FROM accounts WHERE id = NEW.from_account_id);

    -- update to_account user balance
    UPDATE users
    SET balance = balance + NEW.amount
    WHERE id = (SELECT user_id FROM accounts WHERE id = NEW.to_account_id);
END;

-- update from_account and to_account balance after transfer amount update
CREATE TRIGGER IF NOT EXISTS update_from_account_and_to_account_balance_after_transfer_amount_update
AFTER UPDATE ON transfers
FOR EACH ROW
WHEN NEW.amount != OLD.amount
BEGIN
    -- update from_account balance
    UPDATE accounts
    SET balance = balance - (NEW.amount - OLD.amount)
    WHERE id = NEW.from_account_id;

    -- update to_account balance
    UPDATE accounts
    SET balance = balance + (NEW.amount - OLD.amount)
    WHERE id = NEW.to_account_id;

    -- update from_account user balance
    UPDATE users
    SET balance = balance - (NEW.amount - OLD.amount)
    WHERE id = (SELECT user_id FROM accounts WHERE id = NEW.from_account_id);

    -- update to_account user balance
    UPDATE users
    SET balance = balance + (NEW.amount - OLD.amount)
    WHERE id = (SELECT user_id FROM accounts WHERE id = NEW.to_account_id);
END;

-- update from_account and to_account balance after transfer delete
CREATE TRIGGER IF NOT EXISTS update_from_account_and_to_account_balance_after_transfer_delete
AFTER DELETE ON transfers
FOR EACH ROW
BEGIN
    -- update from_account balance
    UPDATE accounts
    SET balance = balance + OLD.amount
    WHERE id = OLD.from_account_id;

    -- update to_account balance
    UPDATE accounts
    SET balance = balance - OLD.amount
    WHERE id = OLD.to_account_id;

    -- update from_account user balance
    UPDATE users
    SET balance = balance + OLD.amount
    WHERE id = (SELECT user_id FROM accounts WHERE id = OLD.from_account_id);

    -- update to_account user balance
    UPDATE users
    SET balance = balance - OLD.amount
    WHERE id = (SELECT user_id FROM accounts WHERE id = OLD.to_account_id);
END;

--------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------

--
-- expense_sources table
--

-- the expense_sources table is used to store the user expense sources
-- the espense sources are used to categorize the user expenses
-- the name field and created_at field are unique for prevent duplicate expense sources or conflicts in the transactions
CREATE TABLE IF NOT EXISTS expense_sources (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at TEXT NOT NULL DEFAULT (datetime('now', 'localtime')) UNIQUE,
    updated_at TEXT NOT NULL DEFAULT (datetime('now', 'localtime')),
    name TEXT NOT NULL UNIQUE, -- expense source name, must be unique
    description TEXT NOT NULL, -- expense source description
    user_id INTEGER NOT NULL DEFAULT 1,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

--
-- expense_sources table indexes
--

CREATE INDEX IF NOT EXISTS expense_sources_id_index ON expense_sources (id);
CREATE INDEX IF NOT EXISTS expense_sources_name_index ON expense_sources (name);
CREATE INDEX IF NOT EXISTS expense_sources_user_id_index ON expense_sources (user_id);

--
-- expense_sources table triggers
--

-- update expense_source updateded_at after expense_source update
CREATE TRIGGER IF NOT EXISTS update_expense_source_updated_at_after_expense_source_update
AFTER UPDATE ON expense_sources
FOR EACH ROW
BEGIN
    UPDATE expense_sources
    SET updated_at = datetime('now', 'localtime')
    WHERE id = NEW.id;
END;

------------------------------------------------------------------------------------
------------------------------------------------------------------------------------
------------------------------------------------------------------------------------

--
-- expenses table
--

-- the expenses table is used to store the user expenses
-- the expenses are used to store the user money expenses
-- the table has triggers to update the account balance after expense insert, update and delete
-- every time a user spends money, an expense is inserted in the expenses table and the account balance is updated
-- the created_at field is unique for prevent duplicate expenses, or conflicts in the transactions
CREATE TABLE IF NOT EXISTS expenses (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at TEXT NOT NULL DEFAULT (datetime('now', 'localtime')) UNIQUE,
    updated_at TEXT NOT NULL DEFAULT (datetime('now', 'localtime')),
    date TEXT NOT NULL DEFAULT (datetime('now', 'localtime')), -- expense date
    description TEXT, -- expense description
    amount REAL NOT NULL DEFAULT 0.0, -- expense amount
    account_id INTEGER NOT NULL, -- account id, the account that pays the expense
    expense_source_id INTEGER NOT NULL, -- expense source id, the expense source that categorizes the expense
    FOREIGN KEY (account_id) REFERENCES accounts(id),
    FOREIGN KEY (expense_source_id) REFERENCES expense_sources(id)
);

-- expenses table indexes
CREATE INDEX IF NOT EXISTS expenses_id_index ON expenses (id);
CREATE INDEX IF NOT EXISTS expenses_date_index ON expenses (date);
CREATE INDEX IF NOT EXISTS expenses_amount_index ON expenses (amount);
CREATE INDEX IF NOT EXISTS expenses_account_id_index ON expenses (account_id);
CREATE INDEX IF NOT EXISTS expenses_expense_source_id_index ON expenses (expense_source_id);

-- expenses table triggers
-- update expense updateded_at after expense update
CREATE TRIGGER IF NOT EXISTS update_expense_updated_at_after_expense_update
AFTER UPDATE ON expenses
FOR EACH ROW
BEGIN
    UPDATE expenses
    SET updated_at = datetime('now', 'localtime')
    WHERE id = NEW.id;
END;

-- update account balance after expense insert
CREATE TRIGGER IF NOT EXISTS update_account_balance_after_expense_insert
AFTER INSERT ON expenses
FOR EACH ROW
BEGIN
    -- update account balance
    UPDATE accounts
    SET balance = balance - NEW.amount
    WHERE id = NEW.account_id;

    -- update user balance
    UPDATE users
    SET balance = balance - NEW.amount
    WHERE id = (SELECT user_id FROM accounts WHERE id = NEW.account_id);
END;

-- update account balance after expense amount update
CREATE TRIGGER IF NOT EXISTS update_account_balance_after_expense_amount_update
AFTER UPDATE ON expenses
FOR EACH ROW
WHEN NEW.amount != OLD.amount
BEGIN
    -- update account balance
    UPDATE accounts
    SET balance = balance - (NEW.amount - OLD.amount)
    WHERE id = NEW.account_id;

    -- update user balance
    UPDATE users
    SET balance = balance - (NEW.amount - OLD.amount)
    WHERE id = (SELECT user_id FROM accounts WHERE id = NEW.account_id);
END;

-- update account balance after expense delete
CREATE TRIGGER IF NOT EXISTS update_account_balance_after_expense_delete
AFTER DELETE ON expenses
FOR EACH ROW
BEGIN
    -- update account balance
    UPDATE accounts
    SET balance = balance + OLD.amount
    WHERE id = OLD.account_id;

    -- update user balance
    UPDATE users
    SET balance = balance + OLD.amount
    WHERE id = (SELECT user_id FROM accounts WHERE id = OLD.account_id);
END;

--------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------

--
-- rents table
--

-- the rents table is used to store the user rents
-- the rents are used to create an expense that is paid periodically
-- the period field is used to store the rent period, it can be DAY, WEEK, MONTH, YEAR
-- in the app there is a script that checks the rents and creates the expenses for the rents if there aren't expenses that match the period based on the reference_expense_id
-- the reference_expense_id field is used to store the reference expense id, the reference expense is the expense that is used to create the rent expenses overtime. 
-- the date of the reference expense and the period of the rent are used to calculate when create a new expense for the rent
-- the reference expense date is used as 'rent start date' and the rent period is used to calculate the next rent date
-- the table has triggers to update the account balance after rent insert, update and delete
-- the created_at field and the reference_expense_id field are unique for prevent duplicate rents, or conflicts in the transactions
CREATE TABLE IF NOT EXISTS rents (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at TEXT NOT NULL DEFAULT (datetime('now', 'localtime')) UNIQUE,
    updated_at TEXT NOT NULL DEFAULT (datetime('now', 'localtime')),
    period TEXT NOT NULL DEFAULT 'MONTH', -- rent period, it can be DAY, WEEK, MONTH, YEAR
    enabled INTEGER NOT NULL DEFAULT 1, -- rent enabled, if 0 the rent is disabled if 1 the rent is enabled
    reference_expense_id INTEGER NOT NULL UNIQUE, -- reference expense id, the reference expense that is used to create the rent expenses
    FOREIGN KEY (reference_expense_id) REFERENCES expenses(id)
);

--
-- rents table indexes
--

CREATE INDEX IF NOT EXISTS rents_id_index ON rents (id);
CREATE INDEX IF NOT EXISTS rents_period_index ON rents (period);
CREATE INDEX IF NOT EXISTS rents_reference_expense_id_index ON rents (reference_expense_id);

--
-- rents table triggers
--

-- update rent updateded_at after rent update
CREATE TRIGGER IF NOT EXISTS update_rent_updated_at_after_rent_update
AFTER UPDATE ON rents
FOR EACH ROW
BEGIN
    UPDATE rents
    SET updated_at = datetime('now', 'localtime')
    WHERE id = NEW.id;
END;

--------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------

--
-- income_sources table
--

-- the income_sources table is used to store the user income sources
-- the income sources are used to categorize the user incomes
-- the name field and created_at field are unique for prevent duplicate income sources or conflicts in the transactions
CREATE TABLE IF NOT EXISTS income_sources (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at TEXT NOT NULL DEFAULT (datetime('now', 'localtime')) UNIQUE,
    updated_at TEXT NOT NULL DEFAULT (datetime('now', 'localtime')),
    name TEXT NOT NULL UNIQUE, -- income source name, must be unique
    description TEXT NOT NULL, -- income source description
    user_id INTEGER NOT NULL DEFAULT 1, -- user id, the user that owns the income source
    FOREIGN KEY (user_id) REFERENCES users(id)
);

--
-- income_sources table indexes
--

CREATE INDEX IF NOT EXISTS income_sources_id_index ON income_sources (id);
CREATE INDEX IF NOT EXISTS income_sources_name_index ON income_sources (name);
CREATE INDEX IF NOT EXISTS income_sources_user_id_index ON income_sources (user_id);

--
-- income_sources table triggers
--

-- update income_source updateded_at after income_source update
CREATE TRIGGER IF NOT EXISTS update_income_source_updated_at_after_income_source_update
AFTER UPDATE ON income_sources
FOR EACH ROW
BEGIN
    UPDATE income_sources
    SET updated_at = datetime('now', 'localtime')
    WHERE id = NEW.id;
END;

--------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------

--
-- incomes table
--

-- the incomes table is used to store the user incomes
-- the incomes are used to store the user money incomes
-- the table has triggers to update the account balance after income insert, update and delete
-- every time a user earns money, an income is inserted in the incomes table and the account balance is updated
-- the created_at field is unique for prevent duplicate incomes, or conflicts in the transactions
CREATE TABLE IF NOT EXISTS incomes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at TEXT NOT NULL DEFAULT (datetime('now', 'localtime')) UNIQUE,
    updated_at TEXT NOT NULL DEFAULT (datetime('now', 'localtime')),
    date TEXT NOT NULL DEFAULT (datetime('now', 'localtime')), -- income date
    description TEXT, -- income description
    amount REAL NOT NULL DEFAULT 0.0, -- income amount
    account_id INTEGER NOT NULL, -- account id, the account that receives the income
    income_source_id INTEGER NOT NULL, -- income source id, the income source that categorizes the income
    FOREIGN KEY (account_id) REFERENCES accounts(id),
    FOREIGN KEY (income_source_id) REFERENCES income_sources(id)
);

--
-- incomes table indexes
--

CREATE INDEX IF NOT EXISTS incomes_id_index ON incomes (id);
CREATE INDEX IF NOT EXISTS incomes_date_index ON incomes (date);
CREATE INDEX IF NOT EXISTS incomes_amount_index ON incomes (amount);
CREATE INDEX IF NOT EXISTS incomes_account_id_index ON incomes (account_id);
CREATE INDEX IF NOT EXISTS incomes_income_source_id_index ON incomes (income_source_id);

--
-- incomes table triggers
--

-- update income updateded_at after income update
CREATE TRIGGER IF NOT EXISTS update_income_updated_at_after_income_update
AFTER UPDATE ON incomes
FOR EACH ROW
BEGIN
    UPDATE incomes
    SET updated_at = datetime('now', 'localtime')
    WHERE id = NEW.id;
END;

-- update account balance after income insert
CREATE TRIGGER IF NOT EXISTS update_account_balance_after_income_insert
AFTER INSERT ON incomes
FOR EACH ROW
BEGIN
    -- update account balance
    UPDATE accounts
    SET balance = balance + NEW.amount
    WHERE id = NEW.account_id;

    -- update user balance
    UPDATE users
    SET balance = balance + NEW.amount
    WHERE id = (SELECT user_id FROM accounts WHERE id = NEW.account_id);
END;

-- update account balance after income amount update
CREATE TRIGGER IF NOT EXISTS update_account_balance_after_income_amount_update
AFTER UPDATE ON incomes
FOR EACH ROW
WHEN NEW.amount != OLD.amount
BEGIN
    -- update account balance
    UPDATE accounts
    SET balance = balance + (NEW.amount - OLD.amount)
    WHERE id = NEW.account_id;

    -- update user balance
    UPDATE users
    SET balance = balance + (NEW.amount - OLD.amount)
    WHERE id = (SELECT user_id FROM accounts WHERE id = NEW.account_id);
END;

-- update account balance after income delete
CREATE TRIGGER IF NOT EXISTS update_account_balance_after_income_delete
AFTER DELETE ON incomes
FOR EACH ROW
BEGIN
    -- update account balance
    UPDATE accounts
    SET balance = balance - OLD.amount
    WHERE id = OLD.account_id;

    -- update user balance
    UPDATE users
    SET balance = balance - OLD.amount
    WHERE id = (SELECT user_id FROM accounts WHERE id = OLD.account_id);
END;

--------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------

--
-- balance_recalc_logs table
--

-- the balance_recalc_logs table is used to store the user balance recalculation logs
-- every time a log is added, the accounts and user balance are recalculated
-- the log simply store the date of the recalculation and the user id
-- the date it is used to understand when the last recalculation was executed and if must be executed again
-- the created_at field is unique for prevent duplicate logs, or conflicts in the transactions
CREATE TABLE IF NOT EXISTS balance_recalc_logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at TEXT NOT NULL DEFAULT (datetime('now', 'localtime')) UNIQUE,
    updated_at TEXT NOT NULL DEFAULT (datetime('now', 'localtime')),
    user_id INTEGER NOT NULL DEFAULT 1, -- user id, the user that owns the balance recalculation log
    FOREIGN KEY (user_id) REFERENCES users(id)
);

--
-- balance_recalc_logs table indexes
--
CREATE INDEX IF NOT EXISTS balance_recalc_logs_id_index ON balance_recalc_logs (id);
CREATE INDEX IF NOT EXISTS balance_recalc_logs_created_at_index ON balance_recalc_logs (created_at);
CREATE INDEX IF NOT EXISTS balance_recalc_logs_user_id_index ON balance_recalc_logs (user_id);

--
-- balance_recalc_logs table triggers
--

-- update balance_recalc_log updateded_at after balance_recalc_log update
CREATE TRIGGER IF NOT EXISTS update_balance_recalc_log_updated_at_after_balance_recalc_log_update
AFTER UPDATE ON balance_recalc_logs
FOR EACH ROW
BEGIN
    UPDATE balance_recalc_logs
    SET updated_at = datetime('now', 'localtime')
    WHERE id = NEW.id;
END;

-- update accounts and user balance after balance_recalc_log insert
CREATE TRIGGER IF NOT EXISTS update_accounts_and_user_balance_after_balance_recalc_log_insert
AFTER INSERT ON balance_recalc_logs
FOR EACH ROW
BEGIN
    -- update accounts balance
    UPDATE accounts
    SET balance = (
        (SELECT COALESCE(SUM(amount), 0) FROM transfers WHERE to_account_id = accounts.id) +
        (SELECT COALESCE(SUM(amount), 0) FROM incomes WHERE account_id = accounts.id) -
        (SELECT COALESCE(SUM(amount), 0) FROM expenses WHERE account_id = accounts.id) -
        (SELECT COALESCE(SUM(amount), 0) FROM transfers WHERE from_account_id = accounts.id)
    )
    WHERE user_id = NEW.user_id;

    -- update user balance
    UPDATE users
    SET balance = (
        SELECT COALESCE(SUM(balance), 0)
        FROM accounts 
        WHERE user_id = NEW.user_id
    )
    WHERE id = NEW.user_id;
END;

--------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------

--
-- areas table
--

-- the areas table is used to store the user areas
-- the areas are used to categorize the user skills, goals, quests missions, rituals, notes, etc.
-- the name field and created_at field are unique for prevent duplicate areas or conflicts in the transactions
CREATE TABLE IF NOT EXISTS areas (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at TEXT NOT NULL DEFAULT (datetime('now', 'localtime')) UNIQUE,
    updated_at TEXT NOT NULL DEFAULT (datetime('now', 'localtime')),
    name TEXT NOT NULL UNIQUE, -- area name, must be unique
    description TEXT NOT NULL, -- area description
    user_id INTEGER NOT NULL DEFAULT 1, -- user id, the user that owns the area
    FOREIGN KEY (user_id) REFERENCES users(id)
);

--
-- areas table indexes
--

CREATE INDEX IF NOT EXISTS areas_id_index ON areas (id);
CREATE INDEX IF NOT EXISTS areas_name_index ON areas (name);
CREATE INDEX IF NOT EXISTS areas_user_id_index ON areas (user_id);

--
-- areas table triggers
--

-- update area updateded_at after area update
CREATE TRIGGER IF NOT EXISTS update_area_updated_at_after_area_update
AFTER UPDATE ON areas
FOR EACH ROW
BEGIN
    UPDATE areas
    SET updated_at = datetime('now', 'localtime')
    WHERE id = NEW.id;
END;

--------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------

--
-- skills table
--

-- the skills table is used to store the user skills
-- the skills are used to store the user abilities, the user can level up the skills by gaining xp
-- the xp field is used to store the skill experience points, the level field is used to store the skill level
-- the next_level_xp field is used to store the skill next level xp, the area_id field is used to store the area id that the skill belongs to
-- the name field and created_at field are unique for prevent duplicate skills or conflicts in the transactions
CREATE TABLE IF NOT EXISTS skills (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at TEXT NOT NULL DEFAULT (datetime('now', 'localtime')) UNIQUE,
    updated_at TEXT NOT NULL DEFAULT (datetime('now', 'localtime')),
    name TEXT NOT NULL UNIQUE, -- skill name, must be unique
    description TEXT NOT NULL, -- skill description
    xp INTEGER NOT NULL DEFAULT 0, -- skill experience points, they are increased by missions, rituals and goals
    next_level_xp INTEGER NOT NULL DEFAULT 50, -- skill next level xp, it is increased by 50 after level up
    level INTEGER NOT NULL DEFAULT 1, -- skill level, it is increased by 1 after level up
    area_id INTEGER NOT NULL, -- area id, the area that the skill belongs to
    FOREIGN KEY (area_id) REFERENCES areas(id)
);

--
-- skills table indexes
--

CREATE INDEX IF NOT EXISTS skills_id_index ON skills (id);
CREATE INDEX IF NOT EXISTS skills_name_index ON skills (name);
CREATE INDEX IF NOT EXISTS skills_level_index ON skills (level);
CREATE INDEX IF NOT EXISTS skills_area_id_index ON skills (area_id);

--
-- skills table triggers
--

-- update skill updateded_at after skill update
CREATE TRIGGER IF NOT EXISTS update_skill_updated_at_after_skill_update
AFTER UPDATE ON skills
FOR EACH ROW
BEGIN
    UPDATE skills
    SET updated_at = datetime('now', 'localtime')
    WHERE id = NEW.id;
END;

--------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------

--
-- goals table
--

-- the goals table is used to store the user goals
-- the goals are used to store the user objectives, the user can complete the goals to gain xp and coins
-- the goals can haves a period, it can be WEEK, MONTH or YEAR
-- the period and the reference_date together refers to the period that the goal is related, for example if the period is WEEK the goal is related to the week of the reference_date
-- the status field is used to store the goal status, it can be NEW, TODO, NEXT, WAITING, DOING, DONE or OVERDUE
-- the created_at field and name field are unique for prevent duplicate goals or conflicts in the transactions
CREATE TABLE IF NOT EXISTS goals (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at TEXT NOT NULL DEFAULT (datetime('now', 'localtime')) UNIQUE,
    updated_at TEXT NOT NULL DEFAULT (datetime('now', 'localtime')),
    reference_date TEXT NOT NULL DEFAULT (datetime('now', 'localtime')),
    period TEXT NOT NULL DEFAULT 'WEEK', -- goal period, it can be WEEK, MONTH, YEAR
    name TEXT NOT NULL UNIQUE, -- goal name, must be unique
    description TEXT NOT NULL, -- goal description
    status TEXT NOT NULL DEFAULT 'NEW', -- goal status, it can be NEW, TODO, NEXT, WAITING, DOING, DONE or OVERDUE
    xp_reward INTEGER NOT NULL DEFAULT 0, -- goal xp reward, it is increased by completing the goal
    coins_reward INTEGER NOT NULL DEFAULT 0, -- goal coins reward, it is increased by completing the goal
    area_id INTEGER NOT NULL, -- area id, the area that the goal belongs to
    FOREIGN KEY (area_id) REFERENCES areas(id)
);

--
-- goals table indexes
--

CREATE INDEX IF NOT EXISTS goals_id_index ON goals (id);
CREATE INDEX IF NOT EXISTS goals_name_index ON goals (name);
CREATE INDEX IF NOT EXISTS goals_reference_date_index ON goals (reference_date);
CREATE INDEX IF NOT EXISTS goals_period_index ON goals (period);
CREATE INDEX IF NOT EXISTS goals_status_index ON goals (status);
CREATE INDEX IF NOT EXISTS goals_area_id_index ON goals (area_id);

--
-- goals table triggers
--

-- update goal updateded_at after goal update
CREATE TRIGGER IF NOT EXISTS update_goal_updated_at_after_goal_update
AFTER UPDATE ON goals
FOR EACH ROW
BEGIN
    UPDATE goals
    SET updated_at = datetime('now', 'localtime')
    WHERE id = NEW.id;
END;

-- update user xp and coins after goal completion
CREATE TRIGGER IF NOT EXISTS update_user_xp_and_coins_after_goal_completion
AFTER UPDATE ON goals
FOR EACH ROW
WHEN NEW.status = 'DONE' AND OLD.status != 'DONE'
BEGIN
    -- update user xp and coins
    UPDATE users
    SET xp = xp + NEW.xp_reward,
        coins = coins + NEW.coins_reward
    WHERE id = NEW.user_id;

    -- update user stats after level up (if xp >= next_level_xp)
    UPDATE users
    SET level = level + 1,
        next_level_xp = next_level_xp + 50,
        max_hp = max_hp + 10,
        hp = hp + 10,
        max_pp = max_pp + 5,
        pp = pp + 5,
        xp = xp - next_level_xp
    WHERE xp >= next_level_xp;
END;

--------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------

--
-- quests table
--

-- the quest table is used to store the user quests
-- a quest is a big mission that can be divided into smaller missions
CREATE TABLE IF NOT EXISTS quests (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at TEXT NOT NULL DEFAULT (datetime('now', 'localtime')) UNIQUE,
    updated_at TEXT NOT NULL DEFAULT (datetime('now', 'localtime')),
    start_date TEXT NOT NULL DEFAULT (datetime('now', 'localtime')), -- quest start date
    due_date TEXT NOT NULL DEFAULT (datetime('now', 'localtime')), -- quest due date
    name TEXT NOT NULL UNIQUE, -- quest name, must be unique
    description TEXT NOT NULL, -- quest description
    status TEXT NOT NULL DEFAULT 'NEW', -- quest status, it can be NEW, TODO, DOING, DONE or OVERDUE
    xp_reward INTEGER NOT NULL DEFAULT 0, -- quest xp reward, it is increased by completing the quest
    karma_reward INTEGER NOT NULL DEFAULT 0, -- quest karma reward, it is increased by completing the quest
    area_id INTEGER NOT NULL, -- area id, the area that the quest belongs to
    FOREIGN KEY (area_id) REFERENCES areas(id)
);

--
-- quests table indexes
--

CREATE INDEX IF NOT EXISTS quests_id_index ON quests (id);
CREATE INDEX IF NOT EXISTS quests_name_index ON quests (name);
CREATE INDEX IF NOT EXISTS quests_start_date_index ON quests (start_date);
CREATE INDEX IF NOT EXISTS quests_due_date_index ON quests (due_date);
CREATE INDEX IF NOT EXISTS quests_status_index ON quests (status);
CREATE INDEX IF NOT EXISTS quests_area_id_index ON quests (area_id);


--
-- quests table triggers
--

-- update quest updateded_at after quest update
CREATE TRIGGER IF NOT EXISTS update_quest_updated_at_after_quest_update
AFTER UPDATE ON quests
FOR EACH ROW
BEGIN
    UPDATE quests
    SET updated_at = datetime('now', 'localtime')
    WHERE id = NEW.id;
END;

-- distribute quest xp_reward to missions after quest insert
CREATE TRIGGER IF NOT EXISTS distribute_quest_xp_reward_to_missions_after_quest_insert
AFTER INSERT ON quests
FOR EACH ROW
BEGIN
    -- update missions xp_reward
    UPDATE missions
    SET xp_reward = CASE
        WHEN (SELECT COUNT(*) FROM missions WHERE quest_id = NEW.id) > 0 THEN ROUND(NEW.xp_reward / (SELECT COUNT(*) FROM missions WHERE quest_id = NEW.id))
        ELSE 0
    END
    WHERE quest_id = NEW.id;
END;

-- distribute quest xp_reward to missions after quest update
CREATE TRIGGER IF NOT EXISTS distribute_quest_xp_reward_to_missions_after_quest_update
AFTER UPDATE ON quests
FOR EACH ROW
BEGIN
    -- update missions xp_reward
    UPDATE missions
    SET xp_reward = CASE
        WHEN (SELECT COUNT(*) FROM missions WHERE quest_id = NEW.id) > 0 THEN ROUND(NEW.xp_reward / (SELECT COUNT(*) FROM missions WHERE quest_id = NEW.id))
        ELSE 0
    END
    WHERE quest_id = NEW.id;
END;

-- add karma to user after quest completion
CREATE TRIGGER IF NOT EXISTS add_karma_to_user_after_quest_completion
AFTER UPDATE ON quests
FOR EACH ROW
WHEN NEW.status = 'DONE' AND OLD.status != 'DONE'
BEGIN
    -- update user karma
    UPDATE users
    SET karma = karma + NEW.karma_reward
    WHERE id = NEW.user_id;
END;

--------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------

--
-- missions table
--

-- the missions table is used to store the user missions
-- the missions are used to store the user activities, the user can complete the missions to gain xp, a mission consumes pp
-- if the completition of a mission cunsume more pp than the user has, the user loose hp as a penalty. the amount of hp lost is equal to the difference between the pp cost and the user pp
-- a mission can be linked to a quest and/or a skill
-- if the mission is linked to a quest, the quest xp_reward is distributed to the missions linked to the quest
-- the mission status can be NEW, TODO, DOING, DONE or OVERDUE
-- if a mission is not completed until the due_date, the mission status is changed to OVERDUE
-- the created_at field is unique for prevent duplicate missions or conflicts in the transactions
CREATE TABLE IF NOT EXISTS missions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at TEXT NOT NULL DEFAULT (datetime('now', 'localtime')) UNIQUE,
    updated_at TEXT NOT NULL DEFAULT (datetime('now', 'localtime')),
    due_date TEXT NOT NULL DEFAULT (datetime('now', 'localtime')), -- mission due date
    name TEXT NOT NULL, -- mission name
    description TEXT NOT NULL, -- mission description
    status TEXT NOT NULL DEFAULT 'NEW', -- mission status, it can be NEW, TODO, DOING, DONE or OVERDUE
    xp_reward INTEGER NOT NULL DEFAULT 0, -- mission xp reward, it is increased by completing the mission
    pp_cost INTEGER NOT NULL DEFAULT 0, -- mission pp cost, it is decreased by completing the mission
    user_id INTEGER NOT NULL DEFAULT 1, -- user id, the user that owns the mission
    quest_id INTEGER, -- quest id, the quest that the mission belongs to
    skill_id INTEGER, -- skill id, the skill that the mission belongs to
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (quest_id) REFERENCES quests(id),
    FOREIGN KEY (skill_id) REFERENCES skills(id)
);

--
-- missions table indexes
--

CREATE INDEX IF NOT EXISTS missions_id_index ON missions (id);
CREATE INDEX IF NOT EXISTS missions_name_index ON missions (name);
CREATE INDEX IF NOT EXISTS missions_due_date_index ON missions (due_date);
CREATE INDEX IF NOT EXISTS missions_status_index ON missions (status);
CREATE INDEX IF NOT EXISTS missions_user_id_index ON missions (user_id);
CREATE INDEX IF NOT EXISTS missions_pp_cost_index ON missions (pp_cost);
CREATE INDEX IF NOT EXISTS missions_quest_id_index ON missions (quest_id);
CREATE INDEX IF NOT EXISTS missions_skill_id_index ON missions (skill_id);

--
-- missions table triggers
--

-- distribute quest xp_reward to missions after mission insert
CREATE TRIGGER IF NOT EXISTS distribute_quest_xp_reward_to_missions_after_mission_insert
AFTER INSERT ON missions
FOR EACH ROW
BEGIN
    -- update missions xp_reward
    UPDATE missions
    SET xp_reward = CASE
        WHEN (SELECT COUNT(*) FROM missions WHERE quest_id = NEW.quest_id) > 0 THEN ROUND((SELECT xp_reward FROM quests WHERE id = NEW.quest_id) / (SELECT COUNT(*) FROM missions WHERE quest_id = NEW.quest_id))
        ELSE 0
    END
    WHERE quest_id = NEW.quest_id;
END;

-- distribute quest xp_reward to missions after mission update
CREATE TRIGGER IF NOT EXISTS distribute_quest_xp_reward_to_missions_after_mission_update
AFTER UPDATE ON missions
FOR EACH ROW
BEGIN
    -- update missions xp_reward
    UPDATE missions
    SET xp_reward = CASE
        WHEN (SELECT COUNT(*) FROM missions WHERE quest_id = NEW.quest_id) > 0 THEN ROUND((SELECT xp_reward FROM quests WHERE id = NEW.quest_id) / (SELECT COUNT(*) FROM missions WHERE quest_id = NEW.quest_id))
        ELSE 0
    END
    WHERE quest_id = NEW.quest_id;
END;

-- distribute quest xp_reward to missions after mission delete
CREATE TRIGGER IF NOT EXISTS distribute_quest_xp_reward_to_missions_after_mission_delete
AFTER DELETE ON missions
FOR EACH ROW
BEGIN
    -- update missions xp_reward
    UPDATE missions
    SET xp_reward = CASE
        WHEN (SELECT COUNT(*) FROM missions WHERE quest_id = OLD.quest_id) > 0 THEN ROUND((SELECT xp_reward FROM quests WHERE id = OLD.quest_id) / (SELECT COUNT(*) FROM missions WHERE quest_id = OLD.quest_id))
        ELSE 0
    END
    WHERE quest_id = OLD.quest_id;
END;

-- update user and skill xp after mission completion
CREATE TRIGGER IF NOT EXISTS update_user_xp_after_mission_completion
AFTER UPDATE ON missions
FOR EACH ROW
WHEN NEW.status = 'DONE' AND OLD.status != 'DONE'
BEGIN
    -- update user xp
    UPDATE users
    SET xp = xp + NEW.xp_reward,
        hp = CASE
            WHEN pp >= NEW.pp_cost THEN hp
            ELSE hp - (NEW.pp_cost - pp)
        END,
        pp = CASE
            WHEN pp >= NEW.pp_cost THEN pp - NEW.pp_cost
            ELSE 0
        END
    WHERE id = NEW.user_id;

    -- update user stats after level up (if xp >= next_level_xp)
    UPDATE users
    SET level = level + 1,
        next_level_xp = next_level_xp + 50,
        max_hp = max_hp + 10,
        hp = hp + 10,
        max_pp = max_pp + 5,
        pp = pp + 5,
        xp = xp - next_level_xp
    WHERE xp >= next_level_xp;

    -- update skill xp
    UPDATE skills
    SET xp = xp + NEW.xp_reward
    WHERE id = NEW.skill_id;

    -- update skill stats after level up (if xp >= next_level_xp)
    UPDATE skills
    SET level = level + 1,
        next_level_xp = next_level_xp + 50,
        xp = xp - next_level_xp
    WHERE xp >= next_level_xp;
END;

--------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------

--
-- rituals table
--

-- the rituals table is used to store the user rituals
-- the rituals are used to store the user routines, the user can complete the rituals to gain xp, a ritual consumes pp
-- the ritual period can be DAY, WEEK, MONTH or YEAR
-- the period is used with the start_date to calculate when to create a new ritual log
-- rituals can be enabled or disabled, if a ritual is disabled the auto ritual log creation is disabled
-- the ritual must be linked to a skill
-- the created_at and name fields are unique for prevent duplicate rituals or conflicts in the transactions
CREATE TABLE IF NOT EXISTS rituals (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at TEXT NOT NULL DEFAULT (datetime('now', 'localtime')) UNIQUE,
    updated_at TEXT NOT NULL DEFAULT (datetime('now', 'localtime')),
    start_date TEXT NOT NULL DEFAULT (datetime('now', 'localtime')), -- ritual start date
    period TEXT NOT NULL DEFAULT 'DAY', -- ritual period, it can be DAY, WEEK, MONTH, YEAR
    enabled INTEGER NOT NULL DEFAULT 1, -- ritual enabled, if 0 the ritual is disabled if 1 the ritual is enabled
    name TEXT NOT NULL UNIQUE, -- ritual name, must be unique
    description TEXT NOT NULL, -- ritual description
    xp_reward INTEGER NOT NULL DEFAULT 0, -- ritual xp reward, it is increased by completing the ritual
    pp_cost INTEGER NOT NULL DEFAULT 0, -- ritual pp cost, it is decreased by completing the ritual
    user_id INTEGER NOT NULL DEFAULT 1, -- user id, the user that owns the ritual
    skill_id INTEGER NOT NULL, -- skill id, the skill that the ritual belongs to
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (skill_id) REFERENCES skills(id)
);

--
-- rituals table indexes
--

CREATE INDEX IF NOT EXISTS rituals_id_index ON rituals (id);
CREATE INDEX IF NOT EXISTS rituals_start_date_index ON rituals (start_date);
CREATE INDEX IF NOT EXISTS rituals_period_index ON rituals (period);
CREATE INDEX IF NOT EXISTS rituals_enabled_index ON rituals (enabled);
CREATE INDEX IF NOT EXISTS rituals_name_index ON rituals (name);
CREATE INDEX IF NOT EXISTS rituals_pp_cost_index ON rituals (pp_cost);
CREATE INDEX IF NOT EXISTS rituals_user_id_index ON rituals (user_id);
CREATE INDEX IF NOT EXISTS rituals_skill_id_index ON rituals (skill_id);

--
-- rituals table triggers
--

-- update ritual updateded_at after ritual update
CREATE TRIGGER IF NOT EXISTS update_ritual_updated_at_after_ritual_update
AFTER UPDATE ON rituals
FOR EACH ROW
BEGIN
    UPDATE rituals
    SET updated_at = datetime('now', 'localtime')
    WHERE id = NEW.id;
END;

--------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------

--
-- rituals_logs table
--

-- the rituals_logs table is used to store the user rituals logs
-- the rituals logs are used to store the user rituals completions
-- the rituals logs are created automatically based on the ritual period and start_date
-- the ritual log status can be 0 (not completed) or 1 (completed)
-- every time a ritual log is updated to completed, the user gains xp increased by the ritual xp_reward
-- every time a ritual log is updated to completed, the user loses pp decreased by the ritual pp_cost
CREATE TABLE IF NOT EXISTS rituals_logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at TEXT NOT NULL DEFAULT (datetime('now', 'localtime')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now', 'localtime')),
    date TEXT NOT NULL DEFAULT (datetime('now', 'localtime')), -- ritual log date
    status INTEGER NOT NULL DEFAULT 0, -- ritual log status, it can be 0 (not completed) or 1 (completed)
    ritual_id INTEGER NOT NULL, -- ritual id, the ritual that the ritual log belongs to
    FOREIGN KEY (ritual_id) REFERENCES rituals(id)
);

--
-- rituals_logs table indexes
--

CREATE INDEX IF NOT EXISTS rituals_logs_id_index ON rituals_logs (id);
CREATE INDEX IF NOT EXISTS rituals_logs_date_index ON rituals_logs (date);
CREATE INDEX IF NOT EXISTS rituals_logs_status_index ON rituals_logs (status);
CREATE INDEX IF NOT EXISTS rituals_logs_ritual_id_index ON rituals_logs (ritual_id);

--
-- rituals_logs table triggers
--

-- update updated_at after ritual log update
CREATE TRIGGER IF NOT EXISTS update_ritual_log_updated_at_after_ritual_log_update
AFTER UPDATE ON rituals_logs
FOR EACH ROW
BEGIN
    UPDATE rituals_logs
    SET updated_at = datetime('now', 'localtime')
    WHERE id = NEW.id;
END;

-- update user and skill xp after ritual_log completion
CREATE TRIGGER IF NOT EXISTS update_user_xp_after_ritual_log_completion
AFTER UPDATE ON rituals_logs
FOR EACH ROW
WHEN NEW.status = 1 AND OLD.status != 1
BEGIN
    -- update user xp
    UPDATE users
    SET xp = xp + (SELECT xp_reward FROM rituals WHERE id = NEW.ritual_id),
        hp = CASE
            WHEN pp >= (SELECT pp_cost FROM rituals WHERE id = NEW.ritual_id) THEN hp
            ELSE hp - ((SELECT pp_cost FROM rituals WHERE id = NEW.ritual_id) - pp)
        END,
        pp = CASE
            WHEN pp >= (SELECT pp_cost FROM rituals WHERE id = NEW.ritual_id) THEN pp - (SELECT pp_cost FROM rituals WHERE id = NEW.ritual_id)
            ELSE 0
        END
    WHERE id = (SELECT user_id FROM rituals WHERE id = NEW.ritual_id);

    -- update user stats after level up (if xp >= next_level_xp)
    UPDATE users
    SET level = level + 1,
        next_level_xp = next_level_xp + 50,
        max_hp = max_hp + 10,
        hp = hp + 10,
        max_pp = max_pp + 5,
        pp = pp + 5,
        xp = xp - next_level_xp
    WHERE xp >= next_level_xp;

    -- update skill xp
    UPDATE skills
    SET xp = xp + (SELECT xp_reward FROM rituals WHERE id = NEW.ritual_id)
    WHERE id = (SELECT skill_id FROM rituals WHERE id = NEW.ritual_id);

    -- update skill stats after level up (if xp >= next_level_xp)
    UPDATE skills
    SET level = level + 1,
        next_level_xp = next_level_xp + 50,
        xp = xp - next_level_xp
    WHERE xp >= next_level_xp;
END;

--------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------

--
-- day_check_ins table
--

-- the day_check_ins table is used to store the user daily check ins
-- the daily check in are part of the journaling process, the user can check in the mood, energy, sleep and main mission of the day
-- the check in is considered to be the first thing that the user does in the day
-- the main_mission_id is the mission that the user wants to complete in the day and is linked to the mission table
-- the created_at field and date field are unique for prevent duplicate check ins or conflicts in the transactions
CREATE TABLE IF NOT EXISTS day_check_ins (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at TEXT NOT NULL DEFAULT (datetime('now', 'localtime')) UNIQUE,
    updated_at TEXT NOT NULL DEFAULT (datetime('now', 'localtime')),
    date TEXT NOT NULL DEFAULT (datetime('now', 'localtime')) UNIQUE, -- check in date
    mood TEXT NOT NULL DEFAULT 'NEUTRAL', -- check in mood, it can be GREAT, GOOD, NEUTRAL, BAD or TERRIBLE
    energy TEXT NOT NULL DEFAULT 'NEUTRAL', -- check in energy, it can be GREAT, GOOD, NEUTRAL, BAD or TERRIBLE
    sleep TEXT NOT NULL DEFAULT 'NEUTRAL', -- check in sleep, it can be GREAT, GOOD, NEUTRAL, BAD or TERRIBLE
    main_mission_id INTEGER NOT NULL, -- main mission id, the mission that the user wants to complete in the day
    user_id INTEGER NOT NULL DEFAULT 1, -- user id, the user that owns the check in
    FOREIGN KEY (main_mission_id) REFERENCES missions(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

--
-- day_check_ins table indexes
--

CREATE INDEX IF NOT EXISTS day_check_ins_id_index ON day_check_ins (id);
CREATE INDEX IF NOT EXISTS day_check_ins_date_index ON day_check_ins (date);
CREATE INDEX IF NOT EXISTS day_check_ins_mood_index ON day_check_ins (mood);
CREATE INDEX IF NOT EXISTS day_check_ins_energy_index ON day_check_ins (energy);
CREATE INDEX IF NOT EXISTS day_check_ins_sleep_index ON day_check_ins (sleep);
CREATE INDEX IF NOT EXISTS day_check_ins_main_mission_id_index ON day_check_ins (main_mission_id);
CREATE INDEX IF NOT EXISTS day_check_ins_user_id_index ON day_check_ins (user_id);

--
-- day_check_ins table triggers
--

-- update check in updateded_at after check in update
CREATE TRIGGER IF NOT EXISTS update_day_check_in_updated_at_after_day_check_in_update
AFTER UPDATE ON day_check_ins
FOR EACH ROW
BEGIN
    UPDATE day_check_ins
    SET updated_at = datetime('now', 'localtime')
    WHERE id = NEW.id;
END;

--------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------

--
-- day_check_outs table
--

-- the day_check_outs table is used to store the user daily check outs
-- the daily check out are part of the journaling process, the user can check out the content, story of the day, gratitude and to improve
-- the check out is considered to be the last thing that the user does in the day
-- the created_at field and date field are unique for prevent duplicate check outs or conflicts in the transactions
CREATE TABLE IF NOT EXISTS day_check_outs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at TEXT NOT NULL DEFAULT (datetime('now', 'localtime')) UNIQUE,
    updated_at TEXT NOT NULL DEFAULT (datetime('now', 'localtime')),
    date TEXT NOT NULL DEFAULT (datetime('now', 'localtime')) UNIQUE, -- check out date
    content TEXT NOT NULL, -- check out content, free text to write about the day
    story_of_the_day TEXT NOT NULL, -- check out story of the day, the most memorable thing that happened in the day
    gratitude TEXT NOT NULL, -- check out gratitude, the things that the user is grateful for
    to_improve TEXT NOT NULL, -- check out to improve, the things that the user wants to improve
    user_id INTEGER NOT NULL DEFAULT 1, -- user id, the user that owns the check out
    FOREIGN KEY (user_id) REFERENCES users(id)
);

--
-- day_check_outs table indexes
--

CREATE INDEX IF NOT EXISTS day_check_outs_id_index ON day_check_outs (id);
CREATE INDEX IF NOT EXISTS day_check_outs_date_index ON day_check_outs (date);
CREATE INDEX IF NOT EXISTS day_check_outs_user_id_index ON day_check_outs (user_id);

--
-- day_check_outs table triggers
--

-- update check out updateded_at after check out update
CREATE TRIGGER IF NOT EXISTS update_day_check_out_updated_at_after_day_check_out_update
AFTER UPDATE ON day_check_outs
FOR EACH ROW
BEGIN
    UPDATE day_check_outs
    SET updated_at = datetime('now', 'localtime')
    WHERE id = NEW.id;
END;

--------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------

--
-- week_check_ins table
--

-- the week_check_ins table is used to store the user weekly check ins
-- the weekly check in are part of the journaling process, the user can check in the scope and main goal of the week
-- the check in is considered to be the first thing that the user does in the week
-- the main_goal_id is the goal that the user wants to complete in the week and is linked to the goal table
-- the created_at field and date field are unique for prevent duplicate check ins or conflicts in the transactions
CREATE TABLE IF NOT EXISTS week_check_ins (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at TEXT NOT NULL DEFAULT (datetime('now', 'localtime')) UNIQUE,
    updated_at TEXT NOT NULL DEFAULT (datetime('now', 'localtime')),
    date TEXT NOT NULL DEFAULT (datetime('now', 'localtime')) UNIQUE, -- check in date
    scope TEXT NOT NULL, -- check in scope, the user scope for the week
    main_goal_id INTEGER NOT NULL, -- main goal id, the goal that the user wants to complete in the week
    user_id INTEGER NOT NULL DEFAULT 1, -- user id, the user that owns the check in
    FOREIGN KEY (main_goal_id) REFERENCES goals(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

--
-- week_check_ins table indexes
--

CREATE INDEX IF NOT EXISTS week_check_ins_id_index ON week_check_ins (id);
CREATE INDEX IF NOT EXISTS week_check_ins_date_index ON week_check_ins (date);
CREATE INDEX IF NOT EXISTS week_check_ins_main_goal_id_index ON week_check_ins (main_goal_id);
CREATE INDEX IF NOT EXISTS week_check_ins_user_id_index ON week_check_ins (user_id);

--
-- week_check_ins table triggers
--

-- update check in updateded_at after check in update
CREATE TRIGGER IF NOT EXISTS update_week_check_in_updated_at_after_week_check_in_update
AFTER UPDATE ON week_check_ins
FOR EACH ROW
BEGIN
    UPDATE week_check_ins
    SET updated_at = datetime('now', 'localtime')
    WHERE id = NEW.id;
END;

--------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------

--
-- week_check_outs table
--

-- the week_check_outs table is used to store the user weekly check outs
-- the weekly check out are part of the journaling process, the user can check out the content, story of the week, wins, fails and to improve
-- the check out is considered to be the last thing that the user does in the week
-- the created_at field and date field are unique for prevent duplicate check outs or conflicts in the transactions
CREATE TABLE IF NOT EXISTS week_check_outs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at TEXT NOT NULL DEFAULT (datetime('now', 'localtime')) UNIQUE,
    updated_at TEXT NOT NULL DEFAULT (datetime('now', 'localtime')),
    date TEXT NOT NULL DEFAULT (datetime('now', 'localtime')) UNIQUE, -- check out date
    content TEXT NOT NULL, -- check out content, free text to write about the week
    story_of_the_week TEXT NOT NULL, -- check out story of the week, the most memorable thing that happened in the week
    wins TEXT NOT NULL, -- check out wins, the things that the user achieved in the week
    fails TEXT NOT NULL, -- check out fails, the things that the user failed in the week
    to_improve TEXT NOT NULL, -- check out to improve, the things that the user wants to improve
    user_id INTEGER NOT NULL DEFAULT 1, -- user id, the user that owns the check out
    FOREIGN KEY (user_id) REFERENCES users(id)
);


--
-- week_check_outs table indexes
--

CREATE INDEX IF NOT EXISTS week_check_outs_id_index ON week_check_outs (id);
CREATE INDEX IF NOT EXISTS week_check_outs_date_index ON week_check_outs (date);
CREATE INDEX IF NOT EXISTS week_check_outs_user_id_index ON week_check_outs (user_id);

--
-- week_check_outs table triggers
--

-- update check out updateded_at after check out update
CREATE TRIGGER IF NOT EXISTS update_week_check_out_updated_at_after_week_check_out_update
AFTER UPDATE ON week_check_outs
FOR EACH ROW
BEGIN
    UPDATE week_check_outs
    SET updated_at = datetime('now', 'localtime')
    WHERE id = NEW.id;
END;

--------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------

--
-- month_check_ins table
--

-- the month_check_ins table is used to store the user monthly check ins
-- the monthly check in are part of the journaling process, the user can check in the scope and main goal of the month
-- the check in is considered to be the first thing that the user does in the month
-- the main_goal_id is the goal that the user wants to complete in the month and is linked to the goal table
-- the created_at field and date field are unique for prevent duplicate check ins or conflicts in the transactions
CREATE TABLE IF NOT EXISTS month_check_ins (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at TEXT NOT NULL DEFAULT (datetime('now', 'localtime')) UNIQUE,
    updated_at TEXT NOT NULL DEFAULT (datetime('now', 'localtime')),
    date TEXT NOT NULL DEFAULT (datetime('now', 'localtime')) UNIQUE, -- check in date
    scope TEXT NOT NULL, -- check in scope, the user scope for the month, the user can define the month theme
    main_goal_id INTEGER NOT NULL, -- main goal id, the goal that the user wants to complete in the month
    user_id INTEGER NOT NULL DEFAULT 1, -- user id, the user that owns the check in
    FOREIGN KEY (main_goal_id) REFERENCES goals(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

--
-- month_check_ins table indexes
--

CREATE INDEX IF NOT EXISTS month_check_ins_id_index ON month_check_ins (id);
CREATE INDEX IF NOT EXISTS month_check_ins_date_index ON month_check_ins (date);
CREATE INDEX IF NOT EXISTS month_check_ins_main_goal_id_index ON month_check_ins (main_goal_id);
CREATE INDEX IF NOT EXISTS month_check_ins_user_id_index ON month_check_ins (user_id);

--
-- month_check_ins table triggers
--

-- update check in updateded_at after check in update
CREATE TRIGGER IF NOT EXISTS update_month_check_in_updated_at_after_month_check_in_update
AFTER UPDATE ON month_check_ins
FOR EACH ROW
BEGIN
    UPDATE month_check_ins
    SET updated_at = datetime('now', 'localtime')
    WHERE id = NEW.id;
END;

--------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------

--
-- month_check_outs table
--

-- the month_check_outs table is used to store the user monthly check outs
-- the monthly check out are part of the journaling process, the user can check out the content, story of the month, wins, fails and to improve
-- the check out is considered to be the last thing that the user does in the month
-- the created_at field and date field are unique for prevent duplicate check outs or conflicts in the transactions
CREATE TABLE IF NOT EXISTS month_check_outs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at TEXT NOT NULL DEFAULT (datetime('now', 'localtime')) UNIQUE,
    updated_at TEXT NOT NULL DEFAULT (datetime('now', 'localtime')),
    date TEXT NOT NULL DEFAULT (datetime('now', 'localtime')) UNIQUE, -- check out date
    content TEXT NOT NULL, -- check out content, free text to write about the month
    story_of_the_month TEXT NOT NULL, -- check out story of the month, the most memorable thing that happened in the month
    wins TEXT NOT NULL, -- check out wins, the things that the user achieved in the month
    fails TEXT NOT NULL, -- check out fails, the things that the user failed in the month
    to_improve TEXT NOT NULL, -- check out to improve, the things that the user wants to improve
    user_id INTEGER NOT NULL DEFAULT 1, -- user id, the user that owns the check out
    FOREIGN KEY (user_id) REFERENCES users(id)
);

--
-- month_check_outs table indexes
--

CREATE INDEX IF NOT EXISTS month_check_outs_id_index ON month_check_outs (id);
CREATE INDEX IF NOT EXISTS month_check_outs_date_index ON month_check_outs (date);
CREATE INDEX IF NOT EXISTS month_check_outs_user_id_index ON month_check_outs (user_id);

--
-- month_check_outs table triggers
--

-- update check out updateded_at after check out update
CREATE TRIGGER IF NOT EXISTS update_month_check_out_updated_at_after_month_check_out_update
AFTER UPDATE ON month_check_outs
FOR EACH ROW
BEGIN
    UPDATE month_check_outs
    SET updated_at = datetime('now', 'localtime')
    WHERE id = NEW.id;
END;

--------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------

--
-- year_check_ins table
--

-- the year_check_ins table is used to store the user yearly check ins
-- the yearly check in are part of the journaling process, the user can check in the scope, who am i and who am i becoming
-- the check in is considered to be the first thing that the user does in the year
-- the main_goal_id is the goal that the user wants to complete in the year and is linked to the goal table
-- the created_at field and date field are unique for prevent duplicate check ins or conflicts in the transactions
CREATE TABLE IF NOT EXISTS year_check_ins (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at TEXT NOT NULL DEFAULT (datetime('now', 'localtime')) UNIQUE,
    updated_at TEXT NOT NULL DEFAULT (datetime('now', 'localtime')),
    date TEXT NOT NULL DEFAULT (datetime('now', 'localtime')) UNIQUE, -- check in date
    scope TEXT NOT NULL, -- check in scope, the user scope for the year, the user can define the year theme
    who_am_i TEXT NOT NULL, -- check in who am i, the user can define who am i in the year
    who_am_i_becoming TEXT NOT NULL, -- check in who am i becoming, the user can define who am i becoming in the year
    main_goal_id INTEGER NOT NULL, -- main goal id, the goal that the user wants to complete in the year
    user_id INTEGER NOT NULL DEFAULT 1, -- user id, the user that owns the check in
    FOREIGN KEY (main_goal_id) REFERENCES goals(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

--
-- year_check_ins table indexes
--

CREATE INDEX IF NOT EXISTS year_check_ins_id_index ON year_check_ins (id);
CREATE INDEX IF NOT EXISTS year_check_ins_date_index ON year_check_ins (date);
CREATE INDEX IF NOT EXISTS year_check_ins_main_goal_id_index ON year_check_ins (main_goal_id);
CREATE INDEX IF NOT EXISTS year_check_ins_user_id_index ON year_check_ins (user_id);

--
-- year_check_ins table triggers
--

-- update check in updateded_at after check in update
CREATE TRIGGER IF NOT EXISTS update_year_check_in_updated_at_after_year_check_in_update
AFTER UPDATE ON year_check_ins
FOR EACH ROW
BEGIN
    UPDATE year_check_ins
    SET updated_at = datetime('now', 'localtime')
    WHERE id = NEW.id;
END;

--------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------

--
-- year_check_outs table
--

-- the year_check_outs table is used to store the user yearly check outs
-- the yearly check out are part of the journaling process, the user can check out the content, story of the year, strong points, weak points, fails and to improve
-- the check out is considered to be the last thing that the user does in the year
-- the created_at field and date field are unique for prevent duplicate check outs or conflicts in the transactions
CREATE TABLE IF NOT EXISTS year_check_outs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at TEXT NOT NULL DEFAULT (datetime('now', 'localtime')) UNIQUE,
    updated_at TEXT NOT NULL DEFAULT (datetime('now', 'localtime')),
    date TEXT NOT NULL DEFAULT (datetime('now', 'localtime')) UNIQUE, -- check out date
    content TEXT NOT NULL, -- check out content, free text to write about the year  
    story_of_the_year TEXT NOT NULL, -- check out story of the year, the most memorable thing that happened in the year
    strong_points TEXT NOT NULL, -- check out strong points, the user strong points in the year
    weak_points TEXT NOT NULL, -- check out weak points, the user weak points in the year
    fails TEXT NOT NULL, -- check out fails, the user fails in the year
    to_improve TEXT NOT NULL, -- check out to improve, the things that the user wants to improve
    user_id INTEGER NOT NULL DEFAULT 1, -- user id, the user that owns the check out
    FOREIGN KEY (user_id) REFERENCES users(id)
);

--
-- year_check_outs table indexes
--

CREATE INDEX IF NOT EXISTS year_check_outs_id_index ON year_check_outs (id);
CREATE INDEX IF NOT EXISTS year_check_outs_date_index ON year_check_outs (date);
CREATE INDEX IF NOT EXISTS year_check_outs_user_id_index ON year_check_outs (user_id);

--
-- year_check_outs table triggers
--

-- update check out updateded_at after check out update
CREATE TRIGGER IF NOT EXISTS update_year_check_out_updated_at_after_year_check_out_update
AFTER UPDATE ON year_check_outs
FOR EACH ROW
BEGIN
    UPDATE year_check_outs
    SET updated_at = datetime('now', 'localtime')
    WHERE id = NEW.id;
END;

--------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------

--
-- journal_entries table
--

-- the journal_entries table is used to store the user journal entries
-- the journal entries are used to store the user thoughts, ideas, concepts, knowledge, experiences, feelings, emotions, reflections, insights, learnings, quotes, references, resources, references or reflections
-- everyjournal entry is linked to a check in or check out and can be linked to multiple notes.
-- every journal entry can be of type day, week, month or year
CREATE TABLE IF NOT EXISTS journal_entries (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at TEXT NOT NULL DEFAULT (datetime('now', 'localtime')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now', 'localtime')),
    date TEXT NOT NULL DEFAULT (datetime('now', 'localtime')), -- journal entry date
    type TEXT NOT NULL DEFAULT 'DAY', -- journal entry type, it can be DAY, WEEK, MONTH or YEAR
    content TEXT NOT NULL, -- journal entry content
    day_check_in_id INTEGER, -- check in id, the day check in that the journal entry belongs to
    day_check_out_id INTEGER, -- check out id, the check out that the journal entry belongs to
    week_check_in_id INTEGER, -- check in id, the week check in that the journal entry belongs to
    week_check_out_id INTEGER, -- check out id, the week check out that the journal entry belongs to
    month_check_in_id INTEGER, -- check in id, the month check in that the journal entry belongs to
    month_check_out_id INTEGER, -- check out id, the month check out that the journal entry belongs to
    year_check_in_id INTEGER, -- check in id, the year check in that the journal entry belongs to
    year_check_out_id INTEGER, -- check out id, the year check out that the journal entry belongs to
    user_id INTEGER NOT NULL DEFAULT 1, -- user id, the user that owns the journal entry
    FOREIGN KEY (day_check_in_id) REFERENCES day_check_ins(id),
    FOREIGN KEY (day_check_out_id) REFERENCES day_check_outs(id),
    FOREIGN KEY (week_check_in_id) REFERENCES week_check_ins(id),
    FOREIGN KEY (week_check_out_id) REFERENCES week_check_outs(id),
    FOREIGN KEY (month_check_in_id) REFERENCES month_check_ins(id),
    FOREIGN KEY (month_check_out_id) REFERENCES month_check_outs(id),
    FOREIGN KEY (year_check_in_id) REFERENCES year_check_ins(id),
    FOREIGN KEY (year_check_out_id) REFERENCES year_check_outs(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

--
-- journal_entries table indexes
--

CREATE INDEX IF NOT EXISTS journal_entries_id_index ON journal_entries (id);
CREATE INDEX IF NOT EXISTS journal_entries_date_index ON journal_entries (date);
CREATE INDEX IF NOT EXISTS journal_entries_type_index ON journal_entries (type);
CREATE INDEX IF NOT EXISTS journal_entries_day_check_in_id_index ON journal_entries (day_check_in_id);
CREATE INDEX IF NOT EXISTS journal_entries_day_check_out_id_index ON journal_entries (day_check_out_id);
CREATE INDEX IF NOT EXISTS journal_entries_week_check_in_id_index ON journal_entries (week_check_in_id);
CREATE INDEX IF NOT EXISTS journal_entries_week_check_out_id_index ON journal_entries (week_check_out_id);
CREATE INDEX IF NOT EXISTS journal_entries_month_check_in_id_index ON journal_entries (month_check_in_id);
CREATE INDEX IF NOT EXISTS journal_entries_month_check_out_id_index ON journal_entries (month_check_out_id);
CREATE INDEX IF NOT EXISTS journal_entries_year_check_in_id_index ON journal_entries (year_check_in_id);
CREATE INDEX IF NOT EXISTS journal_entries_year_check_out_id_index ON journal_entries (year_check_out_id);
CREATE INDEX IF NOT EXISTS journal_entries_user_id_index ON journal_entries (user_id);

--
-- journal_entries table triggers
--

-- update journal entry updateded_at after journal entry update
CREATE TRIGGER IF NOT EXISTS update_journal_entry_updated_at_after_journal_entry_update
AFTER UPDATE ON journal_entries
FOR EACH ROW
BEGIN
    UPDATE journal_entries
    SET updated_at = datetime('now', 'localtime')
    WHERE id = NEW.id;
END;

--------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------

--
-- arguments table
--

-- the arguments table is used to store the user arguments
-- an argument is a topic, asubject or theme that the user wants to study, learn or discuss
-- the argument are used to categorize the user knowledge
-- the argument are linked to areas and can be linked to single or multiple skills
-- the created_at field and name field are unique for prevent duplicate arguments or conflicts in the transactions
CREATE TABLE IF NOT EXISTS arguments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at TEXT NOT NULL DEFAULT (datetime('now', 'localtime')) UNIQUE,
    updated_at TEXT NOT NULL DEFAULT (datetime('now', 'localtime')),
    name TEXT NOT NULL UNIQUE, -- argument name, must be unique
    description TEXT NOT NULL, -- argument description
    area_id INTEGER NOT NULL, -- area id, the area that the argument belongs to
    FOREIGN KEY (area_id) REFERENCES areas(id)
);

--
-- arguments table indexes
--

CREATE INDEX IF NOT EXISTS arguments_id_index ON arguments (id);
CREATE INDEX IF NOT EXISTS arguments_name_index ON arguments (name);
CREATE INDEX IF NOT EXISTS arguments_area_id_index ON arguments (area_id);

--
-- arguments table triggers
--

-- update argument updateded_at after argument update
CREATE TRIGGER IF NOT EXISTS update_argument_updated_at_after_argument_update
AFTER UPDATE ON arguments
FOR EACH ROW
BEGIN
    UPDATE arguments
    SET updated_at = datetime('now', 'localtime')
    WHERE id = NEW.id;
END;

--------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------

--
-- arguments_skills table
--

-- the arguments_skills table is used to store the user arguments skills
-- the arguments skills are used to link the arguments to the skills
CREATE TABLE IF NOT EXISTS arguments_skills (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at TEXT NOT NULL DEFAULT (datetime('now', 'localtime')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now', 'localtime')),
    argument_id INTEGER NOT NULL, -- argument id, the argument that the argument skill belongs to
    skill_id INTEGER NOT NULL, -- skill id, the skill that the argument skill belongs to
    FOREIGN KEY (argument_id) REFERENCES arguments(id),
    FOREIGN KEY (skill_id) REFERENCES skills(id)
);

--
-- arguments_skills table indexes
--

CREATE INDEX IF NOT EXISTS arguments_skills_id_index ON arguments_skills (id);
CREATE INDEX IF NOT EXISTS arguments_skills_argument_id_index ON arguments_skills (argument_id);
CREATE INDEX IF NOT EXISTS arguments_skills_skill_id_index ON arguments_skills (skill_id);

--
-- arguments_skills table triggers
--

-- update argument skill updateded_at after argument skill update
CREATE TRIGGER IF NOT EXISTS update_argument_skill_updated_at_after_argument_skill_update
AFTER UPDATE ON arguments_skills
FOR EACH ROW
BEGIN
    UPDATE arguments_skills
    SET updated_at = datetime('now', 'localtime')
    WHERE id = NEW.id;
END;

--------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------

--
-- notes table
--

-- the notes table is used to store the user notes
-- the notes are used to store the user thoughts, ideas, concepts, knowledge, experiences, feelings, emotions, reflections, insights, learnings, quotes, references, resources, references or reflections
-- the notes are linked to argumetns
-- the created_at and title fields are unique for prevent duplicate notes or conflicts in the transactions
CREATE TABLE IF NOT EXISTS notes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at TEXT NOT NULL DEFAULT (datetime('now', 'localtime')) UNIQUE,
    updated_at TEXT NOT NULL DEFAULT (datetime('now', 'localtime')),
    title TEXT NOT NULL UNIQUE, -- note title, must be unique
    content TEXT NOT NULL, -- note content
    journal_entry_id INTEGER, -- journal entry id, the journal entry that the note belongs to
    argument_id INTEGER NOT NULL, -- argument id, the argument that the note belongs to
    FOREIGN KEY (journal_entry_id) REFERENCES journal_entries(id),
    FOREIGN KEY (argument_id) REFERENCES arguments(id)
);

--
-- notes table indexes
--

CREATE INDEX IF NOT EXISTS notes_id_index ON notes (id);
CREATE INDEX IF NOT EXISTS notes_title_index ON notes (title);
CREATE INDEX IF NOT EXISTS notes_journal_entry_id_index ON notes (journal_entry_id);
CREATE INDEX IF NOT EXISTS notes_argument_id_index ON notes (argument_id);

--
-- notes table triggers
--

-- update note updateded_at after note update
CREATE TRIGGER IF NOT EXISTS update_note_updated_at_after_note_update
AFTER UPDATE ON notes
FOR EACH ROW
BEGIN
    UPDATE notes
    SET updated_at = datetime('now', 'localtime')
    WHERE id = NEW.id;
END;

--------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------

--
-- resources table
--

-- the resources table is used to store the user resources
-- a resource is a link to external content, like articles, books, courses, podcasts, videos, websites or any other type of content
-- resources are used to store the user references, resources, links, sources or tools
-- the resources are linked to an argument anc can be linked to multiple notes
-- the created_at, title and url fields are unique for prevent duplicate resources or conflicts in the transactions
CREATE TABLE IF NOT EXISTS resources (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at TEXT NOT NULL DEFAULT (datetime('now', 'localtime')) UNIQUE,
    updated_at TEXT NOT NULL DEFAULT (datetime('now', 'localtime')),
    title TEXT NOT NULL UNIQUE, -- resource title, must be unique
    description TEXT NOT NULL, -- resource description
    url TEXT NOT NULL UNIQUE, -- resource url
    argument_id INTEGER NOT NULL, -- argument id, the argument that the resource belongs to
    FOREIGN KEY (argument_id) REFERENCES arguments(id)
);

--
-- resources table indexes
--

CREATE INDEX IF NOT EXISTS resources_id_index ON resources (id);
CREATE INDEX IF NOT EXISTS resources_title_index ON resources (title);
CREATE INDEX IF NOT EXISTS resources_url_index ON resources (url);
CREATE INDEX IF NOT EXISTS resources_argument_id_index ON resources (argument_id);

--
-- resources table triggers
--

-- update resource updateded_at after resource update
CREATE TRIGGER IF NOT EXISTS update_resource_updated_at_after_resource_update
AFTER UPDATE ON resources
FOR EACH ROW
BEGIN
    UPDATE resources
    SET updated_at = datetime('now', 'localtime')
    WHERE id = NEW.id;
END;

--------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------

--
-- notes_resources table
--

-- the notes_resources table is used to store the user notes resources
-- the notes resources are used to link the notes to the resources
CREATE TABLE IF NOT EXISTS notes_resources (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at TEXT NOT NULL DEFAULT (datetime('now', 'localtime')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now', 'localtime')),
    note_id INTEGER NOT NULL, -- note id, the note that the note resource belongs to
    resource_id INTEGER NOT NULL, -- resource id, the resource that the note resource belongs to
    FOREIGN KEY (note_id) REFERENCES notes(id),
    FOREIGN KEY (resource_id) REFERENCES resources(id)
);

--
-- notes_resources table indexes
--

CREATE INDEX IF NOT EXISTS notes_resources_id_index ON notes_resources (id);
CREATE INDEX IF NOT EXISTS notes_resources_note_id_index ON notes_resources (note_id);
CREATE INDEX IF NOT EXISTS notes_resources_resource_id_index ON notes_resources (resource_id);

--
-- notes_resources table triggers
--

-- update note resource updateded_at after note resource update
CREATE TRIGGER IF NOT EXISTS update_note_resource_updated_at_after_note_resource_update
AFTER UPDATE ON notes_resources
FOR EACH ROW
BEGIN
    UPDATE notes_resources
    SET updated_at = datetime('now', 'localtime')
    WHERE id = NEW.id;
END;

--------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------

--
-- attachments table
--

-- the attachments table is used to store the user attachments
-- an attachment is a file, image, video, audio or any other type of content that the user wants to store
-- the attachments are linked to an argument and can be linked to multiple notes
-- the created_at, title fields are unique for prevent duplicate attachments or conflicts in the transactions
CREATE TABLE IF NOT EXISTS attachments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at TEXT NOT NULL DEFAULT (datetime('now', 'localtime')) UNIQUE,
    updated_at TEXT NOT NULL DEFAULT (datetime('now', 'localtime')),
    title TEXT NOT NULL UNIQUE, -- attachment title, must be unique
    description TEXT NOT NULL, -- attachment description
    file_name TEXT NOT NULL, -- attachment file name
    content BLOB NOT NULL, -- attachment content in binary format
    mime_type TEXT NOT NULL, -- attachment mime type
    argument_id INTEGER NOT NULL, -- argument id, the argument that the attachment belongs to
    FOREIGN KEY (argument_id) REFERENCES arguments(id)
);

--
-- attachments table indexes
--

CREATE INDEX IF NOT EXISTS attachments_id_index ON attachments (id);
CREATE INDEX IF NOT EXISTS attachments_title_index ON attachments (title);
CREATE INDEX IF NOT EXISTS attachments_file_name_index ON attachments (file_name);
CREATE INDEX IF NOT EXISTS attachments_mime_type_index ON attachments (mime_type);
CREATE INDEX IF NOT EXISTS attachments_argument_id_index ON attachments (argument_id);

--
-- attachments table triggers
--

-- update attachment updateded_at after attachment update
CREATE TRIGGER IF NOT EXISTS update_attachment_updated_at_after_attachment_update
AFTER UPDATE ON attachments
FOR EACH ROW
BEGIN
    UPDATE attachments
    SET updated_at = datetime('now', 'localtime')
    WHERE id = NEW.id;
END;

--------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------

--
-- notes_attachments table
--

-- the notes_attachments table is used to store the user notes attachments
-- the notes attachments are used to link the notes to the attachments
CREATE TABLE IF NOT EXISTS notes_attachments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at TEXT NOT NULL DEFAULT (datetime('now', 'localtime')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now', 'localtime')),
    note_id INTEGER NOT NULL, -- note id, the note that the note attachment belongs to
    attachment_id INTEGER NOT NULL, -- attachment id, the attachment that the note attachment belongs to
    FOREIGN KEY (note_id) REFERENCES notes(id),
    FOREIGN KEY (attachment_id) REFERENCES attachments(id)
);

--
-- notes_attachments table indexes
--

CREATE INDEX IF NOT EXISTS notes_attachments_id_index ON notes_attachments (id);
CREATE INDEX IF NOT EXISTS notes_attachments_note_id_index ON notes_attachments (note_id);
CREATE INDEX IF NOT EXISTS notes_attachments_attachment_id_index ON notes_attachments (attachment_id);

--
-- notes_attachments table triggers
--

-- update note attachment updateded_at after note attachment update
CREATE TRIGGER IF NOT EXISTS update_note_attachment_updated_at_after_note_attachment_update
AFTER UPDATE ON notes_attachments
FOR EACH ROW
BEGIN
    UPDATE notes_attachments
    SET updated_at = datetime('now', 'localtime')
    WHERE id = NEW.id;
END;

--------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------

