package main

import (
	"fmt"
	"os"
	"time"

	"Diplom_Go/pkg/api"
	"Diplom_Go/pkg/db"
	"Diplom_Go/pkg/server"

	_ "modernc.org/sqlite"
)

func main() {

	//step3 - function
	//for test
	now := time.Now()
	//str1, err := api.NextDate(now, "20240113", "d 7")// работает
	//str1, err := api.NextDate(now, "20240113", "y")
	//str1, err := api.NextDate(now, "20240113", "w 7")
	//str1, err := api.NextDate(now, "20240113", "m -2 2,8")
	//str1, err := api.NextDate(now, "20240102", "y")

	/*var task db.Task
	task.Date = "20250716"
	task.Comment = ""
	task.Repeat = "d 1"
	err := api.CheckDate(&task)*/

	str1, err := api.NextDate(now, "20250716", "d 1")
	/*var task db.Task
	task.Comment = ""
	task.Date = "20240108"
	task.Repeat = "d 10"
	task.Title = "Уроки"

	id, err := db.AddTask(&task)
	fmt.Println("    4.1, id = " + string(id))*/

	fmt.Println("str1=" + str1)

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

	//step1 - Web server
	todo_port := os.Getenv("TODO_PORT")
	port := ":7540"
	if len(todo_port) > 0 {
		port = todo_port
	}
	webDir := "./web"

	err = server.Run(port, webDir)
	if err != nil {
		fmt.Println(err)
	}

}
