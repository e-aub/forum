package database

var LikesTable = `CREATE TABLE IF NOT EXISTS likes (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        user_id INTEGER NOT NULL,
        post_id INTEGER NOT NULL,
        target_type TEXT CHECK(target_type IN ('post', 'comment')) NOT NULL,
        type TEXT CHECK(type IN ('like', 'dislike')) NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
        FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE
    );
	
	`
