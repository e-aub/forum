package handlers

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func Creat_Post_Table() {
	file, err := sql.Open("sqlite3", "Data_base/Posts.db")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	_, err = file.Exec(`CREATE TABLE IF NOT EXISTS Posts(
	PostId integer Not null primary key,
	UserId integer not null References Users (userid),
	Title text,
	Content text,
	created_at text,
	updated_at text )`)
	if err != nil {
		panic(err)
	}
}

func New_Post_db(*Posts) error {
	
}
