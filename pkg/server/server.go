package server

import (
	"Diplom_Go/pkg/api"
	"log"
	"net/http"
	"path/filepath"
)

func Run(Port string, webDir string) error {

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
		log.Fatal(err)
		return err
	}
	return nil

}
