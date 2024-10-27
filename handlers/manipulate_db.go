package handlers

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func Creat_TAble_Posts() {
	file, err := sql.Open("sqlite3", "Data_base/data.db")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	_, err = file.Exec(`CREATE TABLE IF NOT EXISTS Posts(
						PostId  INTEGER PRIMARY KEY  AUTOINCREMENT,
						UserId  INTEGER not null,
						Title Text,
						Content Text,
						Created_at Text,
						Updated_at Text,
						FOREIGN KEY(UserId) REFERENCES Users(UserId))`)
	if err != nil {
		panic(err)
	}
}

func Insert_Post(p *Posts) {
	file, err := sql.Open("sqlite3", "Data_base/data.db")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	statement, err := file.Prepare(`INSERT INTO Posts(Title,Content,Created_at , UserId) Values (?,?,?,?)`)
	if err != nil {
		panic(err)
	}
	_, err = statement.Exec(p.Title, p.Content, p.Created_At, p.UserId)
	if err != nil {
		panic(err)
	}
}

func Update_Post(p *Posts){
	
}