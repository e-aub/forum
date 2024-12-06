package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"forum/internal/utils"
	"net/http"
	"os"
)

func InsertOrUpdateReactionHandler(w http.ResponseWriter, r *http.Request, db *sql.DB, userID int) {
	r.Header.Add("content-type", "application/json")
	reactionType := r.URL.Query().Get("type")
	target := r.URL.Query().Get("target")
	id := r.URL.Query().Get("target_id")
	fmt.Println(reactionType, target, id)

	if target != "" && reactionType != "" {
		var insertQuery string
		switch target {
		case "post":
			insertQuery = `INSERT INTO reactions (type_id, user_id, post_id) 
			VALUES (?, ?, ?)
			ON CONFLICT (user_id, post_id) DO UPDATE SET type_id = EXCLUDED.type_id;
			`
		case "comment":
			insertQuery = `INSERT INTO reactions (type_id, user_id, comment_id) 
			VALUES (?, ?, ?)
			ON CONFLICT (user_id, comment_id) DO UPDATE SET type_id = EXCLUDED.type_id;
			`
		}
		_, err := db.Exec(insertQuery, reactionType, userID, id)
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
	target := r.URL.Query().Get("target")
	id := r.URL.Query().Get("target_id")

	if target != "" && id != "" {
		var deleteQuery string
		switch target {
		case "post":
			deleteQuery = `DELETE FROM reactions WHERE user_id = ? AND post_id = ?`
		case "comment":
			deleteQuery = `DELETE FROM reactions WHERE user_id = ? AND comment_id = ?`
		}
		_, err := db.Exec(deleteQuery, userID, id)
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
	target := r.URL.Query().Get("target")
	targetID := r.URL.Query().Get("target_id")

	// Validate required parameters
	if target == "" || targetID == "" {
		utils.RespondWithJSON(w, http.StatusBadRequest, `{"error": "Target and target_id are required"}`)
		return
	}

	// Build the query based on the target
	var query string
	switch target {
	case "post":
		query = `
        SELECT 
            COALESCE(GROUP_CONCAT(CASE WHEN reaction_type.name = 'like' THEN reactions.user_id ELSE NULL END), '') AS liked_by,
            COALESCE(GROUP_CONCAT(CASE WHEN reaction_type.name = 'dislike' THEN reactions.user_id ELSE NULL END), '') AS disliked_by,
            COALESCE(CASE 
                WHEN reactions.user_id = ? AND reaction_type.name = 'like' THEN 'liked'
                WHEN reactions.user_id = ? AND reaction_type.name = 'dislike' THEN 'disliked'
                ELSE NULL
            END, '') AS user_reaction
        FROM reactions
        JOIN reaction_type ON reactions.type_id = reaction_type.reaction_id
        WHERE reactions.post_id = ?
        GROUP BY reactions.post_id;`

	case "comment":
		query = `
        SELECT 
            COALESCE(GROUP_CONCAT(CASE WHEN reaction_type.name = 'like' THEN reactions.user_id ELSE NULL END), '') AS liked_by,
            COALESCE(GROUP_CONCAT(CASE WHEN reaction_type.name = 'dislike' THEN reactions.user_id ELSE NULL END), '') AS disliked_by,
            COALESCE(CASE 
                WHEN reactions.user_id = ? AND reaction_type.name = 'like' THEN 'liked'
                WHEN reactions.user_id = ? AND reaction_type.name = 'dislike' THEN 'disliked'
                ELSE NULL
            END, '') AS user_reaction
        FROM reactions
        JOIN reaction_type ON reactions.type_id = reaction_type.reaction_id
        WHERE reactions.comment_id = ?
        GROUP BY reactions.comment_id;`

	default:
		utils.RespondWithJSON(w, http.StatusBadRequest, `{"error": "Invalid target"}`)
		return
	}

	// Execute the query
	rows, err := db.Query(query, userId, userId, targetID)
	if err != nil {
		fmt.Println(err)
		utils.RespondWithJSON(w, http.StatusInternalServerError, `{"error": "Failed to fetch reactions"}`)
		return
	}
	defer rows.Close()

	// Parse results
	var likedBy, dislikedBy []int
	var userReaction string
	if rows.Next() {
		var likedByJSON, dislikedByJSON string
		if err := rows.Scan(&likedByJSON, &dislikedByJSON, &userReaction); err != nil {
			utils.RespondWithJSON(w, http.StatusInternalServerError, `{"error": "Failed to parse reaction data"}`)
			return
		}
		// Convert JSON arrays into slices
		if likedByJSON != "" {
			if err := json.Unmarshal([]byte(likedByJSON), &likedBy); err != nil {
				fmt.Println(err)
				utils.RespondWithJSON(w, http.StatusInternalServerError, `{"error": "Failed to parse liked_by data"}`)
				return
			}
		} else {
			likedBy = []int{} // If the JSON is empty, initialize as an empty slice
		}

		if dislikedByJSON != "" {
			if err := json.Unmarshal([]byte(dislikedByJSON), &dislikedBy); err != nil {
				fmt.Println(err)
				utils.RespondWithJSON(w, http.StatusInternalServerError, `{"error": "Failed to parse disliked_by data"}`)
				return
			}
		} else {
			dislikedBy = []int{} // If the JSON is empty, initialize as an empty slice
		}
	} else {
		// No reactions found
		utils.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
			"liked_by":      []int{},
			"disliked_by":   []int{},
			"user_reaction": nil,
		})
		return
	}
	var response utils.Reaction
	// Prepare and send the response
	response.LikedBy = likedBy
	response.DislikedBy = dislikedBy
	response.UserReaction = userReaction

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
