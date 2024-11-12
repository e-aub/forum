package database

var PostsTable = `CREATE TABLE IF NOT EXISTS posts (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        user_id INTEGER NOT NULL,
        title TEXT NOT NULL,
        content TEXT NOT NULL,
        like_count INTEGER DEFAULT 0,
        dislike_count INTEGER DEFAULT 0,
        created_at DATETIME DEFAULT (strftime('%Y-%m-%d %H:%M:%S', 'now')),
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
    );
	
	`
