package server

import (
	"Diplom_Go/pkg/api"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func Run(Port string, webDir string) error {

	todoPort := os.Getenv("TODO_PORT")
	if len(todoPort) > 0 {
		Port = todoPort
	}
	if !strings.HasPrefix(Port, ":") {
		Port = ":" + Port
	}
	log.Println("set Port = ", Port)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Если запрос к корню, возвращаем index.html
		if r.URL.Path == "/" {
			http.ServeFile(w, r, filepath.Join(webDir, "index.html"))
			return
		}
		// Иначе обслуживаем статические файлы из webDir
		http.ServeFile(w, r, filepath.Join(webDir, r.URL.Path))
	})

	//
	api.Init()

	// Запуск сервера на порту Port
	log.Println("Сервер запущен на http://localhost" + Port)
	err := http.ListenAndServe(Port, nil)
	if err != nil {
		log.Println("Ошибка сервера:")
		log.Fatal(err)

		return err
	}
	return nil

}
