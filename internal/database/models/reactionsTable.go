package database

var ReactionTable = `CREATE TABLE IF NOT EXISTS reactions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    type_id TEXT NOT NULL,
    user_id INTEGER NOT NULL,
    post_id INTEGER,
    comment_id INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (user_id, post_id),
    UNIQUE (user_id, comment_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
    FOREIGN KEY (type_id) REFERENCES reaction_type(reaction_id) ON DELETE CASCADE,
    FOREIGN KEY (comment_id) REFERENCES comments(id) ON DELETE CASCADE,
    CHECK ((post_id IS NOT NULL AND comment_id IS NULL) OR (post_id IS NULL AND comment_id IS NOT NULL))
);
`

var ReactionsTypeTable = `
    CREATE TABLE IF NOT EXISTS reaction_type (
        reaction_id TEXT PRIMARY KEY,
        name TEXT NOT NULL UNIQUE
    );
    INSERT OR IGNORE INTO reaction_type (reaction_id, name) VALUES
        ('like', 'Like'), 
        ('dislike', 'Dislike'), 
        ('love', 'Love'), 
        ('haha', 'Haha'), 
        ('wow', 'Wow'), 
        ('sad', 'Sad'), 
        ('angry', 'Angry');
`
