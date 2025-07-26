package main

import (
	"Diplom_Go/pkg/db"
	"Diplom_Go/pkg/server"
	"log"

	_ "modernc.org/sqlite"
)


func main() {

	//step2 - DataBase SQLite
	dbFile := "scheduler.db"

	err := db.Init(dbFile)
	if err != nil {
		log.Println(err)
	}

	//step1 - Web server
	port := ":7540"
	webDir := "./web"

	err = server.Run(port, webDir)
	if err != nil {
		log.Println(err)
	}

}
