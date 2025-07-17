package db

import (
	"database/sql"
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

func Init(dbFile string) error {

	install := false
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		install = true
	}

	db, err := Open(dbFile)
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
	return nil
}

func Open(dbFile string) (*sql.DB, error) {

	todo_dbfile := os.Getenv("TODO_DBFILE")
	if len(todo_dbfile) > 0 {
		dbFile = todo_dbfile
	}

	db_r, err := sql.Open("sqlite", dbFile)
	if err != nil {

		return nil, err
	}

	return db_r, nil
}
