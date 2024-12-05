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
	id := r.URL.Query().Get("target_id")
	user := r.URL.Query().Get("is_own_react")
	var query string
	var rows *sql.Rows
	switch target {
	case "post":
		if user == "false" || user == "" {
			query = `SELECT reactions.user_id, reactions.post_id, reaction_type.name, reaction_id 
			FROM reactions 
			JOIN reaction_type 
			ON reactions.type_id = reaction_type.reaction_id 
			WHERE reactions.post_id = ?;`
			var err error
			rows, err = utils.QueryRows(db, query, id)
			if err != nil {
				utils.RespondWithJSON(w, http.StatusInternalServerError, `{"error": "Internal Server Error"}`)
				return
			}
			defer rows.Close()
		} else if user == "true" {
			query = `SELECT reactions.user_id, reactions.post_id, reaction_type.name, reaction_id 
			FROM reactions 
			JOIN reaction_type 
			ON reactions.type_id = reaction_type.reaction_id 
			WHERE reactions.post_id = ? AND reactions.user_id = ?;`
			row, err := utils.QueryRow(db, query, id, userId)
			if err != nil && err != sql.ErrNoRows {
				fmt.Println(err)
				utils.RespondWithJSON(w, http.StatusInternalServerError, `{"error": "Internal Server Error"}`)
				return
			}
			reaction := utils.Reaction{}
			err = row.Scan(&reaction.UserId, &reaction.TargetId, &reaction.Name, &reaction.ReactionId)
			if err != nil && err != sql.ErrNoRows {
				fmt.Println(err)

				utils.RespondWithJSON(w, http.StatusInternalServerError, `{"error": "Internal Server Error"}`)
				return
			}
			utils.RespondWithJSON(w, http.StatusOK, reaction)
			return
		} else {
			utils.RespondWithJSON(w, http.StatusBadRequest, `{"error": "Bad Request"}`)
			return
		}
	case "comment":
		if user == "false" || user == "" {
			query = `SELECT reactions.user_id, reactions.comment_id, reaction_type.name, reaction_id 
			FROM reactions 
			JOIN reaction_type 
			ON reactions.type_id = reaction_type.reaction_id 
			WHERE reactions.comment_id = ?;`
			var err error
			rows, err = utils.QueryRows(db, query, id)
			if err != nil {
				utils.RespondWithJSON(w, http.StatusInternalServerError, `{"error": "Internal Server Error"}`)
				return
			}
			defer rows.Close()
		} else if user == "true" {
			query = `SELECT reactions.user_id, reactions.comment_id, reaction_type.name, reaction_id 
			FROM reactions 
			JOIN reaction_type 
			ON reactions.type_id = reaction_type.reaction_id 
			WHERE reactions.comment_id = ? AND reactions.user_id = ?;`
			row, err := utils.QueryRow(db, query, id, userId)
			if err != nil {
				utils.RespondWithJSON(w, http.StatusInternalServerError, `{"error": "Internal Server Error"}`)
				return
			}
			reaction := utils.Reaction{}
			err = row.Scan(&reaction.UserId, &reaction.TargetId, &reaction.Name, &reaction.ReactionId)
			if err != nil && err != sql.ErrNoRows {
				utils.RespondWithJSON(w, http.StatusInternalServerError, `{"error": "Internal Server Error1"}`)
				return
			}
			utils.RespondWithJSON(w, http.StatusOK, reaction)
			return
		} else {
			utils.RespondWithJSON(w, http.StatusBadRequest, `{"error": "Bad Request"}`)
			return
		}
	}
	var reactions []utils.Reaction
	for rows.Next() {
		var reaction utils.Reaction
		err := rows.Scan(&reaction.UserId, &reaction.TargetId, &reaction.Name, &reaction.ReactionId)
		if err != nil {
			utils.RespondWithJSON(w, http.StatusInternalServerError, `{"error": "Internal Server Error"}`)
			return
		}
		reactions = append(reactions, reaction)
	}
	utils.RespondWithJSON(w, http.StatusOK, reactions)
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
