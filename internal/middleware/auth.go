package auth

import "database/sql"

func IsUserRegistered(db *sql.DB, email, username string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = ? OR username = ?);`
	err := db.QueryRow(query, email, username).Scan(&exists)
	return exists, err
}

// Register a new user in the database
func RegisterUser(db *sql.DB, username, email, password string) error {
	insertQuery := `INSERT INTO users (username, email, password) VALUES (?, ?, ?);`
	_, err := db.Exec(insertQuery, username, email, password)
	return err
}

func GetPasswordByUsername(db *sql.DB, username string) (string, error) {
	var password string
	query := `SELECT password FROM users WHERE username = ?;`
	err := db.QueryRow(query, username).Scan(&password)
	if err != nil {
		return "", err

	}
	return password, nil
}
