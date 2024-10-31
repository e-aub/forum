package models

var UsersTable = `CREATE TABLE IF NOT EXISTS users (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	username TEXT NOT NULL UNIQUE,
	email TEXT NOT NULL UNIQUE,
	password TEXT NOT NULL,
	session
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
);`

var SessionsTable = `CREATE TABLE sessions (
    session_id TEXT PRIMARY KEY,      -- UUID if enabled, otherwise a unique identifier (could be INT if UUID is not used)
    user_id INTEGER NOT NULL,         -- Foreign key referencing users table
    expires_at DATETIME NOT NULL,     -- Expiration datetime for the session
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,  -- Timestamp of session creation
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
`
