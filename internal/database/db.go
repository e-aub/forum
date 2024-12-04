package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	database "forum/internal/database/models"
	utils "forum/internal/utils"
)

type Conn_db struct {
	Db *sql.DB
	Mu sync.Mutex
}

func CreateDatabase(dbPath string) *Conn_db {
	conn := new(Conn_db)
	var err error
	conn.Db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("%sError opening database:%s%s\n", utils.Colors["red"], err.Error(), utils.Colors["reset"])
	}
	// Verify the connection
	if err = conn.Db.Ping(); err != nil {
		log.Fatalf("%sError accessing database: %s%s\n", utils.Colors["red"], err.Error(), utils.Colors["reset"])
	} else {
		fmt.Printf("%sDatabase created/opened successfully!%s\n", utils.Colors["green"], utils.Colors["reset"])
	}

	_, err = conn.Db.Exec(`PRAGMA foreign_keys=ON;`)
	if err != nil {
		log.Fatalf("%sError enabling foreign keys: %s%s\n", utils.Colors["red"], err.Error(), utils.Colors["reset"])
	}
	return conn
}

func (conn *Conn_db) CreateTables() {
	_, err := conn.Db.Exec(database.UsersTable + database.SessionsTable + database.ReactionTable + database.ReactionsTypeTable +
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

func (conn *Conn_db) InsertPost(p *utils.Post) (int64, error) {
	conn.Mu.Lock()
	defer conn.Mu.Unlock()

	transaction, err := conn.Db.Begin()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error starting transaction:", err)
		return 0, err
	}
	stmt, err := transaction.Prepare(`INSERT INTO posts(user_id ,title,content) Values (?,?,?);`)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error Adding post:", err)
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(p.UserId, p.Title, p.Content)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error Adding post:", err)
		return 0, err
	}

	lastPostID, err := result.LastInsertId()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error in assigning category to post", err)
		return 0, err
	}

	err = LinkPostWithCategory(transaction, p.Categories, lastPostID, p.UserId)
	if err != nil {
		return 0, err
	}
	err = transaction.Commit()
	if err != nil {
		fmt.Fprintln(os.Stderr, "transaction aborted")
		return 0, err
	}
	defer transaction.Rollback()

	return lastPostID, nil
}

func (conn *Conn_db) ReadPost(Id ...int) (*utils.Post, error) {
	conn.Mu.Lock()
	defer conn.Mu.Unlock()

	Post := &utils.Post{}
	var postId int

	switch len(Id) {
	case 0:
		// Get the ID of the last post created
		query := `SELECT MAX(id) FROM posts`
		row, err := utils.QueryRow(conn.Db, query)
		if err != nil {
			return nil, err
		}
		err = row.Scan(&Post.PostId)
		if err != nil {
			return nil, err
		}
		return Post, nil

	case 1:
		postId = Id[0]
	default:
		return nil, fmt.Errorf("too many arguments")
	}

	statement, err := conn.Db.Begin()
	if err != nil {
		return nil, err
	}
	defer statement.Rollback()
	stmt, err := statement.Prepare(`
        SELECT p.id, p.user_id, p.title, p.content, p.created_at, u.username 
        FROM posts p
        JOIN users u ON p.user_id = u.id
        WHERE p.id = ?
    `)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(postId)
	err = row.Scan(&Post.PostId, &Post.UserId, &Post.Title, &Post.Content, &Post.CreatedAt, &Post.UserName)
	if err != nil {
		return nil, err
	}

	err = statement.Commit()
	if err != nil {
		return nil, err
	}

	return Post, nil
}

func (conn *Conn_db) Get_session(ses string) (int, error) {
	conn.Mu.Lock()
	defer conn.Mu.Unlock()
	var sessionid int
	tx, err := conn.Db.Begin()
	if err != nil {
		return 0, err
	}
	statement, err := tx.Prepare(`SELECT user_id FROM sessions WHERE session_id = ?`)
	if err != nil {
		return 0, err
	}
	err = statement.QueryRow(ses).Scan(&sessionid)
	if err != nil {
		return 0, err
	}
	return sessionid, tx.Commit()
}

func (conn *Conn_db) InsertSession(userData *utils.User) error {
	conn.Mu.Lock()
	defer conn.Mu.Unlock()
	tx, err := conn.Db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec("INSERT INTO sessions (session_id, user_id, expires_at) VALUES (?, ?, ?)", userData.SessionId, userData.UserId, userData.Expiration)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (conn *Conn_db) CreateComment(c *utils.Comment) error {
	conn.Mu.Lock()
	defer conn.Mu.Unlock()

	if c.Content == "" || c.Post_id == 0 || c.User_id == 0 {
		return errors.New("comment DATA issue")
	}

	tx, err := conn.Db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
        INSERT INTO comments (user_id, post_id, content, created_at) 
        VALUES (?, ?, ?, ?)
    `)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(c.User_id, c.Post_id, c.Content, c.Created_at)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	c.Comment_id = int(id)

	return tx.Commit()
}

func (conn *Conn_db) GetComments(postID int, userId int) ([]utils.Comment, error) {
	conn.Mu.Lock()
	defer conn.Mu.Unlock()

	tx, err := conn.Db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
		SELECT comments.id, comments.content, comments.created_at, users.username
		FROM comments
		INNER JOIN users ON comments.user_id = users.id
		WHERE comments.post_id = ?
		ORDER BY comments.created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(postID)
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

	return comments, tx.Commit()
}

func (conn *Conn_db) GetCategoryContentIds(categoryId string) ([]int, error) {
	conn.Mu.Lock()
	defer conn.Mu.Unlock()
	rows, err := utils.QueryRows(conn.Db, "SELECT post_id FROM post_categories WHERE category_id=?", categoryId)
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

func (conn *Conn_db) GetPostCategories(PostId int, userId int) ([]string, error) {
	conn.Mu.Lock()
	defer conn.Mu.Unlock()

	query := `
	SELECT categories.name 
	FROM post_categories
	JOIN categories ON categories.id = post_categories.category_id
	AND post_categories.post_id = ?;
`
	rows, err := utils.QueryRows(conn.Db, query, PostId)
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
