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
	if len(pathParts) != 5 || pathParts[2] != "react" || (pathParts[4] != "like" && pathParts[4] != "dislike") {
		http.Error(w, "invalid url", http.StatusBadRequest)
		return
	}
	reaction := pathParts[4]
	// Convert postId to integer
	postID, err := strconv.Atoi(pathParts[3])
	if err != nil || postID < 0 {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	// Process dislike logic here, e.g., update in database
	err = updateLikeCount(db, postID, userID, reaction) // Define this function to handle your database update
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send JSON response back to the client
	w.Header().Set("Content-Type", "application/json")
}

// updateDislikeCount inserts a dislike for a post in the likes table
func updateLikeCount(db *sql.DB, postID int, userID int, reaction string) error {

	// Check if the user has already disliked this post
	var exists int
	query := `SELECT EXISTS(SELECT 1 FROM likes WHERE post_id = ? AND user_id = ?)`
	err := db.QueryRow(query, postID, userID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check existing reaction 1: %w", err)
	}
	if exists == 1 { // user changed the reaction
		var same int
		query = `SELECT EXISTS(SELECT 1 FROM likes WHERE post_id = ? AND user_id = ? AND type = ? )`
		err := db.QueryRow(query, postID, userID, reaction).Scan(&same)
		if err != nil {
			return fmt.Errorf("failed to check existing reaction 2: %w", err)
		}
		// if liked/dislied twise delete
		// If the user has already reacted the same way (liked/disliked twice), delete the reaction
		if same == 1 { // user removes the reaction
			if reaction == "like" {
				err := removeLike(db, postID)
				if err != nil {
					return err
				}
			} else if reaction == "dislike" {
				err := removeDislike(db, postID)
				if err != nil {
					return err
				}
			}
			deleteQuery := `DELETE FROM likes WHERE user_id = ? AND post_id = ? AND target_type = 'post'`
			_, err := db.Exec(deleteQuery, userID, postID)
			if err != nil {
				return fmt.Errorf("failed to delete reaction: %w", err)
			}
			return nil
		} else {
			// here
			if reaction == "like" {
				removeDislike(db, postID)
			} else if reaction == "dislike" {
				removeLike(db, postID)
			}
			//Update the reaction if an entry exists
			updateQuery := `UPDATE likes SET type = ? WHERE post_id = ? AND user_id = ? AND target_type = 'post'`
			_, err = db.Exec(updateQuery, reaction, postID, userID)
			if err != nil {
				return fmt.Errorf("failed to update reaction 4: %w", err)
			}
			if reaction == "like" {
				err := addLike(db, postID)
				if err != nil {
					return err
				}
			} else if reaction == "dislike" {
				err := addDislike(db, postID)
				if err != nil {
					return err
				}
			}
		}
	} else { // first time user reation
		// Insert a new reaction if no entry exists
		insertQuery := `INSERT INTO likes (user_id, post_id, target_type, type) VALUES (?, ?, 'post', ?)`
		_, err = db.Exec(insertQuery, userID, postID, reaction)
		if err != nil {
			return fmt.Errorf("failed to insert reaction 5: %w", err)
		}
		// Update the like or dislike count
		if reaction == "like" {
			err := addLike(db, postID)
			if err != nil {
				return err
			}
		} else if reaction == "dislike" {
			err := addDislike(db, postID)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func addLike(db *sql.DB, postID int) error {
	updateQuery := `UPDATE posts SET like_count = like_count + 1 WHERE id = ?`
	_, err := db.Exec(updateQuery, postID)
	if err != nil {
		return fmt.Errorf("failed to update reaction: %w", err)
	}
	return nil
}

func addDislike(db *sql.DB, postID int) error {
	updateQuery := `UPDATE posts SET dislike_count = dislike_count + 1 WHERE id = ?`
	_, err := db.Exec(updateQuery, postID)
	if err != nil {
		return fmt.Errorf("failed to update reaction: %w", err)
	}
	return nil
}

func removeLike(db *sql.DB, postID int) error {
	updateQuery := `UPDATE posts SET like_count = like_count - 1 WHERE id = ?`
	_, err := db.Exec(updateQuery, postID)
	if err != nil {
		return fmt.Errorf("failed to update reaction: %w", err)
	}
	return nil
}

func removeDislike(db *sql.DB, postID int) error {
	updateQuery := `UPDATE posts SET dislike_count = dislike_count - 1 WHERE id = ?`
	_, err := db.Exec(updateQuery, postID)
	if err != nil {
		return fmt.Errorf("failed to update reaction: %w", err)
	}
	return nil
}
