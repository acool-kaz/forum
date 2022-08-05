package storage

import (
	"database/sql"
	"fmt"
)

const userTable = `CREATE TABLE IF NOT EXISTS user (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	email TEXT UNIQUE,
	username TEXT UNIQUE,
	hashPassword TEXT,
	posts INT DEFAULT 0,
	likes INT DEFAULT 0,
	comments INT DEFAULT 0,
	session_token TEXT DEFAULT NULL,
	expiresAt DATETIME DEFAULT NULL
);`

const postTable = `CREATE TABLE IF NOT EXISTS post (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	creater TEXT,
	title TEXT,
	description TEXT,
	created_at DATE DEFAULT (datetime('now','localtime')),
	likes INT DEFAULT 0,
	dislikes INT DEFAULT 0,
	FOREIGN KEY (creater) REFERENCES user(username)
);`

const postCategoryTable = `CREATE TABLE IF NOT EXISTS post_category (
	postId INTEGER,
	category TEXT,
	FOREIGN KEY (postId) REFERENCES post(id) ON DELETE CASCADE
);`

const commentTable = `CREATE TABLE IF NOT EXISTS comment (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	postId INTEGER,
	creater TEXT,
	text TEXT,
	likes INT DEFAULT 0,
	dislikes INT DEFAULT 0,
	FOREIGN KEY (postId) REFERENCES post(id) ON DELETE CASCADE
);`

const likesTable = `CREATE TABLE IF NOT EXISTS likes (
	username TEXT,
	postId INTEGER DEFAULT NULL,
	commentId INTEGER DEFAULT NULL,
	FOREIGN KEY (postId) REFERENCES post(id) ON DELETE CASCADE,
	FOREIGN KEY (commentId) REFERENCES comment(id) ON DELETE CASCADE
);`

const dislikeTable = `CREATE TABLE IF NOT EXISTS dislikes (
	username TEXT,
	postId INTEGER DEFAULT NULL,
	commentId INTEGER DEFAULT NULL,
	FOREIGN KEY (postId) REFERENCES post(id) ON DELETE CASCADE,
	FOREIGN KEY (commentId) REFERENCES comment(id) ON DELETE CASCADE
);`

var tables = []string{userTable, postTable, postCategoryTable, commentTable, likesTable, dislikeTable}

func InitDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "forum.db")
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	_, err = db.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		return nil, err
	}
	return db, nil
}

func CreateTables(db *sql.DB) error {
	var err error
	for _, table := range tables {
		_, err = db.Exec(table)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
	return err
}
