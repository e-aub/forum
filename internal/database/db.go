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
	_, err := db.Exec(database.UsersTable + database.SessionsTable + database.ReactionTable + database.ReactionsTypeTable +
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

func InsertPost(p *utils.Post, db *sql.DB, categories []string) (int64, error) {
	transaction, err := db.Begin()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error starting transaction:", err)
		return 0, err
	}
	stmt, err := transaction.Prepare(`INSERT INTO posts(user_id ,title,content) Values (?,?,?);`)
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
		fmt.Fprintln(os.Stderr, "error in assigning category to post", err)
		return 0, err
	}

	categories = append(categories, "2")
	err = LinkPostWithCategory(transaction, categories, lastPostID, p.UserId)
	if err != nil {
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

func UpdatePost(p *utils.Post, db *sql.DB) {
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

func DeletePost(p *utils.Post, db *sql.DB) {
	statement, err := db.Prepare(`DELETE FROM posts WHERE id = ?`)
	if err != nil {
		log.Fatal(err)
	}
	_, err = statement.Exec(p.PostId)
	if err != nil {
		log.Fatal(err)
	}
}

func ReadPost(db *sql.DB, userId int, postId int) (*utils.Post, error) {
	query := `SELECT * FROM posts WHERE id = ?`
	row, err := utils.QueryRow(db, query, postId)
	if err != nil {
		return nil, err
	}
	Post := &utils.Post{}
	err = row.Scan(&Post.PostId, &Post.UserId, &Post.Title, &Post.Content, &Post.Created_At)
	if err != nil {
		return nil, err
	}

	Post.UserName, err = GetUserName(int(Post.UserId), db)
	if err != nil {
		return nil, err
	}
	fmt.Println(Post)
	return Post, nil
}
func isLiked(db *sql.DB, userId int, postId int, target_type string) (bool, bool) {
	// Query to check for likes or dislikes for the given user and post.
	var query string
	if target_type == "post" {
		query = `SELECT type FROM likes WHERE user_id = ? AND post_id = ? AND target_type = ? LIMIT 1`
	} else if target_type == "comment" {
		query = `SELECT type FROM likes WHERE user_id = ? AND comment_id = ? AND target_type = ? LIMIT 1`
	}

	var reactionType string
	err := db.QueryRow(query, userId, postId, target_type).Scan(&reactionType)
	if err != nil {
		if err == sql.ErrNoRows {
			// No interaction found.
			fmt.Println("\033[31m", err, "\033[0m")
			return false, false
		}
		// Log error if needed and handle it appropriately.
		fmt.Println("Error   likes:", err)
		return false, false
	}

	// Determine the type of interaction.
	switch reactionType {
	case "like":
		return true, false
	case "dislike":
		return false, true
	default:
		return false, false
	}
}

func GetLastPostId(db *sql.DB) (int, error) {
	query := `SELECT MAX(id) FROM posts `
	row, err := utils.QueryRow(db, query)
	if err != nil {
		return 0, err
	}
	result := 0
	err = row.Scan(&result)
	if err != nil {
		return 0, err
	}
	return result, nil
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

func GetUserIDByUsername(db *sql.DB, userData *utils.User) error {
	// var userID int64
	err := db.QueryRow("SELECT id FROM users WHERE username = ?", userData.UserName).Scan(userData.UserId)
	if err != nil {
		return err
	}
	return nil
}

func InsertSession(db *sql.DB, sessionID string, userID int64, expiration time.Time) error {
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

func GetComments(postID int, db *sql.DB, userId int) ([]utils.Comment, error) {
	query := `
	SELECT comments.id, comments.content, comments.created_at, users.username, comments.like_count , comments.dislike_count FROM comments
    INNER JOIN users ON comments.user_id = users.id
    WHERE comments.post_id = ?
	ORDER BY comments.created_at DESC;
	`
	rows, err := utils.QueryRows(db, query, postID)
	if err != nil {
		return nil, errors.New(err.Error() + "here 1")
	}
	defer rows.Close()

	var comments []utils.Comment
	for rows.Next() {
		var comment utils.Comment
		err := rows.Scan(&comment.Comment_id, &comment.Content, &comment.Created_at, &comment.User_name, &comment.LikeCount, &comment.DislikeCount)
		if err != nil {
			return nil, errors.New(err.Error() + "here 2")
		}
		comment.Post_id = postID
		if userId <= 0 {
			comment.Clicked, comment.DisClicked = false, false
		} else {
			comment.Clicked, comment.DisClicked = isLiked(db, userId, comment.Comment_id, "comment")
		}

		comments = append(comments, comment)
	}
	if err = rows.Err(); err != nil {
		return nil, errors.New(err.Error() + "here 3")
	}
	return comments, nil
}

func GetCategoryContentIds(db *sql.DB, categoryId string, userId int) ([]int, error) {
	if categoryId == "1" || categoryId == "2" {
		rows, err := utils.QueryRows(db, `SELECT post_id FROM post_categories WHERE category_id=? AND user_id=?`, categoryId, userId)
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
	rows, err := utils.QueryRows(db, "SELECT post_id FROM post_categories WHERE category_id=?", categoryId)
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
	row, err := utils.QueryRow(db, "SELECT username FROM users WHERE id = ?", id)
	if err != nil {
		return "", err
	}
	err = row.Scan(&name)
	if err != nil {
		return "", err
	}
	return name, nil
}

func GetPostCategories(db *sql.DB, PostId int, userId int) ([]string, error) {
	query := `SELECT categories.name 
	FROM post_categories
	JOIN categories ON categories.id = post_categories.category_id
	WHERE (post_categories.category_id = 1 OR post_categories.category_id = 2) 
	AND post_categories.user_id = ? 
	AND post_categories.post_id = ?
	UNION
	SELECT categories.name 
	FROM post_categories
	JOIN categories ON categories.id = post_categories.category_id
	WHERE post_categories.category_id != 1 
	AND post_categories.category_id != 2 
	AND post_categories.post_id = ?;
`
	rows, err := utils.QueryRows(db, query, userId, PostId, PostId)
	if err != nil {
		return nil, err
	}

	var categories []string

	for rows.Next() {
		var category string
		err := rows.Scan(&category)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}

func LinkPostWithCategory(transaction *sql.Tx, categories []string, postId int64, userId int) error {
	for _, category := range categories {
		stmt, err := transaction.Prepare(`INSERT INTO post_categories(user_id, post_id, category_id) VALUES(?, ?, ?);`)
		if err != nil {
			return err
		}
		defer stmt.Close()
		tmp, err := strconv.Atoi(category)
		if err != nil {
			return err
		}
		_, err = stmt.Exec(userId, postId, tmp)
		if err != nil {
			return err
		}
	}
	return nil
}
