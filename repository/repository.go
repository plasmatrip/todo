package repository

import (
	"database/sql"
	"log"
	"os"

	"todo/configs"
)

// var db *sql.DB

var schema = `
	CREATE TABLE scheduler (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date CHAR(8) NOT NULL DEFAULT "",
		title VARCHAR(256) NOT NULL DEFAULT "",
		comment TEXT NOT NULL DEFAULT "",
		repeat VARCHAR(128) NOT NULL DEFAULT ""
	);
	CREATE INDEX schedule_date ON scheduler (date);
`

type Repository struct {
	db *sql.DB
}

func NewToDo() *Repository {
	return &Repository{db: open()}
}

func (d *Repository) Close() {
	d.db.Close()
}

func open() *sql.DB {
	var err error
	db, err := sql.Open("sqlite", configs.DBDir+configs.DBFile)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	if _, err := os.Stat(configs.DBDir + configs.DBFile); err != nil {
		if _, err := os.Stat(configs.DBDir); err != nil {
			if err := os.Mkdir(configs.DBDir, 0755); err != nil {
				log.Fatal(err)
			}
		}
		_, err = db.Exec(schema)
		if err != nil {
			log.Panic(err)
		}
	}
	return db
}
