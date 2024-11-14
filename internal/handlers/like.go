package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Response struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

func ReactHandler(db *sql.DB, userID int, isUser bool, w http.ResponseWriter, r *http.Request) {
	if !isUser {
		http.Error(w, "Unauthorized: user not registered", http.StatusUnauthorized) // 401 status code
		return
	}
	// Parse the postId from the URL
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) != 6 || pathParts[2] != "react" || (pathParts[4] != "like" && pathParts[4] != "dislike") || (pathParts[5] != "post" && pathParts[5] != "comment") {
		http.Error(w, "invalid url", http.StatusBadRequest)
		return
	}
	reaction := pathParts[4]
	target_type := pathParts[5]
	// Convert postId to integer
	postID, err := strconv.Atoi(pathParts[3])
	if err != nil || postID < 0 {
		http.Error(w, "Invalid post ID 0", http.StatusBadRequest)
		return
	}

	// Process dislike logic here, e.g., update in database
	err = updateLikeCount(db, postID, userID, reaction, target_type) // Define this function to handle your database update
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send JSON response back to the client
	w.Header().Set("Content-Type", "application/json")
}

// updateDislikeCount inserts a dislike for a post in the likes table
func updateLikeCount(db *sql.DB, postID int, userID int, reaction string, target_type string) error {
	// Check if the user has already disliked this post
	var exists int
	var query string
	if target_type == "comment" {
		query = `SELECT EXISTS(SELECT 1 FROM likes WHERE comment_id = ? AND user_id = ? AND target_type = ?)`
	} else if target_type == "post" {
		query = `SELECT EXISTS(SELECT 1 FROM likes WHERE post_id = ? AND user_id = ? AND target_type = ?)`
	}
	err := db.QueryRow(query, postID, userID, target_type).Scan(&exists) // does react exist ?
	if err != nil {
		return fmt.Errorf("failed to check existing reaction 1: %w", err)
	}
	if exists == 1 { // user changed the reaction
		var same int
		if target_type == "comment" {
			query = `SELECT EXISTS(SELECT 1 FROM likes WHERE comment_id = ? AND user_id = ? AND type = ? ANd target_type = ? )`
		} else if target_type == "post" {
			query = `SELECT EXISTS(SELECT 1 FROM likes WHERE post_id = ? AND user_id = ? AND type = ? ANd target_type = ? )`
		}
		err := db.QueryRow(query, postID, userID, reaction, target_type).Scan(&same)
		if err != nil {
			return fmt.Errorf("failed to check existing reaction 2: %w", err)
		}
		// if liked/dislied twise delete
		// If the user has already reacted the same way (liked/disliked twice), delete the reaction
		if same == 1 { // user removes the reaction
			if reaction == "like" {
				err := removeLike(db, postID, target_type)
				if err != nil {
					return err
				}
			} else if reaction == "dislike" {
				err := removeDislike(db, postID, target_type)
				if err != nil {
					return err
				}
			}
			var deleteQuery string
			if target_type == "comment" {
				deleteQuery = `DELETE FROM likes WHERE user_id = ? AND comment_id = ? AND target_type = ?`
			} else if target_type == "post" {
				deleteQuery = `DELETE FROM likes WHERE user_id = ? AND post_id = ? AND target_type = ?`
			}
			_, err := db.Exec(deleteQuery, userID, postID, target_type)
			if err != nil {
				return fmt.Errorf("failed to delete reaction: %w", err)
			}
			return nil
		} else {
			// here
			if reaction == "like" {
				removeDislike(db, postID, target_type)
			} else if reaction == "dislike" {
				removeLike(db, postID, target_type)
			}
			//Update the reaction if an entry exists
			var updateQuery string
			if target_type == "comment" {
				updateQuery = `UPDATE likes SET type = ? WHERE comment_id = ? AND user_id = ? AND target_type = ?`
			} else if target_type == "post" {
				updateQuery = `UPDATE likes SET type = ? WHERE post_id = ? AND user_id = ? AND target_type = ?`
			}
			_, err = db.Exec(updateQuery, reaction, postID, userID, target_type)
			if err != nil {
				return fmt.Errorf("failed to update reaction 4: %w", err)
			}
			if reaction == "like" {
				err := addLike(db, postID, target_type)
				if err != nil {
					return err
				}
			} else if reaction == "dislike" {
				err := addDislike(db, postID, target_type)
				if err != nil {
					return err
				}
			}
		}
		// /api/react/17/like/comment
	} else { // first time user reation
		// Insert a new reaction if no entry exists
		var insertQuery string
		if target_type == "comment" {
			insertQuery = `INSERT INTO likes (user_id, comment_id, target_type, type) VALUES (?, ?, ?, ?)`
		} else if target_type == "post" {
			insertQuery = `INSERT INTO likes (user_id, post_id, target_type, type) VALUES (?, ?, ?, ?)`
		}
		_, err = db.Exec(insertQuery, userID, postID, target_type, reaction)
		if err != nil {
			return fmt.Errorf("failed to insert reaction 5: %w", err)
		}
		// Update the like or dislike count
		if reaction == "like" {
			err := addLike(db, postID, target_type)
			if err != nil {
				return err
			}
		} else if reaction == "dislike" {
			err := addDislike(db, postID, target_type)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func addLike(db *sql.DB, postID int, target_type string) error {
	target_type = target_type + "s"
	updateQuery := fmt.Sprintf("UPDATE %s SET like_count = like_count + 1 WHERE id = ?", target_type)
	_, err := db.Exec(updateQuery, postID)
	if err != nil {
		return fmt.Errorf("failed to update reaction: %w", err)
	}
	return nil
}

func addDislike(db *sql.DB, postID int, target_type string) error {
	target_type = target_type + "s"
	updateQuery := fmt.Sprintf("UPDATE %s SET dislike_count = dislike_count + 1 WHERE id = ?", target_type)
	_, err := db.Exec(updateQuery, postID)
	if err != nil {
		return fmt.Errorf("failed to update reaction: %w", err)
	}
	return nil
}

func removeLike(db *sql.DB, postID int, target_type string) error {
	target_type = target_type + "s"
	//	updateQuery := `UPDATE ? SET like_count = like_count - 1 WHERE id = ?`
	updateQuery := fmt.Sprintf("UPDATE %s SET like_count = like_count - 1 WHERE id = ?", target_type)
	_, err := db.Exec(updateQuery, postID)
	if err != nil {
		return fmt.Errorf("failed to update reaction: %w", err)
	}
	return nil
}

func removeDislike(db *sql.DB, postID int, target_type string) error {
	target_type = target_type + "s"
	updateQuery := fmt.Sprintf("UPDATE %s SET dislike_count = dislike_count - 1 WHERE id = ?", target_type)
	_, err := db.Exec(updateQuery, postID)
	if err != nil {
		return fmt.Errorf("failed to update reaction: %w", err)
	}
	return nil
}
