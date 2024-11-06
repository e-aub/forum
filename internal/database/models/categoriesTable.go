package database

type Category struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

var CategoriesTable = `
    CREATE TABLE IF NOT EXISTS categories (
        id VARCHAR(30) PRIMARY KEY NOT NULL UNIQUE,
        name TEXT NOT NULL UNIQUE,
        description TEXT NOT NULL
    );
    INSERT OR IGNORE INTO categories (id, name, description) VALUES
        ('liked_posts', 'Liked Posts', 'This category contains all posts you liked'), 
        ('created_posts', 'Created Posts', 'This category contains all posts you created'), 
        ('music', 'Music', 'Discuss everything related to music, including genres, artists, and concerts'), 
        ('sports', 'Sports', 'Talk about all types of sports, games, and tournaments'), 
        ('movies_tv_shows', 'Movies & TV Shows', 'Share recommendations and discuss your favorite films and series'), 
        ('technology', 'Technology', 'Discuss the latest trends in tech, gadgets, and software'), 
        ('gaming', 'Gaming', 'A place for gamers to discuss games, consoles, and tips'), 
        ('books_literature', 'Books & Literature', 'Share and discover books, authors, and literary genres'), 
        ('travel', 'Travel', 'Exchange travel tips, favorite destinations, and experiences'), 
        ('food_cooking', 'Food & Cooking', 'Discuss recipes, restaurants, and all things culinary');
`

var PostCategoriesTable = `
    CREATE TABLE IF NOT EXISTS post_categories (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        post_id INTEGER NOT NULL,
        category_id VARCHAR(30) NOT NULL,
        FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
        FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE
    );
`
