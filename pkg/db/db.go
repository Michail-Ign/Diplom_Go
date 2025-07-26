package db

import (
	"database/sql"
	"log"
	"os"
)

const schema = `CREATE TABLE scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    title VARCHAR(128) NOT NULL DEFAULT "",
    comment TEXT NOT NULL DEFAULT "", 
    repeat VARCHAR(128) NOT NULL DEFAULT ""
);
CREATE INDEX scheduler_date ON scheduler (date);`

var db *sql.DB
var nameDB string

func Init(dbFile string) error {

	install := false
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		install = true
	}

	nameDB = dbFile
	log.Println("Init parametres nameDB: ", nameDB)

	db, err := Open()
	if err != nil {
		return err
	}
	defer db.Close()

	// если install равен true, после открытия БД требуется выполнить
	// sql-запрос с CREATE TABLE и CREATE INDEX
	if install {
		_, err = db.Exec(schema)
		if err != nil {
			return err
		}
	}
	log.Println("Base init success")
	return nil
}

func Open() (*sql.DB, error) {

	todo_dbfile := os.Getenv("TODO_DBFILE")
	if len(todo_dbfile) > 0 {
		nameDB = todo_dbfile
	}
	log.Println("Open db = ", nameDB)

	db_r, err := sql.Open("sqlite", nameDB)
	if err != nil {

		log.Println("Error Open db = ", err)
		return nil, err
	}

	return db_r, nil
}
