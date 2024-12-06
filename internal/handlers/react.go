package handlers

import (
	"database/sql"
	"fmt"
	"forum/internal/utils"
	"net/http"
	"os"
)

func InsertOrUpdateReactionHandler(w http.ResponseWriter, r *http.Request, db *sql.DB, userID int) {
	r.Header.Add("content-type", "application/json")
	reactionType := r.URL.Query().Get("type")
	targetType := r.URL.Query().Get("target")
	id := r.URL.Query().Get("target_id")

	fmt.Println(reactionType, targetType, id)

	if targetType != "" && reactionType != "" {
		var insertQuery string
		switch targetType {
		case "post":
			insertQuery = `INSERT INTO reactions (reaction_type, user_id, post_id, target_type) 
			VALUES (?, ?, ?, ?)
			ON CONFLICT (user_id, post_id, target_type) DO UPDATE SET reaction_type = EXCLUDED.reaction_type;
			`
		case "comment":
			insertQuery = `INSERT INTO reactions (reaction_type, user_id, comment_id, target_type) 
			VALUES (?, ?, ?, ?)
			ON CONFLICT (user_id, comment_id,target_type) DO UPDATE SET reaction_type = EXCLUDED.reaction_type ;
			`
		}
		_, err := db.Exec(insertQuery, reactionType, userID, id, targetType)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write(nil)
			return
		}
		w.WriteHeader(200)
		w.Write(nil)
	} else {
		utils.RespondWithJSON(w, http.StatusBadRequest, `{"error": "Bad Request"}`)
		return
	}
}

func DeleteReactionHandler(w http.ResponseWriter, r *http.Request, db *sql.DB, userID int) {
	r.Header.Add("content-type", "application/json")
	targetType := r.URL.Query().Get("target")
	id := r.URL.Query().Get("target_id")

	if targetType != "" && id != "" {
		var deleteQuery string
		switch targetType {
		case "post":
			deleteQuery = `DELETE FROM reactions WHERE user_id = ? AND post_id = ? AND target_type = ? `
		case "comment":
			deleteQuery = `DELETE FROM reactions WHERE user_id = ? AND comment_id = ? AND target_type = ? `
		}
		_, err := db.Exec(deleteQuery, userID, id, targetType)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(nil)
			return
		}
		w.WriteHeader(200)
		w.Write(nil)
	} else {
		utils.RespondWithJSON(w, http.StatusBadRequest, `{"error": "Bad Request1"}`)
		return
	}
}

