package database

var LikesTable = `CREATE TABLE IF NOT EXISTS likes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    post_id INTEGER,
    comment_id INTEGER,
    target_type TEXT CHECK(target_type IN ('post', 'comment')) NOT NULL,
    type TEXT CHECK(type IN ('like', 'dislike')) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
    FOREIGN KEY (comment_id) REFERENCES comments(id) ON DELETE CASCADE,
    CHECK (
        (target_type = 'post' AND post_id IS NOT NULL AND comment_id IS NULL) OR 
        (target_type = 'comment' AND comment_id IS NOT NULL AND post_id IS NULL)
    )
);
`
