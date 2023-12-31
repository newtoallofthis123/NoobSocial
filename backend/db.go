package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// DbInstance A Struct to hold the Database Instance
// Any interaction with the database through the application
// will be done through this struct only
type DbInstance struct {
	db *sql.DB
}

// NewDbInstance : Constructs the connection string from
// the environment variables and returns a pointer to the
// DbInstance struct with a connection to the database.
func NewDbInstance() *DbInstance {
	env := GetEnv()
	var connStr = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", env.Host, env.Port, env.User, env.Password, env.Db)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	return &DbInstance{
		db: db,
	}
}

// preStart : The Database Initialization function
func (pq *DbInstance) preStart() {
	query := `
		CREATE TABLE IF NOT EXISTS users (
			username VARCHAR(50) PRIMARY KEY NOT NULL,
		    full_name VARCHAR(225) NOT NULL,
			bio TEXT NOT NULL,
			email TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS sessions (
		    session_id VARCHAR(50) PRIMARY KEY,
		    username VARCHAR(50) REFERENCES users(username) NOT NULL,
		    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		    valid BOOLEAN DEFAULT TRUE
		);

		CREATE TABLE IF NOT EXISTS posts (
			post_id VARCHAR(50) PRIMARY KEY,
			username VARCHAR(50) REFERENCES users(username) NOT NULL,
			content TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS tags (
			tag_id VARCHAR(50) PRIMARY KEY,
			tag_name VARCHAR(50) UNIQUE NOT NULL
		);

		CREATE TABLE IF NOT EXISTS likes (
			post_id VARCHAR(50) REFERENCES posts(post_id) NOT NULL,
			username VARCHAR(50) REFERENCES users(username) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (post_id, username)
		);

		CREATE TABLE IF NOT EXISTS comments (
			comment_id VARCHAR(50) PRIMARY KEY,
			post_id VARCHAR(50) REFERENCES posts(post_id) NOT NULL,
			username VARCHAR(50) REFERENCES users(username) NOT NULL,
			content TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS follows (
			follower VARCHAR(50) REFERENCES users(username) NOT NULL,
			followee VARCHAR(50) REFERENCES users(username) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (follower, followee)
		);
	`

	_, err := pq.db.Exec(query)
	if err != nil {
		panic(err)
	}
}