func GetReactionsHandler(w http.ResponseWriter, r *http.Request, db *sql.DB, userId int) {
	//target := r.URL.Query().Get("target")
	targetID := r.URL.Query().Get("target_id")

	// Prepare query to get likes and dislikes for the target (post or comment)
	var likedBy, dislikedBy []int
	var userReaction string

	// Query for liked users
	likeQuery := `
		SELECT user_id
		FROM reactions
		WHERE post_id = ? AND reaction_type = 'like';`

	// Query for disliked users
	dislikeQuery := `
		SELECT user_id
		FROM reactions
		WHERE post_id = ? AND reaction_type = 'dislike';`

	// Query for user reaction to a post
	userReactionQuery := `
		SELECT reaction_type
		FROM reactions
		WHERE user_id = ? AND post_id = ?;`

	// Execute like query
	rows, err := db.Query(likeQuery, targetID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var userId int
		if err := rows.Scan(&userId); err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		likedBy = append(likedBy, userId)
	}

	// Execute dislike query
	rows, err = db.Query(dislikeQuery, targetID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var userId int
		if err := rows.Scan(&userId); err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		dislikedBy = append(dislikedBy, userId)
	}

	// Execute user reaction query
	err = db.QueryRow(userReactionQuery, userId, targetID).Scan(&userReaction)
	if err != nil && err != sql.ErrNoRows {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Prepare the response
	response := utils.Reaction{
		LikedBy:      likedBy,
		DislikedBy:   dislikedBy,
		UserReaction: userReaction,
	}

	// Send response
	utils.RespondWithJSON(w, http.StatusOK, response)
}

// updateDislikeCount inserts a dislike for a post in the reactions table
// func updateLikeCount(db *sql.DB, postID int, userID int, reaction string, target_type string) error {
// 	// Check if the user has already disliked this post
// 	var exists int
// 	var query string
// 	if target_type == "comment" {
// 		query = `SELECT EXISTS(SELECT 1 FROM reactions WHERE comment_id = ? AND user_id = ? AND target_type = ?)`
// 	} else if target_type == "post" {
// 		query = `SELECT EXISTS(SELECT 1 FROM reactions WHERE post_id = ? AND user_id = ? AND target_type = ?)`
// 	}
// 	err := db.QueryRow(query, postID, userID, target_type).Scan(&exists) // does react exist ?
// 	if err != nil {
// 		return fmt.Errorf("failed to check existing reaction 1: %w", err)
// 	}
// 	if exists == 1 { // user changed the reaction
// 		var same int
// 		if target_type == "comment" {
// 			query = `SELECT EXISTS(SELECT 1 FROM reactions WHERE comment_id = ? AND user_id = ? AND type = ? ANd target_type = ? )`
// 		} else if target_type == "post" {
// 			query = `SELECT EXISTS(SELECT 1 FROM reactions WHERE post_id = ? AND user_id = ? AND type = ? ANd target_type = ? )`
// 		}
// 		err := db.QueryRow(query, postID, userID, reaction, target_type).Scan(&same)
// 		if err != nil {
// 			return fmt.Errorf("failed to check existing reaction 2: %w", err)
// 		}
// 		// if liked/dislied twise delete
// 		// If the user has already reacted the same way (liked/disliked twice), delete the reaction
// 		if same == 1 { // user removes the reaction
// 			if reaction == "like" {
// 				err := removeLike(db, postID, target_type)
// 				if err != nil {
// 					return err
// 				}
// 			} else if reaction == "dislike" {
// 				err := removeDislike(db, postID, target_type)
// 				if err != nil {
// 					return err
// 				}
// 			}
// 			var deleteQuery string
// 			if target_type == "comment" {
// 				deleteQuery = `DELETE FROM reactions WHERE user_id = ? AND comment_id = ? AND target_type = ?`
// 			} else if target_type == "post" {
// 				deleteQuery = `DELETE FROM reactions WHERE user_id = ? AND post_id = ? AND target_type = ?`
// 			}
// 			_, err := db.Exec(deleteQuery, userID, postID, target_type)
// 			if err != nil {
// 				return fmt.Errorf("failed to delete reaction: %w", err)
// 			}
// 			return nil
// 		} else {
// 			// here
// 			if reaction == "like" {
// 				removeDislike(db, postID, target_type)
// 			} else if reaction == "dislike" {
// 				removeLike(db, postID, target_type)
// 			}
// 			//Update the reaction if an entry exists
// 			var updateQuery string
// 			if target_type == "comment" {
// 				updateQuery = `UPDATE reactions SET type = ? WHERE comment_id = ? AND user_id = ? AND target_type = ?`
// 			} else if target_type == "post" {
// 				updateQuery = `UPDATE reactions SET type = ? WHERE post_id = ? AND user_id = ? AND target_type = ?`
// 			}
// 			_, err = db.Exec(updateQuery, reaction, postID, userID, target_type)
// 			if err != nil {
// 				return fmt.Errorf("failed to update reaction 4: %w", err)
// 			}
// 			if reaction == "like" {
// 				err := addLike(db, postID, target_type)
// 				if err != nil {
// 					return err
// 				}
// 			} else if reaction == "dislike" {
// 				err := addDislike(db, postID, target_type)
// 				if err != nil {
// 					return err
// 				}
// 			}
// 		}
// 		// /react17/like/comment
// 	} else { // first time user reation
// 		// Insert a new reaction if no entry exists
// 		var insertQuery string
// 		if target_type == "comment" {
// 			insertQuery = `INSERT INTO reactions (user_id, comment_id, target_type, type) VALUES (?, ?, ?, ?)`
// 		} else if target_type == "post" {
// 			insertQuery = `INSERT INTO reactions (user_id, post_id, target_type, type) VALUES (?, ?, ?, ?)`
// 		}
// 		_, err = db.Exec(insertQuery, userID, postID, target_type, reaction)
// 		if err != nil {
// 			return fmt.Errorf("failed to insert reaction 5: %w", err)
// 		}
// 		// Update the like or dislike count
// 		if reaction == "like" {
// 			err := addLike(db, postID, target_type)
// 			if err != nil {
// 				return err
// 			}
// 		} else if reaction == "dislike" {
// 			err := addDislike(db, postID, target_type)
// 			if err != nil {
// 				return err
// 			}
// 		}
// 	}

// 	return nil
// }

// func addLike(db *sql.DB, postID int, target_type string) error {
// 	target_type = target_type + "s"
// 	updateQuery := fmt.Sprintf("UPDATE %s SET like_count = like_count + 1 WHERE id = ?", target_type)
// 	_, err := db.Exec(updateQuery, postID)
// 	if err != nil {
// 		return fmt.Errorf("failed to update reaction: %w", err)
// 	}
// 	return nil
// }

// func addDislike(db *sql.DB, postID int, target_type string) error {
// 	target_type = target_type + "s"
// 	updateQuery := fmt.Sprintf("UPDATE %s SET dislike_count = dislike_count + 1 WHERE id = ?", target_type)
// 	_, err := db.Exec(updateQuery, postID)
// 	if err != nil {
// 		return fmt.Errorf("failed to update reaction: %w", err)
// 	}
// 	return nil
// }

// func removeLike(db *sql.DB, postID int, target_type string) error {
// 	target_type = target_type + "s"
// 	//	updateQuery := `UPDATE ? SET like_count = like_count - 1 WHERE id = ?`
// 	updateQuery := fmt.Sprintf("UPDATE %s SET like_count = like_count - 1 WHERE id = ?", target_type)
// 	_, err := db.Exec(updateQuery, postID)
// 	if err != nil {
// 		return fmt.Errorf("failed to update reaction: %w", err)
// 	}
// 	return nil
// }

// func removeDislike(db *sql.DB, postID int, target_type string) error {
// 	target_type = target_type + "s"
// 	updateQuery := fmt.Sprintf("UPDATE %s SET dislike_count = dislike_count - 1 WHERE id = ?", target_type)
// 	_, err := db.Exec(updateQuery, postID)
// 	if err != nil {
// 		return fmt.Errorf("failed to update reaction: %w", err)
// 	}
// 	return nil
// }
