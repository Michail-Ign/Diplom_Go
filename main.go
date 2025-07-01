package main

import (
	"database/sql"
	"os"

	_ "modernc.org/sqlite"
)

var db *sql.DB

func main() {

	dbFile := "scheduler.db"
	_, err := os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	}
	// если install равен true, после открытия БД требуется выполнить
	// sql-запрос с CREATE TABLE и CREATE INDEX

	db.Init("scheduler.db")
}
