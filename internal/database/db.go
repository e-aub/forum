package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
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

func CleanupExpiredSessions(db *sql.DB) {
	_, err := db.Exec("DELETE FROM sessions WHERE  expires_at < ?", time.Now())
	if err != nil {
		log.Printf("Error cleaning up expired sessions: %v", err)
	}
}

func Insert_Post(p *utils.Posts, db *sql.DB) (int64, error) {
	transaction, err := db.Begin()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error starting transaction:", err)
		return 0, err
	}
	stmt, err := transaction.Prepare(`
	INSERT INTO posts(user_id ,title,content) Values (?,?,?);
	`)
	if err != nil {
		transaction.Rollback()
		fmt.Fprintln(os.Stderr, "Error Adding post:", err)
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(p.UserId, p.Title, p.Content)
	if err != nil {
		transaction.Rollback()
		fmt.Fprintln(os.Stderr, "Error Adding post:", err)
		return 0, err
	}
	lastPostID, err := result.LastInsertId()
	if err != nil {
		transaction.Rollback()
		fmt.Fprintln(os.Stderr, "error in assigning category to post")
		return 0, err
	}

	stmt1, err := transaction.Prepare(`INSERT INTO post_categories(category_id, post_id) VALUES(?, ?);`)
	if err != nil {
		transaction.Rollback()
		fmt.Fprintln(os.Stderr, "error in assigning category to post")
		return 0, err
	}
	defer stmt1.Close()

	_, err = stmt1.Exec(2, lastPostID)
	if err != nil {
		transaction.Rollback()
		fmt.Fprintln(os.Stderr, "error in assigning category to post")
		transaction.Rollback()
		return 0, err
	}
	err = transaction.Commit()
	if err != nil {
		transaction.Rollback()
		fmt.Fprintln(os.Stderr, "transaction aborted")
		return 0, err

	}
	return lastPostID, nil
}

func Update_Post(p *utils.Posts, db *sql.DB) {
	statement, err := db.Prepare(`UPDATE posts
	SET title=?,
	content = ?,
	updated_at = ?
	WHERE
	id= ?`)
	if err != nil {
		log.Fatal(err)
	}
	_, err = statement.Exec(p.Title, p.Content, p.PostId)
	if err != nil {
		log.Fatal(err)
	}
}

func Delete_Post(p *utils.Posts, db *sql.DB) {
	statement, err := db.Prepare(`DELETE FROM posts WHERE id = ?`)
	if err != nil {
		log.Fatal(err)
	}
	_, err = statement.Exec(p.PostId)
	if err != nil {
		log.Fatal(err)
	}
}

func Read_Post(id int, db *sql.DB) *utils.Posts {
	query := `SELECT * FROM posts WHERE id = ?`
	row := db.QueryRow(query, id)
	Post := &utils.Posts{}
	err := row.Scan(&Post.PostId, &Post.UserId, &Post.Title, &Post.Content, &Post.Created_At)
	if err != nil {
		fmt.Println(err)
	}
	Post.UserName, err = GetUserName(int(Post.UserId), db)
	if err != nil {
		fmt.Println(err)
	}
	return Post
}

func Get_Last(db *sql.DB) int {
	query := `SELECT MAX(id) FROM posts `
	row := db.QueryRow(query)
	result := 0
	_ = row.Scan(&result)
	return result
}

func Get_session(ses string, db *sql.DB) (int, error) {
	var sessionid int
	query := `SELECT user_id FROM sessions WHERE session_id = ?`
	err := db.QueryRow(query, ses).Scan(&sessionid)
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

func CreateComment(c *utils.Comment, db *sql.DB) error {
	if c.Content == "" || c.Post_id == 0 || c.User_id == 0 {
		return errors.New("comment DATA issue")
	}

	query := `
	INSERT INTO comments (user_id, post_id, content, created_at)
	VALUES (?, ?, ?, ?)
	`
	result, err := db.Exec(query, c.User_id, c.Post_id, c.Content, c.Created_at)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	c.Comment_id = int(id)

	return nil
}

func GetComments(postID int, db *sql.DB) ([]utils.Comment, error) {
	query := `
	SELECT comments.id, comments.content, comments.created_at, users.username FROM comments
    INNER JOIN users ON comments.user_id = users.id
    WHERE comments.post_id = ?
	ORDER BY comments.created_at DESC;
	`
	rows, err := db.Query(query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []utils.Comment
	for rows.Next() {
		var comment utils.Comment
		err := rows.Scan(&comment.Comment_id, &comment.Content, &comment.Created_at, &comment.User_name)
		if err != nil {
			return nil, err
		}
		comment.Post_id = postID
		comments = append(comments, comment)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return comments, nil
}

func GetCategoryContent(db *sql.DB, categoryId string) ([]utils.Posts, error) {
	stmt, err := db.Prepare(`SELECT posts.*
	FROM post_categories
	JOIN posts ON post_categories.post_id = posts.id
	WHERE post_categories.category_id = ?
	`)
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query(categoryId)
	if err != nil {
		return nil, err
	}
	var res []utils.Posts
	for rows.Next() {
		var post utils.Posts
		err := rows.Scan(&post.PostId, &post.UserId, &post.Title, &post.Content, &post.Created_At)
		if err != nil {
			return nil, err
		}
		res = append(res, post)
	}
	return res, nil
}

func GetCategoryContentIds(db *sql.DB, categoryId string) ([]int, error) {
	stmt, err := db.Prepare("SELECT post_id FROM post_categories WHERE category_id=?")
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query(categoryId)
	if err != nil {
		return nil, err
	}
	var ids []int
	for rows.Next() {
		tmp := 0
		err := rows.Scan(&tmp)
		if err != nil {
			return nil, err
		}
		ids = append(ids, tmp)
	}
	return ids, nil
}

func GetUserName(id int, db *sql.DB) (string, error) {
	var name string
	err := db.QueryRow("SELECT username FROM users WHERE id = ?", id).Scan(&name)
	if err != nil {
		return "", err
	}
	return name, nil
}

func Get_CategoryofPost(idPost int, db *sql.DB) (string, error) {
	var name string
	stmt, err := db.Prepare(`SELECT categories.name 
							 FROM categories 
							 JOIN post_categories ON post_categories.category_id = categories.id 
							 WHERE post_categories.post_id = ?`)
	if err != nil {
		return "", err
	}
	defer stmt.Close()

	err = stmt.QueryRow(idPost).Scan(&name)
	if err != nil {
		return "", err
	}
	return name, nil
}

func LinkPostWithCategory(db *sql.DB, category string, postId int64) error {
	stmt, err := db.Prepare(`INSERT INTO post_categories(post_id, category_id) VALUES(?, ?)`)
	if err != nil {
		return err
	}
	tmp, _ := strconv.Atoi(category)
	_, err = stmt.Exec(postId, tmp)
	if err != nil {
		return err
	}
	return nil
}
