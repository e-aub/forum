package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	database "forum/internal/database/models"
	utils "forum/internal/utils"
)

func CreateDatabase(dbPath string) *sql.DB {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("%sError opening database:%s%s\n", utils.Colors["red"], err.Error(), utils.Colors["reset"])
	}
	// Verify the connection
	if err = db.Ping(); err != nil {
		log.Fatalf("%sError accessing database: %s%s\n", utils.Colors["red"], err.Error(), utils.Colors["reset"])
	} else {
		fmt.Printf("%sDatabase created/opened successfully!%s\n", utils.Colors["green"], utils.Colors["reset"])
	}

	_, err = db.Exec(`PRAGMA foreign_keys=ON;`)
	if err != nil {
		log.Fatalf("%sError enabling foreign keys: %s%s\n", utils.Colors["red"], err.Error(), utils.Colors["reset"])
	}
	return db
}

func CreateTables(db *sql.DB) {
	_, err := db.Exec(database.UsersTable + database.SessionsTable + database.LikesTable +
		database.CommentsTable + database.PostsTable + database.CategoriesTable + database.PostCategoriesTable)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Created all tables succesfully")
}

func Insert_Post(p *utils.Posts) {
	log.Printf("yess")
	file, err := sql.Open("sqlite3", "db/data.db")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	statement, err := file.Prepare(`INSERT INTO posts(user_id ,title,content,created_at) Values (?,?,?,?)`)
	if err != nil {
		panic(err)
	}
	_, err = statement.Exec(p.UserId, p.Title, p.Content, p.Created_At)
	if err != nil {
		panic(err)
	}
}

func Update_Post(p *utils.Posts) {
	file, err := sql.Open("sqlite3", "db/data.db")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	statement, err := file.Prepare(`UPDATE posts
	SET title=?,
	content = ?,
	updated_at = ?
	WHERE
	id= ?`)
	if err != nil {
		log.Fatal(err)
	}
	_, err = statement.Exec(p.Title, p.Content, p.Updated_At, p.PostId)
	if err != nil {
		log.Fatal(err)
	}
}

func Delete_Post(p *utils.Posts) {
	file, err := sql.Open("sqlite3", "db/data.db")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	statement, err := file.Prepare(`DELETE FROM posts WHERE id = ?`)
	if err != nil {
		log.Fatal(err)
	}
	_, err = statement.Exec(p.PostId)
	if err != nil {
		log.Fatal(err)
	}
}

func Read_Post(id int) *utils.Posts {
	file, err := sql.Open("sqlite3", "db/data.db")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	query := `SELECT * FROM posts WHERE id = ?`
	row := file.QueryRow(query, id)
	Post := &utils.Posts{}
	_ = row.Scan(&Post.PostId, &Post.UserId, &Post.Title, &Post.Content, &Post.Created_At)
	return Post
}

func Get_Last() int {
	file, err := sql.Open("sqlite3", "db/data.db")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	query := `SELECT MAX(id) FROM posts `
	row := file.QueryRow(query)
	result := 0
	_ = row.Scan(&result)
	return result
}

func Get_session(ses string) (float64, error) {
	var sessionid float64
	file, err := sql.Open("sqlite3", "db/data.db")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	query := `SELECT user_id FROM sessions WHERE session_id = ?`
	err = file.QueryRow(query, ses).Scan(&sessionid)
	if err != nil {
		return 0, err
	}
	return sessionid, nil
}

func GetUserIDByUsername(db *sql.DB, username string) (int, error) {
	var userID int
	err := db.QueryRow("SELECT id FROM users WHERE username = ?", username).Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

func InsertSession(db *sql.DB, sessionID string, userID int, expiration time.Time) error {
	_, err := db.Exec("INSERT INTO sessions (session_id, user_id, expires_at) VALUES (?, ?, ?)", sessionID, userID, expiration)
	return err
}
