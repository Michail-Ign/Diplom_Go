package main

import (
	"fmt"
	"os"

	"Diplom_Go/pkg/db"
	"Diplom_Go/pkg/server"

	_ "modernc.org/sqlite"
)

func main() {

	//step1 - Web server
	Port := ":7540"
	webDir := "./web"

	err := server.Run(Port, webDir)
	if err != nil {
		fmt.Println(err)
	}

	//step2 - DataBase SQLite
	todo_dbfile := os.Getenv("TODO_DBFILE")
	dbFile := "scheduler.db"

	if len(todo_dbfile) > 0 {
		dbFile = todo_dbfile
	}

	//_, err := os.Stat(dbFile)

	//var install bool
	//if err != nil {
	//	install = true
	//}

	err = db.Init(dbFile)
	if err != nil {
		fmt.Println(err)
	}
}
