package handlers

import (
	"database/sql"
	"log"

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

func Update_Post(p *Posts) {
	file, err := sql.Open("sqlite3", "Data_base/data.db")
	if err != nil {
		log.Fatal(err)
	}
	statement, err := file.Prepare(`UPDATE Posts
	SET Title=?,
	Content = ?,
	Updated_at = ?
	WHERE
	PostId = ?`)
	if err != nil {
		log.Fatal(err)
	}
	_, err = statement.Exec(p.Title, p.Content, p.Updated_At, p.PostId)
	if err != nil {
		log.Fatal(err)
	}
}

func Delete_Post(p *Posts) {
	file, err := sql.Open("sqlite3", "Data_base/data.db")
	if err != nil {
		log.Fatal(err)
	}
	statement, err := file.Prepare(`DELETE FROM Posts WHERE PostId = ?`)
	if err != nil {
		log.Fatal(err)
	}
	_, err = statement.Exec(p.PostId)
	if err != nil {
		log.Fatal(err)
	}
}

func Read_Post(id int) *Posts {
	file, err := sql.Open("sqlite3", "Data_base/data.db")
	if err != nil {
		log.Fatal(err)
	}
	query := `SELECT * FROM Posts WHERE PostId = ?`
	row := file.QueryRow(query, id)
	Post := &Posts{}
	_ = row.Scan(&Post.PostId, &Post.UserId, &Post.Title, &Post.Content, &Post.Created_At, &Post.Updated_At)
	return Post
}

func Get_Last() int {
	file, err := sql.Open("sqlite3", "Data_base/data.db")
	if err != nil {
		log.Fatal(err)
	}
	query := `SELECT MAX(PostId) FROM Posts `
	row := file.QueryRow(query)
	result := 0
	_ = row.Scan(&result)
	return result
}
