package postgre

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var db *sql.DB

func init() {
	dbUrl := os.Getenv("DB_URL")
	if len(dbUrl) == 0 {
		log.Fatalln("DB_URL is not defined")
	}
	var err error
	db, err = sql.Open("postgres", dbUrl)
	if err != nil {
		panic("Failed to connect to DB")
	}

	err = db.Ping()
	if err != nil {
		panic("Failed to connect to DB")
	}

	_, err = db.Exec(`
	CREATE TABLE users (
		Id 		 TEXT PRIMARY KEY,
		Username TEXT UNIQUE,
		Password TEXT,
		Token 	 TEXT
	);`)

	if err != nil {
		log.Println(err)
	}

	_, err = db.Exec(`
	CREATE TABLE files (
		Id 	   TEXT PRIMARY KEY,
		Data   BYTEA
	);`)

	if err != nil {
		log.Println(err)
	}

	_, err = db.Exec(`
	CREATE TABLE advertisements (
		Id 	   TEXT PRIMARY KEY,
		UserId TEXT,
		FileId   TEXT,
		Url    TEXT,
		CreatedAt TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
		FOREIGN KEY (UserId) REFERENCES users (Id),
		FOREIGN KEY (FileId) REFERENCES files (Id)
	);`)

	if err != nil {
		log.Println(err)
	}

}

func Close() error {
	return db.Close()
}
