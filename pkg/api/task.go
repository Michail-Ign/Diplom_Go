package api

import (
	"fmt"
	"net/http"
)

func taskHandler(w http.ResponseWriter, req *http.Request) {

	fmt.Println("--вызов taskHandler")
	//
	fmt.Println("  req.Method = " + req.Method)

	switch req.Method {
	// обработка других методов будет добавлена на следующих шагах
	case http.MethodPost:
		addTaskHandler(w, req)
	}
}
